package collision

import "github.com/mokiat/go-whiskey/math"

func MakeLine(start, end math.Vec3) Line {
	return Line{
		start:  start,
		end:    end,
		length: end.DecVec3(start).Length(),
	}
}

type Line struct {
	start  math.Vec3
	end    math.Vec3
	length float32
}

func (l Line) Start() math.Vec3 {
	return l.start
}

func (l Line) End() math.Vec3 {
	return l.end
}

func (l Line) Length() float32 {
	return l.length
}

func (l Line) LengthSquared() float32 {
	return l.length * l.length
}

type LineCollision struct {
	intersection math.Vec3
	normal       math.Vec3
	topHeight    float32
	bottomHeight float32
}

func (c LineCollision) Intersection() math.Vec3 {
	return c.intersection
}

func (c LineCollision) Normal() math.Vec3 {
	return c.normal
}

func (c LineCollision) TopHeight() float32 {
	return c.topHeight
}

func (c LineCollision) BottomHeight() float32 {
	return c.bottomHeight
}
