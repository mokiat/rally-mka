package main

import "github.com/mokiat/lacking/data/pack"

func main() {
	packer := pack.NewPacker()
	packer.SetAssetsDir("assets")
	packer.SetResourcesDir("resources")

	// Programs
	packer.Store(
		packer.ProgramAssetFile("forward-albedo.dat").
			WithVertexShader(packer.ShaderResourceFile("forward-albedo.vert")).
			WithFragmentShader(packer.ShaderResourceFile("forward-albedo.frag")),

		packer.ProgramAssetFile("forward-debug.dat").
			WithVertexShader(packer.ShaderResourceFile("forward-debug.vert")).
			WithFragmentShader(packer.ShaderResourceFile("forward-debug.frag")),

		packer.ProgramAssetFile("geometry-skybox.dat").
			WithVertexShader(packer.ShaderResourceFile("geometry-skybox.vert")).
			WithFragmentShader(packer.ShaderResourceFile("geometry-skybox.frag")),

		packer.ProgramAssetFile("lighting-pbr.dat").
			WithVertexShader(packer.ShaderResourceFile("lighting-pbr.vert")).
			WithFragmentShader(packer.ShaderResourceFile("lighting-pbr.frag")),
	)

	// 2D Textures
	packer.Store(
		packer.TwoDTextureAssetFile("loading.dat").
			WithImage(packer.ImageResourceFile("loading.png")),

		packer.TwoDTextureAssetFile("tree.dat").
			WithImage(packer.ImageResourceFile("tree.png")),

		packer.TwoDTextureAssetFile("lamp.dat").
			WithImage(packer.ImageResourceFile("lamp.png")),

		packer.TwoDTextureAssetFile("finish.dat").
			WithImage(packer.ImageResourceFile("finish.png")),

		packer.TwoDTextureAssetFile("hatch_body.dat").
			WithImage(packer.ImageResourceFile("hatch_body.png")),

		packer.TwoDTextureAssetFile("hatch_wheel.dat").
			WithImage(packer.ImageResourceFile("hatch_wheel.png")),

		packer.TwoDTextureAssetFile("suv_body.dat").
			WithImage(packer.ImageResourceFile("suv_body.png")),

		packer.TwoDTextureAssetFile("suv_wheel.dat").
			WithImage(packer.ImageResourceFile("suv_wheel.png")),

		packer.TwoDTextureAssetFile("truck_body.dat").
			WithImage(packer.ImageResourceFile("truck_body.png")),

		packer.TwoDTextureAssetFile("truck_wheel.dat").
			WithImage(packer.ImageResourceFile("truck_wheel.png")),

		packer.TwoDTextureAssetFile("concrete.dat").
			WithImage(packer.ImageResourceFile("concrete.png")),

		packer.TwoDTextureAssetFile("road.dat").
			WithImage(packer.ImageResourceFile("road.png")),

		packer.TwoDTextureAssetFile("barrier.dat").
			WithImage(packer.ImageResourceFile("barrier.png")),

		packer.TwoDTextureAssetFile("grass.dat").
			WithImage(packer.ImageResourceFile("grass.png")),

		packer.TwoDTextureAssetFile("gravel.dat").
			WithImage(packer.ImageResourceFile("gravel.png")),

		packer.TwoDTextureAssetFile("asphalt.dat").
			WithImage(packer.ImageResourceFile("asphalt.png")),
	)

	// Cube Textures
	packer.Store(
		packer.CubeTextureAssetFile("city.dat").
			WithFrontImage(packer.ImageResourceFile("city_front.png")).
			WithBackImage(packer.ImageResourceFile("city_back.png")).
			WithLeftImage(packer.ImageResourceFile("city_left.png")).
			WithRightImage(packer.ImageResourceFile("city_right.png")).
			WithTopImage(packer.ImageResourceFile("city_top.png")).
			WithBottomImage(packer.ImageResourceFile("city_bottom.png")).
			WithDimension(512),
	)
}
