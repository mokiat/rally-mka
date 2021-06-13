package simulation

import (
	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/lacking/app"
	"github.com/mokiat/lacking/async"
	"github.com/mokiat/lacking/game/ecs"
	"github.com/mokiat/lacking/game/graphics"
	"github.com/mokiat/lacking/game/physics"
	"github.com/mokiat/lacking/game/physics/solver"
	"github.com/mokiat/lacking/resource"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/ecscomp"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/ecssys"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/scene"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/scene/car"
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

func NewView(gfxEngine graphics.Engine, physicsEngine *physics.Engine, ecsEngine *ecs.Engine, registry *resource.Registry, gfxWorker *async.Worker) *View {
	return &View{
		gameData: scene.NewData(registry, gfxWorker),

		gfxEngine:     gfxEngine,
		physicsEngine: physicsEngine,
		ecsEngine:     ecsEngine,
	}
}

type View struct {
	gameData *scene.Data

	gfxEngine     graphics.Engine
	physicsEngine *physics.Engine
	ecsEngine     *ecs.Engine

	gfxScene     graphics.Scene
	physicsScene *physics.Scene
	ecsScene     *ecs.Scene

	renderSystem      *ecssys.Renderer
	vehicleSystem     *ecssys.VehicleSystem
	cameraStandSystem *ecssys.CameraStandSystem

	camera graphics.Camera

	freezeFrame bool
}

func (v *View) Load(window app.Window, cb func()) {
	v.gameData.Request().OnSuccess(func(interface{}) {
		window.Schedule(func() error {
			cb()
			return nil
		})
	})
}

func (v *View) Unload(window app.Window) {
	v.gameData.Dismiss()
}

func (v *View) Open(window app.Window) {
	v.gfxScene = v.gfxEngine.CreateScene()
	v.physicsScene = v.physicsEngine.CreateScene(0.015)
	v.ecsScene = v.ecsEngine.CreateScene()

	v.renderSystem = ecssys.NewRenderer(v.ecsScene)
	v.vehicleSystem = ecssys.NewVehicleSystem(v.ecsScene)
	v.cameraStandSystem = ecssys.NewCameraStandSystem(v.ecsScene)

	v.camera = v.gfxScene.CreateCamera()
	v.camera.SetPosition(sprec.NewVec3(0.0, 0.0, 0.0))
	v.camera.SetFoVMode(graphics.FoVModeHorizontalPlus)
	v.camera.SetFoV(sprec.Degrees(66))
	v.camera.SetAutoExposure(true)
	v.camera.SetExposure(1.0)
	v.camera.SetAutoFocus(false)

	v.setupLevel(v.gameData.Level)
}

func (v *View) setupLevel(level *resource.Level) {
	v.gfxScene.Sky().SetBackgroundColor(sprec.NewVec3(0.0, 0.3, 0.8))
	v.gfxScene.Sky().SetSkybox(level.SkyboxTexture.GFXTexture)

	ambientLight := v.gfxScene.CreateAmbientLight()
	ambientLight.SetReflectionTexture(level.AmbientReflectionTexture.GFXTexture)
	ambientLight.SetRefractionTexture(level.AmbientRefractionTexture.GFXTexture)

	sunLight := v.gfxScene.CreateDirectionalLight()
	sunLight.SetRotation(sprec.QuatProd(
		sprec.RotationQuat(sprec.Degrees(225), sprec.BasisYVec3()),
		sprec.RotationQuat(sprec.Degrees(-45), sprec.BasisXVec3()),
	))
	sunLight.SetIntensity(sprec.NewVec3(1.2, 1.2, 1.2))

	for _, staticMesh := range level.StaticMeshes {
		v.gfxScene.CreateMesh(staticMesh.GFXMeshTemplate)
	}

	for _, collisionMesh := range level.CollisionMeshes {
		body := v.physicsScene.CreateBody()
		body.SetPosition(sprec.ZeroVec3())
		body.SetOrientation(sprec.IdentityQuat())
		body.SetStatic(true)
		body.SetRestitutionCoefficient(1.0)
		body.SetCollisionShapes([]physics.CollisionShape{collisionMesh})
	}

	var createModelMesh func(matrix sprec.Mat4, node *resource.Node)
	createModelMesh = func(matrix sprec.Mat4, node *resource.Node) {
		modelMatrix := sprec.Mat4Prod(matrix, node.Matrix)

		gfxMesh := v.gfxScene.CreateMesh(node.Mesh.GFXMeshTemplate)
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

	carModel := v.gameData.CarModel
	targetEntity := v.setupCarDemo(carModel, sprec.NewVec3(0.0, 3.0, 0.0))
	standTarget := targetEntity
	standEntity := v.ecsScene.CreateEntity()
	ecscomp.SetCameraStand(standEntity, &ecscomp.CameraStand{
		Target:         standTarget,
		Camera:         v.camera,
		AnchorPosition: sprec.Vec3Sum(ecscomp.GetPhysics(standTarget).Body.Position(), sprec.NewVec3(0.0, 0.0, -cameraDistance)),
		AnchorDistance: anchorDistance,
		CameraDistance: cameraDistance,
	})
}

func (v *View) setupCarDemo(model *resource.Model, position sprec.Vec3) *ecs.Entity {
	chasis := car.Chassis(model).
		WithName("chasis").
		WithPosition(position).
		Build(v.ecsScene, v.gfxScene, v.physicsScene)
	chasisPhysics := ecscomp.GetPhysics(chasis)

	flWheelRelativePosition := sprec.NewVec3(suspensionWidth, suspensionStart-suspensionLength, 1.07)
	flWheel := car.Wheel(model, car.FrontLeftWheelLocation).
		WithName("front-left-wheel").
		WithPosition(sprec.Vec3Sum(position, flWheelRelativePosition)).
		Build(v.ecsScene, v.gfxScene, v.physicsScene)
	flWheelPhysics := ecscomp.GetPhysics(flWheel)
	v.physicsScene.CreateDoubleBodyConstraint(chasisPhysics.Body, flWheelPhysics.Body,
		solver.NewMatchTranslation().
			SetPrimaryAnchor(flWheelRelativePosition).
			SetIgnoreY(suspensionEnabled),
	)
	v.physicsScene.CreateDoubleBodyConstraint(chasisPhysics.Body, flWheelPhysics.Body, &solver.LimitTranslation{
		MaxY: suspensionStart,
		MinY: suspensionEnd,
	})
	flRotation := solver.NewMatchAxis().
		SetPrimaryAxis(sprec.BasisXVec3()).
		SetSecondaryAxis(sprec.BasisXVec3())
	v.physicsScene.CreateDoubleBodyConstraint(chasisPhysics.Body, flWheelPhysics.Body, flRotation)
	v.physicsScene.CreateDoubleBodyConstraint(chasisPhysics.Body, flWheelPhysics.Body, &solver.Coilover{
		PrimaryAnchor: flWheelRelativePosition,
		FrequencyHz:   suspensionFrequencyHz,
		DampingRatio:  suspensionDampingRatio,
	})

	frWheelRelativePosition := sprec.NewVec3(-suspensionWidth, suspensionStart-suspensionLength, 1.07)
	frWheel := car.Wheel(model, car.FrontRightWheelLocation).
		WithName("front-right-wheel").
		WithPosition(sprec.Vec3Sum(position, frWheelRelativePosition)).
		Build(v.ecsScene, v.gfxScene, v.physicsScene)
	frWheelPhysics := ecscomp.GetPhysics(frWheel)
	v.physicsScene.CreateDoubleBodyConstraint(chasisPhysics.Body, frWheelPhysics.Body,
		solver.NewMatchTranslation().
			SetPrimaryAnchor(frWheelRelativePosition).
			SetIgnoreY(suspensionEnabled),
	)
	v.physicsScene.CreateDoubleBodyConstraint(chasisPhysics.Body, frWheelPhysics.Body, &solver.LimitTranslation{
		MaxY: suspensionStart,
		MinY: suspensionEnd,
	})
	frRotation := solver.NewMatchAxis().
		SetPrimaryAxis(sprec.BasisXVec3()).
		SetSecondaryAxis(sprec.BasisXVec3())
	v.physicsScene.CreateDoubleBodyConstraint(chasisPhysics.Body, frWheelPhysics.Body, frRotation)
	v.physicsScene.CreateDoubleBodyConstraint(chasisPhysics.Body, frWheelPhysics.Body, &solver.Coilover{
		PrimaryAnchor: frWheelRelativePosition,
		FrequencyHz:   suspensionFrequencyHz,
		DampingRatio:  suspensionDampingRatio,
	})

	blWheelRelativePosition := sprec.NewVec3(suspensionWidth, suspensionStart-suspensionLength, -1.56)
	blWheel := car.Wheel(model, car.BackLeftWheelLocation).
		WithName("back-left-wheel").
		WithPosition(sprec.Vec3Sum(position, blWheelRelativePosition)).
		Build(v.ecsScene, v.gfxScene, v.physicsScene)
	blWheelPhysics := ecscomp.GetPhysics(blWheel)
	v.physicsScene.CreateDoubleBodyConstraint(chasisPhysics.Body, blWheelPhysics.Body,
		solver.NewMatchTranslation().
			SetPrimaryAnchor(blWheelRelativePosition).
			SetIgnoreY(suspensionEnabled),
	)
	v.physicsScene.CreateDoubleBodyConstraint(chasisPhysics.Body, blWheelPhysics.Body, &solver.LimitTranslation{
		MaxY: suspensionStart,
		MinY: suspensionEnd,
	})
	v.physicsScene.CreateDoubleBodyConstraint(chasisPhysics.Body, blWheelPhysics.Body,
		solver.NewMatchAxis().
			SetPrimaryAxis(sprec.BasisXVec3()).
			SetSecondaryAxis(sprec.BasisXVec3()),
	)
	v.physicsScene.CreateDoubleBodyConstraint(chasisPhysics.Body, blWheelPhysics.Body, &solver.Coilover{
		PrimaryAnchor: blWheelRelativePosition,
		FrequencyHz:   suspensionFrequencyHz,
		DampingRatio:  suspensionDampingRatio,
	})

	brWheelRelativePosition := sprec.NewVec3(-suspensionWidth, suspensionStart-suspensionLength, -1.56)
	brWheel := car.Wheel(model, car.BackRightWheelLocation).
		WithName("back-right-wheel").
		WithPosition(sprec.Vec3Sum(position, brWheelRelativePosition)).
		Build(v.ecsScene, v.gfxScene, v.physicsScene)
	brWheelPhysics := ecscomp.GetPhysics(brWheel)
	v.physicsScene.CreateDoubleBodyConstraint(chasisPhysics.Body, brWheelPhysics.Body,
		solver.NewMatchTranslation().
			SetPrimaryAnchor(brWheelRelativePosition).
			SetIgnoreY(suspensionEnabled),
	)
	v.physicsScene.CreateDoubleBodyConstraint(chasisPhysics.Body, brWheelPhysics.Body, &solver.LimitTranslation{
		MaxY: suspensionStart,
		MinY: suspensionEnd,
	})
	v.physicsScene.CreateDoubleBodyConstraint(chasisPhysics.Body, brWheelPhysics.Body, solver.NewMatchAxis().
		SetPrimaryAxis(sprec.BasisXVec3()).
		SetSecondaryAxis(sprec.BasisXVec3()),
	)
	v.physicsScene.CreateDoubleBodyConstraint(chasisPhysics.Body, brWheelPhysics.Body, &solver.Coilover{
		PrimaryAnchor: brWheelRelativePosition,
		FrequencyHz:   suspensionFrequencyHz,
		DampingRatio:  suspensionDampingRatio,
	})

	car := v.ecsScene.CreateEntity()
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

func (v *View) Close(window app.Window) {
	v.renderSystem = nil
	v.vehicleSystem = nil
	v.cameraStandSystem = nil

	v.ecsScene.Delete()
	v.physicsScene.Delete()
	v.gfxScene.Delete()
}

func (v *View) OnKeyboardEvent(window app.Window, event app.KeyboardEvent) bool {
	if event.Code == app.KeyCodeF {
		switch event.Type {
		case app.KeyboardEventTypeKeyDown:
			v.freezeFrame = true
			return true
		case app.KeyboardEventTypeKeyUp:
			v.freezeFrame = false
			return true
		}
	}
	return v.vehicleSystem.OnKeyboardEvent(event)
}

func (v *View) Update(window app.Window, elapsedSeconds float32) {
	if v.freezeFrame {
		return
	}

	var gamepad *app.GamepadState
	if state, ok := window.GamepadState(0); ok {
		gamepad = &state
	}

	v.physicsScene.Update(elapsedSeconds)
	v.vehicleSystem.Update(elapsedSeconds, gamepad)
	v.renderSystem.Update()
	v.cameraStandSystem.Update(elapsedSeconds, gamepad)
}

func (v *View) Render(window app.Window, width, height int) {
	v.gfxScene.Render(graphics.NewViewport(0, 0, width, height), v.camera)
}
