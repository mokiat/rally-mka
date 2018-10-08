package imgutil

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/pkg/errors"
)

func OpenImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open image file")
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode image")
	}
	return img, nil
}
