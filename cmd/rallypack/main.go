package main

import (
	"log"

	"github.com/mokiat/lacking/data/asset"
	"github.com/mokiat/lacking/data/pack"
)

const workerCount = 4

func main() {
	packer := pack.NewPacker()

	packer.Schedule(pack.SaveProgramAsset("assets/programs/forward-albedo.dat",
		pack.BuildProgram().
			WithVertexShader(pack.OpenShaderResource("resources/shaders/forward-albedo.vert")).
			WithFragmentShader(pack.OpenShaderResource("resources/shaders/forward-albedo.frag")),
	))

	packer.Schedule(pack.SaveProgramAsset("assets/programs/forward-debug.dat",
		pack.BuildProgram().
			WithVertexShader(pack.OpenShaderResource("resources/shaders/forward-debug.vert")).
			WithFragmentShader(pack.OpenShaderResource("resources/shaders/forward-debug.frag")),
	))

	packer.Schedule(pack.SaveProgramAsset("assets/programs/geometry-skybox.dat",
		pack.BuildProgram().
			WithVertexShader(pack.OpenShaderResource("resources/shaders/geometry-skybox.vert")).
			WithFragmentShader(pack.OpenShaderResource("resources/shaders/geometry-skybox.frag")),
	))

	packer.Schedule(pack.SaveProgramAsset("assets/programs/lighting-pbr.dat",
		pack.BuildProgram().
			WithVertexShader(pack.OpenShaderResource("resources/shaders/lighting-pbr.vert")).
			WithFragmentShader(pack.OpenShaderResource("resources/shaders/lighting-pbr.frag")),
	))

	packer.Schedule(pack.SaveTwoDTextureAsset("assets/textures/twod/loading.dat",
		pack.OpenImageResource("resources/images/loading.png"),
	))

	packer.Schedule(pack.SaveTwoDTextureAsset("assets/textures/twod/concrete.dat",
		pack.OpenImageResource("resources/images/concrete.png"),
	))

	packer.Schedule(pack.SaveTwoDTextureAsset("assets/textures/twod/road.dat",
		pack.OpenImageResource("resources/images/road.png"),
	))

	packer.Schedule(pack.SaveTwoDTextureAsset("assets/textures/twod/barrier.dat",
		pack.OpenImageResource("resources/images/barrier.png"),
	))

	packer.Schedule(pack.SaveTwoDTextureAsset("assets/textures/twod/grass.dat",
		pack.OpenImageResource("resources/images/grass.png"),
	))

	packer.Schedule(pack.SaveTwoDTextureAsset("assets/textures/twod/gravel.dat",
		pack.OpenImageResource("resources/images/gravel.png"),
	))

	packer.Schedule(pack.SaveTwoDTextureAsset("assets/textures/twod/rusty_metal_02_diff_512.dat",
		pack.OpenImageResource("resources/images/rusty_metal_02_diff_512.png"),
	))

	packer.Schedule(pack.SaveTwoDTextureAsset("assets/textures/twod/body.dat",
		pack.OpenImageResource("resources/images/body.png"),
	))

	packer.Schedule(pack.SaveTwoDTextureAsset("assets/textures/twod/wheel.dat",
		pack.OpenImageResource("resources/images/wheel.png"),
	))

	packer.Schedule(pack.SaveTwoDTextureAsset("assets/textures/twod/leafy_tree.dat",
		pack.OpenImageResource("resources/images/leafy_tree.png"),
	))

	srcEquirectangularImage := pack.OpenImageResource("resources/images/syferfontein.hdr")
	skyboxCubeImage := pack.BuildCubeImage().
		WithFrontImage(pack.BuildCubeSideFromEquirectangular(pack.CubeSideFront, srcEquirectangularImage)).
		WithRearImage(pack.BuildCubeSideFromEquirectangular(pack.CubeSideRear, srcEquirectangularImage)).
		WithLeftImage(pack.BuildCubeSideFromEquirectangular(pack.CubeSideLeft, srcEquirectangularImage)).
		WithRightImage(pack.BuildCubeSideFromEquirectangular(pack.CubeSideRight, srcEquirectangularImage)).
		WithTopImage(pack.BuildCubeSideFromEquirectangular(pack.CubeSideTop, srcEquirectangularImage)).
		WithBottomImage(pack.BuildCubeSideFromEquirectangular(pack.CubeSideBottom, srcEquirectangularImage))
	packer.Schedule(pack.SaveCubeTextureAsset("assets/textures/cube/syferfontein.dat", skyboxCubeImage).
		WithFormat(asset.DataFormatRGBA32F),
	)

	reflectionCubeImage := pack.ScaleCubeImage(skyboxCubeImage, 128)
	packer.Schedule(pack.SaveCubeTextureAsset("assets/textures/cube/syferfontein_reflection.dat", reflectionCubeImage).
		WithFormat(asset.DataFormatRGBA32F),
	)

	refractionCubeImage := pack.BuildIrradianceCubeImage(reflectionCubeImage).
		WithSampleCount(50)
	packer.Schedule(pack.SaveCubeTextureAsset("assets/textures/cube/syferfontein_refraction.dat", refractionCubeImage).
		WithFormat(asset.DataFormatRGBA32F),
	)

	packer.Schedule(pack.SaveModelAsset("assets/models/quad.dat",
		pack.OpenGLTFResource("resources/models/quad.gltf"),
	))

	packer.Schedule(pack.SaveModelAsset("assets/models/street_lamp.dat",
		pack.OpenGLTFResource("resources/models/street_lamp.gltf"),
	))

	packer.Schedule(pack.SaveModelAsset("assets/models/suv.dat",
		pack.OpenGLTFResource("resources/models/suv.gltf"),
	))

	packer.Schedule(pack.SaveModelAsset("assets/models/leafy_tree.dat",
		pack.OpenGLTFResource("resources/models/leafy_tree.gltf"),
	))

	packer.Schedule(pack.SaveLevelAsset("assets/levels/forest.dat",
		pack.OpenLevelResource("resources/levels/forest.json"),
	))

	packer.Schedule(pack.SaveLevelAsset("assets/levels/highway.dat",
		pack.OpenLevelResource("resources/levels/highway.json"),
	))

	packer.Schedule(pack.SaveLevelAsset("assets/levels/playground.dat",
		pack.OpenLevelResource("resources/levels/playground.json"),
	))

	if err := packer.Run(workerCount); err != nil {
		log.Fatalf("packer error: %s", err)
	}
}
