package scene

import (
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/stream"
	"github.com/mokiat/rally-mka/internal/engine/graphics"
	"github.com/mokiat/rally-mka/internal/engine/resource"
)

func NewData(registry *resource.Registry) *Data {
	return &Data{
		registry:      registry,
		SkyboxProgram: stream.GetProgram(registry, "skybox"),
		SkyboxTexture: stream.GetCubeTexture(registry, "city"),
		SkyboxMesh:    stream.GetMesh(registry, "skybox"),
	}
}

type Data struct {
	registry      *resource.Registry
	SkyboxProgram stream.ProgramHandle
	SkyboxTexture stream.CubeTextureHandle
	SkyboxMesh    stream.MeshHandle
}

func (d *Data) Request() {
	d.registry.Request(d.SkyboxProgram.Handle)
	d.registry.Request(d.SkyboxTexture.Handle)
	d.registry.Request(d.SkyboxMesh.Handle)
}

func (d *Data) Dismiss() {
	d.registry.Dismiss(d.SkyboxProgram.Handle)
	d.registry.Dismiss(d.SkyboxTexture.Handle)
	d.registry.Dismiss(d.SkyboxMesh.Handle)
}

func (d *Data) IsAvailable() bool {
	return d.SkyboxProgram.IsAvailable() &&
		d.SkyboxTexture.IsAvailable() &&
		d.SkyboxMesh.IsAvailable()
}

type Stage struct {
	skybox *Skybox
}

func NewStage() *Stage {
	return &Stage{}
}

func (s *Stage) Init(data *Data) {
	s.skybox = &Skybox{
		Program: data.SkyboxProgram,
		Texture: data.SkyboxTexture,
		Mesh:    data.SkyboxMesh,
	}
}

func (s *Stage) Render(pipeline *graphics.Pipeline, camera *Camera) {
	if skybox, found := s.getSkybox(camera); found {
		s.renderSky(pipeline, camera, skybox)
	}
}

func (s *Stage) getSkybox(camera *Camera) (*Skybox, bool) {
	// TODO: this implementation should iterate the BSP tree / Octree
	// to find which skybox is relevant for the current camera position
	return s.skybox, s.skybox != nil
}

func (s *Stage) renderSky(pipeline *graphics.Pipeline, camera *Camera, skybox *Skybox) {
	skyboxSequence := pipeline.BeginSequence()
	skyboxSequence.WriteDepth = false
	skyboxSequence.DepthFunc = graphics.DepthFuncLessOrEqual
	skyboxSequence.ProjectionMatrix = camera.ProjectionMatrix()
	skyboxSequence.ViewMatrix = camera.InverseViewMatrix()
	skyboxItem := skyboxSequence.BeginItem()
	skyboxItem.Program = skybox.Program.Get()
	skyboxItem.SkyboxTexture = skybox.Texture.Get()
	skyboxItem.VertexArray = skybox.Mesh.Get().VertexArray
	skyboxItem.IndexCount = skybox.Mesh.Get().SubMeshes[0].IndexCount
	skyboxSequence.EndItem(skyboxItem)
	pipeline.EndSequence(skyboxSequence)
}
