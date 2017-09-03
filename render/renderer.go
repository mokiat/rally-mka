package render

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/mokiat/go-whiskey/math"
)

func NewRenderer() *Renderer {
	return &Renderer{
		skyboxMaterial:   newSkyboxMaterial(),
		textureMaterial:  newTextureMaterial(),
		projectionMatrix: math.IdentityMat4x4(),
		modelMatrix:      math.IdentityMat4x4(),
		viewMatrix:       math.IdentityMat4x4(),
	}
}

type Renderer struct {
	skyboxMaterial  *Material
	textureMaterial *Material

	projectionMatrix math.Mat4x4
	modelMatrix      math.Mat4x4
	viewMatrix       math.Mat4x4
}

func (r *Renderer) Generate() {
	if err := r.skyboxMaterial.Generate(); err != nil {
		panic(err)
	}
	if err := r.textureMaterial.Generate(); err != nil {
		panic(err)
	}
}

func (r *Renderer) SkyboxMaterial() *Material {
	return r.skyboxMaterial
}

func (r *Renderer) TextureMaterial() *Material {
	return r.textureMaterial
}

func (r *Renderer) SetProjectionMatrix(matrix math.Mat4x4) {
	r.projectionMatrix = matrix
}

func (r *Renderer) ProjectionMatrix() math.Mat4x4 {
	return r.projectionMatrix
}

func (r *Renderer) SetModelMatrix(matrix math.Mat4x4) {
	r.modelMatrix = matrix
}

func (r *Renderer) ModelMatrix() math.Mat4x4 {
	return r.modelMatrix
}

func (r *Renderer) SetViewMatrix(matrix math.Mat4x4) {
	r.viewMatrix = matrix
}

func (r *Renderer) ViewMatrix() math.Mat4x4 {
	return r.viewMatrix
}

func (r *Renderer) Render(mesh *Mesh, material *Material) {
	material.program.Use()

	if material.diffuseTextureLocation != -1 {
		gl.ActiveTexture(gl.TEXTURE0)
		mesh.texture.Bind()
		gl.Uniform1i(material.diffuseTextureLocation, 0)
	}

	if material.projectionMatrixLocation != -1 {
		mat := matrixToArray(r.projectionMatrix)
		gl.UniformMatrix4fv(material.projectionMatrixLocation, 1, false, &mat[0])
	}
	if material.modelMatrixLocation != -1 {
		mat := matrixToArray(r.modelMatrix)
		gl.UniformMatrix4fv(material.modelMatrixLocation, 1, false, &mat[0])
	}
	if material.viewMatrixLocation != -1 {
		mat := matrixToArray(r.viewMatrix)
		gl.UniformMatrix4fv(material.viewMatrixLocation, 1, false, &mat[0])
	}

	mesh.vertexBuffer.Bind()
	if int32(material.coordLocation) != -1 {
		gl.EnableVertexAttribArray(material.coordLocation)
		gl.VertexAttribPointer(material.coordLocation, 3, gl.FLOAT, false, mesh.vertexStride*4, gl.PtrOffset(mesh.coordOffset*4))
	}
	if int32(material.normalLocation) != -1 {
		gl.EnableVertexAttribArray(material.normalLocation)
		gl.VertexAttribPointer(material.normalLocation, 3, gl.FLOAT, false, mesh.vertexStride*4, gl.PtrOffset(mesh.normalOffset*4))
	}
	if int32(material.colorLocation) != -1 {
		gl.EnableVertexAttribArray(material.colorLocation)
		gl.VertexAttribPointer(material.colorLocation, 4, gl.FLOAT, false, mesh.vertexStride*4, gl.PtrOffset(mesh.colorOffset*4))
	}
	if int32(material.texCoordLocation) != -1 {
		gl.EnableVertexAttribArray(material.texCoordLocation)
		gl.VertexAttribPointer(material.texCoordLocation, 2, gl.FLOAT, false, mesh.vertexStride*4, gl.PtrOffset(mesh.texCoordOffset*4))
	}

	mesh.indexBuffer.Bind()
	gl.DrawElements(gl.TRIANGLES, int32(mesh.indexCount), gl.UNSIGNED_SHORT, gl.PtrOffset(0))

	if int32(material.texCoordLocation) != -1 {
		gl.DisableVertexAttribArray(material.texCoordLocation)
	}
	if int32(material.colorLocation) != -1 {
		gl.DisableVertexAttribArray(material.colorLocation)
	}
	if int32(material.normalLocation) != -1 {
		gl.DisableVertexAttribArray(material.normalLocation)
	}
	if int32(material.coordLocation) != -1 {
		gl.DisableVertexAttribArray(material.coordLocation)
	}

	gl.UseProgram(0)
}

func matrixToArray(matrix math.Mat4x4) []float32 {
	var values [16]float32
	values[0] = matrix.M11
	values[1] = matrix.M21
	values[2] = matrix.M31
	values[3] = matrix.M41

	values[4] = matrix.M12
	values[5] = matrix.M22
	values[6] = matrix.M32
	values[7] = matrix.M42

	values[8] = matrix.M13
	values[9] = matrix.M23
	values[10] = matrix.M33
	values[11] = matrix.M43

	values[12] = matrix.M14
	values[13] = matrix.M24
	values[14] = matrix.M34
	values[15] = matrix.M44
	return values[:]
}
