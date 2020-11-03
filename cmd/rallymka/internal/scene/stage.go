package scene

import (
	"fmt"
	"time"

	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/lacking/async"
	"github.com/mokiat/lacking/game"
	"github.com/mokiat/lacking/physics"
	"github.com/mokiat/lacking/render"
	"github.com/mokiat/lacking/resource"
	"github.com/mokiat/lacking/shape"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/ecs"
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

func NewStage(gfxWorker *async.Worker) *Stage {
	scene := render.NewScene()
	if err := scene.Init(gfxWorker).Wait().Err; err != nil {
		panic(err) // FIXME
	}
	ecsManager := ecs.NewManager()
	stage := &Stage{
		scene:                scene,
		camera:               render.NewCamera(),
		ecsManager:           ecsManager,
		ecsRenderer:          ecs.NewRenderer(ecsManager, scene),
		ecsVehicleSystem:     ecs.NewVehicleSystem(ecsManager),
		ecsCameraStandSystem: ecs.NewCameraStandSystem(ecsManager),
		physicsEngine:        physics.NewEngine(15 * time.Millisecond),
	}
	scene.SetActiveCamera(stage.camera)
	return stage
}

type Stage struct {
	scene                *render.Scene
	camera               *render.Camera
	ecsManager           *ecs.Manager
	ecsRenderer          *ecs.Renderer
	ecsVehicleSystem     *ecs.VehicleSystem
	ecsCameraStandSystem *ecs.CameraStandSystem
	physicsEngine        *physics.Engine
}

func (s *Stage) Init(data *Data) {
	level := data.Level

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
		s.physicsEngine.AddBody(&physics.Body{
			Position:        sprec.ZeroVec3(),
			Orientation:     sprec.IdentityQuat(),
			IsStatic:        true,
			RestitutionCoef: 1.0,
			CollisionShapes: []shape.Placement{collisionMesh},
		})
	}

	for _, staticEntity := range level.StaticEntities {
		s.scene.Layout().CreateRenderable(staticEntity.Matrix, 100.0, staticEntity.Model)
	}

	carModel := data.CarModel
	targetEntity := s.setupCarDemo(carModel, sprec.NewVec3(0.0, 3.0, 0.0))
	standTarget := targetEntity
	standEntity := s.ecsManager.CreateEntity()
	standEntity.CameraStand = &ecs.CameraStand{
		Target:         standTarget,
		Camera:         s.camera,
		AnchorPosition: sprec.Vec3Sum(standTarget.Physics.Body.Position, sprec.NewVec3(0.0, 0.0, -cameraDistance)),
		AnchorDistance: anchorDistance,
		CameraDistance: cameraDistance,
	}
}

func (s *Stage) setupCarDemo(model *resource.Model, position sprec.Vec3) *ecs.Entity {
	chasis := car.Chassis(model).
		WithName("chasis").
		WithPosition(position).
		Build(s.ecsManager, s.scene)
	s.physicsEngine.AddBody(chasis.Physics.Body)

	suspensionEnabled := true
	suspensionStart := float32(-0.25)
	suspensionEnd := float32(-0.6)
	suspensionWidth := float32(1.0)
	suspensionLength := float32(0.25)
	suspensionFrequencyHz := float32(3.0)
	suspensionDampingRatio := float32(1.0)

	flWheelRelativePosition := sprec.NewVec3(suspensionWidth, suspensionStart-suspensionLength, 1.07)
	flWheel := car.Wheel(model, car.FrontLeftWheelLocation).
		WithName("front-left-wheel").
		WithPosition(sprec.Vec3Sum(position, flWheelRelativePosition)).
		Build(s.ecsManager, s.scene)
	s.physicsEngine.AddBody(flWheel.Physics.Body)
	s.physicsEngine.AddConstraint(physics.MatchTranslationConstraint{
		FirstBody:       chasis.Physics.Body,
		FirstBodyAnchor: flWheelRelativePosition,
		SecondBody:      flWheel.Physics.Body,
		IgnoreY:         suspensionEnabled,
	})
	s.physicsEngine.AddConstraint(physics.LimitTranslationConstraint{
		FirstBody:  chasis.Physics.Body,
		SecondBody: flWheel.Physics.Body,
		MaxY:       suspensionStart,
		MinY:       suspensionEnd,
	})
	flRotation := &physics.MatchAxisConstraint{
		FirstBody:      chasis.Physics.Body,
		FirstBodyAxis:  sprec.BasisXVec3(),
		SecondBody:     flWheel.Physics.Body,
		SecondBodyAxis: sprec.BasisXVec3(),
	}
	s.physicsEngine.AddConstraint(flRotation)
	s.physicsEngine.AddConstraint(&physics.CoiloverConstraint{
		FirstBody:       chasis.Physics.Body,
		FirstBodyAnchor: flWheelRelativePosition,
		SecondBody:      flWheel.Physics.Body,
		FrequencyHz:     suspensionFrequencyHz,
		DampingRatio:    suspensionDampingRatio,
	})

	frWheelRelativePosition := sprec.NewVec3(-suspensionWidth, suspensionStart-suspensionLength, 1.07)
	frWheel := car.Wheel(model, car.FrontRightWheelLocation).
		WithName("front-right-wheel").
		WithPosition(sprec.Vec3Sum(position, frWheelRelativePosition)).
		Build(s.ecsManager, s.scene)
	s.physicsEngine.AddBody(frWheel.Physics.Body)
	s.physicsEngine.AddConstraint(physics.MatchTranslationConstraint{
		FirstBody:       chasis.Physics.Body,
		FirstBodyAnchor: frWheelRelativePosition,
		SecondBody:      frWheel.Physics.Body,
		IgnoreY:         suspensionEnabled,
	})
	s.physicsEngine.AddConstraint(physics.LimitTranslationConstraint{
		FirstBody:  chasis.Physics.Body,
		SecondBody: frWheel.Physics.Body,
		MaxY:       suspensionStart,
		MinY:       suspensionEnd,
	})
	frRotation := &physics.MatchAxisConstraint{
		FirstBody:      chasis.Physics.Body,
		FirstBodyAxis:  sprec.BasisXVec3(),
		SecondBody:     frWheel.Physics.Body,
		SecondBodyAxis: sprec.BasisXVec3(),
	}
	s.physicsEngine.AddConstraint(frRotation)
	s.physicsEngine.AddConstraint(&physics.CoiloverConstraint{
		FirstBody:       chasis.Physics.Body,
		FirstBodyAnchor: frWheelRelativePosition,
		SecondBody:      frWheel.Physics.Body,
		FrequencyHz:     suspensionFrequencyHz,
		DampingRatio:    suspensionDampingRatio,
	})

	blWheelRelativePosition := sprec.NewVec3(suspensionWidth, suspensionStart-suspensionLength, -1.56)
	blWheel := car.Wheel(model, car.BackLeftWheelLocation).
		WithName("back-left-wheel").
		WithPosition(sprec.Vec3Sum(position, blWheelRelativePosition)).
		Build(s.ecsManager, s.scene)
	s.physicsEngine.AddBody(blWheel.Physics.Body)
	s.physicsEngine.AddConstraint(physics.MatchTranslationConstraint{
		FirstBody:       chasis.Physics.Body,
		FirstBodyAnchor: blWheelRelativePosition,
		SecondBody:      blWheel.Physics.Body,
		IgnoreY:         suspensionEnabled,
	})
	s.physicsEngine.AddConstraint(physics.LimitTranslationConstraint{
		FirstBody:  chasis.Physics.Body,
		SecondBody: blWheel.Physics.Body,
		MaxY:       suspensionStart,
		MinY:       suspensionEnd,
	})
	s.physicsEngine.AddConstraint(physics.MatchAxisConstraint{
		FirstBody:      chasis.Physics.Body,
		FirstBodyAxis:  sprec.BasisXVec3(),
		SecondBody:     blWheel.Physics.Body,
		SecondBodyAxis: sprec.BasisXVec3(),
	})
	s.physicsEngine.AddConstraint(&physics.CoiloverConstraint{
		FirstBody:       chasis.Physics.Body,
		FirstBodyAnchor: blWheelRelativePosition,
		SecondBody:      blWheel.Physics.Body,
		FrequencyHz:     suspensionFrequencyHz,
		DampingRatio:    suspensionDampingRatio,
	})

	brWheelRelativePosition := sprec.NewVec3(-suspensionWidth, suspensionStart-suspensionLength, -1.56)
	brWheel := car.Wheel(model, car.BackRightWheelLocation).
		WithName("back-right-wheel").
		WithPosition(sprec.Vec3Sum(position, brWheelRelativePosition)).
		Build(s.ecsManager, s.scene)
	s.physicsEngine.AddBody(brWheel.Physics.Body)
	s.physicsEngine.AddConstraint(physics.MatchTranslationConstraint{
		FirstBody:       chasis.Physics.Body,
		FirstBodyAnchor: brWheelRelativePosition,
		SecondBody:      brWheel.Physics.Body,
		IgnoreY:         suspensionEnabled,
	})
	s.physicsEngine.AddConstraint(physics.LimitTranslationConstraint{
		FirstBody:  chasis.Physics.Body,
		SecondBody: brWheel.Physics.Body,
		MaxY:       suspensionStart,
		MinY:       suspensionEnd,
	})
	s.physicsEngine.AddConstraint(physics.MatchAxisConstraint{
		FirstBody:      chasis.Physics.Body,
		FirstBodyAxis:  sprec.BasisXVec3(),
		SecondBody:     brWheel.Physics.Body,
		SecondBodyAxis: sprec.BasisXVec3(),
	})
	s.physicsEngine.AddConstraint(&physics.CoiloverConstraint{
		FirstBody:       chasis.Physics.Body,
		FirstBodyAnchor: brWheelRelativePosition,
		SecondBody:      brWheel.Physics.Body,
		FrequencyHz:     suspensionFrequencyHz,
		DampingRatio:    suspensionDampingRatio,
	})

	car := s.ecsManager.CreateEntity()
	car.Vehicle = &ecs.Vehicle{
		MaxSteeringAngle: sprec.Degrees(carMaxSteeringAngle),
		SteeringAngle:    sprec.Degrees(0.0),
		Acceleration:     0.0,
		Deceleration:     0.0,
		Chassis: &ecs.Chassis{
			Body: chasis.Physics.Body,
		},
		Wheels: []*ecs.Wheel{
			{
				Body:                 flWheel.Physics.Body,
				RotationConstraint:   flRotation,
				AccelerationVelocity: carFrontAcceleration,
				DecelerationVelocity: carFrontDeceleration,
			},
			{
				Body:                 frWheel.Physics.Body,
				RotationConstraint:   frRotation,
				AccelerationVelocity: carFrontAcceleration,
				DecelerationVelocity: carFrontDeceleration,
			},
			{
				Body:                 blWheel.Physics.Body,
				AccelerationVelocity: carRearAcceleration,
				DecelerationVelocity: carRearDeceleration,
			},
			{
				Body:                 brWheel.Physics.Body,
				AccelerationVelocity: carRearAcceleration,
				DecelerationVelocity: carRearDeceleration,
			},
		},
	}
	car.PlayerControl = &ecs.PlayerControl{}

	return chasis
}

func (s *Stage) Update(ctx game.UpdateContext) {
	s.physicsEngine.Update(ctx.ElapsedTime)
	s.ecsVehicleSystem.Update(ctx)
	s.ecsRenderer.Update(ctx)
	s.ecsCameraStandSystem.Update(ctx)
}

func (s *Stage) Render(ctx game.RenderContext) {
	screenHalfWidth := float32(ctx.WindowSize.Width) / float32(ctx.WindowSize.Height)
	screenHalfHeight := float32(1.0)
	s.camera.SetProjectionMatrix(sprec.PerspectiveMat4(
		-screenHalfWidth, screenHalfWidth, -screenHalfHeight, screenHalfHeight, 1.5, 900.0,
	))
	startTime := time.Now()
	s.scene.Render(ctx)
	elapsedTime := time.Since(startTime)
	fmt.Printf("render duration: %s\n", elapsedTime)
}
