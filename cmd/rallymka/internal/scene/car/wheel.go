package car

import (
	"fmt"

	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/lacking/physics"
	"github.com/mokiat/lacking/render"
	"github.com/mokiat/lacking/resource"
	"github.com/mokiat/lacking/shape"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/ecs"
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
	modifiers []func(entity *ecs.Entity)
}

func (b *WheelBuilder) WithName(name string) *WheelBuilder {
	b.modifiers = append(b.modifiers, func(entity *ecs.Entity) {
		entity.Physics.Body.Name = name
	})
	return b
}

func (b *WheelBuilder) WithPosition(position sprec.Vec3) *WheelBuilder {
	b.modifiers = append(b.modifiers, func(entity *ecs.Entity) {
		entity.Physics.Body.Position = position
	})
	return b
}

func (b *WheelBuilder) Build(ecsManager *ecs.Manager, scene *render.Scene) *ecs.Entity {
	modelNode, _ := b.model.FindNode(fmt.Sprintf("%sWheel", b.location))

	entity := ecsManager.CreateEntity()
	entity.Physics = &ecs.PhysicsComponent{
		Body: &physics.Body{
			Position:          sprec.ZeroVec3(),
			Orientation:       sprec.IdentityQuat(),
			Mass:              wheelMass,
			MomentOfInertia:   physics.SymmetricMomentOfInertia(wheelMomentOfInertia),
			DragFactor:        wheelDragFactor,
			AngularDragFactor: wheelAngularDragFactor,
			RestitutionCoef:   wheelRestitutionCoef,
			// using sphere shape at is easier to do in physics engine at the moment
			CollisionShapes: []shape.Placement{
				{
					Position:    sprec.ZeroVec3(),
					Orientation: sprec.IdentityQuat(),
					Shape:       shape.NewStaticSphere(0.3),
				},
			},
		},
	}
	entity.Render = &ecs.RenderComponent{
		Renderable: scene.Layout().CreateRenderable(sprec.IdentityMat4(), 100.0, &resource.Model{
			Nodes: []*resource.Node{
				modelNode,
			},
		}),
	}
	for _, modifier := range b.modifiers {
		modifier(entity)
	}
	return entity
}
