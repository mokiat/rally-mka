package stream

import (
	"fmt"

	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/rally-mka/internal/data/asset"
	"github.com/mokiat/rally-mka/internal/engine/graphics"
	"github.com/mokiat/rally-mka/internal/engine/physics/collision"
	"github.com/mokiat/rally-mka/internal/engine/resource"
)

const levelResourceType = "level"

func GetLevel(registry *resource.Registry, name string) LevelHandle {
	return LevelHandle{
		Handle: registry.Type(levelResourceType).Resource(name),
	}
}

type LevelHandle struct {
	resource.Handle
}

func (h LevelHandle) Get() *Level {
	return h.Handle.Get().(*Level)
}

func (h LevelHandle) IsAvailable() bool {
	return h.Handle.IsAvailable() && h.Get().isAvailable()
}

type Level struct {
	Waypoints          []sprec.Vec3
	SkyboxTexture      CubeTextureHandle
	CollisionMeshes    []*collision.Mesh
	StartCollisionMesh *collision.Mesh
	StaticMeshes       []*Mesh
	StaticEntities     []*Entity
}

func (l Level) isAvailable() bool {
	if !l.SkyboxTexture.IsAvailable() {
		return false
	}
	for _, mesh := range l.StaticMeshes {
		if !mesh.IsAvailable() {
			return false
		}
	}
	for _, entity := range l.StaticEntities {
		if !entity.IsAvailable() {
			return false
		}
	}
	return true
}

type Entity struct {
	Model  ModelHandle
	Matrix sprec.Mat4
}

func (e Entity) IsAvailable() bool {
	return e.Model.IsAvailable()
}

func NewLevelOperator(locator resource.Locator, gfxWorker *graphics.Worker) *LevelOperator {
	return &LevelOperator{
		locator:   locator,
		gfxWorker: gfxWorker,
	}
}

type LevelOperator struct {
	locator   resource.Locator
	gfxWorker *graphics.Worker
}

func (o *LevelOperator) Register(registry *resource.Registry) {
	registry.RegisterType(levelResourceType, o)
}

func (o *LevelOperator) Allocate(registry *resource.Registry, name string) (resource.Resource, error) {
	in, err := o.locator.Open("assets", "levels", name)
	if err != nil {
		return nil, fmt.Errorf("failed to open level asset %q: %w", name, err)
	}
	defer in.Close()

	levelAsset, err := asset.NewLevelDecoder().Decode(in)
	if err != nil {
		return nil, fmt.Errorf("failed to decode level asset %q: %w", name, err)
	}

	level := &Level{}

	waypoints := make([]sprec.Vec3, len(levelAsset.Waypoints))
	for i, waypointAsset := range levelAsset.Waypoints {
		waypoints[i] = sprec.NewVec3(waypointAsset[0], waypointAsset[1], waypointAsset[2])
	}
	level.Waypoints = waypoints

	skyboxTexture := GetCubeTexture(registry, levelAsset.SkyboxTexture)
	registry.Request(skyboxTexture.Handle)
	level.SkyboxTexture = skyboxTexture

	convertCollisionMesh := func(collisionMeshAsset asset.LevelCollisionMesh) *collision.Mesh {
		var triangles []collision.Triangle
		for _, triangleAsset := range collisionMeshAsset.Triangles {
			triangles = append(triangles, collision.MakeTriangle(
				sprec.NewVec3(triangleAsset[0][0], triangleAsset[0][1], triangleAsset[0][2]),
				sprec.NewVec3(triangleAsset[1][0], triangleAsset[1][1], triangleAsset[1][2]),
				sprec.NewVec3(triangleAsset[2][0], triangleAsset[2][1], triangleAsset[2][2]),
			))
		}
		return collision.NewMesh(triangles)
	}

	collisionMeshes := make([]*collision.Mesh, len(levelAsset.CollisionMeshes))
	for i, collisionMeshAsset := range levelAsset.CollisionMeshes {
		collisionMeshes[i] = convertCollisionMesh(collisionMeshAsset)
	}
	level.CollisionMeshes = collisionMeshes

	level.StartCollisionMesh = convertCollisionMesh(levelAsset.StartCollisionMesh)

	staticMeshes := make([]*Mesh, len(levelAsset.StaticMeshes))
	for i, staticMeshAsset := range levelAsset.StaticMeshes {
		staticMesh, err := AllocateMesh(registry, o.gfxWorker, &staticMeshAsset)
		if err != nil {
			return nil, fmt.Errorf("failed to allocate mesh: %w", err)
		}
		staticMeshes[i] = staticMesh
	}
	level.StaticMeshes = staticMeshes

	staticEntities := make([]*Entity, len(levelAsset.StaticEntities))
	for i, staticEntityAsset := range levelAsset.StaticEntities {
		model := GetModel(registry, staticEntityAsset.Model)
		registry.Request(model.Handle)
		staticEntities[i] = &Entity{
			Model:  model,
			Matrix: floatArrayToMatrix(staticEntityAsset.Matrix),
		}
	}
	level.StaticEntities = staticEntities

	return level, nil
}

func (o *LevelOperator) Release(registry *resource.Registry, resource resource.Resource) error {
	level := resource.(*Level)

	for _, staticEntity := range level.StaticEntities {
		registry.Dismiss(staticEntity.Model.Handle)
	}
	for _, staticMesh := range level.StaticMeshes {
		if err := ReleaseMesh(registry, o.gfxWorker, staticMesh); err != nil {
			return fmt.Errorf("failed to release mesh: %w", err)
		}
	}
	registry.Dismiss(level.SkyboxTexture.Handle)
	return nil
}
