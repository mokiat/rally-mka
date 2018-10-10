package entities

import (
	"fmt"

	"github.com/mokiat/go-whiskey/math"
	"github.com/mokiat/rally-mka/collision"
	"github.com/mokiat/rally-mka/render"
)

type Map interface {
	Load(path string) error
	Generate()
	Draw(renderer *render.Renderer)
	CheckCollisionGround(line collision.Line) (collision.LineCollision, bool)
	CheckCollisionWall(line collision.Line) (collision.LineCollision, bool)
}

func NewMap() Map {
	return &gameMap{}
}

type gameMap struct {
	walls     []Wall
	grounds   []Ground
	dummies   []Dummy
	waypoints []math.Vec3
}

func (m *gameMap) Load(path string) error {
	model := NewExtendedModel()
	if err := model.Load(path); err != nil {
		return err
	}
	m.loadGrounds(model)
	m.loadWalls(model)
	m.loadDummies(model)
	m.loadWaypoints(model)
	return nil
}

func (m *gameMap) loadGrounds(model *ExtendedModel) {
	searchIndex := model.FindObjectIndex("Grounds", 0, true)
	for searchIndex >= 0 {
		object := model.GetObjectIndex(searchIndex)
		m.grounds = append(m.grounds, Ground{
			CollisionMesh: createCollisionMesh(object),
			RenderMesh:    createRenderMesh(model, object),
		})
		searchIndex = model.FindObjectIndex("Grounds", searchIndex+1, true)
	}
}

func (m *gameMap) loadWalls(model *ExtendedModel) {
	searchIndex := model.FindObjectIndex("Walls", 0, true)
	for searchIndex >= 0 {
		object := model.GetObjectIndex(searchIndex)
		m.walls = append(m.walls, Wall{
			CollisionMesh: createCollisionMesh(object),
			RenderMesh:    createRenderMesh(model, object),
		})
		searchIndex = model.FindObjectIndex("Walls", searchIndex+1, true)
	}
}

func (m *gameMap) loadDummies(model *ExtendedModel) {
	searchIndex := model.FindObjectIndex("Dummy", 0, true)
	for searchIndex >= 0 {
		object := model.GetObjectIndex(searchIndex)
		m.dummies = append(m.dummies, Dummy{
			RenderMesh: createRenderMesh(model, object),
		})
		searchIndex = model.FindObjectIndex("Dummy", searchIndex+1, true)
	}
}

func (m *gameMap) loadWaypoints(model *ExtendedModel) {
	waypointID := 1
	searchIndex := model.FindObjectIndex(fmt.Sprintf("Way%d", waypointID), 0, false)
	for searchIndex >= 0 {
		object := model.Objects[searchIndex]
		m.waypoints = append(m.waypoints, model.ObjectCenter(object))
		waypointID++
		searchIndex = model.FindObjectIndex(fmt.Sprintf("Way%d", waypointID), searchIndex+1, false)
	}
}

func (m *gameMap) Generate() {
	for _, ground := range m.grounds {
		ground.RenderMesh.Generate()
	}
	for _, wall := range m.walls {
		wall.RenderMesh.Generate()
	}
	for _, dummy := range m.dummies {
		dummy.RenderMesh.Generate()
	}
}

func (m *gameMap) CheckCollisionGround(line collision.Line) (bestCollision collision.LineCollision, found bool) {
	closestDistance := line.LengthSquared()
	for _, ground := range m.grounds {
		if lineCollision, ok := ground.CollisionMesh.LineCollision(line); ok {
			found = true
			distanceVector := lineCollision.Intersection().DecVec3(line.Start())
			distance := distanceVector.LengthSquared()
			if distance < closestDistance {
				closestDistance = distance
				bestCollision = lineCollision
			}
		}
	}
	return
}

func (m *gameMap) CheckCollisionWall(line collision.Line) (bestCollision collision.LineCollision, found bool) {
	closestDistance := line.LengthSquared()
	for _, wall := range m.walls {
		if lineCollision, ok := wall.CollisionMesh.LineCollision(line); ok {
			found = true
			distanceVector := lineCollision.Intersection().DecVec3(line.Start())
			distance := distanceVector.LengthSquared()
			if distance < closestDistance {
				closestDistance = distance
				bestCollision = lineCollision
			}
		}
	}
	return
}

func (m *gameMap) Draw(renderer *render.Renderer) {
	for _, ground := range m.grounds {
		renderer.Render(ground.RenderMesh, renderer.TextureMaterial())
	}
	for _, wall := range m.walls {
		renderer.Render(wall.RenderMesh, renderer.TextureMaterial())
	}
	for _, dummy := range m.dummies {
		renderer.Render(dummy.RenderMesh, renderer.TextureMaterial())
	}
}
