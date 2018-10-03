package game

import (
	"math/rand"
	"path/filepath"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/mokiat/go-whiskey/math"
	"github.com/mokiat/rally-mka/entities"
	"github.com/mokiat/rally-mka/render"
	"github.com/mokiat/rally-mka/scene"
)

const lapCount = 3
const cameraDistance = 6.0

var tracks = [...]string{
	"tracks/forest/track.m3d",
	"tracks/highway/track.m3d",
}
var cars = []string{
	"cars/hatch/car.m3d",
	"cars/suv/car.m3d",
	"cars/truck/car.m3d",
}

type Controller interface {
	InitScene()
	ResizeScene(int, int)
	UpdateScene()
	RenderScene()
	SetFrame(forward, back, left, right, brake bool)
}

func NewController(assetsDir string) Controller {
	return &controller{
		assetsDir: assetsDir,
		renderer:  render.NewRenderer(),
		gameMap:   entities.NewMap(),
		carMine:   entities.NewCarExtendedModel(),
		camera:    scene.NewCamera(),
	}
}

type controller struct {
	assetsDir string

	renderer *render.Renderer
	gameMap  entities.Map
	carMine  *entities.CarExtendedModel

	cameraAnchor math.Vec3
	camera       *scene.Camera

	goForward bool
	goBack    bool
	goLeft    bool
	goRight   bool
	goBrake   bool
}

func (r *controller) InitScene() {
	var vertexArrayID uint32
	gl.GenVertexArrays(1, &vertexArrayID)
	gl.BindVertexArray(vertexArrayID)

	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.CULL_FACE)

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	r.renderer.Generate()

	r.cameraAnchor = math.Vec3{
		X: 0.0,
		Y: 3.0,
		Z: -cameraDistance,
	}

	rand := rand.New(rand.NewSource(time.Now().Unix()))
	track := filepath.Join(r.assetsDir, tracks[rand.Intn(len(tracks))])
	car := filepath.Join(r.assetsDir, cars[rand.Intn(len(cars))])

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
}

func (r *controller) ResizeScene(width, height int) {
	gl.Viewport(0, 0, int32(width), int32(height))
	screenHalfWidth := float32(width) / float32(height)
	screenHalfHeight := float32(1.0)
	r.renderer.SetProjectionMatrix(math.PerspectiveMat4x4(-screenHalfWidth, screenHalfWidth, -screenHalfHeight, screenHalfHeight, 1.0, 300.0))
}

func (r *controller) UpdateScene() {
	if r.carMine.Laps <= lapCount {
		r.carMine.Frame(r.goForward, r.goBack, r.goLeft, r.goRight, r.goBrake, r.gameMap)
	} else {
		r.carMine.Frame(false, false, false, false, false, r.gameMap)
	}

	// we use a camera anchor to achieve the smooth effect of a
	// camera following the car
	anchorVector := r.cameraAnchor.DecVec3(r.carMine.Position)
	anchorVector = anchorVector.Resize(cameraDistance)
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

func (r *controller) RenderScene() {
	r.renderer.SetViewMatrix(r.camera.InverseViewMatrix())

	gl.Clear(gl.DEPTH_BUFFER_BIT)
	r.gameMap.Draw(r.renderer)
	r.carMine.DrawMe(r.renderer)
}

func (r *controller) SetFrame(forward, back, left, right, brake bool) {
	r.goForward = forward
	r.goBack = back
	r.goLeft = left
	r.goRight = right
	r.goBrake = brake
}
