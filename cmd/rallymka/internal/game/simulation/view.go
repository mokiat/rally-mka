package simulation

import "github.com/mokiat/rally-mka/internal/engine/graphics"

type View struct{}

func (v *View) Resize(width, height int) {
}

func (v *View) Update(elapsedSeconds float32) {
}

func (v *View) Render(pipeline *graphics.Pipeline) {
}
