package scene

import (
	"time"

	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/ecs"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/ecs/constraint"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/stream"
	"github.com/mokiat/rally-mka/internal/engine/collision"
	"github.com/mokiat/rally-mka/internal/engine/graphics"
)

const (
	framebufferWidth  = int32(1024)
	framebufferHeight = int32(576)
)

const (
	carDropHeight      = 1.6
	entityVisualHeight = 5.0
	anchorDistance     = 4.0
	cameraDistance     = 8.0
)

type CarInput struct {
	Forward   bool
	Backward  bool
	TurnLeft  bool
	TurnRight bool
	Handbrake bool
}

func NewStage() *Stage {
	ecsManager := ecs.NewManager()
	stage := &Stage{
		ecsManager:           ecsManager,
		ecsRenderer:          ecs.NewRenderer(ecsManager),
		ecsCameraStandSystem: ecs.NewCameraStandSystem(ecsManager),
		ecsPhysicsSystem:     ecs.NewPhysicsSystem(ecsManager),
		screenFramebuffer:    &graphics.Framebuffer{},
	}
	stage.ecsVehicleSystem = ecs.NewVehicleSystem(ecsManager, stage)
	return stage
}

type Stage struct {
	ecsManager           *ecs.Manager
	ecsRenderer          *ecs.Renderer
	ecsVehicleSystem     *ecs.VehicleSystem
	ecsCameraStandSystem *ecs.CameraStandSystem
	ecsPhysicsSystem     *ecs.PhysicsSystem

	geometryFramebuffer *graphics.Framebuffer
	lightingFramebuffer *graphics.Framebuffer
	screenFramebuffer   *graphics.Framebuffer
	lightingProgram     *graphics.Program
	quadMesh            *stream.Mesh

	collisionMeshes []*collision.Mesh
}

func (s *Stage) Init(data *Data, camera *ecs.Camera) {
	level := data.Level.Get()

	s.geometryFramebuffer = data.GeometryFramebuffer
	s.lightingFramebuffer = data.LightingFramebuffer

	s.lightingProgram = data.DeferredLightingProgram.Get()
	s.quadMesh = data.QuadMesh.Get()

	for _, staticMesh := range level.StaticMeshes {
		entity := s.ecsManager.CreateEntity()
		entity.Transform = &ecs.TransformComponent{
			Position:    sprec.ZeroVec3(),
			Orientation: ecs.NewOrientation(),
		}
		entity.RenderMesh = &ecs.RenderMesh{
			GeomProgram: data.TerrainProgram.Get(),
			Mesh:        staticMesh,
		}
	}

	s.collisionMeshes = level.CollisionMeshes

	for _, staticEntity := range level.StaticEntities {
		entity := s.ecsManager.CreateEntity()
		entity.Transform = &ecs.TransformComponent{
			Position: staticEntity.Matrix.Translation(),
			Orientation: ecs.Orientation{
				VectorX: staticEntity.Matrix.OrientationX(),
				VectorY: staticEntity.Matrix.OrientationY(),
				VectorZ: staticEntity.Matrix.OrientationZ(),
			},
		}
		entity.RenderModel = &ecs.RenderModel{
			GeomProgram: data.EntityProgram.Get(),
			Model:       staticEntity.Model.Get(),
		}
	}

	carProgram := data.CarProgram.Get()
	carModel := data.CarModel.Get()
	// s.spawnCar(carProgram, carModel)

	flWheelNode, _ := carModel.FindNode("wheel_front_right")
	topEntity := s.ecsManager.CreateEntity()
	// topEntity.Debug = &ecs.DebugComponent{
	// 	Name: "top-entity",
	// }
	topEntity.Transform = &ecs.TransformComponent{
		Position:    sprec.NewVec3(0.0, entityVisualHeight, 0.0),
		Orientation: ecs.NewOrientation(),
	}
	topEntity.Motion = &ecs.MotionComponent{
		Mass: 30.0,
		MomentOfInertia: sprec.NewMat3(
			0.8, 0.0, 0.0,
			0.0, 0.8, 0.0,
			0.0, 0.0, 0.8,
		),
		DragFactor:        10.0,
		AngularDragFactor: 0.0,
	}
	topEntity.RenderMesh = &ecs.RenderMesh{
		GeomProgram: carProgram,
		Mesh:        flWheelNode.Mesh,
	}
	s.ecsPhysicsSystem.AddConstraint(constraint.FixedPosition{
		Entity:   topEntity,
		Position: sprec.NewVec3(0.0, entityVisualHeight, 0.0),
	})

	middleEntity := s.ecsManager.CreateEntity()
	middleEntity.Transform = &ecs.TransformComponent{
		Position:    sprec.NewVec3(1.4, entityVisualHeight, 0.0),
		Orientation: ecs.NewOrientation(),
	}
	middleEntity.Motion = &ecs.MotionComponent{
		Mass: 30.0,
		MomentOfInertia: sprec.NewMat3(
			0.8, 0.0, 0.0,
			0.0, 0.8, 0.0,
			0.0, 0.0, 0.8,
		),
		DragFactor:        10.0,
		AngularDragFactor: 0.0,
	}
	middleEntity.RenderMesh = &ecs.RenderMesh{
		GeomProgram: carProgram,
		Mesh:        flWheelNode.Mesh,
	}
	s.ecsPhysicsSystem.AddConstraint(constraint.Rope{
		First:        topEntity,
		FirstAnchor:  sprec.NewVec3(0.2, 0.0, 0.0),
		Second:       middleEntity,
		SecondAnchor: sprec.NewVec3(-0.2, 0.0, 0.0),
		Length:       1.0,
	})

	bottomEntity := s.ecsManager.CreateEntity()
	bottomEntity.Transform = &ecs.TransformComponent{
		Position:    sprec.NewVec3(2.8, entityVisualHeight, 0.0),
		Orientation: ecs.NewOrientation(),
	}
	bottomEntity.Motion = &ecs.MotionComponent{
		Mass: 30.0,
		MomentOfInertia: sprec.NewMat3(
			0.8, 0.0, 0.0,
			0.0, 0.8, 0.0,
			0.0, 0.0, 0.8,
		),
		// AngularVelocity:   sprec.NewVec3(100.0, 0.0, 0.0),
		DragFactor:        10.0,
		AngularDragFactor: 0.0,
	}
	bottomEntity.RenderMesh = &ecs.RenderMesh{
		GeomProgram: carProgram,
		Mesh:        flWheelNode.Mesh,
	}
	s.ecsPhysicsSystem.AddConstraint(constraint.Rope{
		First:        middleEntity,
		FirstAnchor:  sprec.NewVec3(0.2, 0.0, 0.0),
		Second:       bottomEntity,
		SecondAnchor: sprec.NewVec3(-0.2, 0.0, 0.0),
		Length:       1.0,
	})

	// s.ecsPhysicsSystem.AddConstraint(ecs.FixedConstraint{

	// })
	// anchorEntity := s.ecsManager.CreateEntity()
	// anchorEntity.FixedContraint = &ecs.FixedContraintComponent{

	standEntity := s.ecsManager.CreateEntity()
	standEntity.CameraStand = &ecs.CameraStand{
		Target:         topEntity,
		Camera:         camera,
		AnchorPosition: sprec.NewVec3(0.0, entityVisualHeight, -cameraDistance),
		AnchorDistance: anchorDistance,
		CameraDistance: cameraDistance,
	}

	{
		entity := s.ecsManager.CreateEntity()
		entity.RenderSkybox = &ecs.RenderSkybox{
			Program: data.SkyboxProgram.Get(),
			Texture: level.SkyboxTexture.Get(),
			Mesh:    data.SkyboxMesh.Get(),
		}
	}
}

func (s *Stage) spawnCar(carProgram *graphics.Program, carModel *stream.Model) *ecs.Entity {
	bodyNode, _ := carModel.FindNode("body")
	flWheelNode, _ := carModel.FindNode("wheel_front_left")
	frWheelNode, _ := carModel.FindNode("wheel_front_right")
	blWheelNode, _ := carModel.FindNode("wheel_back_left")
	brWheelNode, _ := carModel.FindNode("wheel_back_right")

	createWheelEntity := func(node *stream.Node, isDriven bool) *ecs.Entity {
		wheelEntity := s.ecsManager.CreateEntity()
		wheelEntity.Transform = &ecs.TransformComponent{
			Position:    sprec.ZeroVec3(),
			Orientation: ecs.NewOrientation(),
		}
		wheelEntity.Wheel = &ecs.Wheel{
			IsDriven:       isDriven,
			Length:         0.4,
			Radius:         0.3,
			AnchorPosition: sprec.Vec3Diff(node.Matrix.Translation(), bodyNode.Matrix.Translation()),
		}
		wheelEntity.RenderMesh = &ecs.RenderMesh{
			GeomProgram: carProgram,
			Mesh:        node.Mesh,
		}
		return wheelEntity
	}

	carEntity := s.ecsManager.CreateEntity()
	carEntity.Transform = &ecs.TransformComponent{
		Position:    sprec.ZeroVec3(),
		Orientation: ecs.NewOrientation(),
	}
	carEntity.RenderMesh = &ecs.RenderMesh{
		GeomProgram: carProgram,
		Mesh:        bodyNode.Mesh,
	}
	carEntity.Input = &ecs.Input{}
	carEntity.Vehicle = &ecs.Vehicle{
		Position:    sprec.NewVec3(0.0, carDropHeight, 0.0),
		Orientation: ecs.NewOrientation(),

		FLWheel: createWheelEntity(flWheelNode, true),
		FRWheel: createWheelEntity(frWheelNode, true),
		BLWheel: createWheelEntity(blWheelNode, false),
		BRWheel: createWheelEntity(brWheelNode, false),
	}
	return carEntity
}

func (s *Stage) Resize(width, height int) {
	s.screenFramebuffer.Width = int32(width)
	s.screenFramebuffer.Height = int32(height)
}

func (s *Stage) Update(elapsedTime time.Duration, camera *ecs.Camera, input ecs.CarInput) {
	s.ecsPhysicsSystem.Update(elapsedTime)
	s.ecsVehicleSystem.Update(elapsedTime, input)
	s.ecsCameraStandSystem.Update()
}

func (s *Stage) Render(pipeline *graphics.Pipeline, camera *ecs.Camera) {
	geometrySequence := pipeline.BeginSequence()
	geometrySequence.TargetFramebuffer = s.geometryFramebuffer
	geometrySequence.BackgroundColor = sprec.NewVec4(0.0, 0.6, 1.0, 1.0)
	geometrySequence.ClearColor = true
	geometrySequence.ClearDepth = true
	geometrySequence.DepthFunc = graphics.DepthFuncLessOrEqual
	geometrySequence.ProjectionMatrix = camera.ProjectionMatrix()
	geometrySequence.ViewMatrix = camera.InverseViewMatrix()
	s.ecsRenderer.Render(geometrySequence)
	pipeline.EndSequence(geometrySequence)

	lightingSequence := pipeline.BeginSequence()
	lightingSequence.SourceFramebuffer = s.geometryFramebuffer
	lightingSequence.TargetFramebuffer = s.lightingFramebuffer
	lightingSequence.BlitFramebufferDepth = true
	lightingSequence.ClearColor = true
	// FIXME: this is only for directional... Will need sub-sequences
	lightingSequence.TestDepth = false
	quadItem := lightingSequence.BeginItem()
	quadItem.Program = s.lightingProgram
	quadItem.VertexArray = s.quadMesh.VertexArray
	quadItem.IndexCount = s.quadMesh.SubMeshes[0].IndexCount
	lightingSequence.EndItem(quadItem)
	pipeline.EndSequence(lightingSequence)

	screenSequence := pipeline.BeginSequence()
	screenSequence.SourceFramebuffer = s.lightingFramebuffer
	screenSequence.TargetFramebuffer = s.screenFramebuffer
	screenSequence.BlitFramebufferColor = true
	screenSequence.BlitFramebufferSmooth = true
	pipeline.EndSequence(screenSequence)
}

func (s *Stage) CheckCollision(line collision.Line) (bestCollision collision.LineCollision, found bool) {
	closestDistance := line.LengthSquared()
	for _, mesh := range s.collisionMeshes {
		if lineCollision, ok := mesh.LineCollision(line); ok {
			found = true
			distanceVector := sprec.Vec3Diff(lineCollision.Intersection(), line.Start())
			if distance := distanceVector.SqrLength(); distance < closestDistance {
				closestDistance = distance
				bestCollision = lineCollision
			}
		}
	}
	return
}
