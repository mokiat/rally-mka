package stream

import (
	"fmt"

	"github.com/mokiat/lacking/data/asset"
	"github.com/mokiat/lacking/graphics"
	"github.com/mokiat/rally-mka/internal/engine/resource"
)

const cubeTextureResourceType = "cube_texture"

func GetCubeTexture(registry *resource.Registry, name string) CubeTextureHandle {
	return CubeTextureHandle{
		Handle: registry.Type(cubeTextureResourceType).Resource(name),
	}
}

type CubeTextureHandle struct {
	resource.Handle
}

func (h CubeTextureHandle) Get() *graphics.CubeTexture {
	return h.Handle.Get().(*graphics.CubeTexture)
}

func NewCubeTextureOperator(locator resource.Locator, gfxWorker *graphics.Worker) *CubeTextureOperator {
	return &CubeTextureOperator{
		locator:   locator,
		gfxWorker: gfxWorker,
	}
}

type CubeTextureOperator struct {
	locator   resource.Locator
	gfxWorker *graphics.Worker
}

func (o *CubeTextureOperator) Register(registry *resource.Registry) {
	registry.RegisterType(cubeTextureResourceType, o)
}

func (o *CubeTextureOperator) Allocate(registry *resource.Registry, name string) (resource.Resource, error) {
	in, err := o.locator.Open("assets", "textures", "cube", name)
	if err != nil {
		return nil, fmt.Errorf("failed to open cube texture asset %q: %w", name, err)
	}
	defer in.Close()

	texAsset := new(asset.CubeTexture)
	if err := asset.DecodeCubeTexture(in, texAsset); err != nil {
		return nil, fmt.Errorf("failed to decode cube texture asset %q: %w", name, err)
	}

	texture := &graphics.CubeTexture{}

	gfxTask := o.gfxWorker.Schedule(func() error {
		return texture.Allocate(graphics.CubeTextureData{
			Dimension:      int32(texAsset.Dimension),
			FrontSideData:  texAsset.Sides[asset.TextureSideFront].Data,
			BackSideData:   texAsset.Sides[asset.TextureSideBack].Data,
			LeftSideData:   texAsset.Sides[asset.TextureSideLeft].Data,
			RightSideData:  texAsset.Sides[asset.TextureSideRight].Data,
			TopSideData:    texAsset.Sides[asset.TextureSideTop].Data,
			BottomSideData: texAsset.Sides[asset.TextureSideBottom].Data,
		})
	})
	if err := gfxTask.Wait(); err != nil {
		return nil, fmt.Errorf("failed to allocate gfx cube texture: %w", err)
	}
	return texture, nil
}

func (o *CubeTextureOperator) Release(registry *resource.Registry, resource resource.Resource) error {
	texture := resource.(*graphics.CubeTexture)

	gfxTask := o.gfxWorker.Schedule(func() error {
		return texture.Release()
	})
	if err := gfxTask.Wait(); err != nil {
		return fmt.Errorf("failed to release gfx cube texture: %w", err)
	}
	return nil
}
