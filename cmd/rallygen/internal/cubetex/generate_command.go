package cubetex

import (
	"fmt"
	"os"

	cli "github.com/urfave/cli/v2"

	"github.com/mokiat/lacking/data/asset"
	"github.com/mokiat/rally-mka/internal/data/resource"
)

func GenerateCommand() *cli.Command {
	return &cli.Command{
		Name:      "cubetex",
		Usage:     "generates a cube texture",
		UsageText: "cubetex <front-image> <back-image> <left-image> <right-image> <top-image> <bottom-image> <output-file>",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:  "dimension",
				Usage: "specifies a custom dimension of the texture",
			},
		},
		Action: func(c *cli.Context) error {
			if c.Args().Len() < 7 {
				return fmt.Errorf("insufficient number of arguments: expected 7 got %d", c.Args().Len())
			}
			action := &generateCubeTextureAction{
				FrontFile:  c.Args().Get(0),
				BackFile:   c.Args().Get(1),
				LeftFile:   c.Args().Get(2),
				RightFile:  c.Args().Get(3),
				TopFile:    c.Args().Get(4),
				BottomFile: c.Args().Get(5),
				OutputFile: c.Args().Get(6),
				Dimension:  c.Int("dimension"),
			}
			return action.Run()
		},
	}
}

type generateCubeTextureAction struct {
	FrontFile  string
	BackFile   string
	LeftFile   string
	RightFile  string
	TopFile    string
	BottomFile string
	OutputFile string
	Dimension  int
}

func (a *generateCubeTextureAction) Run() error {
	frontImg, err := a.readResourceImage(a.FrontFile)
	if err != nil {
		return fmt.Errorf("failed to read front image: %w", err)
	}
	if !frontImg.IsSquare() {
		return fmt.Errorf("front image is not a square")
	}

	backImg, err := a.readResourceImage(a.BackFile)
	if err != nil {
		return fmt.Errorf("failed to read back image: %w", err)
	}
	if !backImg.IsSquare() {
		return fmt.Errorf("back image is not a square")
	}

	leftImg, err := a.readResourceImage(a.LeftFile)
	if err != nil {
		return fmt.Errorf("failed to read left image: %w", err)
	}
	if !leftImg.IsSquare() {
		return fmt.Errorf("left image is not a square")
	}

	rightImg, err := a.readResourceImage(a.RightFile)
	if err != nil {
		return fmt.Errorf("failed to read right image: %w", err)
	}
	if !rightImg.IsSquare() {
		return fmt.Errorf("right image is not a square")
	}

	topImg, err := a.readResourceImage(a.TopFile)
	if err != nil {
		return fmt.Errorf("failed to read top image: %w", err)
	}
	if !topImg.IsSquare() {
		return fmt.Errorf("top image is not a square")
	}

	bottomImg, err := a.readResourceImage(a.BottomFile)
	if err != nil {
		return fmt.Errorf("failed to read bottom image: %w", err)
	}
	if !bottomImg.IsSquare() {
		return fmt.Errorf("bottom image is not a square")
	}

	areSameDimension := frontImg.Width() == backImg.Width() &&
		frontImg.Width() == leftImg.Width() &&
		frontImg.Width() == rightImg.Width() &&
		frontImg.Width() == topImg.Width() &&
		frontImg.Width() == bottomImg.Width()
	if !areSameDimension {
		return fmt.Errorf("images are not of the same size")
	}

	if a.Dimension > 0 {
		frontImg.Scale(a.Dimension, a.Dimension)
		backImg.Scale(a.Dimension, a.Dimension)
		leftImg.Scale(a.Dimension, a.Dimension)
		rightImg.Scale(a.Dimension, a.Dimension)
		topImg.Scale(a.Dimension, a.Dimension)
		bottomImg.Scale(a.Dimension, a.Dimension)
	}

	assetTexture := &asset.CubeTexture{
		Dimension: uint16(frontImg.Width()),
	}
	assetTexture.Sides[asset.TextureSideFront] = asset.CubeTextureSide{
		Data: frontImg.RGBData(),
	}
	assetTexture.Sides[asset.TextureSideBack] = asset.CubeTextureSide{
		Data: backImg.RGBData(),
	}
	assetTexture.Sides[asset.TextureSideLeft] = asset.CubeTextureSide{
		Data: leftImg.RGBData(),
	}
	assetTexture.Sides[asset.TextureSideRight] = asset.CubeTextureSide{
		Data: rightImg.RGBData(),
	}
	assetTexture.Sides[asset.TextureSideTop] = asset.CubeTextureSide{
		Data: topImg.RGBData(),
	}
	assetTexture.Sides[asset.TextureSideBottom] = asset.CubeTextureSide{
		Data: bottomImg.RGBData(),
	}

	if err := a.writeAssetCubeTexture(a.OutputFile, assetTexture); err != nil {
		return fmt.Errorf("failed to write texture: %w", err)
	}
	return nil
}

func (a *generateCubeTextureAction) readResourceImage(path string) (*resource.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	image, err := resource.NewImageDecoder().Decode(file)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}
	return image, nil
}

func (a *generateCubeTextureAction) writeAssetCubeTexture(path string, texture *asset.CubeTexture) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	if err := asset.EncodeCubeTexture(file, texture); err != nil {
		return fmt.Errorf("failed to encode texture: %w", err)
	}
	return nil
}
