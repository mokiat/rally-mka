package scene

import (
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/stream"
)

type Skybox struct {
	Program stream.ProgramHandle
	Texture stream.CubeTextureHandle
	Mesh    stream.MeshHandle
}
