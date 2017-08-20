package m3d

import "io"
import "github.com/mokiat/rally-mka/data/storage"

type Model struct {
	Objects []*Object
}

func NewModel() *Model {
	return &Model{}
}

func (m *Model) Load(in io.Reader) error {
	reader := storage.NewTypedReader(in)
	count, err := reader.ReadUInt16()
	if err != nil {
		return err
	}
	m.Objects = make([]*Object, count)
	for i := range m.Objects {
		object := NewObject()
		if err := object.load(reader); err != nil {
			return err
		}
		m.Objects[i] = object
	}

	return nil
}

func (m *Model) Save(out io.Writer) error {
	writer := storage.NewTypedWriter(out)
	if err := writer.WriteUInt16(uint16(len(m.Objects))); err != nil {
		return err
	}
	for _, obj := range m.Objects {
		if err := obj.save(writer); err != nil {
			return err
		}
	}
	return nil
}
