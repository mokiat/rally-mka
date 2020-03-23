package render

import (
	"github.com/go-gl/gl/v4.1-core/gl"

	"github.com/mokiat/go-whiskey/math"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/scene"
)

func NewRenderer(assetsDir string) *Renderer {
	return &Renderer{
		textureMaterial:  newTextureMaterial(),
		projectionMatrix: math.IdentityMat4x4(),
		modelMatrix:      math.IdentityMat4x4(),
		viewMatrix:       math.IdentityMat4x4(),
		matrixCache:      make([]float32, 16),
	}
}

type Renderer struct {
	textureMaterial *Material

	projectionMatrix math.Mat4x4
	modelMatrix      math.Mat4x4
	viewMatrix       math.Mat4x4

	matrixCache []float32
}

func (r *Renderer) Generate() {
	if err := r.textureMaterial.Generate(); err != nil {
		panic(err)
	}
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
		gl.UniformMatrix4fv(material.projectionMatrixLocation, 1, false, r.matrixToArray(r.projectionMatrix))
	}
	if material.modelMatrixLocation != -1 {
		gl.UniformMatrix4fv(material.modelMatrixLocation, 1, false, r.matrixToArray(r.modelMatrix))
	}
	if material.viewMatrixLocation != -1 {
		gl.UniformMatrix4fv(material.viewMatrixLocation, 1, false, r.matrixToArray(r.viewMatrix))
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

func (r *Renderer) RenderScene(stage *scene.Stage, camera *scene.Camera) {
	r.SetViewMatrix(camera.InverseViewMatrix())
}

// NOTE: Use this method only as short-lived function argument
// subsequent calls will reuse the same float32 array
func (r *Renderer) matrixToArray(matrix math.Mat4x4) *float32 {
	r.matrixCache[0] = matrix.M11
	r.matrixCache[1] = matrix.M21
	r.matrixCache[2] = matrix.M31
	r.matrixCache[3] = matrix.M41

	r.matrixCache[4] = matrix.M12
	r.matrixCache[5] = matrix.M22
	r.matrixCache[6] = matrix.M32
	r.matrixCache[7] = matrix.M42

	r.matrixCache[8] = matrix.M13
	r.matrixCache[9] = matrix.M23
	r.matrixCache[10] = matrix.M33
	r.matrixCache[11] = matrix.M43

	r.matrixCache[12] = matrix.M14
	r.matrixCache[13] = matrix.M24
	r.matrixCache[14] = matrix.M34
	r.matrixCache[15] = matrix.M44
	return &r.matrixCache[0]
}
