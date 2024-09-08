package main

import (
	"github.com/mokiat/lacking/game/asset/dsl"
	"github.com/mokiat/lacking/game/asset/mdl"
)

// Day Scene
var _ = func() any {

	// Using GIMP:
	// 1. Scale height to %50
	// 2. Convert to float32
	// 3. Exposure: ~ -6
	skyImage := dsl.CubeImageFromEquirectangular(
		dsl.OpenImage("resources/images/day.hdr"),
	)

	smallerSkyImage := dsl.ResizedCubeImage(skyImage, dsl.Const(512))
	skyTexture := dsl.CreateCubeTexture(smallerSkyImage)

	reflectionCubeImage := dsl.ResizedCubeImage(skyImage, dsl.Const(128))
	reflectionTexture := dsl.CreateCubeTexture(reflectionCubeImage)

	refractionCubeImage := dsl.IrradianceCubeImage(reflectionCubeImage, dsl.SetSampleCount(dsl.Const(50)))
	refractionTexture := dsl.CreateCubeTexture(refractionCubeImage)

	skyMaterial := dsl.CreateTextureSkyMaterial(
		dsl.CreateSampler(skyTexture,
			dsl.SetWrapMode(dsl.Const(mdl.WrapModeClamp)),
			dsl.SetFilterMode(dsl.Const(mdl.FilterModeLinear)),
			dsl.SetMipmapping(dsl.Const(false)),
		),
	)

	sky := dsl.CreateSky(skyMaterial)

	ambientLight := dsl.CreateAmbientLight(
		dsl.SetReflectionTexture(reflectionTexture),
		dsl.SetRefractionTexture(refractionTexture),
	)

	// TODO: Reference Forest Level so that only this model
	// needs to be loaded.

	return dsl.CreateModel("Forest-Day",
		dsl.AddNode(dsl.CreateNode("sky",
			dsl.SetTarget(sky),
		)),
		dsl.AddNode(dsl.CreateNode("AmbientLight",
			dsl.SetTarget(ambientLight),
		)),
	)
}()

// Night Scene
var _ = func() any {

	// Using GIMP:
	// 1. Scale height to %50
	// 2. Convert to float32
	// 3. Exposure: ~ -6
	skyImage := dsl.CubeImageFromEquirectangular(
		dsl.OpenImage("resources/images/night.exr"),
	)

	smallerSkyImage := dsl.ResizedCubeImage(skyImage, dsl.Const(512))
	skyTexture := dsl.CreateCubeTexture(smallerSkyImage)

	reflectionCubeImage := dsl.ResizedCubeImage(skyImage, dsl.Const(128))
	reflectionTexture := dsl.CreateCubeTexture(reflectionCubeImage)

	refractionCubeImage := dsl.IrradianceCubeImage(reflectionCubeImage, dsl.SetSampleCount(dsl.Const(50)))
	refractionTexture := dsl.CreateCubeTexture(refractionCubeImage)

	skyMaterial := dsl.CreateTextureSkyMaterial(
		dsl.CreateSampler(skyTexture,
			dsl.SetWrapMode(dsl.Const(mdl.WrapModeClamp)),
			dsl.SetFilterMode(dsl.Const(mdl.FilterModeLinear)),
			dsl.SetMipmapping(dsl.Const(false)),
		),
	)

	sky := dsl.CreateSky(skyMaterial)

	ambientLight := dsl.CreateAmbientLight(
		dsl.SetReflectionTexture(reflectionTexture),
		dsl.SetRefractionTexture(refractionTexture),
	)

	// TODO: Reference Forest Level so that only this model
	// needs to be loaded.

	return dsl.CreateModel("Forest-Night",
		dsl.AddNode(dsl.CreateNode("Sky",
			dsl.SetTarget(sky),
		)),
		dsl.AddNode(dsl.CreateNode("AmbientLight",
			dsl.SetTarget(ambientLight),
		)),
	)
}()

// Forest Level
var _ = func() any {

	return dsl.CreateModel("Forest",
		dsl.AppendModel(dsl.OpenGLTFModel("resources/models/forest.glb")),
	)
}()
