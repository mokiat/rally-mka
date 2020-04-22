package constraint

import (
	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/ecs"
)

type CopyRotation struct {
	ecs.NilConstraint
	Target       *ecs.Entity
	TargetOffset sprec.Quat
	Entity       *ecs.Entity
}

func (t CopyRotation) ApplyCorrectionImpulses() {
	result := t.Calculate()
	if sprec.Abs(result.DriftY) > 0.0001 {
		result.JacobianY.Apply(t.Target, t.Entity)
	}
	if sprec.Abs(result.DriftZ) > 0.0001 {
		result.JacobianZ.Apply(t.Target, t.Entity)
	}
}

func (t CopyRotation) ApplyCorrectionTranslations() {
	result := t.Calculate()
	if sprec.Abs(result.DriftY) > 0.0001 {
		result.JacobianY.ApplyNudge(t.Target, t.Entity, result.DriftY)
	}
	if sprec.Abs(result.DriftZ) > 0.0001 {
		result.JacobianZ.ApplyNudge(t.Target, t.Entity, result.DriftZ)
	}
}

func (t CopyRotation) Calculate() CopyRotationResult {
	targetTransformComp := t.Target.Transform
	entityTransformComp := t.Entity.Transform

	targetRadiusY := sprec.QuatProd(targetTransformComp.Orientation, t.TargetOffset).OrientationY()
	entityRadiusY := entityTransformComp.Orientation.OrientationY()
	deltaPositionY := sprec.Vec3Diff(entityRadiusY, targetRadiusY)
	normalY := sprec.BasisXVec3()
	if deltaPositionY.SqrLength() > 0.000001 {
		normalY = sprec.UnitVec3(deltaPositionY)
	}
	targetRadiusZ := sprec.QuatProd(targetTransformComp.Orientation, t.TargetOffset).OrientationZ()
	entityRadiusZ := entityTransformComp.Orientation.OrientationZ()
	deltaPositionZ := sprec.Vec3Diff(entityRadiusZ, targetRadiusZ)
	normalZ := sprec.BasisXVec3()
	if deltaPositionZ.SqrLength() > 0.000001 {
		normalZ = sprec.UnitVec3(deltaPositionZ)
	}
	return CopyRotationResult{
		JacobianY: ecs.DoubleEntityJacobian{
			SlopeVelocityFirst: sprec.ZeroVec3(),
			SlopeAngularVelocityFirst: sprec.NewVec3(
				-(normalY.Z*targetRadiusY.Y - normalY.Y*targetRadiusY.Z),
				-(normalY.X*targetRadiusY.Z - normalY.Z*targetRadiusY.X),
				-(normalY.Y*targetRadiusY.X - normalY.X*targetRadiusY.Y),
			),
			SlopeVelocitySecond: sprec.ZeroVec3(),
			SlopeAngularVelocitySecond: sprec.NewVec3(
				(normalY.Z*entityRadiusY.Y - normalY.Y*entityRadiusY.Z),
				(normalY.X*entityRadiusY.Z - normalY.Z*entityRadiusY.X),
				(normalY.Y*entityRadiusY.X - normalY.X*entityRadiusY.Y),
			),
		},
		DriftY: deltaPositionY.Length(),
		JacobianZ: ecs.DoubleEntityJacobian{
			SlopeVelocityFirst: sprec.ZeroVec3(),
			SlopeAngularVelocityFirst: sprec.NewVec3(
				-(normalZ.Z*targetRadiusZ.Y - normalZ.Y*targetRadiusZ.Z),
				-(normalZ.X*targetRadiusZ.Z - normalZ.Z*targetRadiusZ.X),
				-(normalZ.Y*targetRadiusZ.X - normalZ.X*targetRadiusZ.Y),
			),
			SlopeVelocitySecond: sprec.ZeroVec3(),
			SlopeAngularVelocitySecond: sprec.NewVec3(
				(normalZ.Z*entityRadiusZ.Y - normalZ.Y*entityRadiusZ.Z),
				(normalZ.X*entityRadiusZ.Z - normalZ.Z*entityRadiusZ.X),
				(normalZ.Y*entityRadiusZ.X - normalZ.X*entityRadiusZ.Y),
			),
		},
		DriftZ: deltaPositionZ.Length(),
	}
}

type CopyRotationResult struct {
	JacobianY ecs.DoubleEntityJacobian
	DriftY    float32
	JacobianZ ecs.DoubleEntityJacobian
	DriftZ    float32
}
