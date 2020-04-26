package collision

import (
	"github.com/mokiat/gomath/sprec"
)

func MakeLine(start, end sprec.Vec3) Line {
	return Line{
		start:  start,
		end:    end,
		length: sprec.Vec3Diff(end, start).Length(),
	}
}

type Line struct {
	start  sprec.Vec3
	end    sprec.Vec3
	length float32
}

func (l Line) Start() sprec.Vec3 {
	return l.start
}

func (l Line) End() sprec.Vec3 {
	return l.end
}

func (l Line) Length() float32 {
	return l.length
}

func (l Line) LengthSquared() float32 {
	return l.length * l.length
}

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
