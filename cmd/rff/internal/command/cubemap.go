package command

import (
	"os"

	"github.com/mokiat/rally-mka/cmd/rff/internal/imgutil"
	"github.com/mokiat/rally-mka/data/cubemap"
	"github.com/pkg/errors"
	cli "gopkg.in/urfave/cli.v1"
)

func GenerateCubemap() cli.Command {
	return cli.Command{
		Name:      "generate-cubemap",
		Usage:     "generates a cubemap texture",
		UsageText: "generate-cubemap <front-image> <back-image> <left-image> <right-image> <top-image> <bottom-image> <output-file>",
		Action: func(c *cli.Context) error {
			if len(c.Args()) < 7 {
				return errors.New("insufficient arguments")
			}
			action := generateCubemapAction{
				FrontFile:  c.Args().Get(0),
				BackFile:   c.Args().Get(1),
				LeftFile:   c.Args().Get(2),
				RightFile:  c.Args().Get(3),
				TopFile:    c.Args().Get(4),
				BottomFile: c.Args().Get(5),
				OutputFile: c.Args().Get(6),
			}
			return action.Run()
		},
	}
}

type generateCubemapAction struct {
	FrontFile  string
	BackFile   string
	LeftFile   string
	RightFile  string
	TopFile    string
	BottomFile string
	OutputFile string
}

func (a generateCubemapAction) Run() error {
	frontImg, err := imgutil.OpenImage(a.FrontFile)
	if err != nil {
		return errors.Wrap(err, "failed to load front image")
	}
	if !imgutil.IsSquareImage(frontImg) {
		return errors.New("front image is not a square")
	}

	backImg, err := imgutil.OpenImage(a.BackFile)
	if err != nil {
		return errors.Wrap(err, "failed to load back image")
	}
	if !imgutil.IsSquareImage(backImg) {
		return errors.New("back image is not a square")
	}

	leftImg, err := imgutil.OpenImage(a.LeftFile)
	if err != nil {
		return errors.Wrap(err, "failed to load left image")
	}
	if !imgutil.IsSquareImage(leftImg) {
		return errors.New("left image is not a square")
	}

	rightImg, err := imgutil.OpenImage(a.RightFile)
	if err != nil {
		return errors.Wrap(err, "failed to load right image")
	}
	if !imgutil.IsSquareImage(rightImg) {
		return errors.New("right image is not a square")
	}

	topImg, err := imgutil.OpenImage(a.TopFile)
	if err != nil {
		return errors.Wrap(err, "failed to load top image")
	}
	if !imgutil.IsSquareImage(topImg) {
		return errors.New("top image is not a square")
	}

	bottomImg, err := imgutil.OpenImage(a.BottomFile)
	if err != nil {
		return errors.Wrap(err, "failed to load bottom image")
	}
	if !imgutil.IsSquareImage(frontImg) {
		return errors.New("bottom image is not a square")
	}

	areSameDimension := frontImg.Bounds().Dx() == backImg.Bounds().Dx() &&
		frontImg.Bounds().Dx() == leftImg.Bounds().Dx() &&
		frontImg.Bounds().Dx() == rightImg.Bounds().Dx() &&
		frontImg.Bounds().Dx() == topImg.Bounds().Dx() &&
		frontImg.Bounds().Dx() == bottomImg.Bounds().Dx()
	if !areSameDimension {
		return errors.New("images are not of the same size")
	}

	texture := &cubemap.Texture{
		Format:    cubemap.DataFormatRGBA,
		Dimension: uint16(frontImg.Bounds().Dx()),
	}
	texture.Sides[cubemap.SideFront] = cubemap.TextureSide{
		Data: imgutil.ExtractImageData(frontImg),
	}
	texture.Sides[cubemap.SideBack] = cubemap.TextureSide{
		Data: imgutil.ExtractImageData(backImg),
	}
	texture.Sides[cubemap.SideLeft] = cubemap.TextureSide{
		Data: imgutil.ExtractImageData(leftImg),
	}
	texture.Sides[cubemap.SideRight] = cubemap.TextureSide{
		Data: imgutil.ExtractImageData(rightImg),
	}
	texture.Sides[cubemap.SideTop] = cubemap.TextureSide{
		Data: imgutil.ExtractImageData(topImg),
	}
	texture.Sides[cubemap.SideBottom] = cubemap.TextureSide{
		Data: imgutil.ExtractImageData(bottomImg),
	}

	if err := saveCubemap(a.OutputFile, texture); err != nil {
		return errors.Wrap(err, "failed to store cubemap")
	}
	return nil
}

func saveCubemap(path string, texture *cubemap.Texture) error {
	file, err := os.Create(path)
	if err != nil {
		return errors.Wrap(err, "failed to create file")
	}
	defer file.Close()

	encoder := &cubemap.Encoder{
		CompressData: false, // compression not supported yet
	}
	if err := encoder.Encode(file, texture); err != nil {
		return errors.Wrap(err, "failed to encode cubemap")
	}
	return nil
}
