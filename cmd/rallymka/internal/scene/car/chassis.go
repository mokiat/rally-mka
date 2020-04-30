package car

import (
	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/ecs"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/stream"
	"github.com/mokiat/rally-mka/internal/engine/graphics"
	"github.com/mokiat/rally-mka/internal/engine/physics"
)

const (
	chassisRadius = 2
	chassisMass   = 1300.0 / 5.0
	// chassisMass              = 1300.0 / 10.0
	chassisMomentOfInertia   = chassisMass * chassisRadius * chassisRadius / 2.0
	chassisDragFactor        = 0.0 // 0.5 * 6.8 * 1.0
	chassisAngularDragFactor = 0.0 // 0.5 * 6.8 * 1.0
	chassisRestitutionCoef   = 0.3
)

func Chassis(program *graphics.Program, model *stream.Model) *ChassisBuilder {
	return &ChassisBuilder{
		program: program,
		model:   model,
	}
}

type ChassisBuilder struct {
	program   *graphics.Program
	model     *stream.Model
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
	bodyNode, _ := b.model.FindNode("body")

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
			CollisionShape: physics.BoxShape{
				MinX: -0.8,
				MaxX: 0.8,
				MinY: -0.4,
				MaxY: 0.8,
				MinZ: -2.2,
				MaxZ: 1.5,
			},
		},
	}
	entity.Render = &ecs.RenderComponent{
		GeomProgram: b.program,
		Mesh:        bodyNode.Mesh,
		Matrix:      sprec.IdentityMat4(),
	}
	for _, modifier := range b.modifiers {
		modifier(entity)
	}
	return entity
}
