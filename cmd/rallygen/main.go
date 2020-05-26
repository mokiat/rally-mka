package main

import (
	"log"
	"os"

	cli "github.com/urfave/cli/v2"

	"github.com/mokiat/rally-mka/cmd/rallygen/internal/level"
	"github.com/mokiat/rally-mka/cmd/rallygen/internal/mesh"
	"github.com/mokiat/rally-mka/cmd/rallygen/internal/model"
)

func main() {
	app := cli.NewApp()
	app.Name = "rally-mka file format tool"
	app.Usage = "usage"
	app.Commands = []*cli.Command{
		level.GenerateCommand(),
		model.GenerateCommand(),
		mesh.GenerateCommand(),
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal("error: ", err.Error())
	}
}
