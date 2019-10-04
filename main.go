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
	errBaseURLRequired  = errors.New("-u, --base-url argument is required")
	errFileRequired     = errors.New("-f, --file argument is required required")
	errURLRequired      = errors.New("-u, --url argument is required")
	errRequestsRequired = errors.New("-n, --requests argument is required")

	baseURL          string
	filePath         string
	concurrencyLevel int
	noOfRequests     int
	promExport       bool
	promGwURL        string
	httpHeader       string
)

type cassowary struct {
	fileMode         bool
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
	if baseURL == "" {
		return errURLRequired
	}
	if noOfRequests == 0 {
		return errRequestsRequired
	}
	if concurrencyLevel == 0 {
		concurrencyLevel = 1
	}
	if promGwURL == "" {
		promExport = false
	} else {
		promExport = true
	}

	cass := &cassowary{
		fileMode:         false,
		baseURL:          baseURL,
		concurrencyLevel: concurrencyLevel,
		requests:         noOfRequests,
		promExport:       promExport,
		promURL:          promGwURL,
	}

	cass.coordinate()

	return nil
}

func validateRunFile(c *cli.Context) error {
	if baseURL == "" {
		return errBaseURLRequired
	}
	if filePath == "" {
		return errFileRequired
	}
	if concurrencyLevel == 0 {
		concurrencyLevel = 1
	}
	if promGwURL == "" {
		promExport = false
	} else {
		promExport = true
	}

	cass := &cassowary{
		fileMode:         true,
		baseURL:          baseURL,
		concurrencyLevel: concurrencyLevel,
		requests:         noOfRequests,
		promExport:       promExport,
		promURL:          promGwURL,
	}

	cass.coordinate()

	return nil
}

func main() {
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
					Usage:       "specify `FILE` path, local or www, containing the url suffixes",
					Destination: &filePath,
				},
				cli.StringFlag{
					Name:        "p, prompushgwurl",
					Usage:       "specify prometheus push gateway url to send metrics (optional)",
					Destination: &promGwURL,
				},
				cli.StringFlag{
					Name:        "h, header",
					Usage:       "Add Arbitrary header line, eg. 'Host: www.example.com'",
					Destination: &httpHeader,
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
				cli.StringFlag{
					Name:        "h, header",
					Usage:       "Add Arbitrary header line, eg. 'Host: www.example.com'",
					Destination: &httpHeader,
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
