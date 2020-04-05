package twodtex

import (
	"fmt"
	"os"

	cli "github.com/urfave/cli/v2"

	"github.com/mokiat/rally-mka/internal/data/asset"
	"github.com/mokiat/rally-mka/internal/data/resource"
)

func GenerateCommand() *cli.Command {
	return &cli.Command{
		Name:      "twodtex",
		Usage:     "generates a two dimensional texture",
		UsageText: "twodtex <image> <output-file>",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:  "width",
				Usage: "specifies a custom width of the texture",
			},
			&cli.IntFlag{
				Name:  "height",
				Usage: "specifies a custom height of the texture",
			},
		},
		Action: func(c *cli.Context) error {
			if c.Args().Len() < 2 {
				return fmt.Errorf("insufficient number of arguments: expected 2 got %d", c.Args().Len())
			}
			action := &generateTwoDTextureAction{
				ImageFile:  c.Args().Get(0),
				OutputFile: c.Args().Get(1),
				Width:      c.Int("width"),
				Height:     c.Int("height"),
			}
			return action.Run()
		},
	}
}

type generateTwoDTextureAction struct {
	ImageFile  string
	OutputFile string
	Width      int
	Height     int
}

func (a *generateTwoDTextureAction) Run() error {
	resImage, err := a.readResourceImage(a.ImageFile)
	if err != nil {
		return fmt.Errorf("failed to read image: %w", err)
	}

	if a.Width > 0 && a.Height > 0 {
		resImage.Scale(a.Width, a.Height)
	}

	texture := &asset.TwoDTexture{
		Format: asset.TextureFormatRGBA,
		Width:  uint16(resImage.Width()),
		Height: uint16(resImage.Height()),
		Data:   resImage.RGBAData(),
	}

	if err := a.writeAssetTexture(a.OutputFile, texture); err != nil {
		return fmt.Errorf("failed to write texture: %w", err)
	}
	return nil
}

func (a *generateTwoDTextureAction) readResourceImage(path string) (*resource.Image, error) {
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

func (a *generateTwoDTextureAction) writeAssetTexture(path string, texture *asset.TwoDTexture) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	if err := asset.NewTwoDTextureEncoder().Encode(file, texture); err != nil {
		return fmt.Errorf("failed to encode texture: %w", err)
	}
	return nil
}
