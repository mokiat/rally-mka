package level

import (
	"fmt"
	"os"

	cli "github.com/urfave/cli/v2"

	"github.com/mokiat/lacking/data/asset"
	"github.com/mokiat/rally-mka/cmd/rallygen/internal/mesh"
	"github.com/mokiat/rally-mka/internal/data/resource"
)

func GenerateCommand() *cli.Command {
	return &cli.Command{
		Name:      "level",
		Usage:     "generates a level",
		UsageText: "level <input-file> <output-file>",
		Action: func(c *cli.Context) error {
			if c.Args().Len() < 2 {
				return fmt.Errorf("insufficient number of arguments: expected 2 got %d", c.Args().Len())
			}
			action := &generateLevelAction{
				InputFile:  c.Args().Get(0),
				OutputFile: c.Args().Get(1),
			}
			return action.Run()
		},
	}
}

type generateLevelAction struct {
	InputFile  string
	OutputFile string
}

func (a *generateLevelAction) Run() error {
	resLevel, err := a.readResourceLevel(a.InputFile)
	if err != nil {
		return fmt.Errorf("failed to read level: %w", err)
	}

	assetLevel, err := ConvertResourceToAsset(resLevel)
	if err != nil {
		return fmt.Errorf("failed to convert level: %w", err)
	}

	if err := a.writeAssetLevel(a.OutputFile, assetLevel); err != nil {
		return fmt.Errorf("failed to write level: %w", err)
	}
	return nil
}

func (a *generateLevelAction) readResourceLevel(path string) (*resource.Level, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	level, err := resource.NewLevelDecoder().Decode(file)
	if err != nil {
		return nil, fmt.Errorf("failed to decode level: %w", err)
	}
	return level, nil
}

func (a *generateLevelAction) writeAssetLevel(path string, level *asset.Level) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	if err := asset.EncodeLevel(file, level); err != nil {
		return fmt.Errorf("failed to encode level: %w", err)
	}
	return nil
}

func ConvertResourceToAsset(resLevel *resource.Level) (*asset.Level, error) {
	assetLevel := &asset.Level{
		SkyboxTexture: resLevel.SkyboxTexture,
	}

	assetLevel.StaticEntities = make([]asset.LevelEntity, len(resLevel.StaticEntities))
	for i, resEntity := range resLevel.StaticEntities {
		assetLevel.StaticEntities[i] = asset.LevelEntity{
			Model:  resEntity.Model,
			Matrix: resEntity.Matrix,
		}
	}

	assetLevel.StaticMeshes = make([]asset.Mesh, len(resLevel.StaticMeshes))
	for i, resMesh := range resLevel.StaticMeshes {
		mesh, err := mesh.ConvertResourceToAsset(&resMesh)
		if err != nil {
			return nil, fmt.Errorf("failed to convert mesh: %w", err)
		}
		assetLevel.StaticMeshes[i] = *mesh
	}

	assetLevel.CollisionMeshes = make([]asset.LevelCollisionMesh, len(resLevel.CollisionMeshes))
	for i, resCollisionMesh := range resLevel.CollisionMeshes {
		triangles := make([]asset.Triangle, len(resCollisionMesh.Triangles))
		for j, resTriangle := range resCollisionMesh.Triangles {
			triangles[j] = asset.Triangle{
				asset.Point(resTriangle[0]),
				asset.Point(resTriangle[1]),
				asset.Point(resTriangle[2]),
			}
		}
		assetLevel.CollisionMeshes[i] = asset.LevelCollisionMesh{
			Triangles: triangles,
		}
	}

	return assetLevel, nil
}
