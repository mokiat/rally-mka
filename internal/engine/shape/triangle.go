package shape

import "github.com/mokiat/gomath/sprec"

func NewTriangle(a, b, c sprec.Vec3) Triangle {
	return Triangle{
		A: a,
		B: b,
		C: c,
	}
}

func RotatedTriangle(triangle Triangle, rotation sprec.Quat) Triangle {
	return Triangle{
		A: sprec.QuatVec3Rotation(rotation, triangle.A),
		B: sprec.QuatVec3Rotation(rotation, triangle.B),
		C: sprec.QuatVec3Rotation(rotation, triangle.C),
	}
}

func TranslatedTriangle(triangle Triangle, translation sprec.Vec3) Triangle {
	return Triangle{
		A: sprec.Vec3Sum(triangle.A, translation),
		B: sprec.Vec3Sum(triangle.B, translation),
		C: sprec.Vec3Sum(triangle.C, translation),
	}
}

type Triangle struct {
	A sprec.Vec3
	B sprec.Vec3
	C sprec.Vec3
}

func (t Triangle) Center() sprec.Vec3 {
	a := sprec.Vec3Quot(t.A, 3.0)
	b := sprec.Vec3Quot(t.B, 3.0)
	c := sprec.Vec3Quot(t.C, 3.0)
	return sprec.Vec3Sum(sprec.Vec3Sum(a, b), c)
}

func (t Triangle) Normal() sprec.Vec3 {
	vecAB := sprec.Vec3Diff(t.B, t.A)
	vecAC := sprec.Vec3Diff(t.C, t.A)
	return sprec.UnitVec3(sprec.Vec3Cross(vecAB, vecAC))
}

func (t Triangle) Perimeter() float32 {
	lngAB := sprec.Vec3Diff(t.B, t.A).Length()
	lngBC := sprec.Vec3Diff(t.C, t.B).Length()
	lngCA := sprec.Vec3Diff(t.A, t.C).Length()
	return lngAB + lngBC + lngCA
}

func (t Triangle) Area() float32 {
	vecAB := sprec.Vec3Diff(t.B, t.A)
	vecAC := sprec.Vec3Diff(t.C, t.A)
	return sprec.Vec3Cross(vecAB, vecAC).Length() / 2.0
}

func (t Triangle) IsLookingTowards(direction sprec.Vec3) bool {
	return sprec.Vec3Dot(t.Normal(), direction) > 0.0
}
