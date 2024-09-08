package main

import "github.com/mokiat/lacking/game/asset/dsl"

// SUV
var _ = func() any {

	return dsl.CreateModel("SUV",
		dsl.AppendModel(dsl.OpenGLTFModel("resources/models/suv.glb")),
	)
}()
