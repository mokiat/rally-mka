package view

import (
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
	"github.com/mokiat/lacking/ui/mvc"
	"github.com/mokiat/rally-mka/internal/ui/home"
	"github.com/mokiat/rally-mka/internal/ui/intro"
	"github.com/mokiat/rally-mka/internal/ui/model"
	"github.com/mokiat/rally-mka/internal/ui/play"
)

var Application = co.Define(func(props co.Properties, scope co.Scope) co.Instance {
	var (
		data = co.GetData[*model.Application](props)
	)

	mvc.UseBinding(data, func(ch mvc.Change) bool {
		return mvc.IsChange(ch, model.ChangeActiveView)
	})

	return co.New(mat.Switch, func() {
		co.WithData(mat.SwitchData{
			ChildKey: string(data.ActiveView()),
		})
		co.WithScope(scope)

		co.WithChild(string(model.ViewNameIntro), co.New(intro.View, func() {}))
		co.WithChild(string(model.ViewNameHome), co.New(home.View, func() {}))
		co.WithChild(string(model.ViewNamePlay), co.New(play.View, func() {
			co.WithData(play.ViewData{
				GameData: data.GameData(),
			})
		}))
	})
})
