package scene

import (
	"github.com/mokiat/go-whiskey/math"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/stream"
)

type Entity struct {
	Model  *stream.Model
	Matrix math.Mat4x4
}
