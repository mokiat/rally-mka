package car

import (
	"fmt"

	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/lacking/game/ecs"
	"github.com/mokiat/lacking/game/graphics"
	"github.com/mokiat/lacking/game/physics"
	"github.com/mokiat/lacking/resource"
	"github.com/mokiat/lacking/shape"
	"github.com/mokiat/rally-mka/internal/ecscomp"
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
		physicsComponent := ecscomp.GetPhysics(entity)
		physicsComponent.Body.SetName(name)
	})
	return b
}

func (b *WheelBuilder) WithPosition(position sprec.Vec3) *WheelBuilder {
	b.modifiers = append(b.modifiers, func(entity *ecs.Entity) {
		physicsComponent := ecscomp.GetPhysics(entity)
		physicsComponent.Body.SetPosition(position)
	})
	return b
}

func (b *WheelBuilder) Build(ecsScene *ecs.Scene, gfxScene *graphics.Scene, physicsScene *physics.Scene) *ecs.Entity {
	instance, found := b.model.FindMeshInstance(fmt.Sprintf("%sWheel", b.location))
	if !found {
		panic(fmt.Errorf("mesh instance %q not found", fmt.Sprintf("%sWheel", b.location)))
	}
	definition := instance.MeshDefinition
	modelNode := instance.Node
	// modelNode, _ := b.model.FindNode(fmt.Sprintf("%sWheel", b.location))

	physicsBody := physicsScene.CreateBody()
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

	entity := ecsScene.CreateEntity()
	ecscomp.SetPhysics(entity, &ecscomp.Physics{
		Body: physicsBody,
	})

	gfxMesh := gfxScene.CreateMesh(definition.GFXMeshTemplate)
	gfxMesh.SetPosition(modelNode.Matrix.Translation())
	// TODO: Set Rotation
	// TODO: Set Scale

	ecscomp.SetRender(entity, &ecscomp.Render{
		Mesh: gfxMesh,
	})
	for _, modifier := range b.modifiers {
		modifier(entity)
	}
	return entity
}
