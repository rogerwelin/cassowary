package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/schollz/progressbar"
	"github.com/urfave/cli"
)

var (
	errConcurrencyLevel = errors.New("Error: Concurrency level cannot be set to: 0")
	errRequestNo        = errors.New("Error: No. of request cannot be set to: 0")
)

type cassowary struct {
	fileMode         bool
	isTLS            bool
	inputFile        string
	baseURL          string
	concurrencyLevel int
	requests         int
	promExport       bool
	promURL          string
	client           *http.Client
	bar              *progressbar.ProgressBar
}

func validateRun(c *cli.Context) error {

	prometheusEnabled := false

	if c.Int("concurrency") == 0 {
		return errConcurrencyLevel
	}

	if c.Int("requests") == 0 {
		return errRequestNo
	}

	if c.String("prompushgwurl") != "" {
		prometheusEnabled = true
	}

	cass := &cassowary{
		fileMode:         false,
		baseURL:          c.String("url"),
		concurrencyLevel: c.Int("concurrency"),
		requests:         c.Int("requests"),
		promExport:       prometheusEnabled,
		promURL:          c.String("prompushgwurl"),
	}

	//fmt.Printf("%+v\n", cass)
	cass.coordinate()
	return nil
}

func validateRunFile(c *cli.Context) error {

	prometheusEnabled := false

	if c.Int("concurrency") == 0 {
		return errConcurrencyLevel
	}

	if c.Int("requests") == 0 {
		return errRequestNo
	}

	if c.String("prompushgwurl") != "" {
		prometheusEnabled = true
	}

	cass := &cassowary{
		fileMode:         true,
		inputFile:        c.String("file"),
		baseURL:          c.String("base-url"),
		concurrencyLevel: c.Int("concurrency"),
		requests:         c.Int("requests"),
		promExport:       prometheusEnabled,
		promURL:          c.String("prompushgwurl"),
	}

	cass.coordinate()
	return nil
}

func runCLI(args []string) {
	app := cli.NewApp()
	app.Name = "cassowary - 鹤鸵"
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
					Name:     "u, base-url",
					Usage:    "the base url (absoluteURI) to be used",
					Required: true,
				},
				cli.IntFlag{
					Name:     "c, concurrency",
					Usage:    "number of concurrent users",
					Required: true,
				},
				cli.StringFlag{
					Name:     "f, file",
					Usage:    "specify `FILE` path, local or www, containing the url suffixes",
					Required: true,
				},
				cli.StringFlag{
					Name:  "p, prompushgwurl",
					Usage: "specify prometheus push gateway url to send metrics (optional)",
				},
				cli.StringFlag{
					Name:  "H, header",
					Usage: "Add Arbitrary header line, eg. 'Host: www.example.com'",
				},
			},
			Action: validateRunFile,
		},
		{
			Name:  "run",
			Usage: "start load-test",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:     "u, url",
					Usage:    "the url (absoluteURI) to be used",
					Required: true,
				},
				cli.IntFlag{
					Name:     "c, concurrency",
					Usage:    "number of concurrent users",
					Required: true,
				},
				cli.IntFlag{
					Name:     "n, requests",
					Usage:    "number of requests to perform",
					Required: true,
				},
				cli.StringFlag{
					Name:  "p, prompushgwurl",
					Usage: "specify prometheus push gateway url to send metrics (optional)",
				},
				cli.StringFlag{
					Name:  "H, header",
					Usage: "Add Arbitrary header line, eg. 'Host: www.example.com'",
				},
			},
			Action: validateRun,
		},
	}

	err := app.Run(args)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}
