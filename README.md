# Rally MKA [![Go Report Card](https://goreportcard.com/badge/github.com/mokiat/rally-mka)](https://goreportcard.com/report/github.com/mokiat/rally-mka) [![Go Reference](https://pkg.go.dev/badge/github.com/mokiat/rally-mka@master.svg)](https://pkg.go.dev/github.com/mokiat/rally-mka@master)

Rally MKA is a really old game/demo of mine ported to Go for fun and as a way to experiment with new concepts.

[![Game Screenshot](preview.png)](http://mokiat.com/rally-mka/)

## User's Guide

### Browser

Use the following link to try the game in the browser:
http://mokiat.com/rally-mka/

The preferred browser is [Chrome](https://www.google.com/chrome/), which at the time of writing appears to best support WebGL2, WebAssembly and Game Controllers.

### Desktop

Check the [Releases](https://github.com/mokiat/rally-mka/releases) section for ready-to-use binaries that you can use on your computer.

The requirement is that your OS supports `OpenGL 4.6`.

### Controls

Use the following keyboard keys when playing:

- `Left Arrow` - Steer left
- `Right Arrow` - Steer right
- `Up Arrow` - Accelerate
- `Down Arrow` - Decelerate
- `a`, `s`, `d`, `w` - Rotate camera
- `q`, `e` - Zoom in/out camera
- `Enter` - Flip car

In addition, the game supports using a Game Controller.

## Developer's Guide

This section describes how to setup the project on your machine and compile it yourself.

### Prerequisites

- You need [Go 1.18](https://golang.org/dl/) or newer.
- You need the [Git LFS](https://git-lfs.github.com/) plugin. As the project contains large images and models, this is the official way on how not to clog a repository.
- Follow the instructions on the [GLFW for Go](https://github.com/go-gl/glfw) repository and make sure you can run the [GLFW examples](https://github.com/go-gl/example) on your platform.
- Make sure you have [Task](https://taskfile.dev/) installed, as this project uses Taskfiles.

### Setting Up

1. Clone the repository

   ```sh
   git clone https://github.com/mokiat/rally-mka
   cd rally-mka
   ```

1. Download Go dependencies

   ```sh
   go mod download
   ```

1. Generate game assets

   ```sh
   task pack
   ```

#### Desktop

1. Run the game

   ```sh
   task run
   ```

#### Browser

1. Generate web content

   ```sh
   task webpack
   ```

1. Build web assembly executable

   ```sh
   task wasm
   ```

1. Run HTTP server

   ```sh
   task web
   ```

1. Open http://localhost:8080

## Licensing

### Code

All source code in this project is licensed under [Apache License v2](LICENSE).

### Assets

Assets (images, models, textures, etc.) are distributed under the [Creative Commons Attribution 4.0 International](http://creativecommons.org/licenses/by/4.0/) license.

## Special Thanks

The following projects and individuals have contributed significantly to the project:

- **[The Go Team](https://go.dev/)** for making Go programming language.
- **[GLFW for Go](https://github.com/go-gl/glfw)** for making it possible to use GLFW and OpenGL in Go.
- **[LearnOpenGL](https://learnopengl.com/)** for the amazing tutorials.
- **[Poly Haven](https://polyhaven.com/)** for the excellent free images.
- **[Erin Catto](https://github.com/erincatto)** for all the presentations and articles that were used as reference for the physics engine.
- **[Bo0mer](https://github.com/Bo0mer)** for the panoramic image that was used to generate the original `city` skybox images.
- And everyone else whose repository has been used as a dependency.
