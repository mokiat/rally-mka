package main

import (
	"github.com/mokiat/lacking/data/asset"
	"github.com/mokiat/lacking/data/pack"
)

func main() {
	packer := pack.NewPacker()

	// Programs
	packer.Pipeline(func(p *pack.Pipeline) {
		p.SaveProgramAsset("assets/programs/forward-albedo.dat",
			p.BuildProgram(
				pack.WithVertexShader(p.OpenShaderResource("resources/shaders/forward-albedo.vert")),
				pack.WithFragmentShader(p.OpenShaderResource("resources/shaders/forward-albedo.frag")),
			),
		)

		p.SaveProgramAsset("assets/programs/forward-debug.dat",
			p.BuildProgram(
				pack.WithVertexShader(p.OpenShaderResource("resources/shaders/forward-debug.vert")),
				pack.WithFragmentShader(p.OpenShaderResource("resources/shaders/forward-debug.frag")),
			),
		)

		p.SaveProgramAsset("assets/programs/geometry-skybox.dat",
			p.BuildProgram(
				pack.WithVertexShader(p.OpenShaderResource("resources/shaders/geometry-skybox.vert")),
				pack.WithFragmentShader(p.OpenShaderResource("resources/shaders/geometry-skybox.frag")),
			),
		)

		p.SaveProgramAsset("assets/programs/lighting-pbr.dat",
			p.BuildProgram(
				pack.WithVertexShader(p.OpenShaderResource("resources/shaders/lighting-pbr.vert")),
				pack.WithFragmentShader(p.OpenShaderResource("resources/shaders/lighting-pbr.frag")),
			),
		)
	})

	// TwoD Textures
	packer.Pipeline(func(p *pack.Pipeline) {
		p.SaveTwoDTextureAsset("assets/textures/twod/loading.dat",
			p.OpenImageResource("resources/images/loading.png"),
		)

		p.SaveTwoDTextureAsset("assets/textures/twod/concrete.dat",
			p.OpenImageResource("resources/images/concrete.png"),
		)

		p.SaveTwoDTextureAsset("assets/textures/twod/road.dat",
			p.OpenImageResource("resources/images/road.png"),
		)

		p.SaveTwoDTextureAsset("assets/textures/twod/barrier.dat",
			p.OpenImageResource("resources/images/barrier.png"),
		)

		p.SaveTwoDTextureAsset("assets/textures/twod/grass.dat",
			p.OpenImageResource("resources/images/grass.png"),
		)

		p.SaveTwoDTextureAsset("assets/textures/twod/gravel.dat",
			p.OpenImageResource("resources/images/gravel.png"),
		)

		p.SaveTwoDTextureAsset("assets/textures/twod/rusty_metal_02_diff_512.dat",
			p.OpenImageResource("resources/images/rusty_metal_02_diff_512.png"),
		)

		p.SaveTwoDTextureAsset("assets/textures/twod/body.dat",
			p.OpenImageResource("resources/images/body.png"),
		)

		p.SaveTwoDTextureAsset("assets/textures/twod/wheel.dat",
			p.OpenImageResource("resources/images/wheel.png"),
		)

		p.SaveTwoDTextureAsset("assets/textures/twod/leafy_tree.dat",
			p.OpenImageResource("resources/images/leafy_tree.png"),
		)
	})

	// Cube Textures
	packer.Pipeline(func(p *pack.Pipeline) {
		srcEquirectangularImage := p.OpenImageResource("resources/images/syferfontein.hdr")
		skyboxCubeImage := p.BuildCubeImage(
			pack.WithFrontImage(p.BuildCubeSideFromEquirectangular(pack.CubeSideFront, srcEquirectangularImage)),
			pack.WithRearImage(p.BuildCubeSideFromEquirectangular(pack.CubeSideRear, srcEquirectangularImage)),
			pack.WithLeftImage(p.BuildCubeSideFromEquirectangular(pack.CubeSideLeft, srcEquirectangularImage)),
			pack.WithRightImage(p.BuildCubeSideFromEquirectangular(pack.CubeSideRight, srcEquirectangularImage)),
			pack.WithTopImage(p.BuildCubeSideFromEquirectangular(pack.CubeSideTop, srcEquirectangularImage)),
			pack.WithBottomImage(p.BuildCubeSideFromEquirectangular(pack.CubeSideBottom, srcEquirectangularImage)),
		)

		p.SaveCubeTextureAsset("assets/textures/cube/syferfontein.dat", skyboxCubeImage,
			pack.WithFormat(asset.DataFormatRGBA32F),
		)

		reflectionCubeImage := p.ScaleCubeImage(skyboxCubeImage, 128)
		p.SaveCubeTextureAsset("assets/textures/cube/syferfontein_reflection.dat", reflectionCubeImage,
			pack.WithFormat(asset.DataFormatRGBA32F),
		)

		refractionCubeImage := p.BuildIrradianceCubeImage(reflectionCubeImage,
			pack.WithSampleCount(50),
		)
		p.SaveCubeTextureAsset("assets/textures/cube/syferfontein_refraction.dat", refractionCubeImage,
			pack.WithFormat(asset.DataFormatRGBA32F),
		)
	})

	// Models
	packer.Pipeline(func(p *pack.Pipeline) {
		p.SaveModelAsset("assets/models/quad.dat",
			p.OpenGLTFResource("resources/models/quad.gltf"),
		)

		p.SaveModelAsset("assets/models/street_lamp.dat",
			p.OpenGLTFResource("resources/models/street_lamp.gltf"),
		)

		p.SaveModelAsset("assets/models/suv.dat",
			p.OpenGLTFResource("resources/models/suv.gltf"),
		)

		p.SaveModelAsset("assets/models/leafy_tree.dat",
			p.OpenGLTFResource("resources/models/leafy_tree.gltf"),
		)
	})

	// Levels
	packer.Pipeline(func(p *pack.Pipeline) {
		p.SaveLevelAsset("assets/levels/forest.dat",
			p.OpenLevelResource("resources/levels/forest.json"),
		)

		p.SaveLevelAsset("assets/levels/highway.dat",
			p.OpenLevelResource("resources/levels/highway.json"),
		)

		p.SaveLevelAsset("assets/levels/playground.dat",
			p.OpenLevelResource("resources/levels/playground.json"),
		)
	})

	packer.RunParallel()
}
