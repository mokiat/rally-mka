package main

import (
	"fmt"
	"log"

	"github.com/mokiat/lacking/data/pack"
	"github.com/mokiat/lacking/game/asset"
)

func main() {
	if err := runTool(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func runTool() error {
	registry, err := asset.NewDirRegistry(".")
	if err != nil {
		return fmt.Errorf("failed to create registry: %w", err)
	}

	skyboxDay := ensureResource(registry, "bab99e80-ded1-459a-b00b-6a17afa44046", "cube_texture", "Skybox Day")
	skyboxDayReflection := ensureResource(registry, "eb639f55-d6eb-46d7-bd3b-d52fcaa0bc58", "cube_texture", "Skybox Day Reflection")
	skyboxDayRefraction := ensureResource(registry, "0815fb89-7ae6-4229-b9e2-59610c4fc6bc", "cube_texture", "Skybox Day Refraction")

	skyboxNight := ensureResource(registry, "904431e9-7d78-4cbb-ae46-2e70b1458832", "cube_texture", "Skybox Night")
	skyboxNightReflection := ensureResource(registry, "892ece47-c4c8-4ff0-a9ad-0faa208eee81", "cube_texture", "Skybox Night Reflection")
	skyboxNightRefraction := ensureResource(registry, "e79da33a-5131-4029-b168-ea2c5378c169", "cube_texture", "Skybox Night Refraction")

	modelHomeScreen := ensureResource(registry, "d1aef712-4c5a-45b8-ba6f-0385e071a8f1", "model", "Content: Home Screen")
	modelForest := ensureResource(registry, "5f7bd967-dc4a-4252-b1a5-5721cd299d67", "model", "Forest Ride")
	modelSUV := ensureResource(registry, "eaeb7483-7271-441f-a470-c0a8fa225161", "model", "SUV")

	levelHomeScreen := ensureResource(registry, "80dd9049-c183-4d82-b5b2-6f38ca499055", "scene", "Home Screen")
	levelHomeScreen.AddDependency(modelHomeScreen)

	levelForestDay := ensureResource(registry, "884e6395-2300-47bb-9916-b80e3dc0e086", "scene", "Forest-Day")
	levelForestDay.AddDependency(skyboxDay)
	levelForestDay.AddDependency(skyboxDayReflection)
	levelForestDay.AddDependency(skyboxDayRefraction)
	levelForestDay.AddDependency(modelForest)
	levelForestDay.AddDependency(modelSUV)

	levelForestNight := ensureResource(registry, "a288e44d-3ed5-415b-b9c2-4dbe086dfce2", "scene", "Forest-Night")
	levelForestNight.AddDependency(skyboxNight)
	levelForestNight.AddDependency(skyboxNightReflection)
	levelForestNight.AddDependency(skyboxNightRefraction)
	levelForestNight.AddDependency(modelForest)
	levelForestNight.AddDependency(modelSUV)

	if err := registry.Save(); err != nil {
		return fmt.Errorf("error saving resources: %w", err)
	}

	packer := pack.NewPacker(registry)

	// Cube Textures
	packer.Pipeline(func(p *pack.Pipeline) {
		equirectangularImage := p.OpenImageResource("resources/images/day.hdr")
		cubeImage := p.BuildCubeImage(
			pack.WithFrontImage(p.BuildCubeSideFromEquirectangular(pack.CubeSideFront, equirectangularImage)),
			pack.WithRearImage(p.BuildCubeSideFromEquirectangular(pack.CubeSideRear, equirectangularImage)),
			pack.WithLeftImage(p.BuildCubeSideFromEquirectangular(pack.CubeSideLeft, equirectangularImage)),
			pack.WithRightImage(p.BuildCubeSideFromEquirectangular(pack.CubeSideRight, equirectangularImage)),
			pack.WithTopImage(p.BuildCubeSideFromEquirectangular(pack.CubeSideTop, equirectangularImage)),
			pack.WithBottomImage(p.BuildCubeSideFromEquirectangular(pack.CubeSideBottom, equirectangularImage)),
		)

		smallerCubeImage := p.ScaleCubeImage(cubeImage, 512)
		p.SaveCubeTextureAsset(skyboxDay, smallerCubeImage,
			pack.WithFormat(asset.TexelFormatRGBA16F),
		)

		reflectionCubeImage := p.ScaleCubeImage(cubeImage, 128)
		p.SaveCubeTextureAsset(skyboxDayReflection, reflectionCubeImage,
			pack.WithFormat(asset.TexelFormatRGBA16F),
		)

		refractionCubeImage := p.BuildIrradianceCubeImage(reflectionCubeImage,
			pack.WithSampleCount(50),
		)
		p.SaveCubeTextureAsset(skyboxDayRefraction, refractionCubeImage,
			pack.WithFormat(asset.TexelFormatRGBA16F),
		)
	})

	packer.Pipeline(func(p *pack.Pipeline) {
		// Using GIMP:
		// 1. Scale height to %50
		// 2. Convert to float32
		// 3. Exposure: ~ -6

		equirectangularImage := p.OpenImageResource("resources/images/night.exr")
		cubeImage := p.BuildCubeImage(
			pack.WithFrontImage(p.BuildCubeSideFromEquirectangular(pack.CubeSideFront, equirectangularImage)),
			pack.WithRearImage(p.BuildCubeSideFromEquirectangular(pack.CubeSideRear, equirectangularImage)),
			pack.WithLeftImage(p.BuildCubeSideFromEquirectangular(pack.CubeSideLeft, equirectangularImage)),
			pack.WithRightImage(p.BuildCubeSideFromEquirectangular(pack.CubeSideRight, equirectangularImage)),
			pack.WithTopImage(p.BuildCubeSideFromEquirectangular(pack.CubeSideTop, equirectangularImage)),
			pack.WithBottomImage(p.BuildCubeSideFromEquirectangular(pack.CubeSideBottom, equirectangularImage)),
		)

		smallerCubeImage := p.ScaleCubeImage(cubeImage, 512)
		p.SaveCubeTextureAsset(skyboxNight, smallerCubeImage,
			pack.WithFormat(asset.TexelFormatRGBA16F),
		)

		reflectionCubeImage := p.ScaleCubeImage(cubeImage, 128)
		p.SaveCubeTextureAsset(skyboxNightReflection, reflectionCubeImage,
			pack.WithFormat(asset.TexelFormatRGBA16F),
		)

		refractionCubeImage := p.BuildIrradianceCubeImage(reflectionCubeImage,
			pack.WithSampleCount(50),
		)
		p.SaveCubeTextureAsset(skyboxNightRefraction, refractionCubeImage,
			pack.WithFormat(asset.TexelFormatRGBA16F),
		)
	})

	// Models
	packer.Pipeline(func(p *pack.Pipeline) {
		p.SaveModelAsset(modelHomeScreen,
			p.OpenGLTFResource("resources/models/home-screen.glb"),
		)
	})

	packer.Pipeline(func(p *pack.Pipeline) {
		p.SaveModelAsset(modelSUV,
			p.OpenGLTFResource("resources/models/suv.glb"),
		)
	})

	packer.Pipeline(func(p *pack.Pipeline) {
		p.SaveModelAsset(modelForest,
			p.OpenGLTFResource("resources/models/forest.glb"),
		)
	})

	// Levels
	packer.Pipeline(func(p *pack.Pipeline) {
		p.SaveLevelAsset(levelHomeScreen,
			p.OpenLevelResource("resources/levels/home-screen.json"),
		)
	})

	packer.Pipeline(func(p *pack.Pipeline) {
		p.SaveLevelAsset(levelForestDay,
			p.OpenLevelResource("resources/levels/forest-day.json"),
		)
	})

	packer.Pipeline(func(p *pack.Pipeline) {
		p.SaveLevelAsset(levelForestNight,
			p.OpenLevelResource("resources/levels/forest-night.json"),
		)
	})

	return packer.RunParallel()
}

func ensureResource(registry asset.Registry, id, kind, name string) asset.Resource {
	resource := registry.ResourceByID(id)
	if resource == nil {
		resource = registry.CreateIDResource(id, kind, name)
	} else {
		resource.SetName(name)
	}
	return resource
}
