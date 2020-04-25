package ecs

import "github.com/mokiat/gomath/sprec"

type TransformComponent struct {
	Position    sprec.Vec3
	Orientation sprec.Quat
}

func (c *TransformComponent) Translate(offset sprec.Vec3) {
	c.Position = sprec.Vec3Sum(c.Position, offset)
}

func (c *TransformComponent) Rotate(vector sprec.Vec3) {
	if radians := vector.Length(); sprec.Abs(radians) > 0.00001 {
		c.Orientation = sprec.QuatProd(sprec.RotationQuat(sprec.Radians(radians), vector), c.Orientation)
	}
}

func (c TransformComponent) Matrix() sprec.Mat4 {
	return sprec.TransformationMat4(
		c.Orientation.OrientationX(),
		c.Orientation.OrientationY(),
		c.Orientation.OrientationZ(),
		c.Position,
	)
}
