package main

import (
	"log"
	"os"

	"github.com/mokiat/rally-mka/cmd/rff/internal/command"

	cli "gopkg.in/urfave/cli.v1"
)

func main() {
	app := cli.NewApp()
	app.Name = "rally-mka file format tool"
	app.Usage = "usage"
	app.Commands = []cli.Command{
		command.GenerateCubemap(),
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal("error: ", err.Error())
	}
}
