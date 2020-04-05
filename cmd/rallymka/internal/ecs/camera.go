package ecs

import "github.com/mokiat/gomath/sprec"

type Camera struct {
	projectionMatrix sprec.Mat4
	viewMatrix       sprec.Mat4
}

func NewCamera() *Camera {
	return &Camera{
		projectionMatrix: sprec.IdentityMat4(),
		viewMatrix:       sprec.IdentityMat4(),
	}
}

func (c *Camera) SetProjectionMatrix(matrix sprec.Mat4) {
	c.projectionMatrix = matrix
}

func (c *Camera) ProjectionMatrix() sprec.Mat4 {
	return c.projectionMatrix
}

func (c *Camera) SetViewMatrix(matrix sprec.Mat4) {
	c.viewMatrix = matrix
}

func (c *Camera) ViewMatrix() sprec.Mat4 {
	return c.viewMatrix
}

func (c *Camera) InverseViewMatrix() sprec.Mat4 {
	return sprec.InverseMat4(c.viewMatrix)
}
