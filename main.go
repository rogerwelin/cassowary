package main

import (
	"os"

	"github.com/urfave/cli"
)

type Cassowary struct {
	inputFile        os.File
	baseURL          string
	concurrencyLevel int
	requests         int
	promExport       bool
	promUrl          string
}

func validateRun(c *cli.Context) {
}

func validateRunFile(c *cli.Context) {
}

func main() {
	app := cli.NewApp()
	app.Name = "cassowary"
}
