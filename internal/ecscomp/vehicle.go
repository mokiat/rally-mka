package ecscomp

import (
	"github.com/mokiat/gomath/dprec"
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
	MaxSteeringAngle dprec.Angle
	SteeringAngle    dprec.Angle
	Acceleration     float64
	Deceleration     float64
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
	AccelerationVelocity float64
	DecelerationVelocity float64
}
