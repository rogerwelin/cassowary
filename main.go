package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/urfave/cli"
)

var (
	errBaseURLRequired  = errors.New("-u, --base-url argument is required")
	errFileRequired     = errors.New("-f, --file argument is required required")
	errURLRequired      = errors.New("-u, --url argument is required")
	errRequestsRequired = errors.New("-n, --requests argument is required")

	baseURL          string
	filePath         string
	concurrencyLevel int
	noOfRequests     int
	promGwURL        string
)

type cassowary struct {
	inputFile        string
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
	app.Name = "cassowary - 食火鸡"
	app.HelpName = "cassowary"
	app.UsageText = "cassowary [command] [command options] [arguments...]"
	app.EnableBashCompletion = true
	app.Usage = ""
	app.Commands = []cli.Command{
		{
			Name:  "run-file",
			Usage: "start load test",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "u, base-url",
					Usage:       "the base url (absoluteURI) to be used",
					Destination: &baseURL,
				},
				cli.IntFlag{
					Name:        "c, concurrency",
					Usage:       "number of concurrent users. defaults to 1",
					Destination: &concurrencyLevel,
				},
				cli.StringFlag{
					Name:        "f, file",
					Usage:       "specify `FILE` path, local or www, containing the url suffixes (absolute paths)",
					Destination: &filePath,
				},
				cli.StringFlag{
					Name:        "p, prompushgwurl",
					Usage:       "specify prometheus push gateway url to send metrics (optional)",
					Destination: &promGwURL,
				},
			},
			Action: validateRunFile,
		},
		{
			Name:  "run",
			Usage: "start load-test",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "u, url",
					Usage:       "the url (absoluteURI) to be used",
					Destination: &baseURL,
				},
				cli.IntFlag{
					Name:        "c, concurrency",
					Usage:       "number of concurrent users. defaults to 1",
					Destination: &concurrencyLevel,
				},
				cli.IntFlag{
					Name:        "n, requests",
					Usage:       "number of requests to perform",
					Destination: &noOfRequests,
				},
				cli.StringFlag{
					Name:        "p, prompushgwurl",
					Usage:       "specify prometheus push gateway url to send metrics (optional)",
					Destination: &promGwURL,
				},
			},
			Action: validateRun,
		},
	}

	app.CommandNotFound = func(c *cli.Context, command string) {
		err := cli.ShowAppHelp(c)
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}
