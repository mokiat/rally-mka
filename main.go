package main

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/mokiat/rally-mka/game"
)

func init() {
	runtime.LockOSThread()
}

func main() {
	if err := glfw.Init(); err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	const width = 1024
	const height = 576
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	window, err := glfw.CreateWindow(width, height, "Rally MKA", nil, nil)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()
	window.MakeContextCurrent()
	window.SetInputMode(glfw.CursorMode, glfw.CursorHidden)

	glfw.SwapInterval(1)
	if err := gl.Init(); err != nil {
		panic(err)
	}

	assetsDir := filepath.Join(filepath.Dir(os.Args[0]), "..", "Resources", "assets")
	if !dirExists(assetsDir) {
		assetsDir = "assets"
	}

	controller := game.NewController(assetsDir)
	controller.InitScene()

	window.SetFramebufferSizeCallback(func(w *glfw.Window, width int, height int) {
		controller.ResizeScene(width, height)
	})
	fbWidth, fbHeight := window.GetFramebufferSize()
	controller.ResizeScene(fbWidth, fbHeight)

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
		controller.SetFrame(isForward, isBack, isLeft, isRight, isBrake)

		controller.UpdateScene()
		controller.RenderScene()

		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}
