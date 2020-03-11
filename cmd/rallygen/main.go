package main

import (
	"log"
	"os"

	cli "github.com/urfave/cli/v2"

	"github.com/mokiat/rally-mka/cmd/rallygen/internal/command"
	"github.com/mokiat/rally-mka/cmd/rallygen/internal/program"
)

func main() {
	app := cli.NewApp()
	app.Name = "rally-mka file format tool"
	app.Usage = "usage"
	app.Commands = []*cli.Command{
		program.GenerateCommand(),
		command.GenerateCubemap(),
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal("error: ", err.Error())
	}
}
