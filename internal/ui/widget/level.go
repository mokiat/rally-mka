package widget

import (
	"math"
	"time"

	"github.com/mokiat/gog/opt"
	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/std"
	"github.com/mokiat/rally-mka/internal/game/level"
)

var Level = co.Define(&levelComponent{})

type LevelData struct {
	Board *level.Board
}

type levelComponent struct {
	co.BaseComponent

	images      []*ui.Image
	board       *level.Board
	elapsedTime time.Duration
}

func (c *levelComponent) OnCreate() {
	data := co.GetData[LevelData](c.Properties())
	c.board = data.Board

	c.images = []*ui.Image{
		nil,
		co.OpenImage(c.Scope(), "ui/images/tile-grass.png"),
		co.OpenImage(c.Scope(), "ui/images/tile-road-straight.png"),
		co.OpenImage(c.Scope(), "ui/images/tile-road-corner-smooth.png"),
		co.OpenImage(c.Scope(), "ui/images/tile-road-corner-sharp.png"),
		co.OpenImage(c.Scope(), "ui/images/tile-road-split.png"),
	}
	c.elapsedTime = 0
}

func (c *levelComponent) Render() co.Instance {
	return co.New(std.Element, func() {
		co.WithData(std.ElementData{
			Essence:   c,
			IdealSize: opt.V(ui.NewSize(600, 500)),
		})
		co.WithLayoutData(c.Properties().LayoutData())
	})
}

func (c *levelComponent) OnRender(element *ui.Element, canvas *ui.Canvas) {
	c.elapsedTime += canvas.ElapsedTime()

	drawBounds := canvas.DrawBounds(element, false)

	canvas.Translate(drawBounds.Position)

	board := c.board
	centerPosition := tilePosition(board.Center())
	appearAfter := c.elapsedTime
	appearDuration := 10 * time.Millisecond
	for y := range board.Size() {
		for x := range board.Size() {
			if appearAfter > appearDuration {
				tileCoord := level.C(x, y)
				tile := board.Tile(tileCoord)
				image := c.images[tile.Shape]
				tilePosition := sprec.Vec2Sum(
					sprec.Vec2Diff(tilePosition(tileCoord), centerPosition),
					sprec.NewVec2(290.0, 250.0),
				)
				canvas.Push()
				canvas.Translate(tilePosition)
				canvas.Rotate(-sprec.Degrees(60.0 * float32(tile.Rotation)))
				canvas.Translate(sprec.NewVec2(-32.0, -32.0))
				canvas.Reset()
				canvas.Rectangle(sprec.ZeroVec2(), sprec.NewVec2(64.0, 64.0))
				canvas.Fill(ui.Fill{
					Rule:        ui.FillRuleSimple,
					Color:       ui.White(),
					Image:       image,
					ImageOffset: sprec.ZeroVec2(),
					ImageSize:   sprec.NewVec2(64.0, 64.0),
				})
				canvas.Pop()
			}
			appearAfter -= appearDuration
		}
	}

	if appearAfter < appearDuration {
		element.Invalidate() // force redraw
	}
}

func tilePosition(coord level.Coord) sprec.Vec2 {
	const tileSize = 64.0
	x, y := coord.X, coord.Y
	xShift := tileSize * float64(math.Sqrt(3)) / 2.0
	yShift := tileSize * 3.0 / 4.0
	xOffset := float64(0.0)
	if max(y, -y)%2 == 1 {
		xOffset = xShift / 2.0
	}
	return sprec.Vec2{
		X: float32(x)*float32(xShift) + float32(xOffset),
		Y: float32(y) * float32(yShift),
	}
}
