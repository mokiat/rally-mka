package main

import (
	"fmt"
	"log"

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

	var (
		tex2DConcrete   = ensureResource(registry, "e89d3c68-12ba-42ca-bc04-ccefefcf5720", "twod_texture", "Concrete")
		tex2DRoad       = ensureResource(registry, "b2c7a46f-f2a2-4601-bd10-493a68fc094c", "twod_texture", "Road")
		tex2DBarrier    = ensureResource(registry, "3800657a-4407-4bda-bdff-b57748c002ab", "twod_texture", "Barrier")
		tex2DGrass      = ensureResource(registry, "5905593c-6820-465b-8315-1e17be0a6f72", "twod_texture", "Grass")
		tex2DGravel     = ensureResource(registry, "fee51386-e8f1-4dd7-926e-3ca353589e01", "twod_texture", "Gravel")
		tex2DRustyMetal = ensureResource(registry, "7574ab74-980b-4bbd-aea1-3459998ccc71", "twod_texture", "Rusty Metal")
		tex2DCarBody    = ensureResource(registry, "8ed7f73b-0129-48f8-b187-b7e43eff0294", "twod_texture", "Car Body")
		tex2DCarWheel   = ensureResource(registry, "9d52b91b-0292-4442-b100-29f12402c2c1", "twod_texture", "Car Wheel")
		tex2DLeafyTree  = ensureResource(registry, "01b8ab21-3d72-4b31-869e-37039dee0161", "twod_texture", "Leafy Tree")
	)

	var (
		texCubeSkybox           = ensureResource(registry, "bab99e80-ded1-459a-b00b-6a17afa44046", "cube_texture", "Skybox")
		texCubeSkyboxReflection = ensureResource(registry, "eb639f55-d6eb-46d7-bd3b-d52fcaa0bc58", "cube_texture", "Skybox Reflection")
		texCubeSkyboxRefraction = ensureResource(registry, "0815fb89-7ae6-4229-b9e2-59610c4fc6bc", "cube_texture", "Skybox Refraction")
	)

	var (
		modelStreetLamp = ensureResource(registry, "31cb3900-760d-4179-b5d9-79f8e69be8f6", "model", "Street Lamp")
		modelSUV        = ensureResource(registry, "eaeb7483-7271-441f-a470-c0a8fa225161", "model", "SUV")
		modelLeafyTree  = ensureResource(registry, "2c6e3211-68f8-4b31-beaf-e52af5d3be31", "model", "Leafy Tree")
	)

	modelStreetLamp.AddDependency(tex2DRustyMetal)
	modelSUV.AddDependency(tex2DCarBody)
	modelSUV.AddDependency(tex2DCarWheel)
	modelLeafyTree.AddDependency(tex2DLeafyTree)

	var (
		levelPlayground = ensureResource(registry, "9ca25b5c-ffa0-4224-ad80-a3c4d67930b7", "scene", "Playground")
		levelForest     = ensureResource(registry, "884e6395-2300-47bb-9916-b80e3dc0e086", "scene", "Forest")
		levelHighway    = ensureResource(registry, "acf21108-47ad-44ef-ba21-bf5473bfbaa0", "scene", "Highway")
	)

	levelPlayground.AddDependency(tex2DGrass)
	levelPlayground.AddDependency(texCubeSkybox)
	levelPlayground.AddDependency(texCubeSkyboxReflection)
	levelPlayground.AddDependency(texCubeSkyboxRefraction)
	levelPlayground.AddDependency(modelSUV)
	levelPlayground.AddDependency(modelLeafyTree)
	levelForest.AddDependency(tex2DGravel)
	levelForest.AddDependency(tex2DGrass)
	levelForest.AddDependency(tex2DBarrier)
	levelForest.AddDependency(texCubeSkybox)
	levelForest.AddDependency(texCubeSkyboxReflection)
	levelForest.AddDependency(texCubeSkyboxRefraction)
	levelForest.AddDependency(modelSUV)
	levelForest.AddDependency(modelLeafyTree)
	levelHighway.AddDependency(tex2DRoad)
	levelHighway.AddDependency(tex2DConcrete)
	levelHighway.AddDependency(texCubeSkybox)
	levelHighway.AddDependency(texCubeSkyboxReflection)
	levelHighway.AddDependency(texCubeSkyboxRefraction)
	levelHighway.AddDependency(modelSUV)
	levelHighway.AddDependency(modelStreetLamp)

	if err := registry.Save(); err != nil {
		return fmt.Errorf("error saving resources: %w", err)
	}

	packer := pack.NewPacker(registry)

	// TwoD Textures
	packer.Pipeline(func(p *pack.Pipeline) {
		p.SaveTwoDTextureAsset(tex2DConcrete.ID(),
			p.OpenImageResource("resources/images/concrete.png"),
		)

		p.SaveTwoDTextureAsset(tex2DRoad.ID(),
			p.OpenImageResource("resources/images/road.png"),
		)

		p.SaveTwoDTextureAsset(tex2DBarrier.ID(),
			p.OpenImageResource("resources/images/barrier.png"),
		)

		p.SaveTwoDTextureAsset(tex2DGrass.ID(),
			p.OpenImageResource("resources/images/grass.png"),
		)

		p.SaveTwoDTextureAsset(tex2DGravel.ID(),
			p.OpenImageResource("resources/images/gravel.png"),
		)

		p.SaveTwoDTextureAsset(tex2DRustyMetal.ID(),
			p.OpenImageResource("resources/images/rusty_metal_02_diff_512.png"),
		)

		p.SaveTwoDTextureAsset(tex2DCarBody.ID(),
			p.OpenImageResource("resources/images/body.png"),
		)

		p.SaveTwoDTextureAsset(tex2DCarWheel.ID(),
			p.OpenImageResource("resources/images/wheel.png"),
		)

		p.SaveTwoDTextureAsset(tex2DLeafyTree.ID(),
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

		smallerSkyboxCubeImage := p.ScaleCubeImage(skyboxCubeImage, 512)
		p.SaveCubeTextureAsset(texCubeSkybox.ID(), smallerSkyboxCubeImage,
			pack.WithFormat(gameasset.TexelFormatRGBA16F),
		)

		reflectionCubeImage := p.ScaleCubeImage(skyboxCubeImage, 128)
		p.SaveCubeTextureAsset(texCubeSkyboxReflection.ID(), reflectionCubeImage,
			pack.WithFormat(gameasset.TexelFormatRGBA16F),
		)

		refractionCubeImage := p.BuildIrradianceCubeImage(reflectionCubeImage,
			pack.WithSampleCount(50),
		)
		p.SaveCubeTextureAsset(texCubeSkyboxRefraction.ID(), refractionCubeImage,
			pack.WithFormat(gameasset.TexelFormatRGBA16F),
		)
	})

	// Models
	packer.Pipeline(func(p *pack.Pipeline) {
		p.SaveModelAsset(modelStreetLamp.ID(),
			p.OpenGLTFResource("resources/models/street_lamp.gltf"),
		)

		p.SaveModelAsset(modelSUV.ID(),
			p.OpenGLTFResource("resources/models/suv.gltf"),
		)

		p.SaveModelAsset(modelLeafyTree.ID(),
			p.OpenGLTFResource("resources/models/leafy_tree.gltf"),
		)
	})

	// Levels
	packer.Pipeline(func(p *pack.Pipeline) {
		p.SaveLevelAsset(levelPlayground.ID(),
			p.OpenLevelResource("resources/levels/playground.json"),
		)

		p.SaveLevelAsset(levelForest.ID(),
			p.OpenLevelResource("resources/levels/forest.json"),
		)

		p.SaveLevelAsset(levelHighway.ID(),
			p.OpenLevelResource("resources/levels/highway.json"),
		)
	})

	return packer.RunParallel()
}

func ensureResource(registry gameasset.Registry, id, kind, name string) gameasset.Resource {
	resource := registry.ResourceByID(id)
	if resource == nil {
		resource = registry.CreateIDResource(id, kind, name)
	}
	return resource
}
