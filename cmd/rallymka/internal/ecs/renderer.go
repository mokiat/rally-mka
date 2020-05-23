package ecs

import (
	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/lacking/graphics"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/stream"
)

func NewRenderer(ecsManager *Manager) *Renderer {
	return &Renderer{
		ecsManager: ecsManager,
	}
}

type Renderer struct {
	ecsManager *Manager
}

func (r *Renderer) Render(sequence *graphics.Sequence) {
	r.renderRender(sequence)
	r.renderRenderSkyboxes(sequence)
}

func (r *Renderer) renderRender(sequence *graphics.Sequence) {
	for _, entity := range r.ecsManager.Entities() {
		render := entity.Render
		if render == nil {
			continue
		}
		if entity.Physics != nil {
			body := entity.Physics.Body
			render.Matrix = sprec.TransformationMat4(
				body.Orientation.OrientationX(),
				body.Orientation.OrientationY(),
				body.Orientation.OrientationZ(),
				body.Position,
			)
		}
		if render.Model != nil {
			for _, node := range render.Model.Nodes {
				r.renderModelNode(sequence, render.GeomProgram, render.Matrix, node)
			}
		}
		if render.Mesh != nil {
			r.renderMesh(sequence, render.GeomProgram, render.Matrix, render.Mesh)
		}
	}
}

func (r *Renderer) renderRenderSkyboxes(sequence *graphics.Sequence) {
	for _, entity := range r.ecsManager.Entities() {
		renderSkyboxComp := entity.RenderSkybox
		if renderSkyboxComp == nil {
			continue
		}
		r.renderSkybox(sequence, renderSkyboxComp)
	}
}

func (r *Renderer) renderModelNode(sequence *graphics.Sequence, program *graphics.Program, parentMatrix sprec.Mat4, node *stream.Node) {
	matrix := sprec.Mat4Prod(parentMatrix, node.Matrix)
	r.renderMesh(sequence, program, matrix, node.Mesh)
	for _, child := range node.Children {
		r.renderModelNode(sequence, program, matrix, child)
	}
}

func (r *Renderer) renderMesh(sequence *graphics.Sequence, program *graphics.Program, modelMatrix sprec.Mat4, mesh *stream.Mesh) {
	for _, subMesh := range mesh.SubMeshes {
		meshItem := sequence.BeginItem()
		meshItem.Program = program
		meshItem.ModelMatrix = modelMatrix
		if subMesh.DiffuseTexture != nil {
			meshItem.DiffuseTexture = subMesh.DiffuseTexture.Get()
		}
		meshItem.VertexArray = mesh.VertexArray
		meshItem.IndexCount = subMesh.IndexCount
		sequence.EndItem(meshItem)
	}
}

func (r *Renderer) renderSkybox(sequence *graphics.Sequence, renderSkybox *RenderSkybox) {
	for _, subMesh := range renderSkybox.Mesh.SubMeshes {
		item := sequence.BeginItem()
		item.Program = renderSkybox.Program
		item.SkyboxTexture = renderSkybox.Texture
		item.VertexArray = renderSkybox.Mesh.VertexArray
		item.IndexCount = subMesh.IndexCount
		sequence.EndItem(item)
	}
}
