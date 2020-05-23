package stream

import (
	"fmt"

	"github.com/mokiat/lacking/graphics"
	"github.com/mokiat/rally-mka/internal/data/asset"
	"github.com/mokiat/rally-mka/internal/engine/resource"
)

const programResourceType = "program"

func GetProgram(registry *resource.Registry, name string) ProgramHandle {
	return ProgramHandle{
		Handle: registry.Type(programResourceType).Resource(name),
	}
}

type ProgramHandle struct {
	resource.Handle
}

func (h ProgramHandle) Get() *graphics.Program {
	return h.Handle.Get().(*graphics.Program)
}

func NewProgramOperator(locator resource.Locator, gfxWorker *graphics.Worker) *ProgramOperator {
	return &ProgramOperator{
		locator:   locator,
		gfxWorker: gfxWorker,
	}
}

type ProgramOperator struct {
	locator   resource.Locator
	gfxWorker *graphics.Worker
}

func (o *ProgramOperator) Register(registry *resource.Registry) {
	registry.RegisterType(programResourceType, o)
}

func (o *ProgramOperator) Allocate(registry *resource.Registry, name string) (resource.Resource, error) {
	in, err := o.locator.Open("assets", "programs", name)
	if err != nil {
		return nil, fmt.Errorf("failed to open program asset %q: %w", name, err)
	}
	defer in.Close()

	programAsset, err := asset.NewProgramDecoder().Decode(in)
	if err != nil {
		return nil, fmt.Errorf("failed to decode program asset %q: %w", name, err)
	}

	program := &graphics.Program{}

	gfxTask := o.gfxWorker.Schedule(func() error {
		return program.Allocate(graphics.ProgramData{
			VertexShaderSourceCode:   programAsset.VertexSourceCode,
			FragmentShaderSourceCode: programAsset.FragmentSourceCode,
		})
	})
	if err := gfxTask.Wait(); err != nil {
		return nil, fmt.Errorf("failed to allocate gfx program: %w", err)
	}
	return program, nil
}

func (o *ProgramOperator) Release(registry *resource.Registry, resource resource.Resource) error {
	program := resource.(*graphics.Program)

	gfxTask := o.gfxWorker.Schedule(func() error {
		return program.Release()
	})
	if err := gfxTask.Wait(); err != nil {
		return fmt.Errorf("failed to release gfx program: %w", err)
	}
	return nil
}
