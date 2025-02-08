# Rally MKA [![Go Report Card](https://goreportcard.com/badge/github.com/mokiat/rally-mka)](https://goreportcard.com/report/github.com/mokiat/rally-mka) [![Go Reference](https://pkg.go.dev/badge/github.com/mokiat/rally-mka@master.svg)](https://pkg.go.dev/github.com/mokiat/rally-mka@master)

Rally MKA is a really old game/demo of mine ported to Go for fun and as a way to experiment with new concepts. It is also a showcase for the [lacking](https://github.com/mokiat/lacking) game engine.

[![Game Screenshot](preview.png)](https://mokiat.itch.io/rally-mka)


## User's Guide

### Browser

You can play the game on [itch.io](https://mokiat.itch.io/rally-mka).

The preferred browser is [Chrome](https://www.google.com/chrome/), which at the time of writing appears to best support WebGL2, WebAssembly and Game Controllers.

### Desktop

Check the [Releases](https://github.com/mokiat/rally-mka/releases) section for ready-to-use binaries that you can use on your computer.

The requirement is that your OS supports `OpenGL 4.6`.

## Developer's Guide

This section describes how to setup the project on your machine and compile it yourself.

### Prerequisites

- You need [Go 1.23](https://golang.org/dl/) or newer.
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

1. Run an HTTP server

   ```sh
   task web
   ```

1. Open http://localhost:8080

## Licensing

### Code

All source code in this project is licensed under [Apache License v2](LICENSE).

### Assets

Assets (images, models, textures, etc.) are distributed under the [Creative Commons Attribution 4.0 International](http://creativecommons.org/licenses/by/4.0/) license.
