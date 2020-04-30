package collision

import (
	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/rally-mka/internal/engine/shape"
)

func NewMesh(triangles []Triangle) *Mesh {
	mesh := &Mesh{
		triangles: triangles,
	}
	mesh.evaluateCenter()
	mesh.evaluateRadius()
	return mesh
}

type Mesh struct {
	center    sprec.Vec3
	radius    float32
	triangles []Triangle
}

func (m *Mesh) Triangles() []Triangle {
	return m.triangles
}

func (m *Mesh) LineCollision(line shape.Line) (bestCollision LineCollision, found bool) {
	if startDistance := sprec.Vec3Diff(m.center, line.A).Length(); startDistance > line.Length()+m.radius {
		return
	}
	if endDistance := sprec.Vec3Diff(m.center, line.B).Length(); endDistance > line.Length()+m.radius {
		return
	}

	closestDistance := line.SqrLength()
	for _, triangle := range m.triangles {
		if lineCollision, ok := triangle.LineCollision(line); ok {
			found = true
			distanceVector := sprec.Vec3Diff(lineCollision.intersection, line.A)
			distance := distanceVector.SqrLength()
			if distance < closestDistance {
				closestDistance = distance
				bestCollision = lineCollision
			}
		}
	}
	return
}

func (m *Mesh) evaluateCenter() {
	m.center = sprec.ZeroVec3()
	count := 0
	for _, triangle := range m.triangles {
		m.center = sprec.Vec3Sum(m.center, triangle.a)
		m.center = sprec.Vec3Sum(m.center, triangle.b)
		m.center = sprec.Vec3Sum(m.center, triangle.c)
		count += 3
	}
	m.center = sprec.Vec3Quot(m.center, float32(count))
}

func (m *Mesh) evaluateRadius() {
	m.radius = 0.0
	for _, triangle := range m.triangles {
		if radius := sprec.Vec3Diff(m.center, triangle.a).Length(); radius > m.radius {
			m.radius = radius
		}
		if radius := sprec.Vec3Diff(m.center, triangle.b).Length(); radius > m.radius {
			m.radius = radius
		}
		if radius := sprec.Vec3Diff(m.center, triangle.c).Length(); radius > m.radius {
			m.radius = radius
		}
	}
}
