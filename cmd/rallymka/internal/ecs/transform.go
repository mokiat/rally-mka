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
		rotationQuat := sprec.RotationQuat(angle, sprec.UnitVec3(vector))
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

func NewOrientation() Orientation {
	return Orientation{
		VectorX: sprec.BasisXVec3(),
		VectorY: sprec.BasisYVec3(),
		VectorZ: sprec.BasisZVec3(),
	}
}

type Orientation struct {
	VectorX sprec.Vec3
	VectorY sprec.Vec3
	VectorZ sprec.Vec3
}

func (o *Orientation) Rotate(rotation sprec.Vec3) {
	length := rotation.Length()
	if length < 0.00001 {
		return
	}
	matrix := sprec.RotationMat4(sprec.Radians(length), rotation.X, rotation.Y, rotation.Z)
	orientationMatrix := sprec.TransformationMat4(
		o.VectorX,
		o.VectorY,
		o.VectorZ,
		sprec.ZeroVec3(),
	)
	result := sprec.Mat4Prod(matrix, orientationMatrix)
	o.VectorX = sprec.NewVec3(result.M11, result.M21, result.M31)
	o.VectorY = sprec.NewVec3(result.M12, result.M22, result.M32)
	o.VectorZ = sprec.NewVec3(result.M13, result.M23, result.M33)
}

func (o Orientation) MulVec3(vec sprec.Vec3) sprec.Vec3 {
	matrix := sprec.TransformationMat4(
		o.VectorX,
		o.VectorY,
		o.VectorZ,
		sprec.ZeroVec3(),
	)
	result := sprec.Mat4Vec4Prod(matrix, sprec.NewVec4(vec.X, vec.Y, vec.Z, 1.0))
	return result.VecXYZ()
}
