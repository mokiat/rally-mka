package m3d

import "github.com/mokiat/go-whiskey/data/storage"

type Face struct {
	IndexA uint32
	IndexB uint32
	IndexC uint32
}

func readFace(reader storage.TypedReader) (face Face, err error) {
	if face.IndexA, err = reader.ReadUInt32(); err != nil {
		return
	}
	if face.IndexB, err = reader.ReadUInt32(); err != nil {
		return
	}
	if face.IndexC, err = reader.ReadUInt32(); err != nil {
		return
	}
	return
}

func writeFace(writer storage.TypedWriter, face Face) error {
	if err := writer.WriteUInt32(face.IndexA); err != nil {
		return err
	}
	if err := writer.WriteUInt32(face.IndexB); err != nil {
		return err
	}
	if err := writer.WriteUInt32(face.IndexC); err != nil {
		return err
	}
	return nil
}
