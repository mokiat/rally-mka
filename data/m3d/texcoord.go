package m3d

import "github.com/mokiat/go-whiskey/data/storage"

type TexCoord struct {
	U float32
	V float32
}

func readTexCoord(reader storage.TypedReader) (texCoord TexCoord, err error) {
	if texCoord.U, err = reader.ReadFloat32(); err != nil {
		return
	}
	if texCoord.V, err = reader.ReadFloat32(); err != nil {
		return
	}
	return
}

func writeTexCoord(writer storage.TypedWriter, texCoord TexCoord) error {
	if err := writer.WriteFloat32(texCoord.U); err != nil {
		return err
	}
	if err := writer.WriteFloat32(texCoord.V); err != nil {
		return err
	}
	return nil
}
