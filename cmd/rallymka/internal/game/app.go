package game

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

const (
	initialWindowWidth  = 1024
	initialWindowHeight = 576
	appName             = "Rally MKA"
)

type Application struct{}

func (a Application) Run() error {
	if err := glfw.Init(); err != nil {
		return fmt.Errorf("failed to initialize glfw: %w", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	window, err := glfw.CreateWindow(initialWindowWidth, initialWindowHeight, appName, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to create glfw window: %w", err)
	}
	defer window.Destroy()
	window.MakeContextCurrent()
	window.SetInputMode(glfw.CursorMode, glfw.CursorHidden)

	glfw.SwapInterval(1)
	if err := gl.Init(); err != nil {
		return fmt.Errorf("failed to initialize opengl: %w", err)
	}

	assetsDir := filepath.Join(filepath.Dir(os.Args[0]), "..", "Resources", "assets")
	if !dirExists(assetsDir) {
		assetsDir = "assets"
	}

	controller := NewController(assetsDir)
	controller.OnGLInit()

	go func() {
		controller.OnInit()
		lastTick := time.Now()
		for currentTime := range time.Tick(16 * time.Millisecond) {
			controller.OnUpdate(float32(currentTime.Sub(lastTick).Seconds()))
			lastTick = currentTime
		}
	}()

	window.SetFramebufferSizeCallback(func(w *glfw.Window, width int, height int) {
		controller.OnGLResize(width, height)
	})
	fbWidth, fbHeight := window.GetFramebufferSize()
	controller.OnGLResize(fbWidth, fbHeight)

	for !window.ShouldClose() {
		isQuit := window.GetKey(glfw.KeyEscape) == glfw.Press
		if isQuit {
			break
		}

		isForward := window.GetKey(glfw.KeyUp) == glfw.Press
		isBack := window.GetKey(glfw.KeyDown) == glfw.Press
		isLeft := window.GetKey(glfw.KeyLeft) == glfw.Press
		isRight := window.GetKey(glfw.KeyRight) == glfw.Press
		isBrake := window.GetKey(glfw.KeyEnter) == glfw.Press
		isFreeze := window.GetKey(glfw.KeyF) == glfw.Press
		controller.SetFrame(isForward, isBack, isLeft, isRight, isBrake)
		controller.SetFreeze(isFreeze)
		controller.OnGLDraw()

		window.SwapBuffers()
		glfw.PollEvents()
	}
	return nil
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}
