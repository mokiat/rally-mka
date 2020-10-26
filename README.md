# Rally MKA

Rally MKA is a really old game of mine ported to Go for fun and as a way to try out new concepts.

![Game Screenshot](preview.png)

**WARNING:** This repository is under development, experimentation and redesign and is not guaranteed to be stable or well documented!

## Getting Started

This section describes how to setup the project and run the game locally.

### Prerequisites

* You need [Go 1.15](https://golang.org/dl/) or newer.
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

Assets (images, models, textures, etc.), except shaders, are distributed under the [Creative Commons Attribution 4.0 International](http://creativecommons.org/licenses/by/4.0/) license.

## Special Thanks

The following projects and individuals have contributed significantly to the project:

* **[GLFW for Go](https://github.com/go-gl/glfw)** for making it possible to use GLFW and OpenGL in Go.
* **[Bo0mer](https://github.com/Bo0mer)** for the panoramic image that was used to generate the individual `city` skybox images.
* **[Erin Catto](https://github.com/erincatto)** for all the presentations and articles that were used as reference.
* **[TextureHeaven](https://texturehaven.com/)** for the free images.
