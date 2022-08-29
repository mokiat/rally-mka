package car

import (
	"fmt"

	"github.com/mokiat/gomath/dprec"
	"github.com/mokiat/lacking/game"
	"github.com/mokiat/lacking/game/physics"
	"github.com/mokiat/lacking/resource"
	"github.com/mokiat/lacking/util/shape"
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

func (b *ChassisBuilder) WithPosition(position dprec.Vec3) *ChassisBuilder {
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

	physicsBodyDef := scene.Physics().Engine().CreateBodyDefinition(physics.BodyDefinitionInfo{
		Mass:                   chassisMass,
		MomentOfInertia:        physics.SymmetricMomentOfInertia(chassisMomentOfInertia),
		DragFactor:             chassisDragFactor,
		AngularDragFactor:      chassisAngularDragFactor,
		RestitutionCoefficient: chassisRestitutionCoef,
		CollisionShapes: []physics.CollisionShape{
			shape.NewPlacement(
				shape.NewStaticBox(1.6, 1.4, 4.0),
				dprec.NewVec3(0.0, 0.3, -0.4),
				dprec.IdentityQuat(),
			),
		},
	})
	physicsBody := scene.Physics().CreateBody(physics.BodyInfo{
		Name:       instance.Name,
		Definition: physicsBodyDef,
		Position:   dprec.ZeroVec3(),
		Rotation:   dprec.IdentityQuat(),
		IsDynamic:  true,
	})

	gfxMesh := scene.Graphics().CreateMesh(definition.GFXMeshTemplate)

	node := game.NewNode()
	node.SetBody(physicsBody)
	node.SetMesh(gfxMesh)
	for _, modifier := range b.modifiers {
		modifier(node)
	}
	scene.Root().AppendChild(node)
	return node
}
