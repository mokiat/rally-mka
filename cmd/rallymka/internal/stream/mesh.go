package stream

import (
	"fmt"

	"github.com/mokiat/rally-mka/cmd/rallymka/internal/graphics"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/resource"
	"github.com/mokiat/rally-mka/internal/data/asset"
)

const meshResourceType = "mesh"

func GetMesh(registry *resource.Registry, name string) *Mesh {
	return registry.ResourceType(meshResourceType).Resource(name).(*Mesh)
}

type Mesh struct {
	*resource.Handle
	gfxVertexArray *graphics.VertexArray
	subMeshes      []SubMesh
}

type SubMesh struct {
	IndexOffset    int
	IndexCount     int
	DiffuseTexture *TwoDTexture
}

func (m *Mesh) Gfx() *graphics.VertexArray {
	return m.gfxVertexArray
}

func (m *Mesh) SubMeshes() []SubMesh {
	return m.subMeshes
}

func NewMeshController(capacity int, gfxWorker *graphics.Worker) MeshController {
	return MeshController{
		meshes:    make([]Mesh, capacity),
		gfxWorker: gfxWorker,
	}
}

type MeshController struct {
	meshes    []Mesh
	gfxWorker *graphics.Worker
}

func (c MeshController) ResourceTypeName() string {
	return meshResourceType
}

func (c MeshController) Init(index int, handle *resource.Handle) resource.Resource {
	c.meshes[index] = Mesh{
		Handle:         handle,
		gfxVertexArray: &graphics.VertexArray{},
		subMeshes:      make([]SubMesh, 0),
	}
	return &c.meshes[index]
}

func (c MeshController) Load(index int, locator resource.Locator, registry *resource.Registry) error {
	mesh := &c.meshes[index]

	in, err := locator.Open("assets", "meshes", mesh.Name())
	if err != nil {
		return fmt.Errorf("failed to open mesh asset %q: %w", mesh.Name(), err)
	}
	defer in.Close()

	meshAsset, err := asset.NewMeshDecoder().Decode(in)
	if err != nil {
		return fmt.Errorf("failed to decode mesh asset %q: %w", mesh.Name(), err)
	}

	gfxTask := func() error {
		return mesh.gfxVertexArray.Allocate(graphics.VertexArrayData{
			VertexData:     meshAsset.VertexData,
			VertexStride:   int32(meshAsset.VertexStride),
			CoordOffset:    int(meshAsset.CoordOffset),
			NormalOffset:   int(meshAsset.NormalOffset),
			TexCoordOffset: int(meshAsset.TexCoordOffset),
			IndexData:      meshAsset.IndexData,
		})
	}
	if err := c.gfxWorker.Run(gfxTask); err != nil {
		return fmt.Errorf("failed to allocate gfx vertex array: %w", err)
	}

	mesh.subMeshes = make([]SubMesh, len(meshAsset.SubMeshes))
	for i := range mesh.subMeshes {
		mesh.subMeshes[i] = SubMesh{
			IndexOffset: int(meshAsset.SubMeshes[i].IndexOffset),
			IndexCount:  int(meshAsset.SubMeshes[i].IndexCount),
		}
		if meshAsset.SubMeshes[i].DiffuseTexture != "" {
			diffuseTextire := GetTwoDTexture(registry, meshAsset.SubMeshes[i].DiffuseTexture)
			diffuseTextire.Request()
			mesh.subMeshes[i].DiffuseTexture = diffuseTextire
		}
	}

	return nil
}

func (c MeshController) Unload(index int) error {
	mesh := &c.meshes[index]

	for _, subMesh := range mesh.subMeshes {
		if subMesh.DiffuseTexture != nil {
			subMesh.DiffuseTexture.Dismiss()
		}
	}
	mesh.subMeshes = make([]SubMesh, 0)

	gfxTask := func() error {
		return mesh.gfxVertexArray.Release()
	}
	if err := c.gfxWorker.Run(gfxTask); err != nil {
		return fmt.Errorf("failed to release gfx vertex array: %w", err)
	}
	return nil
}
