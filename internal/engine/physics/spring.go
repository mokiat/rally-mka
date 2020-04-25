package physics

import "github.com/mokiat/gomath/sprec"

type SpringConstraint struct {
	NilConstraint
	FirstBody       *Body
	FirstBodyAnchor sprec.Vec3
	SecondBody      *Body
	Length          float32
	Stiffness       float32
}

func (c SpringConstraint) ApplyForce() {
	firstRadiusWS := sprec.QuatVec3Rotation(c.FirstBody.Orientation, c.FirstBodyAnchor)
	firstAnchorWS := sprec.Vec3Sum(c.FirstBody.Position, firstRadiusWS)
	deltaPosition := sprec.Vec3Diff(c.SecondBody.Position, firstAnchorWS)
	extendDistance := deltaPosition.Length() - c.Length
	if sprec.Abs(extendDistance) < 0.0001 {
		return
	}
	if deltaPosition.SqrLength() < 0.00001 {
		return
	}
	normal := sprec.UnitVec3(deltaPosition)

	force := sprec.Vec3Prod(normal, extendDistance*c.Stiffness)
	c.FirstBody.ApplyOffsetForce(firstRadiusWS, force)
	c.SecondBody.ApplyForce(sprec.InverseVec3(force))
}
