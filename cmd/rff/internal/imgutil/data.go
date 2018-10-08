package imgutil

import (
	"image"
	"image/color"
)

func ExtractImageData(img image.Image) []byte {
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()
	data := make([]byte, 4*width*height)

	offset := 0
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, a := GetImageColor(img, x, height-y-1).RGBA()
			data[offset+0] = byte(r >> 8)
			data[offset+1] = byte(g >> 8)
			data[offset+2] = byte(b >> 8)
			data[offset+3] = byte(a >> 8)
			offset += 4
		}
	}

	return data
}

func GetImageColor(img image.Image, x, y int) color.Color {
	return img.At(
		img.Bounds().Min.X+x,
		img.Bounds().Min.Y+y,
	)
}
