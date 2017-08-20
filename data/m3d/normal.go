package m3d

import "github.com/mokiat/rally-mka/data/storage"

type Normal struct {
	X float32
	Y float32
	Z float32
}

func readNormal(reader storage.TypedReader) (normal Normal, err error) {
	if normal.X, err = reader.ReadFloat32(); err != nil {
		return
	}
	if normal.Y, err = reader.ReadFloat32(); err != nil {
		return
	}
	if normal.Z, err = reader.ReadFloat32(); err != nil {
		return
	}
	return
}

func writeNormal(writer storage.TypedWriter, normal Normal) error {
	if err := writer.WriteFloat32(normal.X); err != nil {
		return err
	}
	if err := writer.WriteFloat32(normal.Y); err != nil {
		return err
	}
	if err := writer.WriteFloat32(normal.Z); err != nil {
		return err
	}
	return nil
}
