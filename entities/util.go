package entities

import (
	"github.com/mokiat/go-whiskey-gl/buffer"
	"github.com/mokiat/go-whiskey-gl/texture"
	"github.com/mokiat/go-whiskey/math"
	"github.com/mokiat/rally-mka/collision"
	"github.com/mokiat/rally-mka/internal/data/m3d"
	"github.com/mokiat/rally-mka/render"
)

func createCollisionMesh(obj *m3d.Object) *collision.Mesh {
	var triangles []collision.Triangle
	for _, face := range obj.Faces {
		a := vertexToVector(obj.Vertices[face.IndexA])
		b := vertexToVector(obj.Vertices[face.IndexB])
		c := vertexToVector(obj.Vertices[face.IndexC])
		triangles = append(triangles, collision.MakeTriangle(a, b, c))
	}
	return collision.NewMesh(triangles)
}

func vertexToVector(vertex m3d.Vertex) math.Vec3 {
	return math.Vec3{
		X: vertex.X,
		Y: vertex.Y,
		Z: vertex.Z,
	}
}

func createRenderMesh(model *ExtendedModel, object *m3d.Object) *render.Mesh {
	data := &render.MeshData{}
	if object.HasVertices() {
		data.CoordOffset = data.Stride
		data.Stride += 3
	} else {
		data.CoordOffset = render.DisabledOffset
	}
	if object.HasNormals() {
		data.NormalOffset = data.Stride
		data.Stride += 3
	} else {
		data.NormalOffset = render.DisabledOffset
	}
	if object.HasColors() {
		data.ColorOffset = data.Stride
		data.Stride += 4
	} else {
		data.ColorOffset = render.DisabledOffset
	}
	if object.HasTexCoords() {
		data.TexCoordOffset = data.Stride
		data.Stride += 2
	} else {
		data.TexCoordOffset = render.DisabledOffset
	}

	data.VertexData = buffer.DedicatedFloat32DataPlayground(len(object.Vertices) * data.Stride)
	if object.HasVertices() {
		writer := buffer.NewFloat32DataWriter(data.VertexData, data.CoordOffset, data.Stride)
		for _, vertex := range object.Vertices {
			writer.PutValue3(vertex.X, vertex.Y, vertex.Z)
		}
	}
	if object.HasNormals() {
		writer := buffer.NewFloat32DataWriter(data.VertexData, data.NormalOffset, data.Stride)
		for _, normal := range object.Normals {
			writer.PutValue3(normal.X, normal.Y, normal.Z)
		}
	}
	if object.HasColors() {
		writer := buffer.NewFloat32DataWriter(data.VertexData, data.ColorOffset, data.Stride)
		for _, color := range object.Colors {
			writer.PutValue4(color.R, color.G, color.B, color.A)
		}
	}
	if object.HasTexCoords() {
		writer := buffer.NewFloat32DataWriter(data.VertexData, data.TexCoordOffset, data.Stride)
		for _, texCoord := range object.TexCoords {
			writer.PutValue2(texCoord.U, texCoord.V)
		}
	}
	data.IndexData = buffer.DedicatedUInt16DataPlayground(len(object.Faces) * 3)
	{
		writer := buffer.NewUInt16DataWriter(data.IndexData, 1)
		for _, face := range object.Faces {
			writer.PutValue(uint16(face.IndexA))
			writer.PutValue(uint16(face.IndexB))
			writer.PutValue(uint16(face.IndexC))
		}
	}
	if object.HasTexture() {
		tex2d := model.textures[object.Texture]
		rgbaData := texture.DedicatedRGBAFlatDataPlayground(tex2d.Width(), tex2d.Height())
		rgbaData.SetData(tex2d.Data())
		data.TextureData = rgbaData
	}
	return render.NewMesh(data)
}
