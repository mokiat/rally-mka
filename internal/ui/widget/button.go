package widget

import (
	"github.com/mokiat/gog/opt"
	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
)

type ButtonData struct {
	Text string
}

var defaultButtonData = ButtonData{
	Text: "",
}

type ButtonCallbackData struct {
	ClickListener mat.ClickListener
}

var defaultButtonCallbackData = ButtonCallbackData{
	ClickListener: func() {},
}

var Button = co.Define(func(props co.Properties, scope co.Scope) co.Instance {
	var (
		data         = co.GetOptionalData(props, defaultButtonData)
		callbackData = co.GetOptionalCallbackData(props, defaultButtonCallbackData)
	)

	essence := co.UseState(func() *homeButtonEssence {
		return &homeButtonEssence{
			ButtonBaseEssence: mat.NewButtonBaseEssence(callbackData.ClickListener),
		}
	}).Get()
	essence.SetOnClick(callbackData.ClickListener)

	essence.font = co.OpenFont(scope, "mat:///roboto-bold.ttf")
	essence.fontSize = 26
	essence.text = data.Text

	padding := ui.Spacing{
		Left:   5,
		Right:  5,
		Top:    2,
		Bottom: 2,
	}

	txtSize := essence.font.TextSize(essence.text, essence.fontSize)

	return co.New(mat.Element, func() {
		co.WithData(mat.ElementData{
			Essence: essence,
			Padding: padding,
			IdealSize: opt.V(
				ui.NewSize(int(txtSize.X), int(txtSize.Y)).Grow(padding.Size()),
			),
		})
		co.WithLayoutData(props.LayoutData())
		co.WithChildren(props.Children())
	})
})

var _ ui.ElementRenderHandler = (*homeButtonEssence)(nil)

type homeButtonEssence struct {
	*mat.ButtonBaseEssence
	font     *ui.Font
	fontSize float32
	text     string
}

func (e *homeButtonEssence) OnRender(element *ui.Element, canvas *ui.Canvas) {
	var fontColor ui.Color

	switch e.State() {
	case mat.ButtonStateOver:
		fontColor = ui.RGB(0x00, 0xB2, 0x08)
	case mat.ButtonStateDown:
		fontColor = ui.RGB(0x00, 0x33, 0x00)
	default:
		fontColor = ui.White()
	}

	contentArea := element.ContentBounds()
	textPosition := contentArea.Position
	canvas.Reset()
	canvas.FillText(e.text, sprec.NewVec2(
		float32(textPosition.X),
		float32(textPosition.Y),
	), ui.Typography{
		Font:  e.font,
		Size:  e.fontSize,
		Color: fontColor,
	})
}
