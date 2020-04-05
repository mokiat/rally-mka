# Rally MKA

Rally MKA is a really old game of mine ported to Go for fun and as a way to try out new concepts.

![Game Screenshot](preview.png)

**WARNING:** This repository is under heavy refactoring, experimentation and redesign and no branch is currently guaranteed to be stable or well documented! I am hoping to get things sorted out eventually and have feature branches contain experimental stuff.

## Remarks

Being a port of a very old project of mine, some remarks are in order:

* Though some parts of the code I have rewritten entirely, others I have left as they were originally. Hence, nasty variable names and difficult to understand formulas are to be expected. The project was originally written in Delphi, back when 'clean code' wasn't something I was aware of.
* Currently, the physics of the car is whacky. It was meant to be that way when I originally wrote it - I wanted a strange rally type drifting feel to it.
* Though some corrections on the art have been made, it is quite dated.

## Getting Started

This section describes how to setup the project and run the game locally. If you are interested in just running the game, check the [Releases](https://github.com/mokiat/rally-mka/releases) section of the repository.

### Prerequisites

* You need [Go 1.14](https://golang.org/dl/) or newer.
* You need the [Git LFS](https://git-lfs.github.com/) plugin. As the project contains large images and models, this is the official way on how not to clog a repository.

### Setting Up

1. Setup GLFW 3.3 for Go

    Follow the instructions on the [GLFW for Go](https://github.com/go-gl/glfw) repository and make sure you can run the [GLFW examples](https://github.com/go-gl/example).

1. Clone the repository

    ```sh
    git clone https://github.com/mokiat/rally-mka
    cd rally-mka
    ```

    **Note:** As the project uses _go modules_, there is no need to clone or `go get` the project inside your `GOPATH`.

1. Download Go dependencies

    ```sh
    go mod download
    ```

1. Build the rallygen tool

    ```sh
    (cd cmd/rallygen && go install)
    ```

    **Note:** This tool is used to convert raw image and model resources into an optimized format for the game.

1. Generate game assets

    ```sh
    make assets
    ```

1. Run the game

    ```sh
    go run cmd/rallymka/main.go
    ```

## Licensing

### Code

All source code in this project is licensed under [Apache License v2](LICENSE).

### Assets

Assets (images, models, textures, etc.), except shaders, are distributed under the [Creative Commons Attribution 4.0 International](http://creativecommons.org/licenses/by/4.0/) license.

## Special Thanks

The following projects and individuals have contributed significantly to the project:

* **[GLFW for Go](https://github.com/go-gl/glfw)** for making it possible to use GLFW and OpenGL in Go.
* **[Bo0mer](https://github.com/Bo0mer)** for the panoramic image that was used to generate the individual `city` skybox images.
