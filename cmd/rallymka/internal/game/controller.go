package game

import (
	"sync"

	"github.com/go-gl/gl/v4.1-core/gl"

	"github.com/mokiat/rally-mka/cmd/rallymka/internal/game/input"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/game/loading"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/game/simulation"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/stream"
	"github.com/mokiat/rally-mka/internal/engine/graphics"
	"github.com/mokiat/rally-mka/internal/engine/resource"
)

const (
	maxQueuedResources = 64
	maxResources       = 1024
	maxEvents          = 64
)

type View interface {
	Load()
	IsAvailable() bool
	Unload()

	Open()
	Close()

	Resize(width, height int)
	Update(elapsedSeconds float32, actions input.ActionSet)
	Render(pipeline *graphics.Pipeline)
}

func NewController(assetsDir string) *Controller {
	ioWorker := resource.NewWorker(maxQueuedResources)
	go ioWorker.Work() // TODO: Schedule more concurrent routines
	registry := resource.NewRegistry(ioWorker, maxResources, maxEvents)
	gfxWorker := graphics.NewWorker()

	locator := resource.FileLocator{}
	programOperator := stream.NewProgramOperator(locator, gfxWorker)
	programOperator.Register(registry)
	cubeTextureOperator := stream.NewCubeTextureOperator(locator, gfxWorker)
	cubeTextureOperator.Register(registry)
	twodTextureOperator := stream.NewTwoDTextureOperator(locator, gfxWorker)
	twodTextureOperator.Register(registry)
	modelOperator := stream.NewModelOperator(locator, gfxWorker)
	modelOperator.Register(registry)
	meshOperator := stream.NewMeshOperator(locator, gfxWorker)
	meshOperator.Register(registry)
	levelOperator := stream.NewLevelOperator(locator, gfxWorker)
	levelOperator.Register(registry)

	return &Controller{
		lock: &sync.Mutex{},

		input:          &input.Tracker{},
		loadingView:    loading.NewView(registry),
		simulationView: simulation.NewView(registry),
		activeView:     nil,

		registry:        registry,
		gfxResizeEvents: make(chan windowSize, 32),
		gfxWorker:       gfxWorker,
		gfxRenderer:     graphics.NewRenderer(),
	}
}

type Controller struct {
	lock *sync.Mutex

	windowSize     windowSize
	input          *input.Tracker
	activeView     View
	loadingView    View
	simulationView View

	registry        *resource.Registry
	gfxResizeEvents chan windowSize
	gfxWorker       *graphics.Worker
	gfxRenderer     *graphics.Renderer
}

func (c *Controller) Input() *input.Tracker {
	return c.input
}

func (c *Controller) OnInit() {
	c.loadingView.Load()
	c.simulationView.Load()
}

func (c *Controller) OnUpdate(elapsedSeconds float32) {
	c.processEvents()
	c.pickView()
	c.registry.Update()

	if c.activeView != nil {
		c.activeView.Update(elapsedSeconds, c.input.Get())
		pipeline := c.gfxRenderer.BeginPipeline()
		c.activeView.Render(pipeline)
		c.gfxRenderer.EndPipeline(pipeline)
	}
}

func (*Controller) OnGLInit() {
	// TODO: Move in render pipeline / sequence / item
	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.CULL_FACE)
}

func (c *Controller) OnGLResize(width, height int) {
	gl.Viewport(0, 0, int32(width), int32(height))
	c.gfxResizeEvents <- windowSize{
		Width:  width,
		Height: height,
	}
}

func (c *Controller) OnGLDraw() {
	c.gfxWorker.Work()
	c.gfxRenderer.Render()
}

func (c *Controller) processEvents() {
	for event, ok := c.pollResizeEvent(); ok; event, ok = c.pollResizeEvent() {
		c.windowSize = event
	}
}

func (c *Controller) pollResizeEvent() (windowSize, bool) {
	select {
	case event := <-c.gfxResizeEvents:
		return event, true
	default:
		return windowSize{}, false
	}
}

func (c *Controller) pickView() {
	switch c.activeView {
	case nil:
		if c.loadingView.IsAvailable() {
			c.changeView(c.loadingView)
		}
	case c.loadingView:
		if c.simulationView.IsAvailable() {
			c.changeView(c.simulationView)
		}
	}
}

func (c *Controller) changeView(view View) {
	if c.activeView != nil {
		c.activeView.Close()
	}
	c.activeView = view
	if c.activeView != nil {
		c.activeView.Open()
		c.activeView.Resize(c.windowSize.Width, c.windowSize.Height)
	}
}

type windowSize struct {
	Width  int
	Height int
}
