package car

import (
	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/lacking/physics"
	"github.com/mokiat/lacking/resource"
	"github.com/mokiat/lacking/shape"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/ecs"
)

const (
	chassisRadius = 2
	chassisMass   = 1300.0 / 5.0
	// chassisMass              = 1300.0 / 10.0
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
	modifiers []func(entity *ecs.Entity)
}

func (b *ChassisBuilder) WithName(name string) *ChassisBuilder {
	b.modifiers = append(b.modifiers, func(entity *ecs.Entity) {
		entity.Physics.Body.Name = name
	})
	return b
}

func (b *ChassisBuilder) WithPosition(position sprec.Vec3) *ChassisBuilder {
	b.modifiers = append(b.modifiers, func(entity *ecs.Entity) {
		entity.Physics.Body.Position = position
	})
	return b
}

func (b *ChassisBuilder) Build(ecsManager *ecs.Manager) *ecs.Entity {
	bodyNode, _ := b.model.FindNode("Chassis")

	entity := ecsManager.CreateEntity()
	entity.Physics = &ecs.PhysicsComponent{
		Body: &physics.Body{
			Position:          sprec.ZeroVec3(),
			Orientation:       sprec.IdentityQuat(),
			Mass:              chassisMass,
			MomentOfInertia:   physics.SymmetricMomentOfInertia(chassisMomentOfInertia),
			DragFactor:        chassisDragFactor,
			AngularDragFactor: chassisAngularDragFactor,
			RestitutionCoef:   chassisRestitutionCoef,
			CollisionShapes: []shape.Placement{
				{
					Position:    sprec.NewVec3(0.0, 0.3, -0.4),
					Orientation: sprec.IdentityQuat(),
					Shape:       shape.NewStaticBox(1.6, 1.4, 4.0),
				},
			},
		},
	}
	entity.Render = &ecs.RenderComponent{
		Mesh:   bodyNode.Mesh,
		Matrix: sprec.IdentityMat4(),
	}
	for _, modifier := range b.modifiers {
		modifier(entity)
	}
	return entity
}
