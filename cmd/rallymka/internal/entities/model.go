package entities

// import (
// 	"strings"

// 	"github.com/mokiat/go-whiskey/math"
// 	"github.com/mokiat/rally-mka/internal/data/m2d"
// 	"github.com/mokiat/rally-mka/internal/data/m3d"
// )

// type ExtendedModel struct {
// 	*m3d.Model
// 	textures map[string]*m2d.Texture

// 	MinX float32
// 	MaxX float32
// 	MinY float32
// 	MaxY float32
// 	MinZ float32
// 	MaxZ float32
// }

// func NewExtendedModel() *ExtendedModel {
// 	return &ExtendedModel{
// 		Model:    m3d.NewModel(),
// 		textures: make(map[string]*m2d.Texture),
// 	}
// }

// func (m *ExtendedModel) GetObjectIndex(index int) *m3d.Object {
// 	return m.Objects[index]
// }

// func (m *ExtendedModel) FindObjectIndex(name string, startIndex int, isSubstring bool) int {
// 	if startIndex >= len(m.Objects) {
// 		return -2
// 	}
// 	for i := startIndex; i < len(m.Objects); i++ {
// 		object := m.Objects[i]
// 		if isSubstring {
// 			if object.Name == "" {
// 				continue
// 			}
// 			if strings.HasPrefix(object.Name, name) {
// 				return i
// 			}
// 		} else {
// 			if name == object.Name {
// 				return i
// 			}
// 		}
// 	}
// 	return -1
// }

// func (m *ExtendedModel) DoCenter(object *m3d.Object) {
// 	center := m.ObjectCenter(object)
// 	for i := range object.Vertices {
// 		object.Vertices[i].X -= center.X
// 		object.Vertices[i].Y -= center.Y
// 		object.Vertices[i].Z -= center.Z
// 	}
// }

// func (m *ExtendedModel) ObjectCenter(object *m3d.Object) math.Vec3 {
// 	if len(object.Vertices) == 0 {
// 		return math.NullVec3()
// 	}
// 	center := math.NullVec3()
// 	for _, vertex := range object.Vertices {
// 		center = center.IncCoords(vertex.X, vertex.Y, vertex.Z)
// 	}
// 	center = center.Div(float32(len(object.Vertices)))
// 	return center
// }

// func (m *ExtendedModel) EvaluateMinMax() {
// 	for _, object := range m.Objects {
// 		for _, vertex := range object.Vertices {
// 			if vertex.X < m.MinX {
// 				m.MinX = vertex.X
// 			}
// 			if vertex.X > m.MaxX {
// 				m.MaxX = vertex.X
// 			}
// 			if vertex.Y < m.MinY {
// 				m.MinY = vertex.Y
// 			}
// 			if vertex.Y > m.MaxY {
// 				m.MaxY = vertex.Y
// 			}
// 			if vertex.Z < m.MinZ {
// 				m.MinZ = vertex.Z
// 			}
// 			if vertex.Z > m.MaxZ {
// 				m.MaxZ = vertex.Z
// 			}
// 		}
// 	}
// }
