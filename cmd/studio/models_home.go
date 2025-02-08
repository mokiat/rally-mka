package main

import (
	"github.com/mokiat/gomath/dprec"
	"github.com/mokiat/lacking/game/asset/dsl"
	"github.com/mokiat/lacking/game/asset/mdl"
)

// Home Screen Sky
var _ = func() any {
	daySky := dsl.CreateSky(dsl.CreateColorSkyMaterial(
		dsl.RGB(2.0, 2.5, 3.0),
	))

	daySkyImage := dsl.CubeImageFromEquirectangular(
		dsl.OpenImage("resources/images/skybox-day.exr"),
	)
	daySkyImageSmall := dsl.ResizedCubeImage(daySkyImage, dsl.Const(128))

	dayReflectionCubeImages := dsl.ReflectionCubeImages(daySkyImageSmall, dsl.SetSampleCount(dsl.Const(120)))
	dayReflectionTexture := dsl.CreateCubeMipmapTexture(dayReflectionCubeImages,
		dsl.SetMipmapping(dsl.Const(true)),
	)

	dayRefractionCubeImage := dsl.IrradianceCubeImage(daySkyImageSmall, dsl.SetSampleCount(dsl.Const(50)))
	dayRefractionTexture := dsl.CreateCubeTexture(dayRefractionCubeImage)

	dayAmbientLight := dsl.CreateAmbientLight(
		dsl.SetReflectionTexture(dayReflectionTexture),
		dsl.SetRefractionTexture(dayRefractionTexture),
	)

	dayDirectionalLight := dsl.CreateDirectionalLight(
		dsl.SetEmitColor(dsl.RGB(2.5, 2.5, 2.3)),
		dsl.SetCastShadow(dsl.Const(true)),
	)

	nightSky := dsl.CreateSky(dsl.CreateColorSkyMaterial(
		dsl.RGB(0.01, 0.01, 0.01),
	))

	nightSkyImage := dsl.CubeImageFromEquirectangular(
		dsl.OpenImage("resources/images/skybox-night.exr"),
	)
	nightSkyImageSmall := dsl.ResizedCubeImage(nightSkyImage, dsl.Const(128))

	nightReflectionCubeImages := dsl.ReflectionCubeImages(nightSkyImageSmall, dsl.SetSampleCount(dsl.Const(120)))
	nightReflectionTexture := dsl.CreateCubeMipmapTexture(nightReflectionCubeImages,
		dsl.SetMipmapping(dsl.Const(true)),
	)

	nightRefractionCubeImage := dsl.IrradianceCubeImage(nightSkyImageSmall, dsl.SetSampleCount(dsl.Const(50)))
	nightRefractionTexture := dsl.CreateCubeTexture(nightRefractionCubeImage)

	nightAmbientLight := dsl.CreateAmbientLight(
		dsl.SetReflectionTexture(nightReflectionTexture),
		dsl.SetRefractionTexture(nightRefractionTexture),
	)

	nightDirectionalLight := dsl.CreateDirectionalLight(
		dsl.SetEmitColor(dsl.RGB(0.05, 0.05, 0.05)),
		dsl.SetCastShadow(dsl.Const(true)),
	)

	// TODO: Find from inside the model.
	waterTexture := dsl.Create2DTexture(
		dsl.OpenImage("resources/images/water.png"),
		dsl.SetMipmapping(dsl.Const(true)),
	)

	waterfallShader := `
		textures {
			water sampler2D,
		}
		func #fragment() {
			var uv vec2 = #uv
			uv.y += #time * 1.5 + sin(#uv.x * #uv.y) * 0.1

			var color vec4 = sample(water, uv)
			color *= 0.6
			var alpha float = 1.0 - #uv.y + sin(#uv.x * sin(#uv.y * 100.0) * 100.0)
			if alpha < 0.7 {
				discard
			}
			#color = color
		}
	`

	return dsl.CreateModel("HomeScreen",
		dsl.AppendModel(dsl.OpenGLTFModel("resources/models/home.glb")),
		dsl.EditMaterial("Waterfall",
			dsl.Clear(),
			dsl.AddGeometryPass(dsl.CreateMaterialPass(
				dsl.SetShader(dsl.CreateShader(mdl.ShaderTypeGeometry, waterfallShader)),
				dsl.SetCulling(dsl.Const(mdl.CullModeBack)),
			)),
			dsl.BindSampler("water", dsl.CreateSampler(waterTexture,
				dsl.SetWrapMode(dsl.Const(mdl.WrapModeRepeat)),
				dsl.SetFilterMode(dsl.Const(mdl.FilterModeLinear)),
				dsl.SetMipmapping(dsl.Const(true)),
			)),
		),
		dsl.AddNode(dsl.CreateNode("Day-Sky",
			dsl.SetTarget(daySky),
			dsl.AddNode(dsl.CreateNode("Day-AmbientLight",
				dsl.SetTarget(dayAmbientLight),
			)),
			dsl.AddNode(dsl.CreateNode("Day-DirectionalLight",
				dsl.SetTarget(dayDirectionalLight),
				dsl.SetRotation(dsl.Const(dprec.QuatProd(
					dprec.RotationQuat(dprec.Degrees(-140), dprec.BasisYVec3()),
					dprec.RotationQuat(dprec.Degrees(-45), dprec.BasisXVec3()),
				))),
			)),
		)),
		dsl.AddNode(dsl.CreateNode("Night-Sky",
			dsl.SetTarget(nightSky),
			dsl.AddNode(dsl.CreateNode("Night-AmbientLight",
				dsl.SetTarget(nightAmbientLight),
			)),
			dsl.AddNode(dsl.CreateNode("Night-DirectionalLight",
				dsl.SetTarget(nightDirectionalLight),
				dsl.SetRotation(dsl.Const(dprec.QuatProd(
					dprec.RotationQuat(dprec.Degrees(-140), dprec.BasisYVec3()),
					dprec.RotationQuat(dprec.Degrees(-45), dprec.BasisXVec3()),
				))),
			)),
		)),
	)
}()
