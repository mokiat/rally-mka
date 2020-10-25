package main

import (
	"log"
	"os"

	cli "github.com/urfave/cli/v2"

	"github.com/mokiat/rally-mka/cmd/rallygen/internal/level"
)

func main() {
	app := cli.NewApp()
	app.Name = "rally-mka file format tool"
	app.Usage = "usage"
	app.Commands = []*cli.Command{
		level.GenerateCommand(),
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal("error: ", err.Error())
	}
}
