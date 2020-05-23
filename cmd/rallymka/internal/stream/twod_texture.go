package stream

import (
	"fmt"

	"github.com/mokiat/lacking/data/asset"
	"github.com/mokiat/lacking/graphics"
	"github.com/mokiat/rally-mka/internal/engine/resource"
)

const twodTextureResourceType = "twod_texture"

func GetTwoDTexture(registry *resource.Registry, name string) TwoDTextureHandle {
	return TwoDTextureHandle{
		Handle: registry.Type(twodTextureResourceType).Resource(name),
	}
}

type TwoDTextureHandle struct {
	resource.Handle
}

func (h TwoDTextureHandle) Get() *graphics.TwoDTexture {
	return h.Handle.Get().(*graphics.TwoDTexture)
}

func NewTwoDTextureOperator(locator resource.Locator, gfxWorker *graphics.Worker) *TwoDTextureOperator {
	return &TwoDTextureOperator{
		locator:   locator,
		gfxWorker: gfxWorker,
	}
}

type TwoDTextureOperator struct {
	locator   resource.Locator
	gfxWorker *graphics.Worker
}

func (o *TwoDTextureOperator) Register(registry *resource.Registry) {
	registry.RegisterType(twodTextureResourceType, o)
}

func (o *TwoDTextureOperator) Allocate(registry *resource.Registry, name string) (resource.Resource, error) {
	in, err := o.locator.Open("assets", "textures", "twod", name)
	if err != nil {
		return nil, fmt.Errorf("failed to open twod texture asset %q: %w", name, err)
	}
	defer in.Close()

	texAsset := new(asset.TwoDTexture)
	if err := asset.DecodeTwoDTexture(in, texAsset); err != nil {
		return nil, fmt.Errorf("failed to decode twod texture asset %q: %w", name, err)
	}

	texture := &graphics.TwoDTexture{}

	gfxTask := o.gfxWorker.Schedule(func() error {
		return texture.Allocate(graphics.TwoDTextureData{
			Width:  int32(texAsset.Width),
			Height: int32(texAsset.Height),
			Data:   texAsset.Data,
		})
	})
	if err := gfxTask.Wait(); err != nil {
		return nil, fmt.Errorf("failed to allocate two dimensional gfx texture: %w", err)
	}
	return texture, nil
}

func (o *TwoDTextureOperator) Release(registry *resource.Registry, resource resource.Resource) error {
	texture := resource.(*graphics.TwoDTexture)

	gfxTask := o.gfxWorker.Schedule(func() error {
		return texture.Release()
	})
	if err := gfxTask.Wait(); err != nil {
		return fmt.Errorf("failed to release two dimensional gfx texture: %w", err)
	}
	return nil
}
