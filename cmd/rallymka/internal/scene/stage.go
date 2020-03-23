package scene

import (
	"github.com/mokiat/go-whiskey/math"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/stream"
	"github.com/mokiat/rally-mka/internal/engine/graphics"
	"github.com/mokiat/rally-mka/internal/engine/resource"
)

func NewData(registry *resource.Registry) *Data {
	return &Data{
		registry:       registry,
		SkyboxProgram:  stream.GetProgram(registry, "skybox"),
		SkyboxMesh:     stream.GetMesh(registry, "skybox"),
		TerrainProgram: stream.GetProgram(registry, "diffuse"),
		EntityProgram:  stream.GetProgram(registry, "diffuse"),
		Level:          stream.GetLevel(registry, "forest"),
	}
}

type Data struct {
	registry *resource.Registry

	SkyboxProgram  stream.ProgramHandle
	SkyboxMesh     stream.MeshHandle
	TerrainProgram stream.ProgramHandle
	EntityProgram  stream.ProgramHandle
	Level          stream.LevelHandle
}

func (d *Data) Request() {
	d.registry.Request(d.SkyboxProgram.Handle)
	d.registry.Request(d.SkyboxMesh.Handle)
	d.registry.Request(d.TerrainProgram.Handle)
	d.registry.Request(d.EntityProgram.Handle)
	d.registry.Request(d.Level.Handle)
}

func (d *Data) Dismiss() {
	d.registry.Dismiss(d.SkyboxProgram.Handle)
	d.registry.Dismiss(d.SkyboxMesh.Handle)
	d.registry.Dismiss(d.TerrainProgram.Handle)
	d.registry.Dismiss(d.EntityProgram.Handle)
	d.registry.Dismiss(d.Level.Handle)
}

func (d *Data) IsAvailable() bool {
	return d.SkyboxProgram.IsAvailable() &&
		d.SkyboxMesh.IsAvailable() &&
		d.TerrainProgram.IsAvailable() &&
		d.EntityProgram.IsAvailable() &&
		d.Level.IsAvailable()
}

type Stage struct {
	skybox *Skybox

	terrainProgram *graphics.Program
	terrains       []Terrain

	entityProgram *graphics.Program
	entities      []Entity
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

	s.entityProgram = data.EntityProgram.Get()
	s.entities = make([]Entity, len(level.StaticEntities))
	for i, staticEntity := range level.StaticEntities {
		s.entities[i] = Entity{
			Model:  staticEntity.Model.Get(),
			Matrix: staticEntity.Matrix,
		}
	}
}

func (s *Stage) Render(pipeline *graphics.Pipeline, camera *Camera) {
	s.renderScene(pipeline, camera)
	if skybox, found := s.getSkybox(camera); found {
		s.renderSky(pipeline, camera, skybox)
	}
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
	// sequence.ClearColor = true // TODO once old renderer is removed
	// sequence.ClearDepth = true
	sequence.ProjectionMatrix = camera.ProjectionMatrix()
	sequence.ViewMatrix = camera.InverseViewMatrix()
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
