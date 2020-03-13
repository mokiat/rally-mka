package main

import (
	"log"
	"os"

	cli "github.com/urfave/cli/v2"

	"github.com/mokiat/rally-mka/cmd/rallygen/internal/command"
	"github.com/mokiat/rally-mka/cmd/rallygen/internal/cubetex"
	"github.com/mokiat/rally-mka/cmd/rallygen/internal/level"
	"github.com/mokiat/rally-mka/cmd/rallygen/internal/mesh"
	"github.com/mokiat/rally-mka/cmd/rallygen/internal/model"
	"github.com/mokiat/rally-mka/cmd/rallygen/internal/program"
	"github.com/mokiat/rally-mka/cmd/rallygen/internal/twodtex"
)

func main() {
	app := cli.NewApp()
	app.Name = "rally-mka file format tool"
	app.Usage = "usage"
	app.Commands = []*cli.Command{
		level.GenerateCommand(),
		program.GenerateCommand(),
		model.GenerateCommand(),
		mesh.GenerateCommand(),
		cubetex.GenerateCommand(),
		twodtex.GenerateCommand(),
		command.GenerateCubemap(),
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal("error: ", err.Error())
	}
}
