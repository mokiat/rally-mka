package program

import (
	"fmt"
	"os"

	cli "github.com/urfave/cli/v2"

	"github.com/mokiat/rally-mka/internal/data/asset"
	"github.com/mokiat/rally-mka/internal/data/resource"
)

func GenerateCommand() *cli.Command {
	return &cli.Command{
		Name:      "program",
		Usage:     "generates a shader program",
		UsageText: "program <vertex-shader> <fragment-shader> <output-file>",
		Action: func(c *cli.Context) error {
			if c.Args().Len() < 3 {
				return fmt.Errorf("insufficient number of arguments: expected 3 got %d", c.Args().Len())
			}
			action := &generateProgramAction{
				VertexShaderFile:   c.Args().Get(0),
				FragmentShaderFile: c.Args().Get(1),
				OutputFile:         c.Args().Get(2),
			}
			return action.Run()
		},
	}
}

type generateProgramAction struct {
	VertexShaderFile   string
	FragmentShaderFile string
	OutputFile         string
}

func (a *generateProgramAction) Run() error {
	vertexShader, err := a.readShader(a.VertexShaderFile)
	if err != nil {
		return fmt.Errorf("failed to read vertex shader: %w", err)
	}
	fragmentShader, err := a.readShader(a.FragmentShaderFile)
	if err != nil {
		return fmt.Errorf("failed to read fragment shader: %w", err)
	}

	program := &asset.Program{
		VertexSourceCode:   string(vertexShader.SourceCode),
		FragmentSourceCode: string(fragmentShader.SourceCode),
	}
	if err := a.writeProgram(a.OutputFile, program); err != nil {
		return fmt.Errorf("failed to store program: %w", err)
	}
	return nil
}

func (a *generateProgramAction) readShader(path string) (*resource.Shader, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %q: %w", path, err)
	}
	defer file.Close()

	shader, err := resource.NewShaderDecoder().Decode(file)
	if err != nil {
		return nil, fmt.Errorf("failed to decode shader: %w", err)
	}
	return shader, nil
}

func (a *generateProgramAction) writeProgram(path string, program *asset.Program) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file %q: %w", path, err)
	}
	defer file.Close()

	if err := asset.NewProgramEncoder().Encode(file, program); err != nil {
		return fmt.Errorf("failed to encode program: %w", err)
	}
	return nil
}
