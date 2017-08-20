package m2d

import (
	"io"

	"github.com/mokiat/rally-mka/data/storage"
)

type Color struct {
	R byte
	G byte
	B byte
	A byte
}

type Texture struct {
	width  uint16
	height uint16
	data   []byte
}

func NewTexture(width, height int) *Texture {
	return &Texture{
		width:  uint16(width),
		height: uint16(height),
		data:   make([]byte, width*height*4),
	}
}

func (t *Texture) Width() int {
	return int(t.width)
}

func (t *Texture) Height() int {
	return int(t.height)
}

func (t *Texture) Data() []byte {
	return t.data
}

func (t *Texture) SetTexel(x, y int, color Color) {
	offset := (y*int(t.width) + x) * 4
	t.data[offset+0] = color.R
	t.data[offset+1] = color.G
	t.data[offset+2] = color.B
	t.data[offset+3] = color.A
}

func (t *Texture) Texel(x, y int) Color {
	offset := (y*int(t.width) + x) * 4
	return Color{
		R: t.data[offset+0],
		G: t.data[offset+1],
		B: t.data[offset+2],
		A: t.data[offset+3],
	}
}

func (t *Texture) Load(in io.Reader) error {
	reader := storage.NewTypedReader(in)
	var err error
	if t.width, err = reader.ReadUInt16(); err != nil {
		return err
	}
	if t.height, err = reader.ReadUInt16(); err != nil {
		return err
	}
	t.data = make([]byte, int(t.width)*int(t.height)*4)
	if err = reader.ReadBytes(t.data); err != nil {
		return err
	}
	return nil
}

func (t *Texture) Save(out io.Writer) error {
	writer := storage.NewTypedWriter(out)
	if err := writer.WriteUInt16(t.width); err != nil {
		return err
	}
	if err := writer.WriteUInt16(t.height); err != nil {
		return err
	}
	if err := writer.WriteBytes(t.data); err != nil {
		return err
	}
	return nil
}
