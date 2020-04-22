package constraint

import (
	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/ecs"
)

type CopyAxis struct {
	ecs.NilConstraint
	Target       *ecs.Entity
	TargetAxis   sprec.Vec3
	TargetOffset sprec.Quat
	Entity       *ecs.Entity
	EntityAxis   sprec.Vec3
}

func (a CopyAxis) ApplyCorrectionImpulses() {
	result := a.Calculate()
	if sprec.Abs(result.Drift) > 0.0001 {
		result.Jacobian.Apply(a.Target, a.Entity)
	}
}

func (a CopyAxis) ApplyCorrectionTranslations() {
	result := a.Calculate()
	if sprec.Abs(result.Drift) > 0.0001 {
		result.Jacobian.ApplyNudge(a.Target, a.Entity, result.Drift)
	}
}

func (a CopyAxis) Calculate() CopyAxisResult {
	targetTransformComp := a.Target.Transform
	entityTransformComp := a.Entity.Transform

	targetRadius := sprec.QuatProd(targetTransformComp.Orientation, a.TargetOffset).MulVec3(a.TargetAxis)
	entityRadius := entityTransformComp.Orientation.MulVec3(a.EntityAxis)
	deltaPosition := sprec.Vec3Diff(entityRadius, targetRadius)
	normal := sprec.BasisZVec3()
	if deltaPosition.SqrLength() > 0.000001 {
		normal = sprec.UnitVec3(deltaPosition)
	}
	return CopyAxisResult{
		Jacobian: ecs.DoubleEntityJacobian{
			SlopeVelocityFirst: sprec.ZeroVec3(),
			SlopeAngularVelocityFirst: sprec.NewVec3(
				-(normal.Z*targetRadius.Y - normal.Y*targetRadius.Z),
				-(normal.X*targetRadius.Z - normal.Z*targetRadius.X),
				-(normal.Y*targetRadius.X - normal.X*targetRadius.Y),
			),
			SlopeVelocitySecond: sprec.ZeroVec3(),
			SlopeAngularVelocitySecond: sprec.NewVec3(
				(normal.Z*entityRadius.Y - normal.Y*entityRadius.Z),
				(normal.X*entityRadius.Z - normal.Z*entityRadius.X),
				(normal.Y*entityRadius.X - normal.X*entityRadius.Y),
			),
		},
		Drift: deltaPosition.Length(),
	}
}

type CopyAxisResult struct {
	Jacobian ecs.DoubleEntityJacobian
	Drift    float32
}
