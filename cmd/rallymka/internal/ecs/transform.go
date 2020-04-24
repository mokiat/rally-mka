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
	if angle := sprec.Radians(vector.Length()); angle > sprec.Radians(sprec.Pi/10800.0) {
		rotationQuat := sprec.RotationQuat(angle, vector)
		c.Orientation = sprec.QuatProd(rotationQuat, c.Orientation)
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
