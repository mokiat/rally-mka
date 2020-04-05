package model

import (
	"fmt"
	"os"

	cli "github.com/urfave/cli/v2"

	"github.com/mokiat/rally-mka/cmd/rallygen/internal/mesh"
	"github.com/mokiat/rally-mka/internal/data/asset"
	"github.com/mokiat/rally-mka/internal/data/resource"
)

func GenerateCommand() *cli.Command {
	return &cli.Command{
		Name:      "model",
		Usage:     "generates a model",
		UsageText: "model <input-file> <output-file>",
		Action: func(c *cli.Context) error {
			if c.Args().Len() < 2 {
				return fmt.Errorf("insufficient number of arguments: expected 2 got %d", c.Args().Len())
			}
			action := &generateModelAction{
				InputFile:  c.Args().Get(0),
				OutputFile: c.Args().Get(1),
			}
			return action.Run()
		},
	}
}

type generateModelAction struct {
	InputFile  string
	OutputFile string
}

func (a *generateModelAction) Run() error {
	resModel, err := a.readResourceModel(a.InputFile)
	if err != nil {
		return fmt.Errorf("failed to read model: %w", err)
	}

	assetModel, err := ConvertResourceToAsset(resModel)
	if err != nil {
		return fmt.Errorf("failed to convert model: %w", err)
	}

	if err := a.writeAssetModel(a.OutputFile, assetModel); err != nil {
		return fmt.Errorf("failed to write model: %w", err)
	}
	return nil
}

func (a *generateModelAction) readResourceModel(path string) (*resource.Model, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %q: %w", path, err)
	}
	defer file.Close()

	model, err := resource.NewModelDecoder().Decode(file)
	if err != nil {
		return nil, fmt.Errorf("failed to decode model: %w", err)
	}
	return model, nil
}

func (a *generateModelAction) writeAssetModel(path string, model *asset.Model) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file %q: %w", path, err)
	}
	defer file.Close()

	if err := asset.NewModelEncoder().Encode(file, model); err != nil {
		return fmt.Errorf("failed to encode model: %w", err)
	}
	return nil
}

func ConvertResourceToAsset(resModel *resource.Model) (*asset.Model, error) {
	assetModel := &asset.Model{
		Meshes: make([]asset.Mesh, len(resModel.Meshes)),
		Nodes:  make([]asset.Node, len(resModel.Nodes)),
	}
	for i, resMesh := range resModel.Meshes {
		mesh, err := mesh.ConvertResourceToAsset(&resMesh)
		if err != nil {
			return nil, fmt.Errorf("failed to convert mesh: %w", err)
		}
		assetModel.Meshes[i] = *mesh
	}
	for i, resNode := range resModel.Nodes {
		assetModel.Nodes[i] = asset.Node{
			ParentIndex: int16(resNode.ParentIndex),
			Name:        resNode.Name,
			Matrix:      resNode.Matrix,
			MeshIndex:   uint16(resNode.MeshIndex),
		}
	}
	return assetModel, nil
}
