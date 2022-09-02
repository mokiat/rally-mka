package car

import (
	"fmt"

	"github.com/mokiat/gomath/dprec"
	"github.com/mokiat/lacking/game"
	"github.com/mokiat/lacking/game/physics"
	"github.com/mokiat/lacking/util/shape"
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

func Wheel(model *game.Model, location WheelLocation) *WheelBuilder {
	return &WheelBuilder{
		model:    model,
		location: location,
	}
}

type WheelBuilder struct {
	model     *game.Model
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

func (b *WheelBuilder) WithPosition(position dprec.Vec3) *WheelBuilder {
	b.modifiers = append(b.modifiers, func(node *game.Node) {
		body := node.Body()
		body.SetPosition(position)
	})
	return b
}

func (b *WheelBuilder) Build(scene *game.Scene) *game.Node {
	node := b.model.Root().FindNode(fmt.Sprintf("%sWheel", b.location))

	physicsBodyDef := scene.Physics().Engine().CreateBodyDefinition(physics.BodyDefinitionInfo{
		Mass:                   wheelMass,
		MomentOfInertia:        physics.SymmetricMomentOfInertia(wheelMomentOfInertia),
		DragFactor:             wheelDragFactor,
		AngularDragFactor:      wheelAngularDragFactor,
		RestitutionCoefficient: wheelRestitutionCoef,
		CollisionShapes: []physics.CollisionShape{
			// using sphere shape at is easier to do in physics engine at the moment
			shape.NewPlacement(
				shape.NewStaticSphere(0.3),
				dprec.ZeroVec3(),
				dprec.IdentityQuat(),
			),
		},
	})

	physicsBody := scene.Physics().CreateBody(physics.BodyInfo{
		Name:       node.Name(),
		Definition: physicsBodyDef,
		Position:   dprec.ZeroVec3(),
		Rotation:   dprec.IdentityQuat(),
		IsDynamic:  true,
	})

	node.SetBody(physicsBody)
	for _, modifier := range b.modifiers {
		modifier(node)
	}
	return node
}
