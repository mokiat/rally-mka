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

	texCubeSkybox := ensureResource(registry, "bab99e80-ded1-459a-b00b-6a17afa44046", "cube_texture", "Skybox")
	texCubeSkyboxReflection := ensureResource(registry, "eb639f55-d6eb-46d7-bd3b-d52fcaa0bc58", "cube_texture", "Skybox Reflection")
	texCubeSkyboxRefraction := ensureResource(registry, "0815fb89-7ae6-4229-b9e2-59610c4fc6bc", "cube_texture", "Skybox Refraction")

	modelForest := ensureResource(registry, "5f7bd967-dc4a-4252-b1a5-5721cd299d67", "model", "Forest Ride")
	modelForest.AddDependency(texCubeSkybox)
	modelForest.AddDependency(texCubeSkyboxReflection)
	modelForest.AddDependency(texCubeSkyboxRefraction)
	modelSUV := ensureResource(registry, "eaeb7483-7271-441f-a470-c0a8fa225161", "model", "SUV")

	levelForest := ensureResource(registry, "884e6395-2300-47bb-9916-b80e3dc0e086", "scene", "Forest")
	levelForest.AddDependency(modelForest)
	levelForest.AddDependency(modelSUV)

	if err := registry.Save(); err != nil {
		return fmt.Errorf("error saving resources: %w", err)
	}

	packer := pack.NewPacker(registry)

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

		smallerSkyboxCubeImage := p.ScaleCubeImage(skyboxCubeImage, 1024)
		p.SaveCubeTextureAsset(texCubeSkybox, smallerSkyboxCubeImage,
			pack.WithFormat(asset.TexelFormatRGBA16F),
		)

		reflectionCubeImage := p.ScaleCubeImage(skyboxCubeImage, 128)
		p.SaveCubeTextureAsset(texCubeSkyboxReflection, reflectionCubeImage,
			pack.WithFormat(asset.TexelFormatRGBA16F),
		)

		refractionCubeImage := p.BuildIrradianceCubeImage(reflectionCubeImage,
			pack.WithSampleCount(50),
		)
		p.SaveCubeTextureAsset(texCubeSkyboxRefraction, refractionCubeImage,
			pack.WithFormat(asset.TexelFormatRGBA16F),
		)
	})

	// Models
	packer.Pipeline(func(p *pack.Pipeline) {
		p.SaveModelAsset(modelSUV,
			p.OpenGLTFResource("resources/models/suv.glb"),
		)

		p.SaveModelAsset(modelForest,
			p.OpenGLTFResource("resources/models/forest.glb"),
		)
	})

	// Levels
	packer.Pipeline(func(p *pack.Pipeline) {
		p.SaveLevelAsset(levelForest,
			p.OpenLevelResource("resources/levels/forest.json"),
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
