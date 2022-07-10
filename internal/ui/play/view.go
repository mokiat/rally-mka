package play

import (
	"fmt"

	"github.com/mokiat/gomath/sprec"
	lackgame "github.com/mokiat/lacking/game"
	"github.com/mokiat/lacking/game/ecs"
	"github.com/mokiat/lacking/game/graphics"
	"github.com/mokiat/lacking/game/physics"
	"github.com/mokiat/lacking/game/physics/solver"
	"github.com/mokiat/lacking/resource"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
	"github.com/mokiat/lacking/util/optional"
	"github.com/mokiat/rally-mka/internal/ecscomp"
	"github.com/mokiat/rally-mka/internal/ecssys"
	"github.com/mokiat/rally-mka/internal/game"
	"github.com/mokiat/rally-mka/internal/global"
	"github.com/mokiat/rally-mka/internal/scene"
	"github.com/mokiat/rally-mka/internal/scene/car"
)

const (
	correction = float32(0.9)

	anchorDistance = 6.0
	cameraDistance = 12.0 * correction

	carMaxSteeringAngle  = 30
	carFrontAcceleration = 145 * 1
	carRearAcceleration  = 160 * 1

	// FIXME: Currently, too much front brakes cause the car
	// to straighten. This is due to there being more pressure
	// on the outer wheel which causes it to brake more and turn
	// the car to neutral orientation.
	carFrontDeceleration = 250
	carRearDeceleration  = 180

	suspensionEnabled      = true
	suspensionStart        = float32(-0.25) * correction
	suspensionEnd          = float32(-0.6) * correction
	suspensionWidth        = float32(1.0) * correction
	suspensionLength       = float32(0.15) * correction
	suspensionFrequencyHz  = float32(4.0)
	suspensionDampingRatio = float32(1.2)
)

type ViewData struct {
	GameData *scene.Data
}

var View = co.Define(func(props co.Properties, scope co.Scope) co.Instance {
	var (
		context = co.GetContext[global.Context]()
		data    = co.GetData[ViewData](props)
	)

	lifecycle := co.UseState(func() *playLifecycle {
		return &playLifecycle{
			gameController: context.GameController,
			gameData:       data.GameData,
		}
	}).Get()

	speedState := co.UseState(func() float32 {
		return float32(0.0)
	})
	speed := speedState.Get()

	co.Once(func() {
		co.Window(scope).SetCursorVisible(false)
	})

	co.Once(func() {
		context.GameController.OnUpdate = func() {
			carSpeed := lifecycle.car.Body().Velocity().Length() * 3.6
			speedState.Set(carSpeed)
		}
	})

	co.Defer(func() {
		co.Window(scope).SetCursorVisible(true)
	})

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

		co.WithChild("speed-label", co.New(mat.Label, func() {
			co.WithData(mat.LabelData{
				Font:      co.OpenFont(scope, "mat:///roboto-bold.ttf"),
				FontSize:  optional.Value(float32(24.0)),
				FontColor: optional.Value(ui.White()),
				Text:      fmt.Sprintf("speed: %.4f", speed),
			})
			co.WithLayoutData(mat.LayoutData{
				Left:   optional.Value(0),
				Top:    optional.Value(0),
				Width:  optional.Value(200),
				Height: optional.Value(24),
			})
		}))
	})
})

type playLifecycle struct {
	gameController *game.Controller
	gameData       *scene.Data

	gfxScene     *graphics.Scene
	physicsScene *physics.Scene
	ecsScene     *ecs.Scene

	vehicleSystem     *ecssys.VehicleSystem
	cameraStandSystem *ecssys.CameraStandSystem

	camera *graphics.Camera
	car    *lackgame.Node
}

func (h *playLifecycle) init() {
	scene := h.gameController.Scene()

	h.physicsScene = scene.Physics()
	h.gfxScene = scene.Graphics()
	h.ecsScene = scene.ECS()

	h.vehicleSystem = h.gameController.VehicleSystem()
	h.cameraStandSystem = h.gameController.CameraStandSystem()

	h.camera = h.gameController.Camera()
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
	sunLight.SetMatrix(sprec.TRSMat4(
		sprec.ZeroVec3(),
		sprec.QuatProd(
			sprec.RotationQuat(sprec.Degrees(-90), sprec.BasisYVec3()),
			sprec.RotationQuat(sprec.Degrees(-24), sprec.BasisXVec3()),
		),
		sprec.NewVec3(1.0, 1.0, 1.0),
	))
	// sunLight.SetRotation(sprec.QuatProd(
	// 	sprec.RotationQuat(sprec.Degrees(225), sprec.BasisYVec3()),
	// 	sprec.RotationQuat(sprec.Degrees(-45), sprec.BasisXVec3()),
	// ))
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

	for _, staticEntity := range level.StaticEntities {
		for _, instance := range staticEntity.Model.MeshInstances {
			node := instance.Node
			modelMatrix := sprec.Mat4Prod(staticEntity.Matrix, node.AbsoluteMatrix())

			gfxMesh := h.gfxScene.CreateMesh(instance.MeshDefinition.GFXMeshTemplate)
			gfxMesh.SetMatrix(modelMatrix)

		}
	}

	carModel := h.gameData.CarModel
	targetNode := h.setupCarDemo(carModel, sprec.NewVec3(0.0, 3.0, 0.0))

	cameraEntity := h.ecsScene.CreateEntity()

	// cameraBody := h.physicsScene.CreateBody()
	// cameraBody.SetPosition(sprec.NewVec3(0.0, 10.0, -5.0))
	// cameraBody.SetOrientation(sprec.RotationQuat(sprec.Degrees(180), sprec.BasisYVec3()))
	// cameraBody.SetMass(10.0)
	// cameraBody.SetMomentOfInertia(physics.SymmetricMomentOfInertia(100 * 1 * 1 / 5.0))
	// cameraBody.SetDragFactor(0.0)
	// cameraBody.SetAngularDragFactor(0.0)
	// cameraBody.SetRestitutionCoefficient(-0.5)
	// cameraBody.SetCollisionShapes([]physics.CollisionShape{
	// 	shape.Placement{
	// 		Position:    sprec.ZeroVec3(),
	// 		Orientation: sprec.IdentityQuat(),
	// 		Shape:       shape.NewStaticSphere(1.0),
	// 	},
	// })

	// ecscomp.SetPhysics(cameraEntity, &ecscomp.Physics{
	// 	Body: cameraBody,
	// })
	// ecscomp.SetHierarchy(cameraEntity, &ecscomp.Hierarchy{
	// 	Node: hierarchy.NewNode(nil),
	// })
	// ecscomp.SetCamera(cameraEntity, &ecscomp.Camera{
	// 	Camera: h.camera,
	// })

	// carBody := ecscomp.GetPhysics(targetEntity).Body

	// cameraRodConstraint := solver.NewHingedRod().
	// 	SetLength(5.0).
	// 	SetPrimaryAnchor(sprec.NewVec3(0.0, 0.0, -1)).
	// 	SetSecondaryAnchor(sprec.NewVec3(0.0, 0.0, 0.0))
	// h.physicsScene.CreateDoubleBodyConstraint(carBody, cameraBody, cameraRodConstraint)

	// cameraAxisConstraint := solver.NewMatchAxis().
	// 	SetPrimaryAxis(sprec.BasisYVec3()).
	// 	SetSecondaryAxis(sprec.BasisYVec3())
	// h.physicsScene.CreateDoubleBodyConstraint(carBody, cameraBody, cameraAxisConstraint)

	// cameraAxis2Constraint := solver.NewMatchAxis().
	// 	SetPrimaryAxis(sprec.BasisZVec3()).
	// 	SetSecondaryAxis(sprec.InverseVec3(sprec.BasisZVec3()))
	// h.physicsScene.CreateDoubleBodyConstraint(carBody, cameraBody, cameraAxis2Constraint)

	ecscomp.SetCameraStand(cameraEntity, &ecscomp.CameraStand{
		Target:         targetNode,
		Camera:         h.camera,
		AnchorPosition: sprec.Vec3Sum(targetNode.Body().Position(), sprec.NewVec3(0.0, 0.0, -cameraDistance)),
		AnchorDistance: anchorDistance,
		CameraDistance: cameraDistance,
	})

	h.car = targetNode
}

func (h *playLifecycle) setupCarDemo(model *resource.Model, position sprec.Vec3) *lackgame.Node {
	chasis := car.Chassis(model).
		WithName("chasis").
		WithPosition(position).
		Build(h.gameController.Scene())
	chasisBody := chasis.Body()

	flWheelRelativePosition := sprec.NewVec3(suspensionWidth, suspensionStart-suspensionLength, 1.07*correction)
	flWheel := car.Wheel(model, car.FrontLeftWheelLocation).
		WithName("front-left-wheel").
		WithPosition(sprec.Vec3Sum(position, flWheelRelativePosition)).
		Build(h.gameController.Scene())
	flWheelBody := flWheel.Body()
	h.physicsScene.CreateDoubleBodyConstraint(chasisBody, flWheelBody,
		solver.NewMatchTranslation().
			SetPrimaryAnchor(flWheelRelativePosition).
			SetIgnoreY(suspensionEnabled),
	)
	h.physicsScene.CreateDoubleBodyConstraint(chasisBody, flWheelBody, &solver.LimitTranslation{
		MaxY: suspensionStart,
		MinY: suspensionEnd,
	})
	flRotation := solver.NewMatchAxis().
		SetPrimaryAxis(sprec.BasisXVec3()).
		SetSecondaryAxis(sprec.BasisXVec3())
	h.physicsScene.CreateDoubleBodyConstraint(chasisBody, flWheelBody, flRotation)
	h.physicsScene.CreateDoubleBodyConstraint(chasisBody, flWheelBody, &solver.Coilover{
		PrimaryAnchor: flWheelRelativePosition,
		FrequencyHz:   suspensionFrequencyHz,
		DampingRatio:  suspensionDampingRatio,
	})

	frWheelRelativePosition := sprec.NewVec3(-suspensionWidth, suspensionStart-suspensionLength, 1.07*correction)
	frWheel := car.Wheel(model, car.FrontRightWheelLocation).
		WithName("front-right-wheel").
		WithPosition(sprec.Vec3Sum(position, frWheelRelativePosition)).
		Build(h.gameController.Scene())
	frWheelBody := frWheel.Body()
	h.physicsScene.CreateDoubleBodyConstraint(chasisBody, frWheelBody,
		solver.NewMatchTranslation().
			SetPrimaryAnchor(frWheelRelativePosition).
			SetIgnoreY(suspensionEnabled),
	)
	h.physicsScene.CreateDoubleBodyConstraint(chasisBody, frWheelBody, &solver.LimitTranslation{
		MaxY: suspensionStart,
		MinY: suspensionEnd,
	})
	frRotation := solver.NewMatchAxis().
		SetPrimaryAxis(sprec.BasisXVec3()).
		SetSecondaryAxis(sprec.BasisXVec3())
	h.physicsScene.CreateDoubleBodyConstraint(chasisBody, frWheelBody, frRotation)
	h.physicsScene.CreateDoubleBodyConstraint(chasisBody, frWheelBody, &solver.Coilover{
		PrimaryAnchor: frWheelRelativePosition,
		FrequencyHz:   suspensionFrequencyHz,
		DampingRatio:  suspensionDampingRatio,
	})

	h.physicsScene.CreateDoubleBodyConstraint(flWheelBody, frWheelBody, &Differential{})

	blWheelRelativePosition := sprec.NewVec3(suspensionWidth, suspensionStart-suspensionLength, -1.56*correction)
	blWheel := car.Wheel(model, car.BackLeftWheelLocation).
		WithName("back-left-wheel").
		WithPosition(sprec.Vec3Sum(position, blWheelRelativePosition)).
		Build(h.gameController.Scene())
	blWheelBody := blWheel.Body()
	h.physicsScene.CreateDoubleBodyConstraint(chasisBody, blWheelBody,
		solver.NewMatchTranslation().
			SetPrimaryAnchor(blWheelRelativePosition).
			SetIgnoreY(suspensionEnabled),
	)
	h.physicsScene.CreateDoubleBodyConstraint(chasisBody, blWheelBody, &solver.LimitTranslation{
		MaxY: suspensionStart,
		MinY: suspensionEnd,
	})
	h.physicsScene.CreateDoubleBodyConstraint(chasisBody, blWheelBody,
		solver.NewMatchAxis().
			SetPrimaryAxis(sprec.BasisXVec3()).
			SetSecondaryAxis(sprec.BasisXVec3()),
	)
	h.physicsScene.CreateDoubleBodyConstraint(chasisBody, blWheelBody, &solver.Coilover{
		PrimaryAnchor: blWheelRelativePosition,
		FrequencyHz:   suspensionFrequencyHz,
		DampingRatio:  suspensionDampingRatio,
	})

	brWheelRelativePosition := sprec.NewVec3(-suspensionWidth, suspensionStart-suspensionLength, -1.56*correction)
	brWheel := car.Wheel(model, car.BackRightWheelLocation).
		WithName("back-right-wheel").
		WithPosition(sprec.Vec3Sum(position, brWheelRelativePosition)).
		Build(h.gameController.Scene())
	brWheelBody := brWheel.Body()
	h.physicsScene.CreateDoubleBodyConstraint(chasisBody, brWheelBody,
		solver.NewMatchTranslation().
			SetPrimaryAnchor(brWheelRelativePosition).
			SetIgnoreY(suspensionEnabled),
	)
	h.physicsScene.CreateDoubleBodyConstraint(chasisBody, brWheelBody, &solver.LimitTranslation{
		MaxY: suspensionStart,
		MinY: suspensionEnd,
	})
	h.physicsScene.CreateDoubleBodyConstraint(chasisBody, brWheelBody, solver.NewMatchAxis().
		SetPrimaryAxis(sprec.BasisXVec3()).
		SetSecondaryAxis(sprec.BasisXVec3()),
	)
	h.physicsScene.CreateDoubleBodyConstraint(chasisBody, brWheelBody, &solver.Coilover{
		PrimaryAnchor: brWheelRelativePosition,
		FrequencyHz:   suspensionFrequencyHz,
		DampingRatio:  suspensionDampingRatio,
	})

	h.physicsScene.CreateDoubleBodyConstraint(blWheelBody, brWheelBody, &Differential{})

	car := h.ecsScene.CreateEntity()
	ecscomp.SetVehicle(car, &ecscomp.Vehicle{
		MaxSteeringAngle: sprec.Degrees(carMaxSteeringAngle),
		SteeringAngle:    sprec.Degrees(0.0),
		Acceleration:     0.0,
		Deceleration:     0.0,
		Chassis: &ecscomp.Chassis{
			Body: chasisBody,
		},
		Wheels: []*ecscomp.Wheel{
			{
				Body:                 flWheelBody,
				RotationConstraint:   flRotation,
				AccelerationVelocity: carFrontAcceleration,
				DecelerationVelocity: carFrontDeceleration,
			},
			{
				Body:                 frWheelBody,
				RotationConstraint:   frRotation,
				AccelerationVelocity: carFrontAcceleration,
				DecelerationVelocity: carFrontDeceleration,
			},
			{
				Body:                 blWheelBody,
				AccelerationVelocity: carRearAcceleration,
				DecelerationVelocity: carRearDeceleration,
			},
			{
				Body:                 brWheelBody,
				AccelerationVelocity: carRearAcceleration,
				DecelerationVelocity: carRearDeceleration,
			},
		},
	})
	ecscomp.SetPlayerControl(car, &ecscomp.PlayerControl{})

	return chasis
}

var _ physics.DBConstraintSolver = (*Differential)(nil)

type Differential struct {
	physics.NilDBConstraintSolver
}

func (d *Differential) CalculateImpulses(ctx physics.DBSolverContext) physics.DBImpulseSolution {
	firstRotation := sprec.Vec3Dot(ctx.Primary.Orientation().OrientationX(), ctx.Primary.AngularVelocity())
	secondRotation := sprec.Vec3Dot(ctx.Secondary.Orientation().OrientationX(), ctx.Secondary.AngularVelocity())

	const maxDelta = float32(100.0)

	var firstCorrection sprec.Vec3
	if firstRotation > secondRotation+maxDelta {
		firstCorrection = sprec.Vec3Prod(ctx.Primary.Orientation().OrientationX(), secondRotation+maxDelta-firstRotation)
	}

	var secondCorrection sprec.Vec3
	if secondRotation > firstRotation+maxDelta {
		secondCorrection = sprec.Vec3Prod(ctx.Secondary.Orientation().OrientationX(), firstRotation+maxDelta-secondRotation)
	}

	return physics.DBImpulseSolution{
		Primary: physics.SBImpulseSolution{
			Impulse:        sprec.ZeroVec3(),
			AngularImpulse: firstCorrection,
		},
		Secondary: physics.SBImpulseSolution{
			Impulse:        sprec.ZeroVec3(),
			AngularImpulse: secondCorrection,
		},
	}
}
