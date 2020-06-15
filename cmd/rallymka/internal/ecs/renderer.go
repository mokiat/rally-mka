package ecs

import (
	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/lacking/graphics"
	"github.com/mokiat/lacking/resource"
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
				r.renderModelNode(sequence, render.Matrix, node)
			}
		}
		if render.Mesh != nil {
			r.renderMesh(sequence, render.Matrix, render.Mesh)
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

func (r *Renderer) renderModelNode(sequence *graphics.Sequence, parentMatrix sprec.Mat4, node *resource.Node) {
	matrix := sprec.Mat4Prod(parentMatrix, node.Matrix)
	r.renderMesh(sequence, matrix, node.Mesh)
	for _, child := range node.Children {
		r.renderModelNode(sequence, matrix, child)
	}
}

func (r *Renderer) renderMesh(sequence *graphics.Sequence, modelMatrix sprec.Mat4, mesh *resource.Mesh) {
	for _, subMesh := range mesh.SubMeshes {
		meshItem := sequence.BeginItem()
		meshItem.Program = subMesh.Material.Shader.GeometryProgram
		meshItem.Primitive = subMesh.Primitive
		meshItem.ModelMatrix = modelMatrix
		meshItem.BackfaceCulling = subMesh.Material.BackfaceCulling
		meshItem.Metalness = subMesh.Material.Metalness
		if subMesh.Material.MetalnessTexture != nil {
			meshItem.MetalnessTwoDTexture = subMesh.Material.MetalnessTexture.GFXTexture
		}
		meshItem.Roughness = subMesh.Material.Roughness
		if subMesh.Material.RoughnessTexture != nil {
			meshItem.RoughnessTwoDTexture = subMesh.Material.RoughnessTexture.GFXTexture
		}
		meshItem.AlbedoColor = subMesh.Material.AlbedoColor
		if subMesh.Material.AlbedoTexture != nil {
			meshItem.AlbedoTwoDTexture = subMesh.Material.AlbedoTexture.GFXTexture
		}
		meshItem.NormalScale = subMesh.Material.NormalScale
		if subMesh.Material.NormalTexture != nil {
			meshItem.NormalTwoDTexture = subMesh.Material.NormalTexture.GFXTexture
		}
		meshItem.VertexArray = mesh.GFXVertexArray
		meshItem.IndexOffset = subMesh.IndexOffset
		meshItem.IndexCount = subMesh.IndexCount
		sequence.EndItem(meshItem)
	}
}

func (r *Renderer) renderSkybox(sequence *graphics.Sequence, renderSkybox *RenderSkybox) {
	for _, subMesh := range renderSkybox.Mesh.SubMeshes {
		item := sequence.BeginItem()
		item.Program = renderSkybox.Program
		item.AlbedoCubeTexture = renderSkybox.Texture
		item.VertexArray = renderSkybox.Mesh.GFXVertexArray
		item.IndexOffset = subMesh.IndexOffset
		item.IndexCount = subMesh.IndexCount
		sequence.EndItem(item)
	}
}
