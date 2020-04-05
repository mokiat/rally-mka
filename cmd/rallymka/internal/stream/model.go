package stream

import (
	"fmt"

	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/rally-mka/internal/data/asset"
	"github.com/mokiat/rally-mka/internal/engine/graphics"
	"github.com/mokiat/rally-mka/internal/engine/resource"
)

const modelResourceType = "model"

func GetModel(registry *resource.Registry, name string) ModelHandle {
	return ModelHandle{
		Handle: registry.Type(modelResourceType).Resource(name),
	}
}

type ModelHandle struct {
	resource.Handle
}

func (h ModelHandle) Get() *Model {
	return h.Handle.Get().(*Model)
}

func (h ModelHandle) IsAvailable() bool {
	return h.Handle.IsAvailable() && h.Get().IsAvailable()
}

type Model struct {
	Nodes  []*Node
	meshes []*Mesh
}

func (m Model) FindNode(name string) (*Node, bool) {
	for _, node := range m.Nodes {
		if node.Name == name {
			return node, true
		}
		if child, found := node.FindNode(name); found {
			return child, true
		}
	}
	return nil, false
}

func (m Model) IsAvailable() bool {
	for _, mesh := range m.meshes {
		if !mesh.IsAvailable() {
			return false
		}
	}
	return true
}

type Node struct {
	Name     string
	Matrix   sprec.Mat4
	Mesh     *Mesh
	Parent   *Node
	Children []*Node
}

func (n Node) FindNode(name string) (*Node, bool) {
	for _, node := range n.Children {
		if node.Name == name {
			return node, true
		}
		if child, found := node.FindNode(name); found {
			return child, true
		}
	}
	return nil, false
}

func NewModelOperator(locator resource.Locator, gfxWorker *graphics.Worker) *ModelOperator {
	return &ModelOperator{
		locator:   locator,
		gfxWorker: gfxWorker,
	}
}

type ModelOperator struct {
	locator   resource.Locator
	gfxWorker *graphics.Worker
}

func (o *ModelOperator) Register(registry *resource.Registry) {
	registry.RegisterType(modelResourceType, o)
}

func (o *ModelOperator) Allocate(registry *resource.Registry, name string) (resource.Resource, error) {
	in, err := o.locator.Open("assets", "models", name)
	if err != nil {
		return nil, fmt.Errorf("failed to open model asset %q: %w", name, err)
	}
	defer in.Close()

	modelAsset, err := asset.NewModelDecoder().Decode(in)
	if err != nil {
		return nil, fmt.Errorf("failed to decode model asset %q: %w", name, err)
	}

	model := &Model{}

	meshes := make([]*Mesh, len(modelAsset.Meshes))
	for i, meshAsset := range modelAsset.Meshes {
		mesh, err := AllocateMesh(registry, o.gfxWorker, &meshAsset)
		if err != nil {
			return nil, fmt.Errorf("failed to allocate mesh: %w", err)
		}
		meshes[i] = mesh
	}
	model.meshes = meshes

	nodes := make([]*Node, len(modelAsset.Nodes))
	for i := range nodes {
		nodes[i] = &Node{}
	}
	rootNodes := make([]*Node, 0)
	for i, nodeAsset := range modelAsset.Nodes {
		if nodeAsset.ParentIndex != -1 {
			nodes[i].Parent = nodes[nodeAsset.ParentIndex]
			nodes[nodeAsset.ParentIndex].Children = append(nodes[nodeAsset.ParentIndex].Children, nodes[i])
		} else {
			rootNodes = append(rootNodes, nodes[i])
		}
		nodes[i].Name = nodeAsset.Name
		nodes[i].Matrix = floatArrayToMatrix(nodeAsset.Matrix)
		nodes[i].Mesh = meshes[nodeAsset.MeshIndex]
	}
	model.Nodes = rootNodes

	return model, nil
}

func (o *ModelOperator) Release(registry *resource.Registry, resource resource.Resource) error {
	model := resource.(*Model)

	for _, mesh := range model.meshes {
		if err := ReleaseMesh(registry, o.gfxWorker, mesh); err != nil {
			return fmt.Errorf("failed to release mesh: %w", err)
		}
	}

	return nil
}

// TODO: Math library should provide similar method
func floatArrayToMatrix(values [16]float32) sprec.Mat4 {
	var result sprec.Mat4
	result.M11 = values[0]
	result.M21 = values[1]
	result.M31 = values[2]
	result.M41 = values[3]

	result.M12 = values[4]
	result.M22 = values[5]
	result.M32 = values[6]
	result.M42 = values[7]

	result.M13 = values[8]
	result.M23 = values[9]
	result.M33 = values[10]
	result.M43 = values[11]

	result.M14 = values[12]
	result.M24 = values[13]
	result.M34 = values[14]
	result.M44 = values[15]
	return result
}
