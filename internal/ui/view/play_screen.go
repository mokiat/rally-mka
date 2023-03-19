package view

import (
	"fmt"

	"github.com/mokiat/gog/opt"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
	"github.com/mokiat/rally-mka/internal/ui/controller"
	"github.com/mokiat/rally-mka/internal/ui/global"
	"github.com/mokiat/rally-mka/internal/ui/model"
)

type PlayScreenData struct {
	Play *model.Play
}

var PlayScreen = co.Define(func(props co.Properties, scope co.Scope) co.Instance {
	var (
		context    = co.GetContext[global.Context]()
		screenData = co.GetData[PlayScreenData](props)
		playModel  = screenData.Play
	)

	controller := co.UseState(func() *controller.PlayController {
		// FIXME: This may actually panic if there is a third party
		// waiting / reading on this and it happens to match the Get call.
		playData, err := playModel.Data().Get()
		if err != nil {
			panic(fmt.Errorf("failed to get data: %w", err))
		}

		return controller.NewPlayController(
			co.Window(scope).Window,
			context.Engine,
			playData,
		)
	}).Get()

	co.Once(func() {
		controller.Start()
	})

	co.Defer(func() {
		controller.Stop()
	})

	return co.New(mat.Element, func() {
		co.WithData(mat.ElementData{
			Essence:   controller,
			Focusable: opt.V(true),
			Focused:   opt.V(true),
			Layout:    ui.NewFillLayout(),
		})
	})
})
