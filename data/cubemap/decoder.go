package cubemap

import (
	"io"

	"github.com/mokiat/go-whiskey/data/storage"
	"github.com/pkg/errors"
)

type BytesAllocator interface {
	Allocate(size int) []byte
}

type Decoder struct {
	Allocator BytesAllocator
}

func (d Decoder) Decode(in io.Reader) (*Texture, error) {
	reader := storage.NewTypedReader(in)

	version, err := reader.ReadUInt8()
	if err != nil {
		return nil, errors.Wrap(err, "failed to read version")
	}
	if version != Version {
		return nil, errors.Errorf("unsupported version: %d", version)
	}

	compressed, err := reader.ReadBool()
	if err != nil {
		return nil, errors.Wrap(err, "failed to read compression flag")
	}
	if compressed {
		return nil, errors.New("compression not supported")
	}

	format, err := reader.ReadUInt8()
	if err != nil {
		return nil, errors.Wrap(err, "failed to read format")
	}

	dimension, err := reader.ReadUInt16()
	if err != nil {
		return nil, errors.Wrap(err, "failed to read dimension")
	}

	texture := &Texture{
		Format:    DataFormat(format),
		Dimension: dimension,
	}

	for i := range texture.Sides {
		dataSize, err := reader.ReadUInt32()
		if err != nil {
			return nil, errors.Wrap(err, "failed to read data size")
		}

		data := d.allocateBytes(int(dataSize))
		if err := reader.ReadBytes(data); err != nil {
			return nil, errors.Wrap(err, "failed to read data")
		}
		texture.Sides[i].Data = data
	}
	return texture, nil
}

func (d Decoder) allocateBytes(count int) []byte {
	if d.Allocator == nil {
		return make([]byte, count)
	}
	return d.Allocator.Allocate(count)
}
