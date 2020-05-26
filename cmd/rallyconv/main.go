package main

import (
	"fmt"
	"os"

	"github.com/mokiat/rally-mka/internal/data/gltf"
)

func main() {
	gltfFile, err := os.Open("../models/Reference/reference.glb")
	panicOnErr(err)
	defer gltfFile.Close()

	gltfDoc, err := gltf.NewParser().Parse(gltfFile)
	panicOnErr(err)
	fmt.Printf("gltf doc: %+v\n", gltfDoc)

	// glLevel, err := gltf.Open("../models/export/tarnovo.exp2.gltf")
	// if err != nil {
	// 	panic(err)
	// }

	// level := &asset.Level{
	// 	SkyboxTexture: "city",
	// }

	// for _, node := range glLevel.Nodes {
	// 	getMesh := func() *gltf.Mesh {
	// 		return glLevel.Meshes[*node.Mesh]
	// 	}

	// 	getPositionAccessor := func() *gltf.Accessor {
	// 		return glLevel.Accessors[getMesh().Primitives[0].Attributes["POSITION"]]
	// 	}

	// 	// getPosition := func(index int) sprec.Vec3 {
	// 	// 	positionAccessor := *getPositionAccessor()
	// 	// 	positionView := glLevel.BufferViews[*positionAccessor.BufferView]
	// 	// 	buffer := data.Buffer(glLevel.Buffers[0].Data)
	// 	// 	offset := int(positionView.ByteOffset)
	// 	// 	x := buffer.Float32(offset+0*3*4) + float32(node.Translation[0])
	// 	// 	y := buffer.Float32(offset+1*3*4) + float32(node.Translation[1])
	// 	// 	z := buffer.Float32(offset+2*3*4) + float32(node.Translation[2])
	// 	// 	return sprec.NewVec3(x, y, z)
	// 	// }

	// 	getNormalAccessor := func() *gltf.Accessor {
	// 		return glLevel.Accessors[getMesh().Primitives[0].Attributes["NORMAL"]]
	// 	}

	// 	getTexCoordAccessor := func() *gltf.Accessor {
	// 		return glLevel.Accessors[getMesh().Primitives[0].Attributes["TEXCOORD_0"]]
	// 	}

	// 	getVertexCount := func() uint32 {
	// 		accessor := glLevel.Accessors[getMesh().Primitives[0].Attributes["POSITION"]]
	// 		return accessor.Count
	// 	}

	// 	getIndexAccessor := func() *gltf.Accessor {
	// 		return glLevel.Accessors[*getMesh().Primitives[0].Indices]
	// 	}

	// 	getIndexCount := func() uint32 {
	// 		return getIndexAccessor().Count
	// 	}

	// 	if !strings.HasPrefix(node.Name, "Road") {
	// 		continue
	// 	}

	// 	vertexStride := uint32(3*4 + 3*4 + 2*4)
	// 	vertexCount := getVertexCount()
	// 	vertexData := make([]byte, vertexCount*vertexStride)
	// 	positionView := glLevel.BufferViews[*getPositionAccessor().BufferView]
	// 	normalView := glLevel.BufferViews[*getNormalAccessor().BufferView]
	// 	texCoordView := glLevel.BufferViews[*getTexCoordAccessor().BufferView]
	// 	offset := 0
	// 	buffer := data.Buffer(vertexData)
	// 	for i := 0; i < int(vertexCount); i++ {
	// 		positionOffset := int(positionView.ByteOffset) + i*3*4
	// 		copy(vertexData[offset:], glLevel.Buffers[0].Data[positionOffset:positionOffset+3*4])
	// 		buffer.SetFloat32(offset+0*4, buffer.Float32(offset+0*4)+float32(node.Translation[0]))
	// 		buffer.SetFloat32(offset+1*4, buffer.Float32(offset+1*4)+float32(node.Translation[1]))
	// 		buffer.SetFloat32(offset+2*4, buffer.Float32(offset+2*4)+float32(node.Translation[2]))
	// 		offset += 3 * 4

	// 		normalOffset := int(normalView.ByteOffset) + i*3*4
	// 		copy(vertexData[offset:], glLevel.Buffers[0].Data[normalOffset:normalOffset+3*4])
	// 		offset += 3 * 4

	// 		texCoordOffset := int(texCoordView.ByteOffset) + i*2*4
	// 		copy(vertexData[offset:], glLevel.Buffers[0].Data[texCoordOffset:texCoordOffset+2*4])
	// 		offset += 2 * 4
	// 	}

	// 	indexData := make([]byte, getIndexCount()*2)
	// 	indexBuffer := data.Buffer(indexData)
	// 	indexView := *getIndexAccessor().BufferView
	// 	indexOffset := glLevel.BufferViews[indexView].ByteOffset
	// 	indexSize := glLevel.BufferViews[indexView].ByteLength
	// 	copy(indexData, glLevel.Buffers[0].Data[indexOffset:indexOffset+indexSize])

	// 	staticMesh := asset.Mesh{
	// 		VertexData:     vertexData,
	// 		VertexStride:   uint8(vertexStride),
	// 		CoordOffset:    0,
	// 		NormalOffset:   3 * 4,
	// 		TexCoordOffset: 3*4 + 3*4,
	// 		IndexData:      indexData,
	// 		SubMeshes: []asset.SubMesh{
	// 			{
	// 				Name:           "mesh",
	// 				IndexOffset:    0,
	// 				IndexCount:     getIndexCount(),
	// 				DiffuseTexture: "asphalt",
	// 			},
	// 		},
	// 	}
	// 	level.StaticMeshes = append(level.StaticMeshes, staticMesh)

	// 	collisionMesh := asset.LevelCollisionMesh{}

	// 	for i := 0; i < int(getIndexCount()); i += 3 {
	// 		indexA := indexBuffer.UInt16((i + 0) * 2)
	// 		indexB := indexBuffer.UInt16((i + 1) * 2)
	// 		indexC := indexBuffer.UInt16((i + 2) * 2)

	// 		vertexA := sprec.Vec3{
	// 			X: buffer.Float32(int(indexA)*int(vertexStride) + 0),
	// 			Y: buffer.Float32(int(indexA)*int(vertexStride) + 4),
	// 			Z: buffer.Float32(int(indexA)*int(vertexStride) + 8),
	// 		}
	// 		vertexB := sprec.Vec3{
	// 			X: buffer.Float32(int(indexB)*int(vertexStride) + 0),
	// 			Y: buffer.Float32(int(indexB)*int(vertexStride) + 4),
	// 			Z: buffer.Float32(int(indexB)*int(vertexStride) + 8),
	// 		}
	// 		vertexC := sprec.Vec3{
	// 			X: buffer.Float32(int(indexC)*int(vertexStride) + 0),
	// 			Y: buffer.Float32(int(indexC)*int(vertexStride) + 4),
	// 			Z: buffer.Float32(int(indexC)*int(vertexStride) + 8),
	// 		}
	// 		collisionMesh.Triangles = append(collisionMesh.Triangles, asset.Triangle{
	// 			asset.Point{vertexA.X, vertexA.Y, vertexA.Z},
	// 			asset.Point{vertexB.X, vertexB.Y, vertexB.Z},
	// 			asset.Point{vertexC.X, vertexC.Y, vertexC.Z},
	// 		})
	// 	}

	// 	level.CollisionMeshes = append(level.CollisionMeshes, collisionMesh)
	// }

	// if err := saveLevel(level, "assets/levels/tarnovo.dat"); err != nil {
	// 	panic(err)
	// }
}

func saveLevel(level *asset.Level, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := asset.NewLevelEncoder().Encode(file, level); err != nil {
		return err
	}
	return nil
}

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}
