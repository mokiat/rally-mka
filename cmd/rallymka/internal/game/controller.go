package game

import (
	"math/rand"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/mokiat/go-whiskey-gl/texture"
	"github.com/mokiat/go-whiskey/math"

	"github.com/mokiat/rally-mka/cmd/rallymka/internal/entities"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/game/loading"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/render"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/scene"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/stream"
	"github.com/mokiat/rally-mka/internal/data/cubemap"
	"github.com/mokiat/rally-mka/internal/engine/graphics"
	"github.com/mokiat/rally-mka/internal/engine/resource"
)

const lapCount = 3
const cameraDistance = 8.0
const anchorDistance = 4.0

const (
	maxQueuedResources = 64
	maxResources       = 1024
	maxEvents          = 64
)

var skyboxPath = "skyboxes/city.dat"

var tracks = [...]string{
	"tracks/forest/track.m3d",
	"tracks/highway/track.m3d",
}
var cars = []string{
	"cars/hatch/car.m3d",
	"cars/suv/car.m3d",
	"cars/truck/car.m3d",
}

type View interface {
	Resize(width, height int)
	Update(elapsedSeconds float32)
	Render(pipeline *graphics.Pipeline)
}

func NewController(assetsDir string) *Controller {
	ioWorker := resource.NewWorker(maxQueuedResources)
	go ioWorker.Work() // TODO: Schedule more concurrent routines
	registry := resource.NewRegistry(ioWorker, maxResources, maxEvents)

	locator := resource.FileLocator{}
	gfxWorker := graphics.NewWorker()

	programOperator := stream.NewProgramOperator(locator, gfxWorker)
	programOperator.Register(registry)
	cubeTextureOperator := stream.NewCubeTextureOperator(locator, gfxWorker)
	cubeTextureOperator.Register(registry)
	twodTextureOperator := stream.NewTwoDTextureOperator(locator, gfxWorker)
	twodTextureOperator.Register(registry)
	meshOperator := stream.NewMeshOperator(locator, gfxWorker)
	meshOperator.Register(registry)

	return &Controller{
		lock:       &sync.Mutex{},
		assetsDir:  assetsDir,
		renderer:   render.NewRenderer(assetsDir),
		gameMap:    entities.NewMap(),
		carMine:    entities.NewCarExtendedModel(),
		camera:     scene.NewCamera(),
		stage:      scene.NewStage(),
		gfxWorker:  gfxWorker,
		glRenderer: graphics.NewRenderer(),

		registry:   registry,
		activeView: loading.NewView(registry),
	}
}

type Controller struct {
	lock      *sync.Mutex
	assetsDir string

	registry   *resource.Registry
	activeView View

	renderer *render.Renderer
	gameMap  entities.Map
	carMine  *entities.CarExtendedModel

	cameraAnchor math.Vec3
	camera       *scene.Camera
	stage        *scene.Stage

	gfxWorker  *graphics.Worker
	glRenderer *graphics.Renderer

	vertexArrayID uint32

	goForward bool
	goBack    bool
	goLeft    bool
	goRight   bool
	goBrake   bool
	goFreeze  bool
}

func (c *Controller) OnInit() {
	c.lock.Lock()
	defer c.lock.Unlock()

}

func (r *Controller) OnUpdate(elapsedSeconds float32) {
	r.registry.Update()
	r.activeView.Update(elapsedSeconds)

	pipeline := r.glRenderer.BeginPipeline()
	r.activeView.Render(pipeline)
	r.glRenderer.EndPipeline(pipeline)

	r.lock.Lock()
	defer r.lock.Unlock()

	if r.goFreeze {
		return
	}

	if r.carMine.Laps <= lapCount {
		r.carMine.Frame(elapsedSeconds, r.goForward, r.goBack, r.goLeft, r.goRight, r.goBrake, r.gameMap)
	} else {
		r.carMine.Frame(elapsedSeconds, false, false, false, false, false, r.gameMap)
	}

	// we use a camera anchor to achieve the smooth effect of a
	// camera following the car
	anchorVector := r.cameraAnchor.DecVec3(r.carMine.Position)
	anchorVector = anchorVector.Resize(anchorDistance)
	r.cameraAnchor = r.carMine.Position.IncVec3(anchorVector)

	// the following approach of creating the view matrix coordinates will fail
	// if the camera is pointing directly up or down
	cameraVectorZ := anchorVector
	cameraVectorX := math.Vec3CrossProduct(math.BaseVec3Y(), cameraVectorZ)
	cameraVectorY := math.Vec3CrossProduct(cameraVectorZ, cameraVectorX)
	r.camera.SetViewMatrix(math.Mat4x4MulMany(
		math.TranslationMat4x4(
			r.carMine.Position.X,
			r.carMine.Position.Y,
			r.carMine.Position.Z,
		),
		math.VectorMat4x4(
			cameraVectorX.Resize(1.0),
			cameraVectorY.Resize(1.0),
			cameraVectorZ.Resize(1.0),
			math.NullVec3(),
		),
		math.RotationMat4x4(-25.0, 1.0, 0.0, 0.0),
		math.TranslationMat4x4(0.0, 0.0, cameraDistance),
	))
}

func (r *Controller) OnGLInit() {
	// TODO: Move in render pipeline / sequence
	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.CULL_FACE)

	r.lock.Lock()
	defer r.lock.Unlock()

	gl.GenVertexArrays(1, &r.vertexArrayID)
	gl.BindVertexArray(r.vertexArrayID)

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	r.renderer.Generate()

	rand := rand.New(rand.NewSource(time.Now().Unix()))
	track := filepath.Join(r.assetsDir, tracks[rand.Intn(len(tracks))])
	car := filepath.Join(r.assetsDir, cars[rand.Intn(len(cars))])
	skybox := filepath.Join(r.assetsDir, skyboxPath)

	if err := r.gameMap.Load(track); err != nil {
		panic(err)
	}
	r.gameMap.Generate()

	if err := r.carMine.Load(car); err != nil {
		panic(err)
	}
	r.carMine.Generate()
	r.carMine.Position = math.Vec3{
		X: 0.0,
		Y: 0.6,
		Z: 0.0,
	}
	r.cameraAnchor = r.carMine.Position.IncCoords(0.0, 0.0, -anchorDistance)

	skyboxTexture := loadSkybox(skybox)
	r.stage.Sky = &scene.Skybox{
		Texture: skyboxTexture,
	}
}

func (r *Controller) OnGLResize(width, height int) {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.activeView.Resize(width, height)

	gl.Viewport(0, 0, int32(width), int32(height))
	screenHalfWidth := float32(width) / float32(height)
	screenHalfHeight := float32(1.0)
	r.camera.SetProjectionMatrix(math.PerspectiveMat4x4(-screenHalfWidth, screenHalfWidth, -screenHalfHeight, screenHalfHeight, 1.5, 300.0))
	r.renderer.SetProjectionMatrix(r.camera.ProjectionMatrix())
}

func (r *Controller) OnGLDraw() {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.gfxWorker.Work()

	// modern GPUs prefer that you clear all the buffers
	// and it can be faster due to cache state
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	// TODO: This should slowly be moved as part of RenderScene
	// in renderer
	gl.BindVertexArray(r.vertexArrayID)
	r.gameMap.Draw(r.renderer)
	r.carMine.Draw(r.renderer)

	// it is more optimal to render front to back
	// i.e. the skybox should be last and only unoccupied
	// fragments should be drawn
	r.renderer.RenderScene(r.stage, r.camera)

	r.glRenderer.Render()
}

func (r *Controller) SetFrame(forward, back, left, right, brake bool) {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.goForward = forward
	r.goBack = back
	r.goLeft = left
	r.goRight = right
	r.goBrake = brake
}

func (c *Controller) SetFreeze(frozen bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.goFreeze = frozen
}

func loadSkybox(path string) *texture.CubeTexture {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	decoder := &cubemap.Decoder{}
	tex, err := decoder.Decode(file)
	if err != nil {
		panic(err)
	}

	skycubeTexture := texture.DedicatedRGBACubeDataPlayground(int(tex.Dimension))
	skycubeTexture.SetData(texture.CubeSideFront, tex.Sides[cubemap.SideFront].Data)
	skycubeTexture.SetData(texture.CubeSideBack, tex.Sides[cubemap.SideBack].Data)
	skycubeTexture.SetData(texture.CubeSideLeft, tex.Sides[cubemap.SideLeft].Data)
	skycubeTexture.SetData(texture.CubeSideRight, tex.Sides[cubemap.SideRight].Data)
	skycubeTexture.SetData(texture.CubeSideTop, tex.Sides[cubemap.SideTop].Data)
	skycubeTexture.SetData(texture.CubeSideBottom, tex.Sides[cubemap.SideBottom].Data)

	skycubeTex := texture.NewCubeTexture()
	if err := skycubeTex.Allocate(); err != nil {
		panic(err)
	}
	skycubeTex.Bind()
	skycubeTex.CreateData(skycubeTexture)
	return skycubeTex
}
