package car

import (
	"fmt"

	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/lacking/game"
	"github.com/mokiat/lacking/game/physics"
	"github.com/mokiat/lacking/resource"
	"github.com/mokiat/lacking/shape"
)

const (
	chassisRadius            = 2
	chassisMass              = 1300.0 / 5.0
	chassisMomentOfInertia   = chassisMass * chassisRadius * chassisRadius / 5.0
	chassisDragFactor        = 0.0 // 0.5 * 6.8 * 1.0
	chassisAngularDragFactor = 0.0 // 0.5 * 6.8 * 1.0
	chassisRestitutionCoef   = 0.0
)

func Chassis(model *resource.Model) *ChassisBuilder {
	return &ChassisBuilder{
		model: model,
	}
}

type ChassisBuilder struct {
	model     *resource.Model
	modifiers []func(node *game.Node)
}

func (b *ChassisBuilder) WithName(name string) *ChassisBuilder {
	b.modifiers = append(b.modifiers, func(node *game.Node) {
		body := node.Body()
		body.SetName(name)
	})
	return b
}

func (b *ChassisBuilder) WithPosition(position sprec.Vec3) *ChassisBuilder {
	b.modifiers = append(b.modifiers, func(node *game.Node) {
		body := node.Body()
		body.SetPosition(position)
	})
	return b
}

func (b *ChassisBuilder) Build(scene *game.Scene) *game.Node {
	instance, found := b.model.FindMeshInstance("Chassis")
	if !found {
		panic(fmt.Errorf("mesh instance %q not found", "Chassis"))
	}
	definition := instance.MeshDefinition
	bodyNode := instance.Node

	physicsBody := scene.Physics().CreateBody()
	physicsBody.SetPosition(sprec.ZeroVec3())
	physicsBody.SetOrientation(sprec.IdentityQuat())
	physicsBody.SetMass(chassisMass)
	physicsBody.SetMomentOfInertia(physics.SymmetricMomentOfInertia(chassisMomentOfInertia))
	physicsBody.SetDragFactor(chassisDragFactor)
	physicsBody.SetAngularDragFactor(chassisAngularDragFactor)
	physicsBody.SetRestitutionCoefficient(chassisRestitutionCoef)
	physicsBody.SetCollisionShapes([]physics.CollisionShape{
		shape.Placement{
			Position:    sprec.NewVec3(0.0, 0.3, -0.4),
			Orientation: sprec.IdentityQuat(),
			Shape:       shape.NewStaticBox(1.6, 1.4, 4.0),
		},
	})

	gfxMesh := scene.Graphics().CreateMesh(definition.GFXMeshTemplate)
	gfxMesh.SetMatrix(bodyNode.Matrix)

	node := game.NewNode()
	node.SetBody(physicsBody)
	node.SetMesh(gfxMesh)
	for _, modifier := range b.modifiers {
		modifier(node)
	}
	scene.Root().AppendChild(node)
	return node
}
