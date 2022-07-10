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
	wheelRadius            = 0.3
	wheelMass              = 20.0                                        // wheel: ~12kg; rim: ~8kg
	wheelMomentOfInertia   = wheelMass * wheelRadius * wheelRadius / 2.0 // using cylinder as approximation
	wheelDragFactor        = 0.0                                         // 0.5 * 0.3 * 0.8
	wheelAngularDragFactor = 0.0                                         // 0.5 * 0.3 * 0.8
	wheelRestitutionCoef   = 0.0
)

type WheelLocation string

const (
	FrontLeftWheelLocation  WheelLocation = "FL"
	FrontRightWheelLocation WheelLocation = "FR"
	BackLeftWheelLocation   WheelLocation = "BL"
	BackRightWheelLocation  WheelLocation = "BR"
)

func Wheel(model *resource.Model, location WheelLocation) *WheelBuilder {
	return &WheelBuilder{
		model:    model,
		location: location,
	}
}

type WheelBuilder struct {
	model     *resource.Model
	location  WheelLocation
	modifiers []func(node *game.Node)
}

func (b *WheelBuilder) WithName(name string) *WheelBuilder {
	b.modifiers = append(b.modifiers, func(node *game.Node) {
		body := node.Body()
		body.SetName(name)
	})
	return b
}

func (b *WheelBuilder) WithPosition(position sprec.Vec3) *WheelBuilder {
	b.modifiers = append(b.modifiers, func(node *game.Node) {
		body := node.Body()
		body.SetPosition(position)
	})
	return b
}

func (b *WheelBuilder) Build(scene *game.Scene) *game.Node {
	instance, found := b.model.FindMeshInstance(fmt.Sprintf("%sWheel", b.location))
	if !found {
		panic(fmt.Errorf("mesh instance %q not found", fmt.Sprintf("%sWheel", b.location)))
	}
	definition := instance.MeshDefinition
	modelNode := instance.Node

	physicsBody := scene.Physics().CreateBody()
	physicsBody.SetPosition(sprec.ZeroVec3())
	physicsBody.SetOrientation(sprec.IdentityQuat())
	physicsBody.SetMass(wheelMass)
	physicsBody.SetMomentOfInertia(physics.SymmetricMomentOfInertia(wheelMomentOfInertia))
	physicsBody.SetDragFactor(wheelDragFactor)
	physicsBody.SetAngularDragFactor(wheelAngularDragFactor)
	physicsBody.SetRestitutionCoefficient(wheelRestitutionCoef)
	physicsBody.SetCollisionShapes([]physics.CollisionShape{
		// using sphere shape at is easier to do in physics engine at the moment
		shape.Placement{
			Position:    sprec.ZeroVec3(),
			Orientation: sprec.IdentityQuat(),
			Shape:       shape.NewStaticSphere(0.3),
		},
	})

	gfxMesh := scene.Graphics().CreateMesh(definition.GFXMeshTemplate)
	gfxMesh.SetMatrix(modelNode.Matrix)

	node := game.NewNode()
	node.SetBody(physicsBody)
	node.SetMesh(gfxMesh)
	for _, modifier := range b.modifiers {
		modifier(node)
	}
	scene.Root().AppendChild(node)
	return node
}
