package scene

import (
	"github.com/mokiat/go-whiskey/math"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/stream"
	"github.com/mokiat/rally-mka/internal/engine/collision"
	"github.com/mokiat/rally-mka/internal/engine/graphics"
	"github.com/mokiat/rally-mka/internal/engine/resource"
)

const (
	carDropHeight  = 1.6
	anchorDistance = 4.0
	cameraDistance = 8.0
)

type CarInput struct {
	Forward   bool
	Backward  bool
	TurnLeft  bool
	TurnRight bool
	Handbrake bool
}

func NewData(registry *resource.Registry) *Data {
	return &Data{
		registry:       registry,
		SkyboxProgram:  stream.GetProgram(registry, "skybox"),
		SkyboxMesh:     stream.GetMesh(registry, "skybox"),
		TerrainProgram: stream.GetProgram(registry, "diffuse"),
		EntityProgram:  stream.GetProgram(registry, "diffuse"),
		CarProgram:     stream.GetProgram(registry, "diffuse"),
		CarModel:       stream.GetModel(registry, "suv"),
		Level:          stream.GetLevel(registry, "forest"),
	}
}

type Data struct {
	registry *resource.Registry

	SkyboxProgram  stream.ProgramHandle
	SkyboxMesh     stream.MeshHandle
	TerrainProgram stream.ProgramHandle
	EntityProgram  stream.ProgramHandle
	CarProgram     stream.ProgramHandle
	CarModel       stream.ModelHandle
	Level          stream.LevelHandle
}

func (d *Data) Request() {
	d.registry.Request(d.SkyboxProgram.Handle)
	d.registry.Request(d.SkyboxMesh.Handle)
	d.registry.Request(d.TerrainProgram.Handle)
	d.registry.Request(d.EntityProgram.Handle)
	d.registry.Request(d.CarProgram.Handle)
	d.registry.Request(d.CarModel.Handle)
	d.registry.Request(d.Level.Handle)
}

func (d *Data) Dismiss() {
	d.registry.Dismiss(d.SkyboxProgram.Handle)
	d.registry.Dismiss(d.SkyboxMesh.Handle)
	d.registry.Dismiss(d.TerrainProgram.Handle)
	d.registry.Dismiss(d.EntityProgram.Handle)
	d.registry.Dismiss(d.CarProgram.Handle)
	d.registry.Dismiss(d.CarModel.Handle)
	d.registry.Dismiss(d.Level.Handle)
}

func (d *Data) IsAvailable() bool {
	return d.SkyboxProgram.IsAvailable() &&
		d.SkyboxMesh.IsAvailable() &&
		d.TerrainProgram.IsAvailable() &&
		d.EntityProgram.IsAvailable() &&
		d.CarProgram.IsAvailable() &&
		d.CarModel.IsAvailable() &&
		d.Level.IsAvailable()
}

type Stage struct {
	skybox *Skybox

	terrainProgram *graphics.Program
	terrains       []Terrain

	collisionMeshes []*collision.Mesh

	entityProgram *graphics.Program
	entities      []Entity

	carProgram   *graphics.Program
	carModel     *stream.Model
	car          *Car
	cameraAnchor math.Vec3
}

func NewStage() *Stage {
	return &Stage{}
}

func (s *Stage) Init(data *Data) {
	level := data.Level.Get()

	s.skybox = &Skybox{
		Program: data.SkyboxProgram.Get(),
		Texture: level.SkyboxTexture.Get(),
		Mesh:    data.SkyboxMesh.Get(),
	}

	s.terrainProgram = data.TerrainProgram.Get()
	s.terrains = make([]Terrain, len(level.StaticMeshes))
	for i, staticMesh := range level.StaticMeshes {
		s.terrains[i] = Terrain{
			Mesh: staticMesh,
		}
	}

	s.collisionMeshes = level.CollisionMeshes

	s.entityProgram = data.EntityProgram.Get()
	s.entities = make([]Entity, len(level.StaticEntities))
	for i, staticEntity := range level.StaticEntities {
		s.entities[i] = Entity{
			Model:  staticEntity.Model.Get(),
			Matrix: staticEntity.Matrix,
		}
	}

	s.carProgram = data.CarProgram.Get()
	s.carModel = data.CarModel.Get()
	s.car = NewCar(s, s.carModel, math.TranslationMat4x4(0.0, carDropHeight, 0.0))
}

func (s *Stage) Update(elapsedSeconds float32, camera *Camera, input CarInput) {
	s.updateCar(elapsedSeconds, input)

	carPosition := s.car.Position()
	// we use a camera anchor to achieve the smooth effect of a
	// camera following the car
	anchorVector := s.cameraAnchor.DecVec3(carPosition)
	anchorVector = anchorVector.Resize(anchorDistance)
	s.cameraAnchor = carPosition.IncVec3(anchorVector)

	// the following approach of creating the view matrix coordinates will fail
	// if the camera is pointing directly up or down
	cameraVectorZ := anchorVector
	cameraVectorX := math.Vec3CrossProduct(math.BaseVec3Y(), cameraVectorZ)
	cameraVectorY := math.Vec3CrossProduct(cameraVectorZ, cameraVectorX)
	camera.SetViewMatrix(math.Mat4x4MulMany(
		math.TranslationMat4x4(
			carPosition.X,
			carPosition.Y,
			carPosition.Z,
		),
		math.VectorMat4x4(
			cameraVectorX.Resize(1.0),
			cameraVectorY.Resize(1.0),
			cameraVectorZ.Resize(1.0),
			math.NullVec3(),
		),
		math.RotationMat4x4(-25.0, 1.0, 0.0, 0.0),
		math.TranslationMat4x4(0.0, 0.0, cameraDistance),
	))

	// angle = angle + elapsedSeconds*90

	// camera.SetViewMatrix(math.Mat4x4MulMany(
	// 	math.RotationMat4x4(angle, 0.0, 1.0, 0.0),
	// 	math.RotationMat4x4(-20, 1.0, 0.0, 0.0),
	// 	math.TranslationMat4x4(0.0, 0.0, 15.0),
	// ))
}

func (s *Stage) Render(pipeline *graphics.Pipeline, camera *Camera) {
	s.renderScene(pipeline, camera)
	if skybox, found := s.getSkybox(camera); found {
		s.renderSky(pipeline, camera, skybox)
	}
}

func (s *Stage) CheckCollision(line collision.Line) (bestCollision collision.LineCollision, found bool) {
	closestDistance := line.LengthSquared()
	for _, mesh := range s.collisionMeshes {
		if lineCollision, ok := mesh.LineCollision(line); ok {
			found = true
			distanceVector := lineCollision.Intersection().DecVec3(line.Start())
			distance := distanceVector.LengthSquared()
			if distance < closestDistance {
				closestDistance = distance
				bestCollision = lineCollision
			}
		}
	}
	return
}

func (s *Stage) updateCar(elapsedSeconds float32, input CarInput) {
	// TODO: Move constants as part of car descriptor
	const turnSpeed = 100        // FIXME ORIGINAL: 120
	const returnSpeed = 50       // FIXME ORIGINAL: 60
	const maxWheelAngle = 20     // FIXME ORIGINAL: 30
	const maxAcceleration = 0.1  // FIXME ORIGINAL: 0.01
	const maxDeceleration = 0.05 // FIXME ORIGINAL: 0.005

	if input.TurnLeft {
		if s.car.WheelAngle += elapsedSeconds * turnSpeed; s.car.WheelAngle > maxWheelAngle {
			s.car.WheelAngle = maxWheelAngle
		}
	}
	if input.TurnRight {
		if s.car.WheelAngle -= elapsedSeconds * turnSpeed; s.car.WheelAngle < -maxWheelAngle {
			s.car.WheelAngle = -maxWheelAngle
		}
	}
	if input.TurnLeft == input.TurnRight {
		if s.car.WheelAngle > 0.001 {
			if s.car.WheelAngle -= elapsedSeconds * returnSpeed; s.car.WheelAngle < 0.0 {
				s.car.WheelAngle = 0.0
			}
		}
		if s.car.WheelAngle < -0.01 {
			if s.car.WheelAngle += elapsedSeconds * returnSpeed; s.car.WheelAngle > 0.0 {
				s.car.WheelAngle = 0.0
			}
		}
	}

	s.car.Acceleration = 0.0
	if input.Forward {
		s.car.Acceleration = maxAcceleration
	}
	if input.Backward {
		s.car.Acceleration = -maxDeceleration
	}

	s.car.HandbrakePulled = input.Handbrake

	s.car.Update()
}

func (s *Stage) getSkybox(camera *Camera) (*Skybox, bool) {
	// TODO: this implementation should iterate the BSP tree / Octree
	// to find which skybox is relevant for the current camera position
	return s.skybox, s.skybox != nil
}

func (s *Stage) renderScene(pipeline *graphics.Pipeline, camera *Camera) {
	// XXX: Modern GPUs prefer that you clear all the buffers
	// and it can be faster due to cache state
	sequence := pipeline.BeginSequence()
	sequence.BackgroundColor = math.MakeVec4(0.0, 0.6, 1.0, 1.0)
	sequence.ClearColor = true
	sequence.ClearDepth = true
	sequence.ProjectionMatrix = camera.ProjectionMatrix()
	sequence.ViewMatrix = camera.InverseViewMatrix()

	s.renderModel(sequence, s.carProgram, s.car.ModelMatrix, s.carModel)
	for _, entity := range s.entities {
		s.renderModel(sequence, s.entityProgram, entity.Matrix, entity.Model)
	}
	for _, terrain := range s.terrains {
		s.renderMesh(sequence, s.terrainProgram, math.IdentityMat4x4(), terrain.Mesh)
	}

	pipeline.EndSequence(sequence)
}

func (s *Stage) renderModel(sequence *graphics.Sequence, program *graphics.Program, parentModelMatrix math.Mat4x4, model *stream.Model) {
	for _, node := range model.Nodes {
		s.renderModelNode(sequence, program, parentModelMatrix, node)
	}
}

func (s *Stage) renderModelNode(sequence *graphics.Sequence, program *graphics.Program, parentModelMatrix math.Mat4x4, node *stream.Node) {
	modelMatrix := parentModelMatrix.MulMat4x4(node.Matrix)
	s.renderMesh(sequence, program, modelMatrix, node.Mesh)
	for _, child := range node.Children {
		s.renderModelNode(sequence, program, modelMatrix, child)
	}
}

func (s *Stage) renderMesh(sequence *graphics.Sequence, program *graphics.Program, modelMatrix math.Mat4x4, mesh *stream.Mesh) {
	for _, subMesh := range mesh.SubMeshes {
		meshItem := sequence.BeginItem()
		meshItem.Program = program
		meshItem.ModelMatrix = modelMatrix
		if subMesh.DiffuseTexture != nil {
			meshItem.DiffuseTexture = subMesh.DiffuseTexture.Get()
		}
		meshItem.VertexArray = mesh.VertexArray
		meshItem.IndexCount = subMesh.IndexCount
		sequence.EndItem(meshItem)
	}
}

func (s *Stage) renderSky(pipeline *graphics.Pipeline, camera *Camera, skybox *Skybox) {
	skyboxSequence := pipeline.BeginSequence()
	skyboxSequence.WriteDepth = false
	skyboxSequence.DepthFunc = graphics.DepthFuncLessOrEqual
	skyboxSequence.ProjectionMatrix = camera.ProjectionMatrix()
	skyboxSequence.ViewMatrix = camera.InverseViewMatrix()
	skyboxItem := skyboxSequence.BeginItem()
	skyboxItem.Program = skybox.Program
	skyboxItem.SkyboxTexture = skybox.Texture
	skyboxItem.VertexArray = skybox.Mesh.VertexArray
	skyboxItem.IndexCount = skybox.Mesh.SubMeshes[0].IndexCount
	skyboxSequence.EndItem(skyboxItem)
	pipeline.EndSequence(skyboxSequence)
}
