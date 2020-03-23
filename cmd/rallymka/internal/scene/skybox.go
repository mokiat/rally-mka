package scene

import (
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/stream"
	"github.com/mokiat/rally-mka/internal/engine/graphics"
)

type Skybox struct {
	Program *graphics.Program
	Texture *graphics.CubeTexture
	Mesh    *stream.Mesh
}
