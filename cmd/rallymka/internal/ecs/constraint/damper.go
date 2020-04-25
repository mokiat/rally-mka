package constraint

import (
	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/ecs"
	"github.com/mokiat/rally-mka/internal/engine/physics"
)

type Damper struct {
	physics.NilConstraint
	Target               *ecs.Entity
	TargetRelativeOffset sprec.Vec3
	Entity               *ecs.Entity
	Strength             float32
}

func (c Damper) ApplyForce() {
	targetTransformComp := c.Target.Transform
	targetAnchorRelativePosition := sprec.QuatVec3Rotation(targetTransformComp.Orientation, c.TargetRelativeOffset)

	targetMotionComp := c.Target.Motion
	entityMotionComp := c.Entity.Motion

	deltaVelocity := sprec.Vec3Diff(entityMotionComp.Velocity, sprec.Vec3Sum(targetMotionComp.Velocity, sprec.Vec3Cross(targetMotionComp.AngularVelocity, targetAnchorRelativePosition)))
	if deltaVelocity.SqrLength() < 0.00001 {
		return
	}
	jacobian := sprec.UnitVec3(deltaVelocity)

	force := sprec.Vec3Prod(jacobian, deltaVelocity.Length()*c.Strength)
	targetMotionComp.ApplyOffsetForce(targetAnchorRelativePosition, sprec.InverseVec3(force))
	entityMotionComp.ApplyForce(sprec.InverseVec3(force))
}

// func (d Damper) ApplyBaumgarte() {
// 	result := d.Calculate()
// 	result.Jacobian.ApplySoft(d.Target, d.Entity, result.Force, result.Gamma)
// }

// func (d Damper) ApplyNudge() {
// 	result := d.Calculate()
// 	result.Jacobian.ApplySoft(d.Target, d.Entity, result.Force, result.Gamma)
// }

func (d Damper) Calculate() DamperResult {
	targetTransformComp := d.Target.Transform
	entityTransformComp := d.Entity.Transform
	targetRadius := sprec.QuatVec3Rotation(targetTransformComp.Orientation, d.TargetRelativeOffset)
	targetAnchor := sprec.Vec3Sum(targetTransformComp.Position, targetRadius)
	deltaPosition := sprec.Vec3Diff(entityTransformComp.Position, targetAnchor)
	normal := sprec.BasisXVec3()
	if deltaPosition.SqrLength() > 0.000001 {
		normal = sprec.UnitVec3(deltaPosition)
	}

	targetMotionComp := d.Target.Motion
	entityMotionComp := d.Entity.Motion
	targetAnchorVelocity := sprec.Vec3Sum(targetMotionComp.Velocity, sprec.Vec3Cross(targetMotionComp.AngularVelocity, targetRadius))
	deltaVelocity := sprec.Vec3Diff(entityMotionComp.Velocity, targetAnchorVelocity)

	return DamperResult{
		Jacobian: ecs.DoubleEntityJacobian{
			SlopeVelocityFirst: sprec.NewVec3(
				-normal.X,
				-normal.Y,
				-normal.Z,
			),
			SlopeAngularVelocityFirst: sprec.NewVec3(
				-(normal.Z*targetRadius.Y - normal.Y*targetRadius.Z),
				-(normal.X*targetRadius.Z - normal.Z*targetRadius.X),
				-(normal.Y*targetRadius.X - normal.X*targetRadius.Y),
			),
			SlopeVelocitySecond: sprec.NewVec3(
				normal.X,
				normal.Y,
				normal.Z,
			),
			SlopeAngularVelocitySecond: sprec.ZeroVec3(),
		},
		Drift: deltaPosition.Length(),
		Force: deltaVelocity.Length() * d.Strength,
		Gamma: 0.047,
	}
}

type DamperResult struct {
	Jacobian ecs.DoubleEntityJacobian
	Drift    float32
	Force    float32
	Gamma    float32
}
