package ecscomp

import (
	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/lacking/game/ecs"
	"github.com/mokiat/lacking/game/physics"
	"github.com/mokiat/lacking/game/physics/solver"
)

func SetVehicle(entity *ecs.Entity, component *Vehicle) {
	entity.SetComponent(VehicleComponentID, component)
}

func GetVehicle(entity *ecs.Entity) *Vehicle {
	component := entity.Component(VehicleComponentID)
	if component == nil {
		return nil
	}
	return component.(*Vehicle)
}

type Vehicle struct {
	MaxSteeringAngle sprec.Angle
	SteeringAngle    sprec.Angle
	Acceleration     float32
	Deceleration     float32
	Recover          bool

	Chassis *Chassis
	Wheels  []*Wheel
}

type Chassis struct {
	Body *physics.Body
}

type Wheel struct {
	Body                 *physics.Body
	RotationConstraint   *solver.MatchAxis
	AccelerationVelocity float32
	DecelerationVelocity float32
}
