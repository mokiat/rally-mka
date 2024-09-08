package main

import "github.com/mokiat/lacking/game/asset/dsl"

// Home Screen Sky
var _ = func() any {

	daySky := dsl.CreateSky(dsl.CreateColorSkyMaterial(
		dsl.RGB(20.0, 25.0, 30.0),
	))

	nightSky := dsl.CreateSky(dsl.CreateColorSkyMaterial(
		dsl.RGB(0.01, 0.01, 0.01),
	))

	// TODO: Append home screen level so that only this model
	// needs to be loaded.

	return dsl.CreateModel("Home-Screen",
		dsl.AddNode(dsl.CreateNode("Sky-Day",
			dsl.SetTarget(daySky),
		)),
		dsl.AddNode(dsl.CreateNode("Sky-Night",
			dsl.SetTarget(nightSky),
		)),
	)
}()

// Home Screen Level
var _ = func() any {

	return dsl.CreateModel("HomeScreen",
		dsl.AppendModel(dsl.OpenGLTFModel("resources/models/home-screen.glb")),
	)
}()
