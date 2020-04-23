package constraint

import (
	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/ecs"
)

type Chandelier struct {
	ecs.NilConstraint
	Entity       *ecs.Entity
	EntityAnchor sprec.Vec3
	Length       float32
	Fixture      sprec.Vec3
}

func (c Chandelier) ApplyCorrectionImpulses() {
	result := c.Calculate()
	if sprec.Abs(result.Drift) > 0.0001 {
		result.Jacobian.Apply(c.Entity)
	}
}

func (c Chandelier) ApplyCorrectionTranslations() {
	result := c.Calculate()
	if sprec.Abs(result.Drift) > 0.0001 {
		result.Jacobian.ApplyNudge(c.Entity, result.Drift)
	}
}

func (c Chandelier) Calculate() ChandelierResult {
	tranformComp := c.Entity.Transform
	anchorWorld := sprec.Vec3Sum(tranformComp.Position, sprec.QuatVec3Rotation(tranformComp.Orientation, c.EntityAnchor))
	radius := sprec.Vec3Diff(anchorWorld, tranformComp.Position)
	deltaPosition := sprec.Vec3Diff(anchorWorld, c.Fixture)
	normal := sprec.BasisXVec3()
	if deltaPosition.SqrLength() > 0.000001 {
		normal = sprec.UnitVec3(deltaPosition)
	}

	return ChandelierResult{
		Jacobian: ecs.SingleEntityJacobian{
			SlopeVelocity: sprec.NewVec3(
				normal.X,
				normal.Y,
				normal.Z,
			),
			SlopeAngularVelocity: sprec.NewVec3(
				normal.Z*radius.Y-normal.Y*radius.Z,
				normal.X*radius.Z-normal.Z*radius.X,
				normal.Y*radius.X-normal.X*radius.Y,
			),
		},
		Drift: deltaPosition.Length() - c.Length,
	}
}

type ChandelierResult struct {
	Jacobian ecs.SingleEntityJacobian
	Drift    float32
}
