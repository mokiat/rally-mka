package preset

import (
	"iter"

	"github.com/mokiat/lacking/game/ecs"
)

func newECSScene(scene *ecs.Scene) *ECSScene {
	return &ECSScene{
		Scene: scene,

		nodeComponents:         ecs.NewDenseComponentSet[NodeComponent](scene),
		controlledComponents:   ecs.NewTinyComponentSet[ControlledComponent](scene),
		followCameraComponents: ecs.NewTinyComponentSet[FollowCameraComponent](scene),
		carComponents:          ecs.NewTinyComponentSet[CarComponent](scene),
		carKeyboardComponents:  ecs.NewTinyComponentSet[CarKeyboardComponent](scene),
		carGamepadComponents:   ecs.NewTinyComponentSet[CarGamepadComponent](scene),
		carMouseComponents:     ecs.NewTinyComponentSet[CarMouseComponent](scene),
	}
}

type ECSScene struct {
	*ecs.Scene

	nodeComponents         ecs.ComponentSet[NodeComponent]
	controlledComponents   ecs.ComponentSet[ControlledComponent]
	followCameraComponents ecs.ComponentSet[FollowCameraComponent]
	carComponents          ecs.ComponentSet[CarComponent]
	carKeyboardComponents  ecs.ComponentSet[CarKeyboardComponent]
	carGamepadComponents   ecs.ComponentSet[CarGamepadComponent]
	carMouseComponents     ecs.ComponentSet[CarMouseComponent]
}

func (c *ECSScene) SubscribeDelete(cb func(Entity)) *ecs.DeleteSubscription {
	return c.Scene.SubscribeDelete(func(entityID ecs.EntityID) {
		cb(c.Wrap(entityID))
	})
}

func (c *ECSScene) CreateEntity() Entity {
	return c.Wrap(c.Scene.CreateEntity())
}

func (c *ECSScene) Query(conditions ...ecs.Condition) Result {
	return Result{
		Result: c.Scene.Query(conditions...),
		ctx:    c,
	}
}

func (c *ECSScene) Wrap(entityID ecs.EntityID) Entity {
	return Entity{
		EntityID: entityID,
		scene:    c,
	}
}

func (c *ECSScene) HasNodeComponent() ecs.Condition {
	return ecs.HasComponent(c.nodeComponents)
}

func (c *ECSScene) HasControlledComponent() ecs.Condition {
	return ecs.HasComponent(c.controlledComponents)
}

func (c *ECSScene) HasFollowCameraComponent() ecs.Condition {
	return ecs.HasComponent(c.followCameraComponents)
}

func (c *ECSScene) HasCarComponent() ecs.Condition {
	return ecs.HasComponent(c.carComponents)
}

func (c *ECSScene) HasCarKeyboardControl() ecs.Condition {
	return ecs.HasComponent(c.carKeyboardComponents)
}

func (c *ECSScene) HasCarGamepadControl() ecs.Condition {
	return ecs.HasComponent(c.carGamepadComponents)
}

func (c *ECSScene) HasCarMouseControl() ecs.Condition {
	return ecs.HasComponent(c.carMouseComponents)
}

type Result struct {
	*ecs.Result
	ctx *ECSScene
}

func (r Result) Each(cb func(Entity)) {
	r.Result.Each(func(entity ecs.EntityID) {
		cb(r.ctx.Wrap(entity))
	})
}

func (r Result) Iter() iter.Seq[Entity] {
	return func(yield func(Entity) bool) {
		for entity := range r.Result.Iter() {
			if !yield(r.ctx.Wrap(entity)) {
				return
			}
		}
	}
}

type Entity struct {
	ecs.EntityID
	scene *ECSScene
}

func (e Entity) Exists() bool {
	if e.scene == nil {
		return false
	}
	return e.scene.HasEntity(e.EntityID)
}

func (e Entity) Delete() {
	e.scene.DeleteEntity(e.EntityID)
}

func (e Entity) SetNodeComponent(value NodeComponent) {
	e.scene.nodeComponents.Set(e.EntityID, value)
}

func (e Entity) UnsetNodeComponent() {
	e.scene.nodeComponents.Unset(e.EntityID)
}

func (e Entity) RefNodeComponent() *NodeComponent {
	return e.scene.nodeComponents.Ref(e.EntityID)
}

func (e Entity) SetControlledComponent(value ControlledComponent) {
	e.scene.controlledComponents.Set(e.EntityID, value)
}

func (e Entity) UnsetControlledComponent() {
	e.scene.controlledComponents.Unset(e.EntityID)
}

func (e Entity) RefControlledComponent() *ControlledComponent {
	return e.scene.controlledComponents.Ref(e.EntityID)
}

func (e Entity) SetFollowCameraComponent(value FollowCameraComponent) {
	e.scene.followCameraComponents.Set(e.EntityID, value)
}

func (e Entity) UnsetFollowCameraComponent() {
	e.scene.followCameraComponents.Unset(e.EntityID)
}

func (e Entity) RefFollowCameraComponent() *FollowCameraComponent {
	return e.scene.followCameraComponents.Ref(e.EntityID)
}

func (e Entity) SetCarComponent(value CarComponent) {
	e.scene.carComponents.Set(e.EntityID, value)
}

func (e Entity) UnsetCarComponent() {
	e.scene.carComponents.Unset(e.EntityID)
}

func (e Entity) RefCarComponent() *CarComponent {
	return e.scene.carComponents.Ref(e.EntityID)
}

func (e Entity) SetCarKeyboardComponent(value CarKeyboardComponent) {
	e.scene.carKeyboardComponents.Set(e.EntityID, value)
}

func (e Entity) UnsetCarKeyboardComponent() {
	e.scene.carKeyboardComponents.Unset(e.EntityID)
}

func (e Entity) RefCarKeyboardComponent() *CarKeyboardComponent {
	return e.scene.carKeyboardComponents.Ref(e.EntityID)
}

func (e Entity) SetCarGamepadComponent(value CarGamepadComponent) {
	e.scene.carGamepadComponents.Set(e.EntityID, value)
}

func (e Entity) UnsetCarGamepadComponent() {
	e.scene.carGamepadComponents.Unset(e.EntityID)
}

func (e Entity) RefCarGamepadComponent() *CarGamepadComponent {
	return e.scene.carGamepadComponents.Ref(e.EntityID)
}

func (e Entity) SetCarMouseComponent(value CarMouseComponent) {
	e.scene.carMouseComponents.Set(e.EntityID, value)
}

func (e Entity) UnsetCarMouseComponent() {
	e.scene.carMouseComponents.Unset(e.EntityID)
}

func (e Entity) RefCarMouseComponent() *CarMouseComponent {
	return e.scene.carMouseComponents.Ref(e.EntityID)
}
