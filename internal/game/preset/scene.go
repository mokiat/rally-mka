package preset

import (
	"github.com/mokiat/lacking/game"
)

type GameScene = game.Scene

func NewScene(gameScene *GameScene) *Scene {
	result := &Scene{
		GameScene: gameScene,
		ECSScene:  newECSScene(gameScene.ECS()),
	}
	return result
}

type Scene struct {
	*GameScene
	*ECSScene
}

func (s *Scene) ECS() *ECSScene {
	return s.ECSScene
}

func (s *Scene) onEntityDeleted(entity Entity) {
	if nodeComp := entity.RefNodeComponent(); nodeComp != nil {
		s.Hierarchy().DeleteNode(nodeComp.Node.ID())
	}
}
