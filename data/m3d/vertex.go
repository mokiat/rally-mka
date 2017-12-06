package m3d

import "github.com/mokiat/go-whiskey/data/storage"

type Vertex struct {
	X float32
	Y float32
	Z float32
}

func readVertex(reader storage.TypedReader) (vertex Vertex, err error) {
	if vertex.X, err = reader.ReadFloat32(); err != nil {
		return
	}
	if vertex.Y, err = reader.ReadFloat32(); err != nil {
		return
	}
	if vertex.Z, err = reader.ReadFloat32(); err != nil {
		return
	}
	return
}

func writeVertex(writer storage.TypedWriter, vertex Vertex) error {
	if err := writer.WriteFloat32(vertex.X); err != nil {
		return err
	}
	if err := writer.WriteFloat32(vertex.Y); err != nil {
		return err
	}
	if err := writer.WriteFloat32(vertex.Z); err != nil {
		return err
	}
	return nil
}
