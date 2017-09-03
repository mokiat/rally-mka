package game

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/mokiat/go-whiskey/math"
	"github.com/mokiat/rally-mka/entities"
	"github.com/mokiat/rally-mka/render"
)

const lapCount = 3
const cameraDistance = 100.0

type Controller interface {
	InitScene()
	ResizeScene(int, int)
	UpdateScene()
	RenderScene()
	SetFrame(forward, back, left, right, brake bool)
}

func NewController() Controller {
	return &controller{
		renderer: render.NewRenderer(),
		gameMap:  entities.NewMap(),
		carMine:  entities.NewCarExtendedModel(),
	}
}

type controller struct {
	renderer *render.Renderer

	cameraPosition math.Vec3

	gameMap entities.Map
	carMine *entities.CarExtendedModel

	goForward bool
	goBack    bool
	goLeft    bool
	goRight   bool
	goBrake   bool
}

func (r *controller) InitScene() {
	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.CULL_FACE)

	gl.Enable(gl.ALPHA_TEST)
	gl.AlphaFunc(gl.GEQUAL, 0.8)

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	r.renderer.Generate()

	r.cameraPosition = math.Vec3{
		X: 0.0,
		Y: 50.0,
		Z: -cameraDistance,
	}

	const track = "assets/tracks/forest/track.m3d"
	const car = "assets/cars/suv/car.m3d"

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
		Y: 10.0,
		Z: 0.0,
	}
}

func (r *controller) ResizeScene(width, height int) {
	// gl.Viewport(0, 0, int32(width), int32(height))
	screenHalfWidth := float32(width) / float32(height)
	screenHalfHeight := float32(1.0)
	r.renderer.SetProjectionMatrix(math.PerspectiveMat4x4(-screenHalfWidth, screenHalfWidth, -screenHalfHeight, screenHalfHeight, 1.0, 5000.0))
	r.renderer.SetModelMatrix(math.IdentityMat4x4())
	r.renderer.SetViewMatrix(math.IdentityMat4x4())
}

func (r *controller) UpdateScene() {
	if r.carMine.Laps <= lapCount {
		r.carMine.Frame(r.goForward, r.goBack, r.goLeft, r.goRight, r.goBrake, r.gameMap)
	} else {
		r.carMine.Frame(false, false, false, false, false, r.gameMap)
	}

	cameraVector := r.cameraPosition.DecVec3(r.carMine.Position)
	cameraVector = cameraVector.Resize(cameraDistance)
	r.cameraPosition = r.carMine.Position.IncVec3(cameraVector)
}

func (r *controller) RenderScene() {
	gl.Clear(gl.DEPTH_BUFFER_BIT)

	r.renderer.SetViewMatrix(math.TranslationMat4x4(0.0, 0.0, -cameraDistance))

	cameraVector := r.cameraPosition.DecVec3(r.carMine.Position)
	koef := math.Vec2{
		Y: cameraVector.Y,
		X: math.Sqrt32(cameraDistance*cameraDistance - cameraVector.Y*cameraVector.Y),
	}
	koef = koef.Resize(1.0)

	var angleX float32
	if math.Abs32(koef.X) > 0.0000001 {
		angleX = ((math.Atan32(koef.Y/koef.X)/math.Pi)*180.0 + 90.0*(1.0-math.Signum32(koef.X)))
	} else {
		angleX = 90.0 * math.Signum32(koef.Y)
	}
	r.renderer.SetViewMatrix(r.renderer.ViewMatrix().MulMat4x4(math.RotationMat4x4(25.0+angleX, 1.0, 0.0, 0.0)))

	koef.X = cameraVector.Z
	koef.Y = -cameraVector.X
	koef = koef.Resize(1.0)

	var angleY float32
	if math.Abs32(koef.X) > 0.0000001 {
		angleY = ((math.Atan32(koef.Y/koef.X)/math.Pi)*180.0 - 90.0*(1.0-math.Signum32(koef.X)))
	} else {
		angleY = 0.0
	}
	r.renderer.SetViewMatrix(r.renderer.ViewMatrix().MulMat4x4(math.RotationMat4x4(angleY, 0.0, 1.0, 0.0)))

	r.renderer.SetModelMatrix(math.TranslationMat4x4(-r.carMine.Position.X, -r.carMine.Position.Y, -r.carMine.Position.Z))
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
