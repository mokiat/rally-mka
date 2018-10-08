package imgutil

import "image"

func IsSquareImage(img image.Image) bool {
	return img.Bounds().Dx() == img.Bounds().Dy()
}
