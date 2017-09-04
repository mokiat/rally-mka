package main

import (
	"os"
	"path"
	"runtime"

	"github.com/go-gl/gl/v2.1/gl"
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

	const width = 800
	const height = 600
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	window, err := glfw.CreateWindow(width, height, "Rally MKA", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		panic(err)
	}

	assetsDir := path.Join(path.Dir(os.Args[0]), "assets")
	if !dirExists(assetsDir) {
		assetsDir = "assets"
	}

	controller := game.NewController(assetsDir)
	controller.InitScene()
	controller.ResizeScene(width, height)

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
