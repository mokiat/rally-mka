package scene

import (
	"time"

	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/ecs"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/ecs/constraint"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/ecs/system"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/scene/car"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/stream"
	"github.com/mokiat/rally-mka/internal/data"
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

const maxDebugLines = 1024

var arrayTask *graphics.Task

func NewStage(gfxWorker *graphics.Worker) *Stage {
	indexData := make([]byte, maxDebugLines*2)
	for i := 0; i < maxDebugLines; i++ {
		data.Buffer(indexData).SetUInt16(i*2, uint16(i))
	}
	vertexData := make([]byte, maxDebugLines*4*7*2)
	debugVertexArrayData := graphics.VertexArrayData{
		VertexData:   vertexData,
		VertexStride: 4 * 7,
		CoordOffset:  0,
		ColorOffset:  4 * 3,
		IndexData:    indexData,
	}
	debugVertexArray := &graphics.VertexArray{}
	arrayTask = gfxWorker.Schedule(func() error {
		if err := debugVertexArray.Allocate(debugVertexArrayData); err != nil {
			panic(err)
		}
		return nil
	}) // FIXME: Race condition

	ecsManager := ecs.NewManager()
	stage := &Stage{
		ecsManager:           ecsManager,
		ecsRenderer:          ecs.NewRenderer(ecsManager),
		ecsCarSystem:         system.NewCarSystem(ecsManager),
		ecsCameraStandSystem: ecs.NewCameraStandSystem(ecsManager),
		ecsPhysicsSystem:     ecs.NewPhysicsSystem(ecsManager, 15*time.Millisecond),
		screenFramebuffer:    &graphics.Framebuffer{},
		debugVertexArray:     debugVertexArray,
		debugVertexArrayData: debugVertexArrayData,
	}
	return stage
}

type Stage struct {
	ecsManager           *ecs.Manager
	ecsRenderer          *ecs.Renderer
	ecsCarSystem         *system.CarSystem
	ecsCameraStandSystem *ecs.CameraStandSystem
	ecsPhysicsSystem     *ecs.PhysicsSystem

	geometryFramebuffer *graphics.Framebuffer
	lightingFramebuffer *graphics.Framebuffer
	screenFramebuffer   *graphics.Framebuffer
	lightingProgram     *graphics.Program
	quadMesh            *stream.Mesh

	debugProgram         *graphics.Program
	debugVertexArray     *graphics.VertexArray
	debugVertexArrayData graphics.VertexArrayData

	collisionMeshes []*collision.Mesh
}

func (s *Stage) Init(data *Data, camera *ecs.Camera) {
	level := data.Level.Get()

	s.debugProgram = data.DebugProgram.Get()

	s.geometryFramebuffer = data.GeometryFramebuffer
	s.lightingFramebuffer = data.LightingFramebuffer

	s.lightingProgram = data.DeferredLightingProgram.Get()
	s.quadMesh = data.QuadMesh.Get()

	for _, staticMesh := range level.StaticMeshes {
		entity := s.ecsManager.CreateEntity()
		entity.Transform = &ecs.TransformComponent{
			Position:    sprec.ZeroVec3(),
			Orientation: sprec.IdentityQuat(),
		}
		entity.RenderMesh = &ecs.RenderMesh{
			GeomProgram: data.TerrainProgram.Get(),
			Mesh:        staticMesh,
		}
	}

	s.collisionMeshes = level.CollisionMeshes
	for _, collisionMesh := range level.CollisionMeshes {
		entity := s.ecsManager.CreateEntity()
		entity.Transform = &ecs.TransformComponent{
			Position:    sprec.ZeroVec3(),
			Orientation: sprec.IdentityQuat(),
		}
		entity.Collision = &ecs.CollisionComponent{
			RestitutionCoef: 1.0,
			CollisionShape: ecs.MeshShape{
				Mesh: collisionMesh,
			},
		}
	}

	for _, staticEntity := range level.StaticEntities {
		entity := s.ecsManager.CreateEntity()
		entity.Transform = &ecs.TransformComponent{
			Position:    staticEntity.Matrix.Translation(),
			Orientation: sprec.IdentityQuat(),
			// FIXME
			// Orientation: ecs.Orientation{
			// 	VectorX: staticEntity.Matrix.OrientationX(),
			// 	VectorY: staticEntity.Matrix.OrientationY(),
			// 	VectorZ: staticEntity.Matrix.OrientationZ(),
			// },
		}
		entity.RenderModel = &ecs.RenderModel{
			GeomProgram: data.EntityProgram.Get(),
			Model:       staticEntity.Model.Get(),
		}
	}

	carProgram := data.CarProgram.Get()
	carModel := data.CarModel.Get()

	// ----------------------------------------------

	// targetEntity := s.spawnCar(carProgram, carModel)

	// ----------------------------------------------

	// targetEntity :=
	// 	s.setupChandelierDemo(carProgram, carModel, sprec.NewVec3(0.0, 10.0, 0.0))

	// ----------------------------------------------

	// targetEntity :=
	// 	s.setupCoiloverDemo(carProgram, carModel, sprec.NewVec3(0.0, 10.0, -5.0))

	// ----------------------------------------------

	// targetEntity :=
	// 	s.setupRodDemo(carProgram, carModel, sprec.NewVec3(0.0, 10.0, 5.0))

	// ----------------------------------------------

	targetEntity :=
		s.setupCarDemo(carProgram, carModel, sprec.NewVec3(0.0, 2.0, 10.0))

	// ----------------------------------------------

	standTarget := targetEntity
	standEntity := s.ecsManager.CreateEntity()
	standEntity.CameraStand = &ecs.CameraStand{
		Target:         standTarget,
		Camera:         camera,
		AnchorPosition: sprec.Vec3Sum(standTarget.Transform.Position, sprec.NewVec3(0.0, 0.0, -cameraDistance)),
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

func (s *Stage) setupChandelierDemo(program *graphics.Program, model *stream.Model, position sprec.Vec3) *ecs.Entity {
	fakeFixtureTire := car.Tire(program, model, car.FrontRightTireLocation).
		WithPosition(position).
		Build(s.ecsManager)
	s.ecsPhysicsSystem.AddConstraint(constraint.FixedTranslation{
		Entity:   fakeFixtureTire,
		Position: position,
	})

	playTire := car.Tire(program, model, car.FrontRightTireLocation).
		WithPosition(sprec.Vec3Sum(position, sprec.NewVec3(-2.3, 0.0, 0.0))).
		Build(s.ecsManager)
	s.ecsPhysicsSystem.AddConstraint(constraint.Chandelier{
		Entity:       playTire,
		EntityAnchor: sprec.NewVec3(0.3, 0.0, 0.0),
		Length:       2.0,
		Fixture:      position,
	})

	return fakeFixtureTire
}

func (s *Stage) setupRodDemo(program *graphics.Program, model *stream.Model, position sprec.Vec3) *ecs.Entity {
	topTirePosition := position
	topTire := car.Tire(program, model, car.FrontRightTireLocation).
		WithPosition(topTirePosition).
		Build(s.ecsManager)
	s.ecsPhysicsSystem.AddConstraint(constraint.FixedTranslation{
		Entity:   topTire,
		Position: topTirePosition,
	})

	middleTirePosition := sprec.Vec3Sum(topTirePosition, sprec.NewVec3(1.4, 0.0, 0.0))
	middleTire := car.Tire(program, model, car.FrontRightTireLocation).
		WithPosition(middleTirePosition).
		Build(s.ecsManager)
	s.ecsPhysicsSystem.AddConstraint(constraint.HingedRod{
		First:        topTire,
		FirstAnchor:  sprec.NewVec3(0.2, 0.0, 0.0),
		Second:       middleTire,
		SecondAnchor: sprec.NewVec3(-0.2, 0.0, 0.0),
		Length:       1.0,
	})

	bottomTirePosition := sprec.Vec3Sum(middleTirePosition, sprec.NewVec3(1.4, 0.0, 0.0))
	bottomTire := car.Tire(program, model, car.FrontRightTireLocation).
		WithPosition(bottomTirePosition).
		Build(s.ecsManager)
	s.ecsPhysicsSystem.AddConstraint(constraint.HingedRod{
		First:        middleTire,
		FirstAnchor:  sprec.NewVec3(0.2, 0.0, 0.0),
		Second:       bottomTire,
		SecondAnchor: sprec.NewVec3(-0.2, 0.0, 0.0),
		Length:       1.0,
	})

	return topTire
}

func (s *Stage) setupCoiloverDemo(program *graphics.Program, model *stream.Model, position sprec.Vec3) *ecs.Entity {
	fixtureTire := car.Tire(program, model, car.FrontRightTireLocation).
		WithPosition(position).
		Build(s.ecsManager)
	s.ecsPhysicsSystem.AddConstraint(constraint.FixedTranslation{
		Entity:   fixtureTire,
		Position: position,
	})

	fallingTire := car.Tire(program, model, car.FrontRightTireLocation).
		WithPosition(sprec.Vec3Sum(position, sprec.NewVec3(0.0, -0.8, 0.0))).
		Build(s.ecsManager)
	s.ecsPhysicsSystem.AddConstraint(constraint.Spring{
		Target:    fixtureTire,
		Entity:    fallingTire,
		Length:    1.0,
		Stiffness: 100.0,
	})
	s.ecsPhysicsSystem.AddConstraint(constraint.Damper{
		Target:   fixtureTire,
		Entity:   fallingTire,
		Strength: 20.0,
	})
	return fixtureTire
}

func (s *Stage) setupCarDemo(program *graphics.Program, model *stream.Model, position sprec.Vec3) *ecs.Entity {
	chasis := car.Chassis(program, model).
		WithPosition(position).
		Build(s.ecsManager)
	// chasis.Motion.AngularVelocity = sprec.NewVec3(0.0, -0.5, 0.0)
	// chasis.Motion.AngularVelocity = sprec.NewVec3(0.0, 0.0, 1.0)
	// s.ecsPhysicsSystem.AddConstraint(constraint.FixedTranslation{
	// 	Entity:   chasis,
	// 	Position: position,
	// })

	// suspensionLength := float32(1.0)
	// suspensionStiffness := float32(5000.0)
	// suspensionDampness := float32(0.4)
	suspensionLength := float32(0.5)
	suspensionStiffness := float32(6000.0)
	suspensionDampness := float32(0.9)

	flTireRelativePosition := sprec.NewVec3(1.0, -0.6-suspensionLength/2.0, 1.25)
	flTire := car.Tire(program, model, car.FrontLeftTireLocation).
		WithPosition(sprec.Vec3Sum(position, flTireRelativePosition)).
		Build(s.ecsManager)
	s.ecsPhysicsSystem.AddConstraint(constraint.CopyTranslation{
		Target:         chasis,
		Entity:         flTire,
		RelativeOffset: flTireRelativePosition,
		SkipY:          true,
	})
	flRotation := &constraint.CopyAxis{
		Target:       chasis,
		TargetOffset: sprec.IdentityQuat(),
		TargetAxis:   sprec.BasisXVec3(),
		Entity:       flTire,
		EntityAxis:   sprec.BasisXVec3(),
	}
	s.ecsPhysicsSystem.AddConstraint(flRotation)
	flSpringAttachmentRelativePosition := sprec.Vec3Sum(flTireRelativePosition, sprec.NewVec3(0.0, suspensionLength, 0.0))
	s.ecsPhysicsSystem.AddConstraint(constraint.Spring{
		Target:               chasis,
		TargetRelativeOffset: flSpringAttachmentRelativePosition,
		Entity:               flTire,
		Stiffness:            suspensionStiffness,
		Length:               suspensionLength,
	})
	s.ecsPhysicsSystem.AddConstraint(constraint.Damper{
		Target:               chasis,
		TargetRelativeOffset: flSpringAttachmentRelativePosition,
		Entity:               flTire,
		Strength:             suspensionDampness,
	})

	frTireRelativePosition := sprec.NewVec3(-1.0, -0.6-suspensionLength/2.0, 1.25)
	frTire := car.Tire(program, model, car.FrontRightTireLocation).
		WithPosition(sprec.Vec3Sum(position, frTireRelativePosition)).
		Build(s.ecsManager)
	s.ecsPhysicsSystem.AddConstraint(constraint.CopyTranslation{
		Target:         chasis,
		Entity:         frTire,
		RelativeOffset: frTireRelativePosition,
		SkipY:          true,
	})
	frRotation := &constraint.CopyAxis{
		Target:       chasis,
		TargetOffset: sprec.IdentityQuat(),
		TargetAxis:   sprec.BasisXVec3(),
		Entity:       frTire,
		EntityAxis:   sprec.BasisXVec3(),
	}
	s.ecsPhysicsSystem.AddConstraint(frRotation)
	frSpringAttachmentRelativePosition := sprec.Vec3Sum(frTireRelativePosition, sprec.NewVec3(0.0, suspensionLength, 0.0))
	s.ecsPhysicsSystem.AddConstraint(constraint.Spring{
		Target:               chasis,
		TargetRelativeOffset: frSpringAttachmentRelativePosition,
		Entity:               frTire,
		Stiffness:            suspensionStiffness,
		Length:               suspensionLength,
	})
	s.ecsPhysicsSystem.AddConstraint(constraint.Damper{
		Target:               chasis,
		TargetRelativeOffset: frSpringAttachmentRelativePosition,
		Entity:               frTire,
		Strength:             suspensionDampness,
	})

	blTireRelativePosition := sprec.NewVec3(1.0, -0.6-suspensionLength/2.0, -1.45)
	blTire := car.Tire(program, model, car.BackLeftTireLocation).
		WithPosition(sprec.Vec3Sum(position, blTireRelativePosition)).
		Build(s.ecsManager)
	s.ecsPhysicsSystem.AddConstraint(constraint.CopyTranslation{
		Target:         chasis,
		Entity:         blTire,
		RelativeOffset: blTireRelativePosition,
		SkipY:          true,
	})
	s.ecsPhysicsSystem.AddConstraint(constraint.CopyAxis{
		Target:       chasis,
		TargetAxis:   sprec.BasisXVec3(),
		TargetOffset: sprec.IdentityQuat(),
		Entity:       blTire,
		EntityAxis:   sprec.BasisXVec3(),
	})
	blSpringAttachmentRelativePosition := sprec.Vec3Sum(blTireRelativePosition, sprec.NewVec3(0.0, suspensionLength, 0.0))
	s.ecsPhysicsSystem.AddConstraint(constraint.Spring{
		Target:               chasis,
		TargetRelativeOffset: blSpringAttachmentRelativePosition,
		Entity:               blTire,
		Stiffness:            suspensionStiffness,
		Length:               suspensionLength,
	})
	s.ecsPhysicsSystem.AddConstraint(constraint.Damper{
		Target:               chasis,
		TargetRelativeOffset: blSpringAttachmentRelativePosition,
		Entity:               blTire,
		Strength:             suspensionDampness,
	})

	brTireRelativePosition := sprec.NewVec3(-1.0, -0.6-suspensionLength/2.0, -1.45)
	brTire := car.Tire(program, model, car.BackRightTireLocation).
		WithPosition(sprec.Vec3Sum(position, brTireRelativePosition)).
		Build(s.ecsManager)
	s.ecsPhysicsSystem.AddConstraint(constraint.CopyTranslation{
		Target:         chasis,
		Entity:         brTire,
		RelativeOffset: brTireRelativePosition,
		SkipY:          true,
	})
	s.ecsPhysicsSystem.AddConstraint(constraint.CopyAxis{
		Target:       chasis,
		TargetAxis:   sprec.BasisXVec3(),
		TargetOffset: sprec.IdentityQuat(),
		Entity:       brTire,
		EntityAxis:   sprec.BasisXVec3(),
	})
	brSpringAttachmentRelativePosition := sprec.Vec3Sum(brTireRelativePosition, sprec.NewVec3(0.0, suspensionLength, 0.0))
	s.ecsPhysicsSystem.AddConstraint(constraint.Spring{
		Target:               chasis,
		TargetRelativeOffset: brSpringAttachmentRelativePosition,
		Entity:               brTire,
		Stiffness:            suspensionStiffness,
		Length:               suspensionLength,
	})
	s.ecsPhysicsSystem.AddConstraint(constraint.Damper{
		Target:               chasis,
		TargetRelativeOffset: brSpringAttachmentRelativePosition,
		Entity:               brTire,
		Strength:             suspensionDampness,
	})

	car := s.ecsManager.CreateEntity()
	car.Car = &ecs.Car{
		Body:            chasis,
		FLWheelRotation: flRotation,
		FLWheel:         flTire,
		FRWheelRotation: frRotation,
		FRWheel:         frTire,
		BLWheel:         blTire,
		BRWheel:         brTire,
	}

	return chasis
}

func (s *Stage) Resize(width, height int) {
	s.screenFramebuffer.Width = int32(width)
	s.screenFramebuffer.Height = int32(height)
}

func (s *Stage) Update(elapsedTime time.Duration, camera *ecs.Camera, input ecs.CarInput) {
	s.ecsPhysicsSystem.Update(elapsedTime)
	s.ecsCarSystem.Update(elapsedTime, input)
	s.ecsCameraStandSystem.Update()
}

func (s *Stage) Render(pipeline *graphics.Pipeline, camera *ecs.Camera) {
	if !arrayTask.Done() {
		panic("NOT DONE!")
	}

	pipeline.SchedulePreRender(func() {
		if err := s.debugVertexArray.Update(s.debugVertexArrayData); err != nil {
			panic(err)
		}
	})

	geometrySequence := pipeline.BeginSequence()
	geometrySequence.TargetFramebuffer = s.geometryFramebuffer
	geometrySequence.BackgroundColor = sprec.NewVec4(0.0, 0.6, 1.0, 1.0)
	geometrySequence.ClearColor = true
	geometrySequence.ClearDepth = true
	geometrySequence.DepthFunc = graphics.DepthFuncLessOrEqual
	geometrySequence.ProjectionMatrix = camera.ProjectionMatrix()
	geometrySequence.ViewMatrix = camera.InverseViewMatrix()
	s.renderDebugLines(geometrySequence, s.ecsPhysicsSystem.GetDebug())
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

func (s *Stage) renderDebugLines(sequence *graphics.Sequence, lines []ecs.DebugLine) {
	for i, line := range lines {
		vertexStride := 4 * 7 * 2
		data.Buffer(s.debugVertexArrayData.VertexData).SetFloat32(vertexStride*i+0, line.A.X)
		data.Buffer(s.debugVertexArrayData.VertexData).SetFloat32(vertexStride*i+4, line.A.Y)
		data.Buffer(s.debugVertexArrayData.VertexData).SetFloat32(vertexStride*i+8, line.A.Z)
		data.Buffer(s.debugVertexArrayData.VertexData).SetFloat32(vertexStride*i+12, line.Color.X)
		data.Buffer(s.debugVertexArrayData.VertexData).SetFloat32(vertexStride*i+16, line.Color.Y)
		data.Buffer(s.debugVertexArrayData.VertexData).SetFloat32(vertexStride*i+20, line.Color.Z)
		data.Buffer(s.debugVertexArrayData.VertexData).SetFloat32(vertexStride*i+24, line.Color.W)
		data.Buffer(s.debugVertexArrayData.VertexData).SetFloat32(vertexStride*i+28, line.B.X)
		data.Buffer(s.debugVertexArrayData.VertexData).SetFloat32(vertexStride*i+32, line.B.Y)
		data.Buffer(s.debugVertexArrayData.VertexData).SetFloat32(vertexStride*i+36, line.B.Z)
		data.Buffer(s.debugVertexArrayData.VertexData).SetFloat32(vertexStride*i+40, line.Color.X)
		data.Buffer(s.debugVertexArrayData.VertexData).SetFloat32(vertexStride*i+44, line.Color.Y)
		data.Buffer(s.debugVertexArrayData.VertexData).SetFloat32(vertexStride*i+48, line.Color.Z)
		data.Buffer(s.debugVertexArrayData.VertexData).SetFloat32(vertexStride*i+52, line.Color.W)
	}

	item := sequence.BeginItem()
	item.Primitive = graphics.RenderPrimitiveLines
	item.Program = s.debugProgram
	item.ModelMatrix = sprec.IdentityMat4()
	item.VertexArray = s.debugVertexArray
	item.IndexCount = int32(len(lines) * 2)
	sequence.EndItem(item)
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
