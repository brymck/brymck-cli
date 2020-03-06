package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/brymck/brymck-cli/cmd/cli/commands"
)

func main() {
	app := &cli.App{Commands: []*cli.Command{
		commands.GetCalendarCommand(),
		commands.GetRiskCommand(),
		commands.GetSecuritiesCommand(),
	}}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
