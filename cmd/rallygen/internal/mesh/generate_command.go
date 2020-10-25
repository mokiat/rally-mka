package mesh

import (
	"errors"
	"fmt"

	"github.com/mokiat/lacking/data"
	"github.com/mokiat/lacking/data/asset"
	"github.com/mokiat/rally-mka/internal/data/resource"
)

func ConvertResourceToAsset(resMesh *resource.Mesh) (*asset.Mesh, error) {
	if len(resMesh.Coords) == 0 {
		return nil, fmt.Errorf("missing coords")
	}
	if len(resMesh.SubMeshes) == 0 {
		return nil, fmt.Errorf("missing sub meshes")
	}
	if len(resMesh.Coords)%3 != 0 {
		return nil, errors.New("coords values not multiple of three")
	}
	vertexCount := len(resMesh.Coords) / 3
	layout := evaluateMeshVertexLayout(resMesh)
	vertexData := data.Buffer(make([]byte, vertexCount*layout.Stride))
	if len(resMesh.Coords) > 0 {
		offset := layout.CoordOffset
		for index := 0; index < len(resMesh.Coords); index += 3 {
			vertexData.SetFloat32(offset+0, resMesh.Coords[index+0])
			vertexData.SetFloat32(offset+4, resMesh.Coords[index+1])
			vertexData.SetFloat32(offset+8, resMesh.Coords[index+2])
			offset += layout.Stride
		}
	}
	if len(resMesh.Normals) > 0 {
		offset := layout.NormalOffset
		for index := 0; index < len(resMesh.Normals); index += 3 {
			vertexData.SetFloat32(offset+0, resMesh.Normals[index+0])
			vertexData.SetFloat32(offset+4, resMesh.Normals[index+1])
			vertexData.SetFloat32(offset+8, resMesh.Normals[index+2])
			offset += layout.Stride
		}
	}
	if len(resMesh.TexCoords) > 0 {
		offset := layout.TexCoordOffset
		for index := 0; index < len(resMesh.TexCoords); index += 2 {
			vertexData.SetFloat32(offset+0, resMesh.TexCoords[index+0])
			vertexData.SetFloat32(offset+4, resMesh.TexCoords[index+1])
			offset += layout.Stride
		}
	}
	indexCount := len(resMesh.Indices)
	indexData := data.Buffer(make([]byte, indexCount*2))
	for i, index := range resMesh.Indices {
		indexData.SetUInt16(i*2, uint16(index))
	}
	subMeshes := make([]asset.SubMesh, len(resMesh.SubMeshes))
	for i, resSubMesh := range resMesh.SubMeshes {
		subMeshes[i] = asset.SubMesh{
			Primitive:   asset.PrimitiveTriangles,
			IndexOffset: uint32(resSubMesh.IndexOffset * 2),
			IndexCount:  uint32(resSubMesh.IndexCount),
			Material: asset.Material{
				Type:         "pbr",
				Metalness:    0.0,
				Roughness:    0.5,
				ColorTexture: resSubMesh.DiffuseTexture,
			},
		}
	}
	return &asset.Mesh{
		VertexData: vertexData,
		VertexLayout: asset.VertexLayout{
			CoordOffset:    int16(layout.CoordOffset),
			CoordStride:    int16(layout.Stride),
			NormalOffset:   int16(layout.NormalOffset),
			NormalStride:   int16(layout.Stride),
			TangentOffset:  asset.UnspecifiedOffset,
			TangentStride:  int16(layout.Stride),
			TexCoordOffset: int16(layout.TexCoordOffset),
			TexCoordStride: int16(layout.Stride),
			ColorOffset:    asset.UnspecifiedOffset,
			ColorStride:    int16(layout.Stride),
		},
		IndexData: indexData,
		SubMeshes: subMeshes,
	}, nil
}

func evaluateMeshVertexLayout(mesh *resource.Mesh) meshVertexLayout {
	var layout meshVertexLayout
	if len(mesh.Coords) > 0 {
		layout.CoordOffset = layout.Stride
		layout.Stride += 3 * 4
	} else {
		layout.CoordOffset = int(asset.UnspecifiedOffset)
	}
	if len(mesh.Normals) > 0 {
		layout.NormalOffset = layout.Stride
		layout.Stride += 3 * 4
	} else {
		layout.NormalOffset = int(asset.UnspecifiedOffset)
	}
	if len(mesh.TexCoords) > 0 {
		layout.TexCoordOffset += layout.Stride
		layout.Stride += 2 * 4
	} else {
		layout.TexCoordOffset = int(asset.UnspecifiedOffset)
	}
	return layout
}

type meshVertexLayout struct {
	Stride         int
	CoordOffset    int
	NormalOffset   int
	TexCoordOffset int
}
