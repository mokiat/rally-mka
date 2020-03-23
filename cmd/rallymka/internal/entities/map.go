package entities

import (
	"github.com/mokiat/rally-mka/internal/engine/collision"
)

type Map interface {
	Load(path string) error
	CheckCollisionGround(line collision.Line) (collision.LineCollision, bool)
	CheckCollisionWall(line collision.Line) (collision.LineCollision, bool)
}

func NewMap() Map {
	return &gameMap{}
}

type gameMap struct {
	walls   []Wall
	grounds []Ground
}

func (m *gameMap) Load(path string) error {
	model := NewExtendedModel()
	if err := model.Load(path); err != nil {
		return err
	}
	m.loadGrounds(model)
	m.loadWalls(model)
	return nil
}

func (m *gameMap) loadGrounds(model *ExtendedModel) {
	searchIndex := model.FindObjectIndex("Grounds", 0, true)
	for searchIndex >= 0 {
		object := model.GetObjectIndex(searchIndex)
		m.grounds = append(m.grounds, Ground{
			CollisionMesh: createCollisionMesh(object),
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
		})
		searchIndex = model.FindObjectIndex("Walls", searchIndex+1, true)
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
