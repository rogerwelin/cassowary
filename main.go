package main

import (
	"os"

	"github.com/urfave/cli"
)

type cassowary struct {
	inputFile        os.File
	baseURL          string
	concurrencyLevel int
	requests         int
	promExport       bool
	promURL          string
}

func validateRun(c *cli.Context) {
}

func validateRunFile(c *cli.Context) {
}

func main() {
	app := cli.NewApp()
	app.Name = "cassowary"
}
