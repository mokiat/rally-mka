package car

import (
	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/lacking/game/ecs"
	"github.com/mokiat/lacking/game/graphics"
	"github.com/mokiat/lacking/game/physics"
	"github.com/mokiat/lacking/resource"
	"github.com/mokiat/lacking/shape"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/ecscomp"
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
	modifiers []func(entity *ecs.Entity)
}

func (b *ChassisBuilder) WithName(name string) *ChassisBuilder {
	b.modifiers = append(b.modifiers, func(entity *ecs.Entity) {
		physicsComponent := ecscomp.GetPhysics(entity)
		physicsComponent.Body.SetName(name)
	})
	return b
}

func (b *ChassisBuilder) WithPosition(position sprec.Vec3) *ChassisBuilder {
	b.modifiers = append(b.modifiers, func(entity *ecs.Entity) {
		physicsComponent := ecscomp.GetPhysics(entity)
		physicsComponent.Body.SetPosition(position)
	})
	return b
}

func (b *ChassisBuilder) Build(ecsScene *ecs.Scene, gfxScene graphics.Scene, physicsScene *physics.Scene) *ecs.Entity {
	// bodyNode, _ := b.model.FindNode("Chassis")

	physicsBody := physicsScene.CreateBody()
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

	entity := ecsScene.CreateEntity()
	ecscomp.SetPhysics(entity, &ecscomp.Physics{
		Body: physicsBody,
	})
	// TODO
	// ecscomp.SetRender(entity, &ecscomp.Render{
	// 	Renderable: scene.Layout().CreateRenderable(sprec.IdentityMat4(), 100.0, &resource.Model{
	// 		Nodes: []*resource.Node{
	// 			bodyNode,
	// 		},
	// 	}),
	// })

	for _, modifier := range b.modifiers {
		modifier(entity)
	}
	return entity
}
