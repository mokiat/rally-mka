package main

import (
	"github.com/mokiat/lacking/game/asset/dsl"
	"github.com/mokiat/lacking/game/asset/mdl"
)

var _ = func() any {
	model := dsl.OpenGLTFModel("resources/models/tiles.glb")

	// TODO: Find from inside the model.
	grassTexture := dsl.Create2DTexture(
		dsl.OpenImage("resources/images/grass.png"),
		dsl.SetMipmapping(dsl.Const(true)),
	)

	// TODO: Find from inside the model.
	dirtTexture := dsl.Create2DTexture(
		dsl.OpenImage("resources/images/dirt.png"),
		dsl.SetMipmapping(dsl.Const(true)),
	)

	noiseTexture := dsl.Create2DTexture(
		dsl.OpenImage("resources/images/perlin.png"),
		dsl.SetMipmapping(dsl.Const(true)),
	)

	tilingShader := `
		texture terrain sampler2D
		texture noise sampler2D

		func #fragment() {
			var noiseScale vec2
			noiseScale.x = 0.1
			noiseScale.y = 0.1
			var noiseSample vec4 = sample(noise, #varyingUV * noiseScale)

			var cs float = 0.8660 // cos 30
			var sn float = 0.5 // sin 30
			var uv2 vec2
			uv2.x = #varyingUV.x * cs - #varyingUV.y * sn + 0.5
			uv2.y = #varyingUV.x * sn + #varyingUV.y * cs + 0.5

			var color1 vec4 = sample(terrain, #varyingUV)
			var color2 vec4 = sample(terrain, uv2)

			#color = mix(color1, color2, smoothstep(0.0, 1.0, noiseSample.r + 0.3))
		}
	`

	return dsl.Save("Tiles.dat", dsl.CreateModel(
		dsl.AppendModel(model),
		dsl.EditMaterial("Grass",
			dsl.Clear(),
			dsl.AddGeometryPass(dsl.CreateMaterialPass(
				dsl.SetShader(dsl.CreateShader(mdl.ShaderTypeGeometry, tilingShader)),
				dsl.SetCulling(dsl.Const(mdl.CullModeBack)),
			)),
			dsl.BindSampler("terrain", dsl.CreateSampler(grassTexture,
				dsl.SetWrapMode(dsl.Const(mdl.WrapModeRepeat)),
				dsl.SetFilterMode(dsl.Const(mdl.FilterModeLinear)),
				dsl.SetMipmapping(dsl.Const(true)),
			)),
			dsl.BindSampler("noise", dsl.CreateSampler(noiseTexture,
				dsl.SetWrapMode(dsl.Const(mdl.WrapModeRepeat)),
				dsl.SetFilterMode(dsl.Const(mdl.FilterModeNearest)),
				dsl.SetMipmapping(dsl.Const(true)),
			)),
		),
		dsl.EditMaterial("Dirt",
			dsl.Clear(),
			dsl.AddGeometryPass(dsl.CreateMaterialPass(
				dsl.SetShader(dsl.CreateShader(mdl.ShaderTypeGeometry, tilingShader)),
				dsl.SetCulling(dsl.Const(mdl.CullModeBack)),
			)),
			dsl.BindSampler("terrain", dsl.CreateSampler(dirtTexture,
				dsl.SetWrapMode(dsl.Const(mdl.WrapModeRepeat)),
				dsl.SetFilterMode(dsl.Const(mdl.FilterModeLinear)),
				dsl.SetMipmapping(dsl.Const(true)),
			)),
			dsl.BindSampler("noise", dsl.CreateSampler(noiseTexture,
				dsl.SetWrapMode(dsl.Const(mdl.WrapModeRepeat)),
				dsl.SetFilterMode(dsl.Const(mdl.FilterModeNearest)),
				dsl.SetMipmapping(dsl.Const(true)),
			)),
		),
	))
}()
