package stream

import (
	"fmt"

	"github.com/mokiat/rally-mka/internal/data/asset"
	"github.com/mokiat/rally-mka/internal/engine/graphics"
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

type Mesh struct {
	VertexArray *graphics.VertexArray
	SubMeshes   []SubMesh
}

type SubMesh struct {
	IndexOffset    int
	IndexCount     int32
	DiffuseTexture *TwoDTextureHandle
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
	mesh := &Mesh{
		VertexArray: &graphics.VertexArray{},
	}

	in, err := o.locator.Open("assets", "meshes", name)
	if err != nil {
		return nil, fmt.Errorf("failed to open mesh asset %q: %w", name, err)
	}
	defer in.Close()

	meshAsset, err := asset.NewMeshDecoder().Decode(in)
	if err != nil {
		return nil, fmt.Errorf("failed to decode mesh asset %q: %w", name, err)
	}

	gfxTask := o.gfxWorker.Schedule(func() error {
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

func (o *MeshOperator) Release(registry *resource.Registry, resource resource.Resource) error {
	mesh := resource.(*Mesh)

	for _, subMesh := range mesh.SubMeshes {
		if subMesh.DiffuseTexture != nil {
			registry.Dismiss(subMesh.DiffuseTexture.Handle)
		}
	}

	gfxTask := o.gfxWorker.Schedule(func() error {
		return mesh.VertexArray.Release()
	})
	if err := gfxTask.Wait(); err != nil {
		return fmt.Errorf("failed to release gfx vertex array: %w", err)
	}
	return nil
}