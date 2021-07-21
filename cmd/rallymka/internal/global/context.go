package global

import (
	"github.com/mokiat/lacking/game/graphics"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/game"
)

type Context struct {
	GFXEngine      graphics.Engine
	GameController *game.Controller
}
