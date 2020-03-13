package scene

import "github.com/mokiat/go-whiskey/math"

type Camera struct {
	projectionMatrix math.Mat4x4
	viewMatrix       math.Mat4x4
}

func NewCamera() *Camera {
	return &Camera{
		projectionMatrix: math.IdentityMat4x4(),
		viewMatrix:       math.IdentityMat4x4(),
	}
}

func (c *Camera) SetProjectionMatrix(matrix math.Mat4x4) {
	c.projectionMatrix = matrix
}

func (c *Camera) ProjectionMatrix() math.Mat4x4 {
	return c.projectionMatrix
}

func (c *Camera) SetViewMatrix(matrix math.Mat4x4) {
	c.viewMatrix = matrix
}

func (c *Camera) ViewMatrix() math.Mat4x4 {
	return c.viewMatrix
}

func (c *Camera) InverseViewMatrix() math.Mat4x4 {
	return c.viewMatrix.Inverse()
}
