package scene

import (
	"github.com/mokiat/lacking/async"
	"github.com/mokiat/lacking/graphics"
	"github.com/mokiat/lacking/resource"
)

func NewData(registry *resource.Registry, gfxWorker *graphics.Worker) *Data {
	return &Data{
		registry:  registry,
		gfxWorker: gfxWorker,

		GeometryFramebuffer: &graphics.Framebuffer{},
		LightingFramebuffer: &graphics.Framebuffer{},
	}
}

type Data struct {
	registry  *resource.Registry
	gfxWorker *graphics.Worker

	SkyboxProgram           *resource.Program
	SkyboxMesh              *resource.Mesh
	CarModel                *resource.Model
	Level                   *resource.Level
	DeferredLightingProgram *resource.Program
	QuadMesh                *resource.Mesh
	DebugProgram            *resource.Program
	GeometryFramebuffer     *graphics.Framebuffer
	LightingFramebuffer     *graphics.Framebuffer

	loadOutcome async.Outcome
	gfxTask     *graphics.Task
}

func (d *Data) Request() {
	d.loadOutcome = async.NewCompositeOutcome(
		d.registry.LoadProgram("geometry-skybox").OnSuccess(resource.InjectProgram(&d.SkyboxProgram)),
		d.registry.LoadMesh("skybox").OnSuccess(resource.InjectMesh(&d.SkyboxMesh)),
		d.registry.LoadModel("suv").OnSuccess(resource.InjectModel(&d.CarModel)),
		d.registry.LoadLevel("forest").OnSuccess(resource.InjectLevel(&d.Level)),
		d.registry.LoadProgram("lighting-pbr").OnSuccess(resource.InjectProgram(&d.DeferredLightingProgram)),
		d.registry.LoadMesh("quad").OnSuccess(resource.InjectMesh(&d.QuadMesh)),
		d.registry.LoadProgram("forward-debug").OnSuccess(resource.InjectProgram(&d.DebugProgram)),
	)

	d.gfxTask = d.gfxWorker.Schedule(func() error {
		geometryFramebufferData := graphics.FramebufferData{
			Width:               framebufferWidth,
			Height:              framebufferHeight,
			HasAlbedoAttachment: true,
			HasNormalAttachment: true,
			HasDepthAttachment:  true,
		}
		if err := d.GeometryFramebuffer.Allocate(geometryFramebufferData); err != nil {
			return err
		}
		lightingFramebufferData := graphics.FramebufferData{
			Width:               framebufferWidth,
			Height:              framebufferHeight,
			HasAlbedoAttachment: true,
			HasDepthAttachment:  true,
		}
		if err := d.LightingFramebuffer.Allocate(lightingFramebufferData); err != nil {
			return err
		}
		return nil
	})
}

func (d *Data) Dismiss() {
	d.registry.UnloadProgram(d.SkyboxProgram)
	d.registry.UnloadMesh(d.SkyboxMesh)
	d.registry.UnloadModel(d.CarModel)
	d.registry.UnloadLevel(d.Level)
	d.registry.UnloadProgram(d.DeferredLightingProgram)
	d.registry.UnloadMesh(d.QuadMesh)
	d.registry.UnloadProgram(d.DebugProgram)

	geometryFramebuffer := d.GeometryFramebuffer
	lightingFramebuffer := d.LightingFramebuffer
	d.gfxWorker.Schedule(func() error {
		if err := geometryFramebuffer.Release(); err != nil {
			return err
		}
		if err := lightingFramebuffer.Release(); err != nil {
			return err
		}
		return nil
	})
	d.gfxTask = nil
}

func (d *Data) IsAvailable() bool {
	return d.loadOutcome.IsAvailable() && (d.gfxTask != nil && d.gfxTask.Done())
}
