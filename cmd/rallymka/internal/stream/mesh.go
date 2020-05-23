package stream

import (
	"fmt"

	"github.com/mokiat/lacking/data/asset"
	"github.com/mokiat/lacking/graphics"
	"github.com/mokiat/rally-mka/internal/engine/resource"
)

const meshResourceType = "mesh"

func GetMesh(registry *resource.Registry, name string) MeshHandle {
	return MeshHandle{
		Handle: registry.Type(meshResourceType).Resource(name),
	}
}

type MeshHandle struct {
	resource.Handle
}

func (h MeshHandle) Get() *Mesh {
	return h.Handle.Get().(*Mesh)
}

func (h MeshHandle) IsAvailable() bool {
	return h.Handle.IsAvailable() && h.Get().IsAvailable()
}

type Mesh struct {
	VertexArray *graphics.VertexArray
	SubMeshes   []SubMesh
}

func (m Mesh) IsAvailable() bool {
	for _, subMesh := range m.SubMeshes {
		if !subMesh.IsAvailable() {
			return false
		}
	}
	return true
}

type SubMesh struct {
	IndexOffset    int
	IndexCount     int32
	DiffuseTexture *TwoDTextureHandle
}

func (m SubMesh) IsAvailable() bool {
	if m.DiffuseTexture != nil {
		return m.DiffuseTexture.IsAvailable()
	}
	return true
}

func NewMeshOperator(locator resource.Locator, gfxWorker *graphics.Worker) *MeshOperator {
	return &MeshOperator{
		locator:   locator,
		gfxWorker: gfxWorker,
	}
}

type MeshOperator struct {
	locator   resource.Locator
	gfxWorker *graphics.Worker
}

func (o *MeshOperator) Register(registry *resource.Registry) {
	registry.RegisterType(meshResourceType, o)
}

func (o *MeshOperator) Allocate(registry *resource.Registry, name string) (resource.Resource, error) {
	in, err := o.locator.Open("assets", "meshes", name)
	if err != nil {
		return nil, fmt.Errorf("failed to open mesh asset %q: %w", name, err)
	}
	defer in.Close()

	meshAsset := new(asset.Mesh)
	if err := asset.DecodeMesh(in, meshAsset); err != nil {
		return nil, fmt.Errorf("failed to decode mesh asset %q: %w", name, err)
	}

	return AllocateMesh(registry, o.gfxWorker, meshAsset)
}

func (o *MeshOperator) Release(registry *resource.Registry, resource resource.Resource) error {
	mesh := resource.(*Mesh)
	return ReleaseMesh(registry, o.gfxWorker, mesh)
}

func AllocateMesh(registry *resource.Registry, gfxWorker *graphics.Worker, meshAsset *asset.Mesh) (*Mesh, error) {
	mesh := &Mesh{
		VertexArray: &graphics.VertexArray{},
	}

	gfxTask := gfxWorker.Schedule(func() error {
		return mesh.VertexArray.Allocate(graphics.VertexArrayData{
			VertexData:     meshAsset.VertexData,
			VertexStride:   int32(meshAsset.VertexStride),
			CoordOffset:    int(meshAsset.CoordOffset),
			NormalOffset:   int(meshAsset.NormalOffset),
			TexCoordOffset: int(meshAsset.TexCoordOffset),
			IndexData:      meshAsset.IndexData,
		})
	})
	if err := gfxTask.Wait(); err != nil {
		return nil, fmt.Errorf("failed to allocate gfx vertex array: %w", err)
	}

	mesh.SubMeshes = make([]SubMesh, len(meshAsset.SubMeshes))
	for i := range mesh.SubMeshes {
		subMeshAsset := meshAsset.SubMeshes[i]
		subMesh := SubMesh{
			IndexOffset: int(subMeshAsset.IndexOffset),
			IndexCount:  int32(subMeshAsset.IndexCount),
		}
		if subMeshAsset.DiffuseTexture != "" {
			diffuseTexture := GetTwoDTexture(registry, subMeshAsset.DiffuseTexture)
			registry.Request(diffuseTexture.Handle)
			subMesh.DiffuseTexture = &diffuseTexture
		}
		mesh.SubMeshes[i] = subMesh
	}
	return mesh, nil
}

func ReleaseMesh(registry *resource.Registry, gfxWorker *graphics.Worker, mesh *Mesh) error {
	for _, subMesh := range mesh.SubMeshes {
		if subMesh.DiffuseTexture != nil {
			registry.Dismiss(subMesh.DiffuseTexture.Handle)
		}
	}

	gfxTask := gfxWorker.Schedule(func() error {
		return mesh.VertexArray.Release()
	})
	if err := gfxTask.Wait(); err != nil {
		return fmt.Errorf("failed to release gfx vertex array: %w", err)
	}
	return nil
}
