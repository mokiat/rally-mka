package collision

import (
	"github.com/mokiat/gomath/sprec"
)

func MakeTriangle(a, b, c sprec.Vec3) Triangle {
	return Triangle{
		a:      a,
		b:      b,
		c:      c,
		normal: getNormal(a, b, c),
	}
}

type Triangle struct {
	a      sprec.Vec3
	b      sprec.Vec3
	c      sprec.Vec3
	normal sprec.Vec3
}

func (t Triangle) BoudingSphereRadius() float32 {
	center := t.Center()
	lngA := sprec.Vec3Diff(t.a, center).Length()
	lngB := sprec.Vec3Diff(t.b, center).Length()
	lngC := sprec.Vec3Diff(t.c, center).Length()
	if lngA > lngB && lngA > lngC {
		return lngA
	}
	if lngB > lngC {
		return lngB
	}
	return lngC
}

func (t Triangle) Surface() float32 {
	lngAB := sprec.Vec3Diff(t.b, t.a).Length()
	lngBC := sprec.Vec3Diff(t.c, t.b).Length()
	lngCA := sprec.Vec3Diff(t.a, t.c).Length()
	semiperimeter := (lngAB + lngBC + lngCA) / 2.0
	return sprec.Sqrt(semiperimeter * (semiperimeter - lngAB) * (semiperimeter - lngBC) * (semiperimeter - lngCA))
}

func (t Triangle) A() sprec.Vec3 {
	return t.a
}

func (t Triangle) B() sprec.Vec3 {
	return t.b
}

func (t Triangle) C() sprec.Vec3 {
	return t.c
}

func (t Triangle) Normal() sprec.Vec3 {
	return t.normal
}

func (t Triangle) Center() sprec.Vec3 {
	return sprec.Vec3Quot(sprec.Vec3Sum(sprec.Vec3Sum(t.a, t.b), t.c), 3)
}

func (t Triangle) Contains(point sprec.Vec3) bool {
	return isInTriangle(point, t.a, t.b, t.c, t.normal)
}

func (t Triangle) LineCollision(line Line) (LineCollision, bool) {
	surfaceToStart := sprec.Vec3Diff(line.start, t.a)
	surfaceToEnd := sprec.Vec3Diff(line.end, t.a)

	topHeight := sprec.Vec3Dot(t.normal, surfaceToStart)
	bottomHeight := sprec.Vec3Dot(t.normal, surfaceToEnd)
	if topHeight < 0 || bottomHeight > 0 {
		return LineCollision{}, false
	}

	height := topHeight - bottomHeight
	if sprec.Abs(height) < 0.00001 {
		return LineCollision{}, false
	}

	factor := topHeight / height
	delta := sprec.Vec3Diff(line.end, line.start)
	intersection := sprec.Vec3Sum(sprec.Vec3Prod(delta, factor), line.start)

	if !isInTriangle(intersection, t.a, t.b, t.c, t.normal) {
		return LineCollision{}, false
	}
	return LineCollision{
		intersection: intersection,
		normal:       t.normal,
		topHeight:    topHeight,
		bottomHeight: bottomHeight,
	}, true
}

func isInTriangle(vertex, a, b, c, normal sprec.Vec3) bool {
	return isCounterClockwise(a, b, vertex, normal) &&
		isCounterClockwise(b, c, vertex, normal) &&
		isCounterClockwise(c, a, vertex, normal)
}

func isCounterClockwise(a, b, c, normal sprec.Vec3) bool {
	evaluatedNormal := getNormal(a, b, c)
	return sprec.Vec3Dot(normal, evaluatedNormal) > 0.0
}

func getNormal(a, b, c sprec.Vec3) sprec.Vec3 {
	vector1 := sprec.Vec3Diff(a, c)
	vector2 := sprec.Vec3Diff(b, c)
	direction := sprec.Vec3Cross(vector1, vector2)
	return sprec.ResizedVec3(direction, 1.0)
}
