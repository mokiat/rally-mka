package ecscomp

import (
	"github.com/mokiat/gomath/dprec"
	"github.com/mokiat/lacking/game/ecs"
	"github.com/mokiat/lacking/game/physics"
)

var VehicleComponentID ecs.ComponentTypeID = ecs.NewComponentTypeID()

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
}

type Chassis struct {
	Body *physics.Body
}
