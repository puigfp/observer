package main

import (
	"os"

	"github.com/puigfp/observer/fetch"
	"github.com/puigfp/observer/util"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "observer"
	app.Usage = ""

	app.Commands = []cli.Command{
		fetch.Command,
	}

	err := app.Run(os.Args)
	if err != nil {
		util.ErrorLogger.Fatalln(err)
	}

}
