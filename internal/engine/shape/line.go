package shape

import "github.com/mokiat/gomath/sprec"

func NewLine(a, b sprec.Vec3) Line {
	return Line{
		A: a,
		B: b,
	}
}

func RotatedLine(line Line, rotation sprec.Quat) Line {
	return Line{
		A: sprec.QuatVec3Rotation(rotation, line.A),
		B: sprec.QuatVec3Rotation(rotation, line.B),
	}
}

func TranslatedLine(line Line, translation sprec.Vec3) Line {
	return Line{
		A: sprec.Vec3Sum(line.A, translation),
		B: sprec.Vec3Sum(line.B, translation),
	}
}

type Line struct {
	A sprec.Vec3
	B sprec.Vec3
}

func (l Line) SqrLength() float32 {
	return sprec.Vec3Diff(l.B, l.A).SqrLength()
}

func (l Line) Length() float32 {
	return sprec.Vec3Diff(l.B, l.A).Length()
}
