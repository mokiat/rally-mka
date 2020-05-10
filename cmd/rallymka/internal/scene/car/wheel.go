package car

import (
	"fmt"

	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/ecs"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/stream"
	"github.com/mokiat/rally-mka/internal/engine/graphics"
	"github.com/mokiat/rally-mka/internal/engine/physics"
	"github.com/mokiat/rally-mka/internal/engine/shape"
)

const (
	wheelRadius            = 0.3
	wheelMass              = 20.0                                        // wheel: ~12kg; rim: ~8kg
	wheelMomentOfInertia   = wheelMass * wheelRadius * wheelRadius / 2.0 // using cylinder as approximation
	wheelDragFactor        = 0.0                                         // 0.5 * 0.3 * 0.8
	wheelAngularDragFactor = 0.0                                         // 0.5 * 0.3 * 0.8
	wheelRestitutionCoef   = 0.5
)

type WheelLocation string

const (
	FrontLeftWheelLocation  WheelLocation = "front_left"
	FrontRightWheelLocation WheelLocation = "front_right"
	BackLeftWheelLocation   WheelLocation = "back_left"
	BackRightWheelLocation  WheelLocation = "back_right"
)

func Wheel(program *graphics.Program, model *stream.Model, location WheelLocation) *WheelBuilder {
	return &WheelBuilder{
		program:  program,
		model:    model,
		location: location,
	}
}

type WheelBuilder struct {
	program   *graphics.Program
	model     *stream.Model
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

func (b *WheelBuilder) Build(ecsManager *ecs.Manager) *ecs.Entity {
	modelNode, _ := b.model.FindNode(fmt.Sprintf("wheel_%s", b.location))

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
		GeomProgram: b.program,
		Mesh:        modelNode.Mesh,
		Matrix:      sprec.IdentityMat4(),
	}
	for _, modifier := range b.modifiers {
		modifier(entity)
	}
	return entity
}