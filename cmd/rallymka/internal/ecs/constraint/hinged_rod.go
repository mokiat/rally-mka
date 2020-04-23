package constraint

import (
	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/ecs"
)

type HingedRod struct {
	ecs.NilConstraint
	First        *ecs.Entity
	FirstAnchor  sprec.Vec3
	Second       *ecs.Entity
	SecondAnchor sprec.Vec3
	Length       float32
}

func (r HingedRod) ApplyCorrectionImpulses() {
	result := r.Calculate()
	if sprec.Abs(result.Drift) > 0.0001 {
		result.Jacobian.Apply(r.First, r.Second)
	}
}

func (r HingedRod) ApplyCorrectionTranslations() {
	result := r.Calculate()
	if sprec.Abs(result.Drift) > 0.0001 {
		result.Jacobian.ApplyNudge(r.First, r.Second, result.Drift)
	}
}

func (r HingedRod) Calculate() HingedRodResult {
	firstTransformComp := r.First.Transform
	secondTransformComp := r.Second.Transform
	firstRadius := sprec.QuatVec3Rotation(firstTransformComp.Orientation, r.FirstAnchor)
	secondRadius := sprec.QuatVec3Rotation(secondTransformComp.Orientation, r.SecondAnchor)
	firstAnchorWorld := sprec.Vec3Sum(firstTransformComp.Position, firstRadius)
	secondAnchorWorld := sprec.Vec3Sum(secondTransformComp.Position, secondRadius)
	deltaPosition := sprec.Vec3Diff(secondAnchorWorld, firstAnchorWorld)
	normal := sprec.BasisXVec3()
	if deltaPosition.SqrLength() > 0.000001 {
		normal = sprec.UnitVec3(deltaPosition)
	}

	return HingedRodResult{
		Jacobian: ecs.DoubleEntityJacobian{
			SlopeVelocityFirst: sprec.NewVec3(
				-normal.X,
				-normal.Y,
				-normal.Z,
			),
			SlopeAngularVelocityFirst: sprec.NewVec3(
				-(normal.Z*firstRadius.Y - normal.Y*firstRadius.Z),
				-(normal.X*firstRadius.Z - normal.Z*firstRadius.X),
				-(normal.Y*firstRadius.X - normal.X*firstRadius.Y),
			),
			SlopeVelocitySecond: sprec.NewVec3(
				normal.X,
				normal.Y,
				normal.Z,
			),
			SlopeAngularVelocitySecond: sprec.NewVec3(
				normal.Z*secondRadius.Y-normal.Y*secondRadius.Z,
				normal.X*secondRadius.Z-normal.Z*secondRadius.X,
				normal.Y*secondRadius.X-normal.X*secondRadius.Y,
			),
		},
		Drift: deltaPosition.Length() - r.Length,
	}
}

type HingedRodResult struct {
	Jacobian ecs.DoubleEntityJacobian
	Drift    float32
}
