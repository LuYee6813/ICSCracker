package main

import (
	"log"
	"os"

	"ICSCracker/cmd"
)

func main() {
	cmd.PrintAsciiArt()

	app := cmd.SetupCLI()
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
