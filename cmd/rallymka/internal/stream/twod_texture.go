package stream

import (
	"fmt"

	"github.com/mokiat/rally-mka/cmd/rallymka/internal/graphics"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/resource"
	"github.com/mokiat/rally-mka/internal/data/asset"
)

const twodTextureResourceType = "twod_texture"

func GetTwoDTexture(registry *resource.Registry, name string) *TwoDTexture {
	return registry.ResourceType(twodTextureResourceType).Resource(name).(*TwoDTexture)
}

type TwoDTexture struct {
	*resource.Handle
	gfxTexture *graphics.TwoDTexture
}

func (t *TwoDTexture) Gfx() *graphics.TwoDTexture {
	return t.gfxTexture
}

func NewTwoDTextureController(capacity int, gfxWorker *graphics.Worker) TwoDTextureController {
	return TwoDTextureController{
		textures:  make([]TwoDTexture, capacity),
		gfxWorker: gfxWorker,
	}
}

type TwoDTextureController struct {
	textures  []TwoDTexture
	gfxWorker *graphics.Worker
}

func (c TwoDTextureController) ResourceTypeName() string {
	return twodTextureResourceType
}

func (c TwoDTextureController) Init(index int, handle *resource.Handle) resource.Resource {
	c.textures[index] = TwoDTexture{
		Handle:     handle,
		gfxTexture: &graphics.TwoDTexture{},
	}
	return &c.textures[index]
}

func (c TwoDTextureController) Load(index int, locator resource.Locator, registry *resource.Registry) error {
	texture := &c.textures[index]

	in, err := locator.Open("assets", "textures", "twod", texture.Name())
	if err != nil {
		return fmt.Errorf("failed to open twod texture asset %q: %w", texture.Name(), err)
	}
	defer in.Close()

	texAsset, err := asset.NewTwoDTextureDecoder().Decode(in)
	if err != nil {
		return fmt.Errorf("failed to decode twod texture asset %q: %w", texture.Name(), err)
	}

	gfxTask := func() error {
		return texture.gfxTexture.Allocate(graphics.TwoDTextureData{
			Width:  int32(texAsset.Width),
			Height: int32(texAsset.Height),
			Data:   texAsset.Data,
		})
	}
	if err := c.gfxWorker.Run(gfxTask); err != nil {
		return fmt.Errorf("failed to allocate gfx two dimensional texture: %w", err)
	}
	return nil
}

func (c TwoDTextureController) Unload(index int) error {
	texture := &c.textures[index]

	gfxTask := func() error {
		return texture.gfxTexture.Release()
	}
	if err := c.gfxWorker.Run(gfxTask); err != nil {
		return fmt.Errorf("failed to release gfx two dimensional texture: %w", err)
	}
	return nil
}
