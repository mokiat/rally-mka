package m3d

import "github.com/mokiat/go-whiskey/data/storage"

type Color struct {
	R float32
	G float32
	B float32
	A float32
}

func readColor(reader storage.TypedReader) (color Color, err error) {
	if color.R, err = reader.ReadFloat32(); err != nil {
		return
	}
	if color.G, err = reader.ReadFloat32(); err != nil {
		return
	}
	if color.B, err = reader.ReadFloat32(); err != nil {
		return
	}
	if color.A, err = reader.ReadFloat32(); err != nil {
		return
	}
	return
}

func writeColor(writer storage.TypedWriter, color Color) error {
	if err := writer.WriteFloat32(color.R); err != nil {
		return err
	}
	if err := writer.WriteFloat32(color.G); err != nil {
		return err
	}
	if err := writer.WriteFloat32(color.B); err != nil {
		return err
	}
	if err := writer.WriteFloat32(color.A); err != nil {
		return err
	}
	return nil
}
