package constraint

import (
	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/ecs"
)

type CopyTranslation struct {
	ecs.NilConstraint
	Target         *ecs.Entity
	Entity         *ecs.Entity
	RelativeOffset sprec.Vec3
	SkipX          bool
	SkipY          bool
	SkipZ          bool
}

func (t CopyTranslation) ApplyCorrectionImpulses() {
	result := t.Calculate()
	if sprec.Abs(result.Drift) > 0.0001 {
		result.Jacobian.Apply(t.Target, t.Entity)
	}
}

func (t CopyTranslation) ApplyCorrectionTranslations() {
	result := t.Calculate()
	if sprec.Abs(result.Drift) > 0.0001 {
		result.Jacobian.ApplyNudge(t.Target, t.Entity, result.Drift)
	}
}

func (t CopyTranslation) Calculate() CopyTranslationResult {
	targetTransformComp := t.Target.Transform
	entityTransformComp := t.Entity.Transform
	targetRadius := sprec.QuatVec3Rotation(targetTransformComp.Orientation, t.RelativeOffset)
	targetAnchor := sprec.Vec3Sum(targetTransformComp.Position, targetRadius)
	deltaPosition := sprec.Vec3Diff(entityTransformComp.Position, targetAnchor)
	if t.SkipX {
		deltaPosition = sprec.Vec3Diff(deltaPosition, sprec.Vec3Prod(targetTransformComp.Orientation.OrientationX(), sprec.Vec3Dot(deltaPosition, targetTransformComp.Orientation.OrientationX())))
	}
	if t.SkipY {
		deltaPosition = sprec.Vec3Diff(deltaPosition, sprec.Vec3Prod(targetTransformComp.Orientation.OrientationY(), sprec.Vec3Dot(deltaPosition, targetTransformComp.Orientation.OrientationY())))
	}
	if t.SkipZ {
		deltaPosition = sprec.Vec3Diff(deltaPosition, sprec.Vec3Prod(targetTransformComp.Orientation.OrientationZ(), sprec.Vec3Dot(deltaPosition, targetTransformComp.Orientation.OrientationZ())))
	}
	normal := sprec.BasisXVec3()
	if deltaPosition.SqrLength() > 0.000001 {
		normal = sprec.UnitVec3(deltaPosition)
	}

	return CopyTranslationResult{
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
	}
}

type CopyTranslationResult struct {
	Jacobian ecs.DoubleEntityJacobian
	Drift    float32
}
