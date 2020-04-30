package scene

import (
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/stream"
	"github.com/mokiat/rally-mka/internal/engine/graphics"
	"github.com/mokiat/rally-mka/internal/engine/resource"
)

func NewData(registry *resource.Registry, gfxWorker *graphics.Worker) *Data {
	return &Data{
		registry:  registry,
		gfxWorker: gfxWorker,

		SkyboxProgram:  stream.GetProgram(registry, "skybox"),
		SkyboxMesh:     stream.GetMesh(registry, "skybox"),
		TerrainProgram: stream.GetProgram(registry, "deferred-geometry"),
		EntityProgram:  stream.GetProgram(registry, "deferred-geometry"),
		CarProgram:     stream.GetProgram(registry, "deferred-geometry"),
		CarModel:       stream.GetModel(registry, "suv"),
		Level:          stream.GetLevel(registry, "forest"),

		DeferredGeometryProgram: stream.GetProgram(registry, "deferred-geometry"),
		DeferredLightingProgram: stream.GetProgram(registry, "deferred-lighting"),
		QuadMesh:                stream.GetMesh(registry, "quad"),
		DebugProgram:            stream.GetProgram(registry, "debug"),
		GeometryFramebuffer:     &graphics.Framebuffer{},
		LightingFramebuffer:     &graphics.Framebuffer{},
	}
}

type Data struct {
	registry  *resource.Registry
	gfxWorker *graphics.Worker

	SkyboxProgram  stream.ProgramHandle
	SkyboxMesh     stream.MeshHandle
	TerrainProgram stream.ProgramHandle
	EntityProgram  stream.ProgramHandle
	CarProgram     stream.ProgramHandle
	CarModel       stream.ModelHandle
	Level          stream.LevelHandle

	DeferredGeometryProgram stream.ProgramHandle
	DeferredLightingProgram stream.ProgramHandle
	QuadMesh                stream.MeshHandle

	DebugProgram stream.ProgramHandle

	gfxTask             *graphics.Task
	GeometryFramebuffer *graphics.Framebuffer
	LightingFramebuffer *graphics.Framebuffer
}

func (d *Data) Request() {
	d.registry.Request(d.SkyboxProgram.Handle)
	d.registry.Request(d.SkyboxMesh.Handle)
	d.registry.Request(d.TerrainProgram.Handle)
	d.registry.Request(d.EntityProgram.Handle)
	d.registry.Request(d.CarProgram.Handle)
	d.registry.Request(d.CarModel.Handle)
	d.registry.Request(d.Level.Handle)
	d.registry.Request(d.DeferredGeometryProgram.Handle)
	d.registry.Request(d.DeferredLightingProgram.Handle)
	d.registry.Request(d.QuadMesh.Handle)
	d.registry.Request(d.DebugProgram.Handle)
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
	d.registry.Dismiss(d.SkyboxProgram.Handle)
	d.registry.Dismiss(d.SkyboxMesh.Handle)
	d.registry.Dismiss(d.TerrainProgram.Handle)
	d.registry.Dismiss(d.EntityProgram.Handle)
	d.registry.Dismiss(d.CarProgram.Handle)
	d.registry.Dismiss(d.CarModel.Handle)
	d.registry.Dismiss(d.Level.Handle)
	d.registry.Dismiss(d.DeferredGeometryProgram.Handle)
	d.registry.Dismiss(d.DeferredLightingProgram.Handle)
	d.registry.Dismiss(d.QuadMesh.Handle)
	d.registry.Dismiss(d.DebugProgram.Handle)

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
	return d.SkyboxProgram.IsAvailable() &&
		d.SkyboxMesh.IsAvailable() &&
		d.TerrainProgram.IsAvailable() &&
		d.EntityProgram.IsAvailable() &&
		d.CarProgram.IsAvailable() &&
		d.CarModel.IsAvailable() &&
		d.Level.IsAvailable() &&
		d.DeferredGeometryProgram.IsAvailable() &&
		d.DeferredLightingProgram.IsAvailable() &&
		d.QuadMesh.IsAvailable() &&
		d.DebugProgram.IsAvailable() &&
		(d.gfxTask != nil && d.gfxTask.Done())
}
