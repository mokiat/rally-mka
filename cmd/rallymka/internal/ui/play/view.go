package play

import (
	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/lacking/game/ecs"
	"github.com/mokiat/lacking/game/graphics"
	"github.com/mokiat/lacking/game/physics"
	"github.com/mokiat/lacking/game/physics/solver"
	"github.com/mokiat/lacking/resource"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/ecscomp"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/ecssys"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/game"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/global"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/scene"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/scene/car"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/store"
)

const (
	anchorDistance = 6.0
	cameraDistance = 12.0

	carMaxSteeringAngle  = 30
	carFrontAcceleration = 155
	carRearAcceleration  = 160

	// FIXME: Currently, too much front brakes cause the car
	// to straighten. This is due to there being more pressure
	// on the outer wheel which causes it to brake more and turn
	// the car to neutral orientation.
	carFrontDeceleration = 250
	carRearDeceleration  = 180

	suspensionEnabled      = true
	suspensionStart        = float32(-0.25)
	suspensionEnd          = float32(-0.6)
	suspensionWidth        = float32(1.0)
	suspensionLength       = float32(0.25)
	suspensionFrequencyHz  = float32(3.0)
	suspensionDampingRatio = float32(1.0)
)

type ViewData struct {
	GameData *scene.Data
}

var View = co.Connect(co.ShallowCached(co.Define(func(props co.Properties) co.Instance {
	var context global.Context
	co.InjectContext(&context)

	var (
		data      ViewData
		lifecycle *playLifecycle
	)
	props.InjectData(&data)

	co.UseState(func() interface{} {
		return &playLifecycle{
			gameController: context.GameController,
			gameData:       data.GameData,
		}
	}).Inject(&lifecycle)

	co.Once(func() {
		lifecycle.init()
	})

	co.Defer(func() {
		lifecycle.destroy()
	})

	return co.New(mat.Container, func() {
		co.WithData(mat.ContainerData{
			Layout: mat.NewAnchorLayout(mat.AnchorLayoutSettings{}),
		})
	})

})), co.ConnectMapping{

	Data: func(props co.Properties) interface{} {
		var appStore store.Application
		co.InjectStore(&appStore)

		return ViewData{
			GameData: appStore.GameData,
		}
	},
})

type playLifecycle struct {
	gameController *game.Controller
	gameData       *scene.Data

	gfxScene     graphics.Scene
	physicsScene *physics.Scene
	ecsScene     *ecs.Scene

	renderSystem      *ecssys.Renderer
	vehicleSystem     *ecssys.VehicleSystem
	cameraStandSystem *ecssys.CameraStandSystem

	camera graphics.Camera
}

func (h *playLifecycle) init() {
	h.gfxScene = h.gameController.GFXScene()
	h.physicsScene = h.gameController.PhysicsScene()
	h.ecsScene = h.gameController.ECSScene()

	h.renderSystem = h.gameController.RenderSystem()
	h.vehicleSystem = h.gameController.VehicleSystem()
	h.cameraStandSystem = h.gameController.CameraStandSystem()

	h.camera = h.gameController.Camera()
	h.camera.SetPosition(sprec.NewVec3(0.0, 0.0, 0.0))
	h.camera.SetFoVMode(graphics.FoVModeHorizontalPlus)
	h.camera.SetFoV(sprec.Degrees(66))
	h.camera.SetAutoExposure(true)
	h.camera.SetExposure(1.0)
	h.camera.SetAutoFocus(false)

	h.setupLevel(h.gameData.Level)
}

func (h *playLifecycle) destroy() {
	h.gameData.Dismiss()
}

func (h *playLifecycle) setupLevel(level *resource.Level) {
	h.gfxScene.Sky().SetBackgroundColor(sprec.NewVec3(0.0, 0.3, 0.8))
	h.gfxScene.Sky().SetSkybox(level.SkyboxTexture.GFXTexture)

	ambientLight := h.gfxScene.CreateAmbientLight()
	ambientLight.SetReflectionTexture(level.AmbientReflectionTexture.GFXTexture)
	ambientLight.SetRefractionTexture(level.AmbientRefractionTexture.GFXTexture)

	sunLight := h.gfxScene.CreateDirectionalLight()
	sunLight.SetRotation(sprec.QuatProd(
		sprec.RotationQuat(sprec.Degrees(225), sprec.BasisYVec3()),
		sprec.RotationQuat(sprec.Degrees(-45), sprec.BasisXVec3()),
	))
	sunLight.SetIntensity(sprec.NewVec3(1.2, 1.2, 1.2))

	for _, staticMesh := range level.StaticMeshes {
		h.gfxScene.CreateMesh(staticMesh.GFXMeshTemplate)
	}

	for _, collisionMesh := range level.CollisionMeshes {
		body := h.physicsScene.CreateBody()
		body.SetPosition(sprec.ZeroVec3())
		body.SetOrientation(sprec.IdentityQuat())
		body.SetStatic(true)
		body.SetRestitutionCoefficient(1.0)
		body.SetCollisionShapes([]physics.CollisionShape{collisionMesh})
	}

	var createModelMesh func(matrix sprec.Mat4, node *resource.Node)
	createModelMesh = func(matrix sprec.Mat4, node *resource.Node) {
		modelMatrix := sprec.Mat4Prod(matrix, node.Matrix)

		gfxMesh := h.gfxScene.CreateMesh(node.Mesh.GFXMeshTemplate)
		gfxMesh.SetPosition(modelMatrix.Translation())
		// TODO: SetRotation
		// TODO: SetScale

		for _, child := range node.Children {
			createModelMesh(modelMatrix, child)
		}
	}

	for _, staticEntity := range level.StaticEntities {
		for _, node := range staticEntity.Model.Nodes {
			createModelMesh(staticEntity.Matrix, node)
		}
	}

	carModel := h.gameData.CarModel
	targetEntity := h.setupCarDemo(carModel, sprec.NewVec3(0.0, 3.0, 0.0))
	standTarget := targetEntity
	standEntity := h.ecsScene.CreateEntity()
	ecscomp.SetCameraStand(standEntity, &ecscomp.CameraStand{
		Target:         standTarget,
		Camera:         h.camera,
		AnchorPosition: sprec.Vec3Sum(ecscomp.GetPhysics(standTarget).Body.Position(), sprec.NewVec3(0.0, 0.0, -cameraDistance)),
		AnchorDistance: anchorDistance,
		CameraDistance: cameraDistance,
	})
}

func (h *playLifecycle) setupCarDemo(model *resource.Model, position sprec.Vec3) *ecs.Entity {
	chasis := car.Chassis(model).
		WithName("chasis").
		WithPosition(position).
		Build(h.ecsScene, h.gfxScene, h.physicsScene)
	chasisPhysics := ecscomp.GetPhysics(chasis)

	flWheelRelativePosition := sprec.NewVec3(suspensionWidth, suspensionStart-suspensionLength, 1.07)
	flWheel := car.Wheel(model, car.FrontLeftWheelLocation).
		WithName("front-left-wheel").
		WithPosition(sprec.Vec3Sum(position, flWheelRelativePosition)).
		Build(h.ecsScene, h.gfxScene, h.physicsScene)
	flWheelPhysics := ecscomp.GetPhysics(flWheel)
	h.physicsScene.CreateDoubleBodyConstraint(chasisPhysics.Body, flWheelPhysics.Body,
		solver.NewMatchTranslation().
			SetPrimaryAnchor(flWheelRelativePosition).
			SetIgnoreY(suspensionEnabled),
	)
	h.physicsScene.CreateDoubleBodyConstraint(chasisPhysics.Body, flWheelPhysics.Body, &solver.LimitTranslation{
		MaxY: suspensionStart,
		MinY: suspensionEnd,
	})
	flRotation := solver.NewMatchAxis().
		SetPrimaryAxis(sprec.BasisXVec3()).
		SetSecondaryAxis(sprec.BasisXVec3())
	h.physicsScene.CreateDoubleBodyConstraint(chasisPhysics.Body, flWheelPhysics.Body, flRotation)
	h.physicsScene.CreateDoubleBodyConstraint(chasisPhysics.Body, flWheelPhysics.Body, &solver.Coilover{
		PrimaryAnchor: flWheelRelativePosition,
		FrequencyHz:   suspensionFrequencyHz,
		DampingRatio:  suspensionDampingRatio,
	})

	frWheelRelativePosition := sprec.NewVec3(-suspensionWidth, suspensionStart-suspensionLength, 1.07)
	frWheel := car.Wheel(model, car.FrontRightWheelLocation).
		WithName("front-right-wheel").
		WithPosition(sprec.Vec3Sum(position, frWheelRelativePosition)).
		Build(h.ecsScene, h.gfxScene, h.physicsScene)
	frWheelPhysics := ecscomp.GetPhysics(frWheel)
	h.physicsScene.CreateDoubleBodyConstraint(chasisPhysics.Body, frWheelPhysics.Body,
		solver.NewMatchTranslation().
			SetPrimaryAnchor(frWheelRelativePosition).
			SetIgnoreY(suspensionEnabled),
	)
	h.physicsScene.CreateDoubleBodyConstraint(chasisPhysics.Body, frWheelPhysics.Body, &solver.LimitTranslation{
		MaxY: suspensionStart,
		MinY: suspensionEnd,
	})
	frRotation := solver.NewMatchAxis().
		SetPrimaryAxis(sprec.BasisXVec3()).
		SetSecondaryAxis(sprec.BasisXVec3())
	h.physicsScene.CreateDoubleBodyConstraint(chasisPhysics.Body, frWheelPhysics.Body, frRotation)
	h.physicsScene.CreateDoubleBodyConstraint(chasisPhysics.Body, frWheelPhysics.Body, &solver.Coilover{
		PrimaryAnchor: frWheelRelativePosition,
		FrequencyHz:   suspensionFrequencyHz,
		DampingRatio:  suspensionDampingRatio,
	})

	blWheelRelativePosition := sprec.NewVec3(suspensionWidth, suspensionStart-suspensionLength, -1.56)
	blWheel := car.Wheel(model, car.BackLeftWheelLocation).
		WithName("back-left-wheel").
		WithPosition(sprec.Vec3Sum(position, blWheelRelativePosition)).
		Build(h.ecsScene, h.gfxScene, h.physicsScene)
	blWheelPhysics := ecscomp.GetPhysics(blWheel)
	h.physicsScene.CreateDoubleBodyConstraint(chasisPhysics.Body, blWheelPhysics.Body,
		solver.NewMatchTranslation().
			SetPrimaryAnchor(blWheelRelativePosition).
			SetIgnoreY(suspensionEnabled),
	)
	h.physicsScene.CreateDoubleBodyConstraint(chasisPhysics.Body, blWheelPhysics.Body, &solver.LimitTranslation{
		MaxY: suspensionStart,
		MinY: suspensionEnd,
	})
	h.physicsScene.CreateDoubleBodyConstraint(chasisPhysics.Body, blWheelPhysics.Body,
		solver.NewMatchAxis().
			SetPrimaryAxis(sprec.BasisXVec3()).
			SetSecondaryAxis(sprec.BasisXVec3()),
	)
	h.physicsScene.CreateDoubleBodyConstraint(chasisPhysics.Body, blWheelPhysics.Body, &solver.Coilover{
		PrimaryAnchor: blWheelRelativePosition,
		FrequencyHz:   suspensionFrequencyHz,
		DampingRatio:  suspensionDampingRatio,
	})

	brWheelRelativePosition := sprec.NewVec3(-suspensionWidth, suspensionStart-suspensionLength, -1.56)
	brWheel := car.Wheel(model, car.BackRightWheelLocation).
		WithName("back-right-wheel").
		WithPosition(sprec.Vec3Sum(position, brWheelRelativePosition)).
		Build(h.ecsScene, h.gfxScene, h.physicsScene)
	brWheelPhysics := ecscomp.GetPhysics(brWheel)
	h.physicsScene.CreateDoubleBodyConstraint(chasisPhysics.Body, brWheelPhysics.Body,
		solver.NewMatchTranslation().
			SetPrimaryAnchor(brWheelRelativePosition).
			SetIgnoreY(suspensionEnabled),
	)
	h.physicsScene.CreateDoubleBodyConstraint(chasisPhysics.Body, brWheelPhysics.Body, &solver.LimitTranslation{
		MaxY: suspensionStart,
		MinY: suspensionEnd,
	})
	h.physicsScene.CreateDoubleBodyConstraint(chasisPhysics.Body, brWheelPhysics.Body, solver.NewMatchAxis().
		SetPrimaryAxis(sprec.BasisXVec3()).
		SetSecondaryAxis(sprec.BasisXVec3()),
	)
	h.physicsScene.CreateDoubleBodyConstraint(chasisPhysics.Body, brWheelPhysics.Body, &solver.Coilover{
		PrimaryAnchor: brWheelRelativePosition,
		FrequencyHz:   suspensionFrequencyHz,
		DampingRatio:  suspensionDampingRatio,
	})

	car := h.ecsScene.CreateEntity()
	ecscomp.SetVehicle(car, &ecscomp.Vehicle{
		MaxSteeringAngle: sprec.Degrees(carMaxSteeringAngle),
		SteeringAngle:    sprec.Degrees(0.0),
		Acceleration:     0.0,
		Deceleration:     0.0,
		Chassis: &ecscomp.Chassis{
			Body: chasisPhysics.Body,
		},
		Wheels: []*ecscomp.Wheel{
			{
				Body:                 flWheelPhysics.Body,
				RotationConstraint:   flRotation,
				AccelerationVelocity: carFrontAcceleration,
				DecelerationVelocity: carFrontDeceleration,
			},
			{
				Body:                 frWheelPhysics.Body,
				RotationConstraint:   frRotation,
				AccelerationVelocity: carFrontAcceleration,
				DecelerationVelocity: carFrontDeceleration,
			},
			{
				Body:                 blWheelPhysics.Body,
				AccelerationVelocity: carRearAcceleration,
				DecelerationVelocity: carRearDeceleration,
			},
			{
				Body:                 brWheelPhysics.Body,
				AccelerationVelocity: carRearAcceleration,
				DecelerationVelocity: carRearDeceleration,
			},
		},
	})
	ecscomp.SetPlayerControl(car, &ecscomp.PlayerControl{})

	return chasis
}
