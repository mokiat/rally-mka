package stream

import (
	"fmt"

	"github.com/mokiat/rally-mka/cmd/rallymka/internal/graphics"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/resource"
	"github.com/mokiat/rally-mka/internal/data/asset"
)

const cubeTextureResourceType = "cube_texture"

func GetCubeTexture(registry *resource.Registry, name string) *CubeTexture {
	return registry.ResourceType(cubeTextureResourceType).Resource(name).(*CubeTexture)
}

type CubeTexture struct {
	*resource.Handle
	gfxTexture *graphics.CubeTexture
}

func (t *CubeTexture) Gfx() *graphics.CubeTexture {
	return t.gfxTexture
}

func NewCubeTextureController(capacity int, gfxWorker *graphics.Worker) CubeTextureController {
	return CubeTextureController{
		textures:  make([]CubeTexture, capacity),
		gfxWorker: gfxWorker,
	}
}

type CubeTextureController struct {
	textures  []CubeTexture
	gfxWorker *graphics.Worker
}

func (c CubeTextureController) ResourceTypeName() string {
	return cubeTextureResourceType
}

func (c CubeTextureController) Init(index int, handle *resource.Handle) resource.Resource {
	c.textures[index] = CubeTexture{
		Handle:     handle,
		gfxTexture: &graphics.CubeTexture{},
	}
	return &c.textures[index]
}

func (c CubeTextureController) Load(index int, locator resource.Locator, registry *resource.Registry) error {
	texture := &c.textures[index]

	in, err := locator.Open("assets", "textures", "cube", texture.Name())
	if err != nil {
		return fmt.Errorf("failed to open cube texture asset %q: %w", texture.Name(), err)
	}
	defer in.Close()

	texAsset, err := asset.NewCubeTextureDecoder().Decode(in)
	if err != nil {
		return fmt.Errorf("failed to decode cube texture asset %q: %w", texture.Name(), err)
	}

	gfxTask := func() error {
		return texture.gfxTexture.Allocate(graphics.CubeTextureData{
			FrontSideData:  texAsset.Sides[asset.TextureSideFront].Data,
			BackSideData:   texAsset.Sides[asset.TextureSideBack].Data,
			LeftSideData:   texAsset.Sides[asset.TextureSideLeft].Data,
			RightSideData:  texAsset.Sides[asset.TextureSideRight].Data,
			TopSideData:    texAsset.Sides[asset.TextureSideTop].Data,
			BottomSideData: texAsset.Sides[asset.TextureSideBottom].Data,
		})
	}
	if err := c.gfxWorker.Run(gfxTask); err != nil {
		return fmt.Errorf("failed to allocate gfx cube texture: %w", err)
	}
	return nil
}

func (c CubeTextureController) Unload(index int) error {
	texture := &c.textures[index]

	gfxTask := func() error {
		return texture.gfxTexture.Release()
	}
	if err := c.gfxWorker.Run(gfxTask); err != nil {
		return fmt.Errorf("failed to release gfx cube texture: %w", err)
	}
	return nil
}
