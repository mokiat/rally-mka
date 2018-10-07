package cubemap

import (
	"io"

	"github.com/mokiat/go-whiskey/data/storage"
	"github.com/pkg/errors"
)

type Encoder struct {
	CompressData bool
}

func (e Encoder) Encode(w io.Writer, texture *Texture) error {
	if e.CompressData {
		return errors.New("compression not supported")
	}
	writer := storage.NewTypedWriter(w)

	if err := writer.WriteUInt8(Version); err != nil {
		return errors.Wrap(err, "failed to write version")
	}

	if err := writer.WriteBool(e.CompressData); err != nil {
		return errors.Wrap(err, "failed to write compression flag")
	}

	if err := writer.WriteUInt8(uint8(texture.Format)); err != nil {
		return errors.Wrap(err, "failed to write format")
	}

	if err := writer.WriteUInt16(texture.Dimension); err != nil {
		return errors.Wrap(err, "failed to write dimension")
	}

	for _, side := range texture.Sides {
		if err := writer.WriteUInt32(uint32(len(side.Data))); err != nil {
			return errors.Wrap(err, "failed to write data size")
		}

		if err := writer.WriteBytes(side.Data); err != nil {
			return errors.Wrap(err, "failed to wrtie data")
		}
	}
	return nil
}
