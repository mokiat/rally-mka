package scene

import (
	"time"

	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/lacking/game"
	"github.com/mokiat/lacking/graphics"
	"github.com/mokiat/lacking/physics"
	"github.com/mokiat/lacking/resource"
	"github.com/mokiat/lacking/shape"
	"github.com/mokiat/lacking/world"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/ecs"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/scene/car"
	"github.com/mokiat/rally-mka/internal/data"
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

const maxDebugLines = 1024 * 8

var arrayTask *graphics.Task

func NewStage(gfxWorker *graphics.Worker) *Stage {
	indexData := make([]byte, maxDebugLines*2)
	for i := 0; i < maxDebugLines; i++ {
		data.Buffer(indexData).SetUInt16(i*2, uint16(i))
	}
	vertexData := make([]byte, maxDebugLines*4*7*2)
	debugVertexArrayData := graphics.VertexArrayData{
		VertexData: vertexData,
		Layout: graphics.VertexArrayLayout{
			HasCoord:    true,
			CoordOffset: 0,
			CoordStride: 4 * 7,
			HasColor:    true,
			ColorOffset: 4 * 3,
			ColorStride: 4 * 7,
		},
		IndexData: indexData,
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
		ecsCarSystem:         ecs.NewCarSystem(ecsManager),
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
	ecsCarSystem         *ecs.CarSystem
	ecsCameraStandSystem *ecs.CameraStandSystem
	physicsEngine        *physics.Engine

	geometryFramebuffer *graphics.Framebuffer
	screenFramebuffer   *graphics.Framebuffer
	lightingProgram     *graphics.Program
	quadMesh            *resource.Mesh

	debugProgram         *graphics.Program
	debugVertexArray     *graphics.VertexArray
	debugVertexArrayData graphics.VertexArrayData
	debugLines           []DebugLine
}

var targetEntity *ecs.Entity

func (s *Stage) Init(data *Data, camera *world.Camera) {
	level := data.Level

	s.debugProgram = data.DebugProgram.GFXProgram

	s.geometryFramebuffer = data.GeometryFramebuffer

	s.lightingProgram = data.DeferredLightingProgram.GFXProgram
	s.quadMesh = data.QuadMesh

	for _, staticMesh := range level.StaticMeshes {
		entity := s.ecsManager.CreateEntity()
		entity.Render = &ecs.RenderComponent{
			GeomProgram: data.TerrainProgram.GFXProgram,
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
			CollisionShapes: []shape.Placement{collisionMesh},
		})
	}

	for _, staticEntity := range level.StaticEntities {
		entity := s.ecsManager.CreateEntity()
		entity.Render = &ecs.RenderComponent{
			GeomProgram: data.EntityProgram.GFXProgram,
			Model:       staticEntity.Model,
			Matrix:      staticEntity.Matrix,
		}
	}

	carProgram := data.CarProgram.GFXProgram
	carModel := data.CarModel

	// targetEntity =
	// 	s.setupChandelierDemo(carProgram, carModel, sprec.NewVec3(0.0, 10.0, 0.0))

	// targetEntity =
	// 	s.setupRodDemo(carProgram, carModel, sprec.NewVec3(0.0, 10.0, 5.0))

	// targetEntity =
	// 	s.setupCoiloverDemo(carProgram, carModel, sprec.NewVec3(0.0, 10.0, -5.0))

	// targetEntity =
	// 	s.setupCarDemo(carProgram, carModel, sprec.NewVec3(0.0, 141.0, 0.0))

	targetEntity =
		s.setupCarDemo(carProgram, carModel, sprec.NewVec3(0.0, 2.0, 0.0))

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
			Program: data.SkyboxProgram.GFXProgram,
			Texture: level.SkyboxTexture.GFXTexture,
			Mesh:    data.SkyboxMesh,
		}
	}
}

func (s *Stage) setupChandelierDemo(program *graphics.Program, model *resource.Model, position sprec.Vec3) *ecs.Entity {
	fakeFixtureWheel := car.Wheel(program, model, car.FrontRightWheelLocation).
		WithPosition(position).
		Build(s.ecsManager)
	s.physicsEngine.AddBody(fakeFixtureWheel.Physics.Body)
	s.physicsEngine.AddConstraint(physics.FixedTranslationConstraint{
		Fixture: position,
		Body:    fakeFixtureWheel.Physics.Body,
	})

	playWheel := car.Wheel(program, model, car.FrontRightWheelLocation).
		WithPosition(sprec.Vec3Sum(position, sprec.NewVec3(-2.3, 0.0, 0.0))).
		Build(s.ecsManager)
	s.physicsEngine.AddBody(playWheel.Physics.Body)
	s.physicsEngine.AddConstraint(physics.ChandelierConstraint{
		Fixture:    position,
		Body:       playWheel.Physics.Body,
		BodyAnchor: sprec.NewVec3(0.3, 0.0, 0.0),
		Length:     2.0,
	})

	return fakeFixtureWheel
}

func (s *Stage) setupCoiloverDemo(program *graphics.Program, model *resource.Model, position sprec.Vec3) *ecs.Entity {
	fixtureWheel := car.Wheel(program, model, car.FrontRightWheelLocation).
		WithPosition(position).
		Build(s.ecsManager)
	s.physicsEngine.AddBody(fixtureWheel.Physics.Body)
	s.physicsEngine.AddConstraint(physics.FixedTranslationConstraint{
		Fixture: position,
		Body:    fixtureWheel.Physics.Body,
	})

	fallingWheel := car.Wheel(program, model, car.FrontRightWheelLocation).
		WithPosition(sprec.Vec3Sum(position, sprec.NewVec3(0.0, 2.0, 0.0))).
		Build(s.ecsManager)
	s.physicsEngine.AddBody(fallingWheel.Physics.Body)
	s.physicsEngine.AddConstraint(&physics.CoiloverConstraint{
		FirstBody:    fixtureWheel.Physics.Body,
		SecondBody:   fallingWheel.Physics.Body,
		FrequencyHz:  4.5,
		DampingRatio: 0.1,
	})

	return fixtureWheel
}

func (s *Stage) setupRodDemo(program *graphics.Program, model *resource.Model, position sprec.Vec3) *ecs.Entity {
	topWheelPosition := position
	topWheel := car.Wheel(program, model, car.FrontRightWheelLocation).
		WithPosition(topWheelPosition).
		Build(s.ecsManager)
	s.physicsEngine.AddBody(topWheel.Physics.Body)
	s.physicsEngine.AddConstraint(physics.FixedTranslationConstraint{
		Fixture: topWheelPosition,
		Body:    topWheel.Physics.Body,
	})

	middleWheelPosition := sprec.Vec3Sum(topWheelPosition, sprec.NewVec3(1.4, 0.0, 0.0))
	middleWheel := car.Wheel(program, model, car.FrontRightWheelLocation).
		WithPosition(middleWheelPosition).
		Build(s.ecsManager)
	s.physicsEngine.AddBody(middleWheel.Physics.Body)
	s.physicsEngine.AddConstraint(physics.HingedRodConstraint{
		FirstBody:        topWheel.Physics.Body,
		FirstBodyAnchor:  sprec.NewVec3(0.2, 0.0, 0.0),
		SecondBody:       middleWheel.Physics.Body,
		SecondBodyAnchor: sprec.NewVec3(-0.2, 0.0, 0.0),
		Length:           1.0,
	})

	bottomWheelPosition := sprec.Vec3Sum(middleWheelPosition, sprec.NewVec3(1.4, 0.0, 0.0))
	bottomWheel := car.Wheel(program, model, car.FrontRightWheelLocation).
		WithPosition(bottomWheelPosition).
		Build(s.ecsManager)
	s.physicsEngine.AddBody(bottomWheel.Physics.Body)
	s.physicsEngine.AddConstraint(physics.HingedRodConstraint{
		FirstBody:        middleWheel.Physics.Body,
		FirstBodyAnchor:  sprec.NewVec3(0.2, 0.0, 0.0),
		SecondBody:       bottomWheel.Physics.Body,
		SecondBodyAnchor: sprec.NewVec3(-0.2, 0.0, 0.0),
		Length:           1.0,
	})

	return topWheel
}

func (s *Stage) setupCarDemo(program *graphics.Program, model *resource.Model, position sprec.Vec3) *ecs.Entity {
	chasis := car.Chassis(program, model).
		WithName("chasis").
		WithPosition(position).
		Build(s.ecsManager)
	s.physicsEngine.AddBody(chasis.Physics.Body)

	suspensionEnabled := true
	suspensionWidth := float32(1.0)
	suspensionLength := float32(0.3)
	suspensionFrequencyHz := float32(4.5)
	suspensionDampingRatio := float32(1.0)

	flWheelRelativePosition := sprec.NewVec3(suspensionWidth, -0.6-suspensionLength/2.0, 1.25)
	flWheel := car.Wheel(program, model, car.FrontLeftWheelLocation).
		WithName("front-left-wheel").
		WithPosition(sprec.Vec3Sum(position, flWheelRelativePosition)).
		Build(s.ecsManager)
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
		MaxY:       -0.5,
		MinY:       -1.0,
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

	frWheelRelativePosition := sprec.NewVec3(-suspensionWidth, -0.6-suspensionLength/2.0, 1.25)
	frWheel := car.Wheel(program, model, car.FrontRightWheelLocation).
		WithName("front-right-wheel").
		WithPosition(sprec.Vec3Sum(position, frWheelRelativePosition)).
		Build(s.ecsManager)
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
		MaxY:       -0.5,
		MinY:       -1.0,
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

	blWheelRelativePosition := sprec.NewVec3(suspensionWidth, -0.6-suspensionLength/2.0, -1.45)
	blWheel := car.Wheel(program, model, car.BackLeftWheelLocation).
		WithName("back-left-wheel").
		WithPosition(sprec.Vec3Sum(position, blWheelRelativePosition)).
		Build(s.ecsManager)
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
		MaxY:       -0.5,
		MinY:       -1.0,
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

	brWheelRelativePosition := sprec.NewVec3(-suspensionWidth, -0.6-suspensionLength/2.0, -1.45)
	brWheel := car.Wheel(program, model, car.BackRightWheelLocation).
		WithName("back-right-wheel").
		WithPosition(sprec.Vec3Sum(position, brWheelRelativePosition)).
		Build(s.ecsManager)
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
		MaxY:       -0.5,
		MinY:       -1.0,
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
	car.Car = &ecs.Car{
		Chassis:         chasis.Physics.Body,
		FLWheelRotation: flRotation,
		FLWheel:         flWheel.Physics.Body,
		FRWheelRotation: frRotation,
		FRWheel:         frWheel.Physics.Body,
		BLWheel:         blWheel.Physics.Body,
		BRWheel:         brWheel.Physics.Body,
	}
	car.HumanInput = true

	return chasis
}

func (s *Stage) Resize(width, height int) {
	// Initialize new framebuffer instance, to avoid race conditions
	s.screenFramebuffer = &graphics.Framebuffer{
		Width:  int32(width),
		Height: int32(height),
	}
}

func (s *Stage) Update(ctx game.UpdateContext, camera *world.Camera) {
	s.physicsEngine.Update(ctx.ElapsedTime)
	s.ecsCarSystem.Update(ctx)
	s.ecsCameraStandSystem.Update(ctx)
}

func (s *Stage) Render(pipeline *graphics.Pipeline, camera *world.Camera) {
	if !arrayTask.Done() {
		panic("NOT DONE!")
	}

	pipeline.SchedulePreRender(func() {
		// FIXME: Race condition, as vertexArrayData is being modified!
		// Instead:
		// -- Outside closure function --
		// data := pipeline.StagingData(length)
		// SetVec3(data, ....)
		// etc.
		// debugVertexArray := VertexArrayData{VertexData: data}
		// -- Inside closure function
		// s.debugVertexArray.Update(debugVertexArrayData) ...

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
	geometrySequence.ViewMatrix = camera.ViewMatrix()
	s.ecsRenderer.Render(geometrySequence)
	pipeline.EndSequence(geometrySequence)

	lightingSequence := pipeline.BeginSequence()
	lightingSequence.SourceFramebuffer = s.geometryFramebuffer
	lightingSequence.TargetFramebuffer = s.screenFramebuffer
	lightingSequence.BlitFramebufferDepth = true
	lightingSequence.ClearColor = true
	lightingSequence.TestDepth = false
	lightingSequence.WriteDepth = false
	lightingSequence.ProjectionMatrix = camera.ProjectionMatrix()
	lightingSequence.ViewMatrix = camera.ViewMatrix()
	lightingSequence.CameraMatrix = camera.Matrix()
	quadItem := lightingSequence.BeginItem()
	quadItem.Program = s.lightingProgram
	quadItem.VertexArray = s.quadMesh.GFXVertexArray
	quadItem.IndexOffset = s.quadMesh.SubMeshes[0].IndexOffset
	quadItem.IndexCount = s.quadMesh.SubMeshes[0].IndexCount
	lightingSequence.EndItem(quadItem)
	pipeline.EndSequence(lightingSequence)

	// TODO: Move skybox rendering as part of forward sequence
	// for proper lighting
	// It might make sense to experiment with stencil buffer
	// to reduce unnecessary draw fragments during lighting
	// pass

	// TODO: This would be better achieved via subsequences
	// of the lighting sequence
	forwardSequence := pipeline.BeginSequence()
	forwardSequence.SourceFramebuffer = s.geometryFramebuffer
	forwardSequence.TargetFramebuffer = s.screenFramebuffer
	forwardSequence.TestDepth = false
	forwardSequence.WriteDepth = false
	forwardSequence.DepthFunc = graphics.DepthFuncLessOrEqual
	forwardSequence.ProjectionMatrix = camera.ProjectionMatrix()
	forwardSequence.ViewMatrix = camera.ViewMatrix()
	forwardSequence.CameraMatrix = camera.Matrix()
	s.refreshDebugLines()
	// s.renderDebugLines(forwardSequence)
	pipeline.EndSequence(forwardSequence)
}

type DebugLine struct {
	A     sprec.Vec3
	B     sprec.Vec3
	Color sprec.Vec4
}

func (s *Stage) refreshDebugLines() {
	s.debugLines = s.debugLines[:0]
	for _, body := range s.physicsEngine.Bodies() {
		color := sprec.NewVec4(1.0, 1.0, 1.0, 1.0)
		if body.InCollision {
			color = sprec.NewVec4(1.0, 0.0, 0.0, 1.0)
		}
		for _, placement := range body.CollisionShapes {
			placementWS := placement.Transformed(body.Position, body.Orientation)
			s.renderDebugPlacement(placementWS, color)
		}
	}
}

func (s *Stage) renderDebugPlacement(placement shape.Placement, color sprec.Vec4) {
	switch shape := placement.Shape.(type) {
	case shape.StaticSphere:
		s.renderDebugSphere(placement, shape, color)
	case shape.StaticBox:
		s.renderDebugBox(placement, shape, color)
	case shape.StaticMesh:
		s.renderDebugMesh(placement, shape, color)
	}
}

func (s *Stage) renderDebugSphere(placement shape.Placement, sphere shape.StaticSphere, color sprec.Vec4) {
	// FIXME: Draw sphere, not box!
	box := shape.NewStaticBox(sphere.Radius()*2.0, sphere.Radius()*2.0, sphere.Radius()*2.0)
	s.renderDebugBox(placement, box, color)
}

func (s *Stage) renderDebugBox(placement shape.Placement, box shape.StaticBox, color sprec.Vec4) {
	minX := sprec.Vec3Prod(placement.Orientation.OrientationX(), -box.Width()/2.0)
	maxX := sprec.Vec3Prod(placement.Orientation.OrientationX(), box.Width()/2.0)
	minY := sprec.Vec3Prod(placement.Orientation.OrientationY(), -box.Height()/2.0)
	maxY := sprec.Vec3Prod(placement.Orientation.OrientationY(), box.Height()/2.0)
	minZ := sprec.Vec3Prod(placement.Orientation.OrientationZ(), -box.Length()/2.0)
	maxZ := sprec.Vec3Prod(placement.Orientation.OrientationZ(), box.Length()/2.0)

	p1 := sprec.Vec3Sum(sprec.Vec3Sum(sprec.Vec3Sum(placement.Position, minX), minZ), maxY)
	p2 := sprec.Vec3Sum(sprec.Vec3Sum(sprec.Vec3Sum(placement.Position, minX), maxZ), maxY)
	p3 := sprec.Vec3Sum(sprec.Vec3Sum(sprec.Vec3Sum(placement.Position, maxX), maxZ), maxY)
	p4 := sprec.Vec3Sum(sprec.Vec3Sum(sprec.Vec3Sum(placement.Position, maxX), minZ), maxY)
	p5 := sprec.Vec3Sum(sprec.Vec3Sum(sprec.Vec3Sum(placement.Position, minX), minZ), minY)
	p6 := sprec.Vec3Sum(sprec.Vec3Sum(sprec.Vec3Sum(placement.Position, minX), maxZ), minY)
	p7 := sprec.Vec3Sum(sprec.Vec3Sum(sprec.Vec3Sum(placement.Position, maxX), maxZ), minY)
	p8 := sprec.Vec3Sum(sprec.Vec3Sum(sprec.Vec3Sum(placement.Position, maxX), minZ), minY)

	s.addDebugLine(p1, p2, color)
	s.addDebugLine(p2, p3, color)
	s.addDebugLine(p3, p4, color)
	s.addDebugLine(p4, p1, color)

	s.addDebugLine(p5, p6, color)
	s.addDebugLine(p6, p7, color)
	s.addDebugLine(p7, p8, color)
	s.addDebugLine(p8, p5, color)

	s.addDebugLine(p1, p5, color)
	s.addDebugLine(p2, p6, color)
	s.addDebugLine(p3, p7, color)
	s.addDebugLine(p4, p8, color)
}

func (s *Stage) renderDebugMesh(placement shape.Placement, mesh shape.StaticMesh, color sprec.Vec4) {
	for _, triangle := range mesh.Triangles() {
		triangleWS := triangle.Transformed(placement.Position, placement.Orientation)
		s.addDebugLine(triangleWS.A(), triangleWS.B(), color)
		s.addDebugLine(triangleWS.B(), triangleWS.C(), color)
		s.addDebugLine(triangleWS.C(), triangleWS.A(), color)
	}
}

func (s *Stage) addDebugLine(a, b sprec.Vec3, color sprec.Vec4) {
	s.debugLines = append(s.debugLines, DebugLine{
		A:     a,
		B:     b,
		Color: color,
	})
}

func (s *Stage) renderDebugLines(sequence *graphics.Sequence) {
	for i, line := range s.debugLines {
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
	item.IndexCount = int32(len(s.debugLines) * 2)
	sequence.EndItem(item)
}
