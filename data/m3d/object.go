package m3d

import "github.com/mokiat/go-whiskey/data/storage"

type objectFlag uint8

const (
	hasVertices objectFlag = 1 << iota
	hasNormals
	hasColors
	hasTexCoords
)

type Object struct {
	Name      string
	Texture   string
	Vertices  []Vertex
	Normals   []Normal
	Colors    []Color
	TexCoords []TexCoord
	Faces     []Face
}

func NewObject() *Object {
	return &Object{}
}

func (o *Object) HasTexture() bool {
	return o.Texture != ""
}

func (o *Object) HasVertices() bool {
	return o.Vertices != nil
}

func (o *Object) HasNormals() bool {
	return o.Normals != nil
}

func (o *Object) HasColors() bool {
	return o.Colors != nil
}

func (o *Object) HasTexCoords() bool {
	return o.TexCoords != nil
}

func (o *Object) load(reader storage.TypedReader) error {
	var err error
	if o.Name, err = reader.ReadString8(); err != nil {
		return err
	}
	if o.Texture, err = reader.ReadString8(); err != nil {
		return err
	}
	var mask uint8
	if mask, err = reader.ReadUInt8(); err != nil {
		return err
	}
	var count uint32
	if count, err = reader.ReadUInt32(); err != nil {
		return err
	}
	if (objectFlag(mask) & hasVertices) != 0 {
		o.Vertices = make([]Vertex, count)
		for i := range o.Vertices {
			if vertex, err := readVertex(reader); err != nil {
				return err
			} else {
				o.Vertices[i] = vertex
			}
		}
	}
	if (objectFlag(mask) & hasNormals) != 0 {
		o.Normals = make([]Normal, count)
		for i := range o.Normals {
			if normal, err := readNormal(reader); err != nil {
				return err
			} else {
				o.Normals[i] = normal
			}
		}
	}
	if (objectFlag(mask) & hasColors) != 0 {
		o.Colors = make([]Color, count)
		for i := range o.Colors {
			if color, err := readColor(reader); err != nil {
				return err
			} else {
				o.Colors[i] = color
			}
		}
	}
	if (objectFlag(mask) & hasTexCoords) != 0 {
		o.TexCoords = make([]TexCoord, count)
		for i := range o.TexCoords {
			if texCoord, err := readTexCoord(reader); err != nil {
				return err
			} else {
				o.TexCoords[i] = texCoord
			}
		}
	}
	var faceCount uint32
	if faceCount, err = reader.ReadUInt32(); err != nil {
		return err
	}
	o.Faces = make([]Face, faceCount)
	for i := range o.Faces {
		if face, err := readFace(reader); err != nil {
			return err
		} else {
			o.Faces[i] = face
		}
	}
	return nil
}

func (o *Object) save(writer storage.TypedWriter) error {
	if err := writer.WriteString8(o.Name); err != nil {
		return err
	}
	if err := writer.WriteString8(o.Texture); err != nil {
		return err
	}
	mask := objectFlag(0)
	if o.HasVertices() {
		mask |= hasVertices
	}
	if o.HasNormals() {
		mask |= hasNormals
	}
	if o.HasColors() {
		mask |= hasColors
	}
	if o.HasTexCoords() {
		mask |= hasTexCoords
	}
	if err := writer.WriteUInt8(uint8(mask)); err != nil {
		return err
	}
	if err := writer.WriteUInt32(uint32(len(o.Vertices))); err != nil {
		return err
	}
	if o.HasVertices() {
		for _, vertex := range o.Vertices {
			if err := writeVertex(writer, vertex); err != nil {
				return err
			}
		}
	}
	if o.HasNormals() {
		for _, normal := range o.Normals {
			if err := writeNormal(writer, normal); err != nil {
				return err
			}
		}
	}
	if o.HasColors() {
		for _, color := range o.Colors {
			if err := writeColor(writer, color); err != nil {
				return err
			}
		}
	}
	if o.HasTexCoords() {
		for _, texCoord := range o.TexCoords {
			if err := writeTexCoord(writer, texCoord); err != nil {
				return err
			}
		}
	}
	if err := writer.WriteUInt32(uint32(len(o.Faces))); err != nil {
		return err
	}
	for _, face := range o.Faces {
		if err := writeFace(writer, face); err != nil {
			return err
		}
	}
	return nil
}
