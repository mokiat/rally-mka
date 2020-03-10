package collision

import "github.com/mokiat/go-whiskey/math"

func MakeTriangle(a, b, c math.Vec3) Triangle {
	return Triangle{
		a:      a,
		b:      b,
		c:      c,
		normal: getNormal(a, b, c),
	}
}

type Triangle struct {
	a      math.Vec3
	b      math.Vec3
	c      math.Vec3
	normal math.Vec3
}

func (t Triangle) LineCollision(line Line) (LineCollision, bool) {
	surfaceToStart := line.Start().DecVec3(t.a)
	surfaceToEnd := line.End().DecVec3(t.a)

	topHeight := math.Vec3DotProduct(t.normal, surfaceToStart)
	bottomHeight := math.Vec3DotProduct(t.normal, surfaceToEnd)
	if topHeight < 0 || bottomHeight > 0 {
		return LineCollision{}, false
	}

	height := topHeight - bottomHeight
	if math.Abs32(height) < 0.00001 {
		return LineCollision{}, false
	}

	factor := topHeight / height
	delta := line.End().DecVec3(line.Start())
	intersection := delta.Mul(factor).IncVec3(line.Start())

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

func isInTriangle(vertex, a, b, c, normal math.Vec3) bool {
	return isCounterClockwise(a, b, vertex, normal) &&
		isCounterClockwise(b, c, vertex, normal) &&
		isCounterClockwise(c, a, vertex, normal)
}

func isCounterClockwise(a, b, c, normal math.Vec3) bool {
	evaluatedNormal := getNormal(a, b, c)
	return math.Vec3DotProduct(normal, evaluatedNormal) > 0.0
}

func getNormal(a, b, c math.Vec3) math.Vec3 {
	vector1 := a.DecVec3(c)
	vector2 := b.DecVec3(c)
	normal := math.Vec3CrossProduct(vector1, vector2)
	normal = normal.Resize(1.0)
	return normal
}
