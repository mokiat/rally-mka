package car

import (
	"fmt"

	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/ecs"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/stream"
	"github.com/mokiat/rally-mka/internal/engine/graphics"
)

const (
	tireRadius            = 0.3
	tireMass              = 20.0                                     // tire: ~12kg; rim: ~8kg
	tireMomentOfInertia   = tireMass * tireRadius * tireRadius / 2.0 // using cylinder as approximation
	tireDragFactor        = 0.0                                      // 0.5 * 0.3 * 0.8
	tireAngularDragFactor = 0.0                                      // 0.5 * 0.3 * 0.8
	tireRestitutionCoef   = 0.5
)

type TireLocation string

const (
	FrontLeftTireLocation  TireLocation = "front_left"
	FrontRightTireLocation TireLocation = "front_right"
	BackLeftTireLocation   TireLocation = "back_left"
	BackRightTireLocation  TireLocation = "back_right"
)

func Tire(program *graphics.Program, model *stream.Model, location TireLocation) *TireBuilder {
	return &TireBuilder{
		program:  program,
		model:    model,
		location: location,
	}
}

type TireBuilder struct {
	program   *graphics.Program
	model     *stream.Model
	location  TireLocation
	modifiers []func(entity *ecs.Entity)
}

func (b *TireBuilder) WithDebug(name string) *TireBuilder {
	b.modifiers = append(b.modifiers, func(entity *ecs.Entity) {
		entity.Debug = &ecs.DebugComponent{
			Name: name,
		}
	})
	return b
}

func (b *TireBuilder) WithPosition(position sprec.Vec3) *TireBuilder {
	b.modifiers = append(b.modifiers, func(entity *ecs.Entity) {
		entity.Transform.Position = position
	})
	return b
}

func (b *TireBuilder) Build(ecsManager *ecs.Manager) *ecs.Entity {
	modelNode, _ := b.model.FindNode(fmt.Sprintf("wheel_%s", b.location))

	entity := ecsManager.CreateEntity()
	entity.Transform = &ecs.TransformComponent{
		Position:    sprec.ZeroVec3(),
		Orientation: sprec.IdentityQuat(),
	}
	entity.Motion = &ecs.MotionComponent{
		Mass: tireMass,
		MomentOfInertia: sprec.NewMat3(
			tireMomentOfInertia, 0.0, 0.0,
			0.0, tireMomentOfInertia, 0.0,
			0.0, 0.0, tireMomentOfInertia,
		),
		DragFactor:        tireDragFactor,
		AngularDragFactor: tireAngularDragFactor,
	}
	entity.Collision = &ecs.CollisionComponent{
		RestitutionCoef: tireRestitutionCoef,
		CollisionShape: ecs.CylinderShape{
			Length: 0.4,
			Radius: 0.3,
		},
	}
	entity.RenderMesh = &ecs.RenderMesh{
		GeomProgram: b.program,
		Mesh:        modelNode.Mesh,
	}
	for _, modifier := range b.modifiers {
		modifier(entity)
	}
	return entity
}
