package main

import (
	"github.com/mokiat/lacking/game/asset/dsl"
)

// Vehicle
var _ = func() any {
	return dsl.Save("Vehicle.dat", dsl.CreateModel(
		dsl.AppendModel(dsl.OpenGLTFModel("resources/models/vehicle.glb")),
	))
}()
