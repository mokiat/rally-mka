package scene

import (
	"time"

	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/ecs"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/ecs/system"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/scene/car"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/stream"
	"github.com/mokiat/rally-mka/internal/data"
	"github.com/mokiat/rally-mka/internal/engine/graphics"
	"github.com/mokiat/rally-mka/internal/engine/physics"
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
		physicsEngine:        physics.NewEngine(15 * time.Millisecond),
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
	physicsEngine        *physics.Engine

	geometryFramebuffer *graphics.Framebuffer
	lightingFramebuffer *graphics.Framebuffer
	screenFramebuffer   *graphics.Framebuffer
	lightingProgram     *graphics.Program
	quadMesh            *stream.Mesh

	debugProgram         *graphics.Program
	debugVertexArray     *graphics.VertexArray
	debugVertexArrayData graphics.VertexArrayData
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
		entity.Render = &ecs.RenderComponent{
			GeomProgram: data.TerrainProgram.Get(),
			Mesh:        staticMesh,
			Matrix:      sprec.IdentityMat4(),
		}
	}

	for _, collisionMesh := range level.CollisionMeshes {
		s.physicsEngine.AddBody(&physics.Body{
			Position:        sprec.ZeroVec3(),
			Orientation:     sprec.IdentityQuat(),
			IsStatic:        true,
			RestitutionCoef: 1.0,
			CollisionShape: physics.MeshShape{
				Mesh: collisionMesh,
			},
		})
	}

	for _, staticEntity := range level.StaticEntities {
		entity := s.ecsManager.CreateEntity()
		entity.Render = &ecs.RenderComponent{
			GeomProgram: data.EntityProgram.Get(),
			Model:       staticEntity.Model.Get(),
			Matrix:      staticEntity.Matrix,
		}
	}

	carProgram := data.CarProgram.Get()
	carModel := data.CarModel.Get()

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
		AnchorPosition: sprec.Vec3Sum(standTarget.Physics.Body.Position, sprec.NewVec3(0.0, 0.0, -cameraDistance)),
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
	s.physicsEngine.AddBody(fakeFixtureTire.Physics.Body)
	s.physicsEngine.AddConstraint(physics.FixedTranslationConstraint{
		Fixture: position,
		Body:    fakeFixtureTire.Physics.Body,
	})

	playTire := car.Tire(program, model, car.FrontRightTireLocation).
		WithPosition(sprec.Vec3Sum(position, sprec.NewVec3(-2.3, 0.0, 0.0))).
		Build(s.ecsManager)
	s.physicsEngine.AddBody(playTire.Physics.Body)
	s.physicsEngine.AddConstraint(physics.ChandelierConstraint{
		Fixture:    position,
		Body:       playTire.Physics.Body,
		BodyAnchor: sprec.NewVec3(0.3, 0.0, 0.0),
		Length:     2.0,
	})

	return fakeFixtureTire
}

func (s *Stage) setupCoiloverDemo(program *graphics.Program, model *stream.Model, position sprec.Vec3) *ecs.Entity {
	fixtureTire := car.Tire(program, model, car.FrontRightTireLocation).
		WithPosition(position).
		Build(s.ecsManager)
	s.physicsEngine.AddBody(fixtureTire.Physics.Body)
	s.physicsEngine.AddConstraint(physics.FixedTranslationConstraint{
		Fixture: position,
		Body:    fixtureTire.Physics.Body,
	})

	fallingTire := car.Tire(program, model, car.FrontRightTireLocation).
		WithPosition(sprec.Vec3Sum(position, sprec.NewVec3(0.0, -0.8, 0.0))).
		Build(s.ecsManager)
	s.physicsEngine.AddBody(fallingTire.Physics.Body)
	s.physicsEngine.AddConstraint(physics.SpringConstraint{
		FirstBody:  fixtureTire.Physics.Body,
		SecondBody: fallingTire.Physics.Body,
		Length:     1.0,
		Stiffness:  100.0,
	})
	s.physicsEngine.AddConstraint(physics.DamperConstraint{
		FirstBody:  fixtureTire.Physics.Body,
		SecondBody: fallingTire.Physics.Body,
		Strength:   20.0,
	})
	return fixtureTire
}

func (s *Stage) setupRodDemo(program *graphics.Program, model *stream.Model, position sprec.Vec3) *ecs.Entity {
	topTirePosition := position
	topTire := car.Tire(program, model, car.FrontRightTireLocation).
		WithPosition(topTirePosition).
		Build(s.ecsManager)
	s.physicsEngine.AddBody(topTire.Physics.Body)
	s.physicsEngine.AddConstraint(physics.FixedTranslationConstraint{
		Fixture: topTirePosition,
		Body:    topTire.Physics.Body,
	})

	middleTirePosition := sprec.Vec3Sum(topTirePosition, sprec.NewVec3(1.4, 0.0, 0.0))
	middleTire := car.Tire(program, model, car.FrontRightTireLocation).
		WithPosition(middleTirePosition).
		Build(s.ecsManager)
	s.physicsEngine.AddBody(middleTire.Physics.Body)
	s.physicsEngine.AddConstraint(physics.HingedRodConstraint{
		FirstBody:        topTire.Physics.Body,
		FirstBodyAnchor:  sprec.NewVec3(0.2, 0.0, 0.0),
		SecondBody:       middleTire.Physics.Body,
		SecondBodyAnchor: sprec.NewVec3(-0.2, 0.0, 0.0),
		Length:           1.0,
	})

	bottomTirePosition := sprec.Vec3Sum(middleTirePosition, sprec.NewVec3(1.4, 0.0, 0.0))
	bottomTire := car.Tire(program, model, car.FrontRightTireLocation).
		WithPosition(bottomTirePosition).
		Build(s.ecsManager)
	s.physicsEngine.AddBody(bottomTire.Physics.Body)
	s.physicsEngine.AddConstraint(physics.HingedRodConstraint{
		FirstBody:        middleTire.Physics.Body,
		FirstBodyAnchor:  sprec.NewVec3(0.2, 0.0, 0.0),
		SecondBody:       bottomTire.Physics.Body,
		SecondBodyAnchor: sprec.NewVec3(-0.2, 0.0, 0.0),
		Length:           1.0,
	})

	return topTire
}

func (s *Stage) setupCarDemo(program *graphics.Program, model *stream.Model, position sprec.Vec3) *ecs.Entity {
	chasis := car.Chassis(program, model).
		WithPosition(position).
		Build(s.ecsManager)
	s.physicsEngine.AddBody(chasis.Physics.Body)
	// chasis.Motion.AngularVelocity = sprec.NewVec3(0.0, -0.5, 0.0)
	// chasis.Motion.AngularVelocity = sprec.NewVec3(0.0, 0.0, 1.0)
	// s.physicsEngine.AddConstraint(constraint.FixedTranslation{
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
	s.physicsEngine.AddBody(flTire.Physics.Body)
	s.physicsEngine.AddConstraint(physics.MatchTranslationConstraint{
		FirstBody:       chasis.Physics.Body,
		FirstBodyAnchor: flTireRelativePosition,
		SecondBody:      flTire.Physics.Body,
		IgnoreY:         true,
	})
	flRotation := &physics.MatchAxisConstraint{
		FirstBody:      chasis.Physics.Body,
		FirstBodyAxis:  sprec.BasisXVec3(),
		SecondBody:     flTire.Physics.Body,
		SecondBodyAxis: sprec.BasisXVec3(),
	}
	s.physicsEngine.AddConstraint(flRotation)
	flSpringAttachmentRelativePosition := sprec.Vec3Sum(flTireRelativePosition, sprec.NewVec3(0.0, suspensionLength, 0.0))
	s.physicsEngine.AddConstraint(physics.SpringConstraint{
		FirstBody:       chasis.Physics.Body,
		FirstBodyAnchor: flSpringAttachmentRelativePosition,
		SecondBody:      flTire.Physics.Body,
		Length:          suspensionLength,
		Stiffness:       suspensionStiffness,
	})
	s.physicsEngine.AddConstraint(physics.DamperConstraint{
		FirstBody:       chasis.Physics.Body,
		FirstBodyAnchor: flSpringAttachmentRelativePosition,
		SecondBody:      flTire.Physics.Body,
		Strength:        suspensionDampness,
	})

	frTireRelativePosition := sprec.NewVec3(-1.0, -0.6-suspensionLength/2.0, 1.25)
	frTire := car.Tire(program, model, car.FrontRightTireLocation).
		WithPosition(sprec.Vec3Sum(position, frTireRelativePosition)).
		Build(s.ecsManager)
	s.physicsEngine.AddBody(frTire.Physics.Body)
	s.physicsEngine.AddConstraint(physics.MatchTranslationConstraint{
		FirstBody:       chasis.Physics.Body,
		FirstBodyAnchor: frTireRelativePosition,
		SecondBody:      frTire.Physics.Body,
		IgnoreY:         true,
	})
	frRotation := &physics.MatchAxisConstraint{
		FirstBody:      chasis.Physics.Body,
		FirstBodyAxis:  sprec.BasisXVec3(),
		SecondBody:     frTire.Physics.Body,
		SecondBodyAxis: sprec.BasisXVec3(),
	}
	s.physicsEngine.AddConstraint(frRotation)
	frSpringAttachmentRelativePosition := sprec.Vec3Sum(frTireRelativePosition, sprec.NewVec3(0.0, suspensionLength, 0.0))
	s.physicsEngine.AddConstraint(physics.SpringConstraint{
		FirstBody:       chasis.Physics.Body,
		FirstBodyAnchor: frSpringAttachmentRelativePosition,
		SecondBody:      frTire.Physics.Body,
		Length:          suspensionLength,
		Stiffness:       suspensionStiffness,
	})
	s.physicsEngine.AddConstraint(physics.DamperConstraint{
		FirstBody:       chasis.Physics.Body,
		FirstBodyAnchor: frSpringAttachmentRelativePosition,
		SecondBody:      frTire.Physics.Body,
		Strength:        suspensionDampness,
	})

	blTireRelativePosition := sprec.NewVec3(1.0, -0.6-suspensionLength/2.0, -1.45)
	blTire := car.Tire(program, model, car.BackLeftTireLocation).
		WithPosition(sprec.Vec3Sum(position, blTireRelativePosition)).
		Build(s.ecsManager)
	s.physicsEngine.AddBody(blTire.Physics.Body)
	s.physicsEngine.AddConstraint(physics.MatchTranslationConstraint{
		FirstBody:       chasis.Physics.Body,
		FirstBodyAnchor: blTireRelativePosition,
		SecondBody:      blTire.Physics.Body,
		IgnoreY:         true,
	})
	s.physicsEngine.AddConstraint(physics.MatchAxisConstraint{
		FirstBody:      chasis.Physics.Body,
		FirstBodyAxis:  sprec.BasisXVec3(),
		SecondBody:     blTire.Physics.Body,
		SecondBodyAxis: sprec.BasisXVec3(),
	})
	blSpringAttachmentRelativePosition := sprec.Vec3Sum(blTireRelativePosition, sprec.NewVec3(0.0, suspensionLength, 0.0))
	s.physicsEngine.AddConstraint(physics.SpringConstraint{
		FirstBody:       chasis.Physics.Body,
		FirstBodyAnchor: blSpringAttachmentRelativePosition,
		SecondBody:      blTire.Physics.Body,
		Length:          suspensionLength,
		Stiffness:       suspensionStiffness,
	})
	s.physicsEngine.AddConstraint(physics.DamperConstraint{
		FirstBody:       chasis.Physics.Body,
		FirstBodyAnchor: blSpringAttachmentRelativePosition,
		SecondBody:      blTire.Physics.Body,
		Strength:        suspensionDampness,
	})

	brTireRelativePosition := sprec.NewVec3(-1.0, -0.6-suspensionLength/2.0, -1.45)
	brTire := car.Tire(program, model, car.BackRightTireLocation).
		WithPosition(sprec.Vec3Sum(position, brTireRelativePosition)).
		Build(s.ecsManager)
	s.physicsEngine.AddBody(brTire.Physics.Body)
	s.physicsEngine.AddConstraint(physics.MatchTranslationConstraint{
		FirstBody:       chasis.Physics.Body,
		FirstBodyAnchor: brTireRelativePosition,
		SecondBody:      brTire.Physics.Body,
		IgnoreY:         true,
	})
	s.physicsEngine.AddConstraint(physics.MatchAxisConstraint{
		FirstBody:      chasis.Physics.Body,
		FirstBodyAxis:  sprec.BasisXVec3(),
		SecondBody:     brTire.Physics.Body,
		SecondBodyAxis: sprec.BasisXVec3(),
	})
	brSpringAttachmentRelativePosition := sprec.Vec3Sum(brTireRelativePosition, sprec.NewVec3(0.0, suspensionLength, 0.0))
	s.physicsEngine.AddConstraint(physics.SpringConstraint{
		FirstBody:       chasis.Physics.Body,
		FirstBodyAnchor: brSpringAttachmentRelativePosition,
		SecondBody:      brTire.Physics.Body,
		Length:          suspensionLength,
		Stiffness:       suspensionStiffness,
	})
	s.physicsEngine.AddConstraint(physics.DamperConstraint{
		FirstBody:       chasis.Physics.Body,
		FirstBodyAnchor: brSpringAttachmentRelativePosition,
		SecondBody:      brTire.Physics.Body,
		Strength:        suspensionDampness,
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
	s.physicsEngine.Update(elapsedTime)
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
	// s.renderDebugLines(geometrySequence, s.ecsPhysicsSystem.GetDebug())
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

// func (s *Stage) renderDebugLines(sequence *graphics.Sequence, lines []ecs.DebugLine) {
// 	for i, line := range lines {
// 		vertexStride := 4 * 7 * 2
// 		data.Buffer(s.debugVertexArrayData.VertexData).SetFloat32(vertexStride*i+0, line.A.X)
// 		data.Buffer(s.debugVertexArrayData.VertexData).SetFloat32(vertexStride*i+4, line.A.Y)
// 		data.Buffer(s.debugVertexArrayData.VertexData).SetFloat32(vertexStride*i+8, line.A.Z)
// 		data.Buffer(s.debugVertexArrayData.VertexData).SetFloat32(vertexStride*i+12, line.Color.X)
// 		data.Buffer(s.debugVertexArrayData.VertexData).SetFloat32(vertexStride*i+16, line.Color.Y)
// 		data.Buffer(s.debugVertexArrayData.VertexData).SetFloat32(vertexStride*i+20, line.Color.Z)
// 		data.Buffer(s.debugVertexArrayData.VertexData).SetFloat32(vertexStride*i+24, line.Color.W)
// 		data.Buffer(s.debugVertexArrayData.VertexData).SetFloat32(vertexStride*i+28, line.B.X)
// 		data.Buffer(s.debugVertexArrayData.VertexData).SetFloat32(vertexStride*i+32, line.B.Y)
// 		data.Buffer(s.debugVertexArrayData.VertexData).SetFloat32(vertexStride*i+36, line.B.Z)
// 		data.Buffer(s.debugVertexArrayData.VertexData).SetFloat32(vertexStride*i+40, line.Color.X)
// 		data.Buffer(s.debugVertexArrayData.VertexData).SetFloat32(vertexStride*i+44, line.Color.Y)
// 		data.Buffer(s.debugVertexArrayData.VertexData).SetFloat32(vertexStride*i+48, line.Color.Z)
// 		data.Buffer(s.debugVertexArrayData.VertexData).SetFloat32(vertexStride*i+52, line.Color.W)
// 	}

// 	item := sequence.BeginItem()
// 	item.Primitive = graphics.RenderPrimitiveLines
// 	item.Program = s.debugProgram
// 	item.ModelMatrix = sprec.IdentityMat4()
// 	item.VertexArray = s.debugVertexArray
// 	item.IndexCount = int32(len(lines) * 2)
// 	sequence.EndItem(item)
// }
