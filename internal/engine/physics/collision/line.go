package collision

import "github.com/mokiat/gomath/sprec"

func NewLineCollision(intersection, normal sprec.Vec3, topHeight, bottomHeight float32) LineCollision {
	return LineCollision{
		intersection: intersection,
		normal:       normal,
		topHeight:    topHeight,
		bottomHeight: bottomHeight,
	}
}

type LineCollision struct {
	intersection sprec.Vec3
	normal       sprec.Vec3
	topHeight    float32
	bottomHeight float32
}

func (c LineCollision) Intersection() sprec.Vec3 {
	return c.intersection
}

func (c LineCollision) Normal() sprec.Vec3 {
	return c.normal
}

func (c LineCollision) TopHeight() float32 {
	return c.topHeight
}

func (c LineCollision) BottomHeight() float32 {
	return c.bottomHeight
}
