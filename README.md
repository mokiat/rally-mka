# Rally MKA [![Go Report Card](https://goreportcard.com/badge/github.com/mokiat/rally-mka)](https://goreportcard.com/report/github.com/mokiat/rally-mka) [![Go Reference](https://pkg.go.dev/badge/github.com/mokiat/rally-mka@master.svg)](https://pkg.go.dev/github.com/mokiat/rally-mka@master)

Rally MKA is a really old game of mine ported to Go for fun and as a way to try new concepts out.

![Game Screenshot](preview.png)

## User's Guide

Check the [Releases](https://github.com/mokiat/rally-mka/releases) section for ready-to-use binaries.

## Developer's Guide

This section describes how to setup the project on your machine and compile it yourself.

### Prerequisites

* You need [Go 1.16](https://golang.org/dl/) or newer.
* You need the [Git LFS](https://git-lfs.github.com/) plugin. As the project contains large images and models, this is the official way on how not to clog a repository.
* Follow the instructions on the [GLFW for Go](https://github.com/go-gl/glfw) repository and make sure you can run the [GLFW examples](https://github.com/go-gl/example) on your platform.

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
    make assets
    ```

1. Run the game

    ```sh
    make play
    ```

## Licensing

### Code

All source code in this project is licensed under [Apache License v2](LICENSE).

### Assets

Assets (images, models, textures, etc.), are distributed under the [Creative Commons Attribution 4.0 International](http://creativecommons.org/licenses/by/4.0/) license.

## Special Thanks

The following projects and individuals have contributed significantly to the project:

* **[GLFW for Go](https://github.com/go-gl/glfw)** for making it possible to use GLFW and OpenGL in Go.
* **[LearnOpenGL](https://learnopengl.com/)** for the amazing tutorials.
* **[TextureHeaven](https://texturehaven.com/)** for the excellent free images.
* **[Erin Catto](https://github.com/erincatto)** for all the presentations and articles that were used as reference.
* **[Bo0mer](https://github.com/Bo0mer)** for the panoramic image that was used to generate the original `city` skybox images.
