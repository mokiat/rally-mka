package main

import (
	"fmt"
	"log"

	"github.com/mokiat/lacking/data/asset"
	"github.com/mokiat/lacking/data/pack"
	gameasset "github.com/mokiat/lacking/game/asset"
)

func main() {
	if err := runTool(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func runTool() error {
	registry, err := gameasset.NewDirRegistry(".")
	if err != nil {
		return fmt.Errorf("failed to create registry: %w", err)
	}

	resources := []gameasset.Resource{
		// TwoD Textures
		{
			GUID: "c80a8260-4e6d-479b-add9-cb93a86ca0ee",
			Kind: "twod_texture",
			Name: "Loading",
		},
		{
			GUID: "e89d3c68-12ba-42ca-bc04-ccefefcf5720",
			Kind: "twod_texture",
			Name: "Concrete",
		},
		{
			GUID: "b2c7a46f-f2a2-4601-bd10-493a68fc094c",
			Kind: "twod_texture",
			Name: "Road",
		},
		{
			GUID: "3800657a-4407-4bda-bdff-b57748c002ab",
			Kind: "twod_texture",
			Name: "Barrier",
		},
		{
			GUID: "5905593c-6820-465b-8315-1e17be0a6f72",
			Kind: "twod_texture",
			Name: "Grass",
		},
		{
			GUID: "fee51386-e8f1-4dd7-926e-3ca353589e01",
			Kind: "twod_texture",
			Name: "Gravel",
		},
		{
			GUID: "7574ab74-980b-4bbd-aea1-3459998ccc71",
			Kind: "twod_texture",
			Name: "Rusty Metal",
		},
		{
			GUID: "8ed7f73b-0129-48f8-b187-b7e43eff0294",
			Kind: "twod_texture",
			Name: "Car Body",
		},
		{
			GUID: "9d52b91b-0292-4442-b100-29f12402c2c1",
			Kind: "twod_texture",
			Name: "Car Wheel",
		},
		{
			GUID: "01b8ab21-3d72-4b31-869e-37039dee0161",
			Kind: "twod_texture",
			Name: "Leafy Tree",
		},

		// Cube Textures
		{
			GUID: "bab99e80-ded1-459a-b00b-6a17afa44046",
			Kind: "cube_texture",
			Name: "Skybox",
		},
		{
			GUID: "eb639f55-d6eb-46d7-bd3b-d52fcaa0bc58",
			Kind: "cube_texture",
			Name: "Skybox Reflection",
		},
		{
			GUID: "0815fb89-7ae6-4229-b9e2-59610c4fc6bc",
			Kind: "cube_texture",
			Name: "Skybox Refraction",
		},

		// Models
		{
			GUID: "5323c8e7-9eb6-471f-b14b-585d6ad260f4",
			Kind: "model",
			Name: "Quad",
		},
		{
			GUID: "31cb3900-760d-4179-b5d9-79f8e69be8f6",
			Kind: "model",
			Name: "Street Lamp",
		},
		{
			GUID: "eaeb7483-7271-441f-a470-c0a8fa225161",
			Kind: "model",
			Name: "SUV",
		},
		{
			GUID: "2c6e3211-68f8-4b31-beaf-e52af5d3be31",
			Kind: "model",
			Name: "Leafy Tree",
		},

		// Levels
		{
			GUID: "9ca25b5c-ffa0-4224-ad80-a3c4d67930b7",
			Kind: "scene",
			Name: "Playground",
		},
		{
			GUID: "884e6395-2300-47bb-9916-b80e3dc0e086",
			Kind: "scene",
			Name: "Forest",
		},
		{
			GUID: "acf21108-47ad-44ef-ba21-bf5473bfbaa0",
			Kind: "scene",
			Name: "Highway",
		},
	}
	if err := registry.WriteResources(resources); err != nil {
		return fmt.Errorf("failed to write resources: %w", err)
	}

	dependencies := []gameasset.Dependency{
		{
			SourceGUID: "31cb3900-760d-4179-b5d9-79f8e69be8f6", // lamp
			TargetGUID: "7574ab74-980b-4bbd-aea1-3459998ccc71", // rusty metal
		},
		{
			SourceGUID: "2c6e3211-68f8-4b31-beaf-e52af5d3be31", // leafy tree
			TargetGUID: "01b8ab21-3d72-4b31-869e-37039dee0161", // leafy tree
		},
		{
			SourceGUID: "eaeb7483-7271-441f-a470-c0a8fa225161", // suv
			TargetGUID: "8ed7f73b-0129-48f8-b187-b7e43eff0294", // car body
		},
		{
			SourceGUID: "eaeb7483-7271-441f-a470-c0a8fa225161", // suv
			TargetGUID: "9d52b91b-0292-4442-b100-29f12402c2c1", // car wheel
		},
	}
	if err := registry.WriteDependencies(dependencies); err != nil {
		return fmt.Errorf("failed to write dependencies: %w", err)
	}

	packer := pack.NewPacker(registry)

	// TwoD Textures
	packer.Pipeline(func(p *pack.Pipeline) {
		p.SaveTwoDTextureAsset("c80a8260-4e6d-479b-add9-cb93a86ca0ee",
			p.OpenImageResource("resources/images/loading.png"),
		)

		p.SaveTwoDTextureAsset("e89d3c68-12ba-42ca-bc04-ccefefcf5720",
			p.OpenImageResource("resources/images/concrete.png"),
		)

		p.SaveTwoDTextureAsset("b2c7a46f-f2a2-4601-bd10-493a68fc094c",
			p.OpenImageResource("resources/images/road.png"),
		)

		p.SaveTwoDTextureAsset("3800657a-4407-4bda-bdff-b57748c002ab",
			p.OpenImageResource("resources/images/barrier.png"),
		)

		p.SaveTwoDTextureAsset("5905593c-6820-465b-8315-1e17be0a6f72",
			p.OpenImageResource("resources/images/grass.png"),
		)

		p.SaveTwoDTextureAsset("fee51386-e8f1-4dd7-926e-3ca353589e01",
			p.OpenImageResource("resources/images/gravel.png"),
		)

		p.SaveTwoDTextureAsset("7574ab74-980b-4bbd-aea1-3459998ccc71",
			p.OpenImageResource("resources/images/rusty_metal_02_diff_512.png"),
		)

		p.SaveTwoDTextureAsset("8ed7f73b-0129-48f8-b187-b7e43eff0294",
			p.OpenImageResource("resources/images/body.png"),
		)

		p.SaveTwoDTextureAsset("9d52b91b-0292-4442-b100-29f12402c2c1",
			p.OpenImageResource("resources/images/wheel.png"),
		)

		p.SaveTwoDTextureAsset("01b8ab21-3d72-4b31-869e-37039dee0161",
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

		p.SaveCubeTextureAsset("bab99e80-ded1-459a-b00b-6a17afa44046", skyboxCubeImage,
			pack.WithFormat(asset.TexelFormatRGBA32F),
		)

		reflectionCubeImage := p.ScaleCubeImage(skyboxCubeImage, 128)
		p.SaveCubeTextureAsset("eb639f55-d6eb-46d7-bd3b-d52fcaa0bc58", reflectionCubeImage,
			pack.WithFormat(asset.TexelFormatRGBA32F),
		)

		refractionCubeImage := p.BuildIrradianceCubeImage(reflectionCubeImage,
			pack.WithSampleCount(50),
		)
		p.SaveCubeTextureAsset("0815fb89-7ae6-4229-b9e2-59610c4fc6bc", refractionCubeImage,
			pack.WithFormat(asset.TexelFormatRGBA32F),
		)
	})

	// Models
	packer.Pipeline(func(p *pack.Pipeline) {
		p.SaveModelAsset("5323c8e7-9eb6-471f-b14b-585d6ad260f4",
			p.OpenGLTFResource("resources/models/quad.gltf"),
		)

		p.SaveModelAsset("31cb3900-760d-4179-b5d9-79f8e69be8f6",
			p.OpenGLTFResource("resources/models/street_lamp.gltf"),
		)

		p.SaveModelAsset("eaeb7483-7271-441f-a470-c0a8fa225161",
			p.OpenGLTFResource("resources/models/suv.gltf"),
		)

		p.SaveModelAsset("2c6e3211-68f8-4b31-beaf-e52af5d3be31",
			p.OpenGLTFResource("resources/models/leafy_tree.gltf"),
		)
	})

	// Levels
	packer.Pipeline(func(p *pack.Pipeline) {
		p.SaveLevelAsset("9ca25b5c-ffa0-4224-ad80-a3c4d67930b7",
			p.OpenLevelResource("resources/levels/playground.json"),
		)

		p.SaveLevelAsset("884e6395-2300-47bb-9916-b80e3dc0e086",
			p.OpenLevelResource("resources/levels/forest.json"),
		)

		p.SaveLevelAsset("acf21108-47ad-44ef-ba21-bf5473bfbaa0",
			p.OpenLevelResource("resources/levels/highway.json"),
		)
	})

	return packer.RunParallel()
}
