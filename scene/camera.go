package scene

import "github.com/mokiat/go-whiskey/math"

type Camera struct {
	viewMatrix math.Mat4x4
}

func NewCamera() *Camera {
	return &Camera{
		viewMatrix: math.IdentityMat4x4(),
	}
}

func (c *Camera) SetViewMatrix(matrix math.Mat4x4) {
	c.viewMatrix = matrix
}

func (c *Camera) ViewMatrix() math.Mat4x4 {
	return c.viewMatrix
}

func (c *Camera) InverseViewMatrix() math.Mat4x4 {
	return c.viewMatrix.QuickInverse()
}
