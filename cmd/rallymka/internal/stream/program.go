package stream

import (
	"fmt"

	"github.com/mokiat/rally-mka/cmd/rallymka/internal/graphics"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/resource"
	"github.com/mokiat/rally-mka/internal/data/asset"
)

const programResourceType = "program"

func GetProgram(registry *resource.Registry, name string) *Program {
	return registry.ResourceType(programResourceType).Resource(name).(*Program)
}

type Program struct {
	*resource.Handle
	gfxProgram *graphics.Program
}

func (p *Program) Gfx() *graphics.Program {
	return p.gfxProgram
}

func NewProgramController(capacity int, gfxWorker *graphics.Worker) ProgramController {
	return ProgramController{
		programs:  make([]Program, capacity),
		gfxWorker: gfxWorker,
	}
}

type ProgramController struct {
	programs  []Program
	gfxWorker *graphics.Worker
}

func (c ProgramController) ResourceTypeName() string {
	return programResourceType
}

func (c ProgramController) Init(index int, handle *resource.Handle) resource.Resource {
	c.programs[index] = Program{
		Handle:     handle,
		gfxProgram: &graphics.Program{},
	}
	return &c.programs[index]
}

func (c ProgramController) Load(index int, locator resource.Locator, registry *resource.Registry) error {
	program := &c.programs[index]

	in, err := locator.Open("assets", "programs", program.Name())
	if err != nil {
		return fmt.Errorf("failed to open program asset %q: %w", program.Name(), err)
	}
	defer in.Close()

	programAsset, err := asset.NewProgramDecoder().Decode(in)
	if err != nil {
		return fmt.Errorf("failed to decode program asset %q: %w", program.Name(), err)
	}

	gfxTask := func() error {
		return program.gfxProgram.Allocate(graphics.ProgramData{
			VertexShaderSourceCode:   programAsset.VertexSourceCode,
			FragmentShaderSourceCode: programAsset.FragmentSourceCode,
		})
	}
	if err := c.gfxWorker.Run(gfxTask); err != nil {
		return fmt.Errorf("failed to allocate gfx program: %w", err)
	}
	return nil
}

func (c ProgramController) Unload(index int) error {
	program := &c.programs[index]

	gfxTask := func() error {
		return program.gfxProgram.Release()
	}
	if err := c.gfxWorker.Run(gfxTask); err != nil {
		return fmt.Errorf("failed to release gfx program: %w", err)
	}
	return nil
}
