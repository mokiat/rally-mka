package collision

import "github.com/mokiat/go-whiskey/math"

func NewMesh(triangles []Triangle) *Mesh {
	mesh := &Mesh{
		triangles: triangles,
	}
	mesh.evaluateCenter()
	mesh.evaluateRadius()
	return mesh
}

type Mesh struct {
	center    math.Vec3
	radius    float32
	triangles []Triangle
}

func (m *Mesh) Triangles() []Triangle {
	return m.triangles
}

func (m *Mesh) LineCollision(line Line) (bestCollision LineCollision, found bool) {
	if startDistance := line.Start().DecVec3(m.center).Length(); startDistance > line.Length()+m.radius {
		return
	}
	if endDistance := line.End().DecVec3(m.center).Length(); endDistance > line.Length()+m.radius {
		return
	}

	closestDistance := line.LengthSquared()
	for _, triangle := range m.triangles {
		if lineCollision, ok := triangle.LineCollision(line); ok {
			found = true
			distanceVector := lineCollision.Intersection().DecVec3(line.Start())
			distance := distanceVector.LengthSquared()
			if distance < closestDistance {
				closestDistance = distance
				bestCollision = lineCollision
			}
		}
	}
	return
}

func (m *Mesh) evaluateCenter() {
	m.center = math.Vec3{}
	count := 0
	for _, triangle := range m.triangles {
		m.center = m.center.IncVec3(triangle.a)
		m.center = m.center.IncVec3(triangle.b)
		m.center = m.center.IncVec3(triangle.c)
		count += 3
	}
	m.center = m.center.Div(float32(count))
}

func (m *Mesh) evaluateRadius() {
	m.radius = 0.0
	for _, triangle := range m.triangles {
		if radius := triangle.a.DecVec3(m.center).Length(); radius > m.radius {
			m.radius = radius
		}
		if radius := triangle.b.DecVec3(m.center).Length(); radius > m.radius {
			m.radius = radius
		}
		if radius := triangle.c.DecVec3(m.center).Length(); radius > m.radius {
			m.radius = radius
		}
	}
}
