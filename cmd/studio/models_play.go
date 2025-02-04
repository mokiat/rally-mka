package main

import (
	"github.com/mokiat/gomath/dprec"
	"github.com/mokiat/lacking/game/asset/dsl"
	"github.com/mokiat/lacking/game/asset/mdl"
)

// Day Scene
var _ = func() any {

	// Transformation using GIMP:
	// 1. Scale height to %50
	// 2. Convert to float32/linear (Image -> Precision)
	// 3. Exposure to 1.3 (Color -> Exposure)
	skyImage := dsl.CubeImageFromEquirectangular(
		dsl.OpenImage("resources/images/skybox-day.exr"),
	)
	skyImageSmall := dsl.ResizedCubeImage(skyImage, dsl.Const(128))

	smallerSkyImage := dsl.ResizedCubeImage(skyImage, dsl.Const(512))
	skyTexture := dsl.CreateCubeTexture(smallerSkyImage)

	reflectionCubeImages := dsl.ReflectionCubeImages(skyImageSmall, dsl.SetSampleCount(dsl.Const(120)))
	reflectionTexture := dsl.CreateCubeMipmapTexture(reflectionCubeImages,
		dsl.SetMipmapping(dsl.Const(true)),
	)

	refractionCubeImage := dsl.IrradianceCubeImage(skyImageSmall, dsl.SetSampleCount(dsl.Const(50)))
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

	directionalLight := dsl.CreateDirectionalLight(
		dsl.SetEmitColor(dsl.RGB(2.5, 2.5, 2.3)),
		dsl.SetCastShadow(dsl.Const(true)),
	)

	return dsl.CreateModel("PlayScreen-Day",
		dsl.AddNode(dsl.CreateNode("sky",
			dsl.SetTarget(sky),
		)),
		dsl.AddNode(dsl.CreateNode("AmbientLight",
			dsl.SetTarget(ambientLight),
		)),
		dsl.AddNode(dsl.CreateNode("DirectionalLight",
			dsl.SetTarget(directionalLight),
			dsl.SetRotation(dsl.Const(dprec.QuatProd(
				dprec.RotationQuat(dprec.Degrees(-140), dprec.BasisYVec3()),
				dprec.RotationQuat(dprec.Degrees(-45), dprec.BasisXVec3()),
			))),
		)),
	)
}()

// Night Scene
var _ = func() any {

	// Transformation using GIMP:
	// 1. Scale height to %50
	// 2. Convert to float32/linear (Image -> Precision)
	// 3. Exposure to -4 (Color -> Exposure)
	skyImage := dsl.CubeImageFromEquirectangular(
		dsl.OpenImage("resources/images/skybox-night.exr"),
	)
	skyImageSmall := dsl.ResizedCubeImage(skyImage, dsl.Const(128))

	smallerSkyImage := dsl.ResizedCubeImage(skyImage, dsl.Const(512))
	skyTexture := dsl.CreateCubeTexture(smallerSkyImage)

	reflectionCubeImages := dsl.ReflectionCubeImages(skyImageSmall, dsl.SetSampleCount(dsl.Const(120)))
	reflectionTexture := dsl.CreateCubeMipmapTexture(reflectionCubeImages,
		dsl.SetMipmapping(dsl.Const(true)),
	)

	refractionCubeImage := dsl.IrradianceCubeImage(skyImageSmall, dsl.SetSampleCount(dsl.Const(50)))
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

	return dsl.CreateModel("PlayScreen-Night",
		dsl.AddNode(dsl.CreateNode("Sky",
			dsl.SetTarget(sky),
		)),
		dsl.AddNode(dsl.CreateNode("AmbientLight",
			dsl.SetTarget(ambientLight),
		)),
	)
}()
