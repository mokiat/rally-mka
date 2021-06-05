package scene

import (
	"time"

	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/lacking/app"
	"github.com/mokiat/lacking/async"
	"github.com/mokiat/lacking/game/ecs"
	"github.com/mokiat/lacking/game/physics"
	"github.com/mokiat/lacking/game/physics/solver"
	"github.com/mokiat/lacking/graphics"
	"github.com/mokiat/lacking/render"
	"github.com/mokiat/lacking/resource"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/ecscomp"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/ecssys"
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
)

type CarInput struct {
	Forward   bool
	Backward  bool
	TurnLeft  bool
	TurnRight bool
	Handbrake bool
}

func NewStage() *Stage {
	scene := render.NewScene()
	ecsEngine := ecs.NewEngine()
	ecsScene := ecsEngine.CreateScene()
	physicsEngine := physics.NewEngine()
	stage := &Stage{
		scene:                scene,
		camera:               render.NewCamera(),
		ecsScene:             ecsScene,
		ecsRenderer:          ecssys.NewRenderer(ecsScene, scene),
		ecsVehicleSystem:     ecssys.NewVehicleSystem(ecsScene),
		ecsCameraStandSystem: ecssys.NewCameraStandSystem(ecsScene),
		physicsScene:         physicsEngine.CreateScene(0.015),
	}
	return stage
}

type Stage struct {
	scene                *render.Scene
	camera               *render.Camera
	ecsScene             *ecs.Scene
	ecsRenderer          *ecssys.Renderer
	ecsVehicleSystem     *ecssys.VehicleSystem
	ecsCameraStandSystem *ecssys.CameraStandSystem
	physicsScene         *physics.Scene
}

func (s *Stage) Init(gfxWorker *async.Worker, data *Data) {
	level := data.Level

	if err := s.scene.Init(gfxWorker).Wait().Err; err != nil {
		panic(err) // FIXME
	}

	s.scene.SetActiveCamera(s.camera)

	s.scene.Layout().SetSkybox(&render.Skybox{
		SkyboxTexture:            data.Level.SkyboxTexture.GFXTexture,
		AmbientReflectionTexture: data.Level.AmbientReflectionTexture.GFXTexture,
		AmbientRefractionTexture: data.Level.AmbientRefractionTexture.GFXTexture,
	})

	for _, staticMesh := range level.StaticMeshes {
		s.scene.Layout().CreateRenderable(sprec.IdentityMat4(), 100.0, &resource.Model{
			Name: "static",
			Nodes: []*resource.Node{
				{
					Name:   "root",
					Matrix: sprec.IdentityMat4(),
					Mesh:   staticMesh,
				},
			},
		})
	}

	for _, collisionMesh := range level.CollisionMeshes {
		body := s.physicsScene.CreateBody()
		body.SetPosition(sprec.ZeroVec3())
		body.SetOrientation(sprec.IdentityQuat())
		body.SetStatic(true)
		body.SetRestitutionCoefficient(1.0)
		body.SetCollisionShapes([]physics.CollisionShape{collisionMesh})
	}

	for _, staticEntity := range level.StaticEntities {
		s.scene.Layout().CreateRenderable(staticEntity.Matrix, 100.0, staticEntity.Model)
	}

	carModel := data.CarModel
	targetEntity := s.setupCarDemo(carModel, sprec.NewVec3(0.0, 3.0, 0.0))
	standTarget := targetEntity
	standEntity := s.ecsScene.CreateEntity()
	ecscomp.SetCameraStand(standEntity, &ecscomp.CameraStand{
		Target:         standTarget,
		Camera:         s.camera,
		AnchorPosition: sprec.Vec3Sum(ecscomp.GetPhysics(standTarget).Body.Position(), sprec.NewVec3(0.0, 0.0, -cameraDistance)),
		AnchorDistance: anchorDistance,
		CameraDistance: cameraDistance,
	})
}

func (s *Stage) setupCarDemo(model *resource.Model, position sprec.Vec3) *ecs.Entity {
	const (
		suspensionEnabled      = true
		suspensionStart        = float32(-0.25)
		suspensionEnd          = float32(-0.6)
		suspensionWidth        = float32(1.0)
		suspensionLength       = float32(0.25)
		suspensionFrequencyHz  = float32(3.0)
		suspensionDampingRatio = float32(1.0)
	)

	chasis := car.Chassis(model).
		WithName("chasis").
		WithPosition(position).
		Build(s.ecsScene, s.scene, s.physicsScene)
	chasisPhysics := ecscomp.GetPhysics(chasis)

	flWheelRelativePosition := sprec.NewVec3(suspensionWidth, suspensionStart-suspensionLength, 1.07)
	flWheel := car.Wheel(model, car.FrontLeftWheelLocation).
		WithName("front-left-wheel").
		WithPosition(sprec.Vec3Sum(position, flWheelRelativePosition)).
		Build(s.ecsScene, s.scene, s.physicsScene)
	flWheelPhysics := ecscomp.GetPhysics(flWheel)

	s.physicsScene.CreateDoubleBodyConstraint(chasisPhysics.Body, flWheelPhysics.Body,
		solver.NewMatchTranslation().
			SetPrimaryAnchor(flWheelRelativePosition).
			SetIgnoreY(suspensionEnabled),
	)
	s.physicsScene.CreateDoubleBodyConstraint(chasisPhysics.Body, flWheelPhysics.Body, &solver.LimitTranslation{
		MaxY: suspensionStart,
		MinY: suspensionEnd,
	})
	flRotation := solver.NewMatchAxis().
		SetPrimaryAxis(sprec.BasisXVec3()).
		SetSecondaryAxis(sprec.BasisXVec3())
	s.physicsScene.CreateDoubleBodyConstraint(chasisPhysics.Body, flWheelPhysics.Body, flRotation)
	s.physicsScene.CreateDoubleBodyConstraint(chasisPhysics.Body, flWheelPhysics.Body, &solver.Coilover{
		PrimaryAnchor: flWheelRelativePosition,
		FrequencyHz:   suspensionFrequencyHz,
		DampingRatio:  suspensionDampingRatio,
	})

	frWheelRelativePosition := sprec.NewVec3(-suspensionWidth, suspensionStart-suspensionLength, 1.07)
	frWheel := car.Wheel(model, car.FrontRightWheelLocation).
		WithName("front-right-wheel").
		WithPosition(sprec.Vec3Sum(position, frWheelRelativePosition)).
		Build(s.ecsScene, s.scene, s.physicsScene)
	frWheelPhysics := ecscomp.GetPhysics(frWheel)
	s.physicsScene.CreateDoubleBodyConstraint(chasisPhysics.Body, frWheelPhysics.Body,
		solver.NewMatchTranslation().
			SetPrimaryAnchor(frWheelRelativePosition).
			SetIgnoreY(suspensionEnabled),
	)
	s.physicsScene.CreateDoubleBodyConstraint(chasisPhysics.Body, frWheelPhysics.Body, &solver.LimitTranslation{
		MaxY: suspensionStart,
		MinY: suspensionEnd,
	})
	frRotation := solver.NewMatchAxis().
		SetPrimaryAxis(sprec.BasisXVec3()).
		SetSecondaryAxis(sprec.BasisXVec3())
	s.physicsScene.CreateDoubleBodyConstraint(chasisPhysics.Body, frWheelPhysics.Body, frRotation)
	s.physicsScene.CreateDoubleBodyConstraint(chasisPhysics.Body, frWheelPhysics.Body, &solver.Coilover{
		PrimaryAnchor: frWheelRelativePosition,
		FrequencyHz:   suspensionFrequencyHz,
		DampingRatio:  suspensionDampingRatio,
	})

	blWheelRelativePosition := sprec.NewVec3(suspensionWidth, suspensionStart-suspensionLength, -1.56)
	blWheel := car.Wheel(model, car.BackLeftWheelLocation).
		WithName("back-left-wheel").
		WithPosition(sprec.Vec3Sum(position, blWheelRelativePosition)).
		Build(s.ecsScene, s.scene, s.physicsScene)
	blWheelPhysics := ecscomp.GetPhysics(blWheel)
	s.physicsScene.CreateDoubleBodyConstraint(chasisPhysics.Body, blWheelPhysics.Body,
		solver.NewMatchTranslation().
			SetPrimaryAnchor(blWheelRelativePosition).
			SetIgnoreY(suspensionEnabled),
	)
	s.physicsScene.CreateDoubleBodyConstraint(chasisPhysics.Body, blWheelPhysics.Body, &solver.LimitTranslation{
		MaxY: suspensionStart,
		MinY: suspensionEnd,
	})
	s.physicsScene.CreateDoubleBodyConstraint(chasisPhysics.Body, blWheelPhysics.Body,
		solver.NewMatchAxis().
			SetPrimaryAxis(sprec.BasisXVec3()).
			SetSecondaryAxis(sprec.BasisXVec3()),
	)
	s.physicsScene.CreateDoubleBodyConstraint(chasisPhysics.Body, blWheelPhysics.Body, &solver.Coilover{
		PrimaryAnchor: blWheelRelativePosition,
		FrequencyHz:   suspensionFrequencyHz,
		DampingRatio:  suspensionDampingRatio,
	})

	brWheelRelativePosition := sprec.NewVec3(-suspensionWidth, suspensionStart-suspensionLength, -1.56)
	brWheel := car.Wheel(model, car.BackRightWheelLocation).
		WithName("back-right-wheel").
		WithPosition(sprec.Vec3Sum(position, brWheelRelativePosition)).
		Build(s.ecsScene, s.scene, s.physicsScene)
	brWheelPhysics := ecscomp.GetPhysics(brWheel)
	s.physicsScene.CreateDoubleBodyConstraint(chasisPhysics.Body, brWheelPhysics.Body,
		solver.NewMatchTranslation().
			SetPrimaryAnchor(brWheelRelativePosition).
			SetIgnoreY(suspensionEnabled),
	)
	s.physicsScene.CreateDoubleBodyConstraint(chasisPhysics.Body, brWheelPhysics.Body, &solver.LimitTranslation{
		MaxY: suspensionStart,
		MinY: suspensionEnd,
	})
	s.physicsScene.CreateDoubleBodyConstraint(chasisPhysics.Body, brWheelPhysics.Body, solver.NewMatchAxis().
		SetPrimaryAxis(sprec.BasisXVec3()).
		SetSecondaryAxis(sprec.BasisXVec3()),
	)
	s.physicsScene.CreateDoubleBodyConstraint(chasisPhysics.Body, brWheelPhysics.Body, &solver.Coilover{
		PrimaryAnchor: brWheelRelativePosition,
		FrequencyHz:   suspensionFrequencyHz,
		DampingRatio:  suspensionDampingRatio,
	})

	car := s.ecsScene.CreateEntity()
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

func (s *Stage) OnKeyboardEvent(event app.KeyboardEvent) bool {
	return s.ecsVehicleSystem.OnKeyboardEvent(event)
}

func (s *Stage) Update(window app.Window, elapsedTime time.Duration) {
	var gamepad *app.GamepadState
	if state, ok := window.GamepadState(0); ok {
		gamepad = &state
	}

	s.physicsScene.Update(float32(elapsedTime.Seconds()))
	s.ecsVehicleSystem.Update(elapsedTime, gamepad)
	s.ecsRenderer.Update()
	s.ecsCameraStandSystem.Update(elapsedTime, gamepad)
}

func (s *Stage) Render(width, height int, pipeline *graphics.Pipeline) {
	screenHalfWidth := float32(width) / float32(height)
	screenHalfHeight := float32(1.0)
	s.camera.SetProjectionMatrix(sprec.PerspectiveMat4(
		-screenHalfWidth, screenHalfWidth, -screenHalfHeight, screenHalfHeight, 1.5, 900.0,
	))
	s.scene.Render(width, height, pipeline)
}
