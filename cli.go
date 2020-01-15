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
	errNotValidURL      = errors.New("Error: Not a valid URL. Must have the following format: http{s}://{host}")
	errNotValidHeader   = errors.New("Error: Not a valid header value. Did you forget : ?")
)

type cassowary struct {
	fileMode         bool
	isTLS            bool
	inputFile        string
	baseURL          string
	concurrencyLevel int
	requests         int
	exportMetrics    bool
	// The filename which json metrics should written
	// to if `exportMetrics` is true, otherwise it defaults to "out.json".
	exportMetricsFile string
	promExport        bool
	promURL           string
	requestHeader     []string
	client            *http.Client
	bar               *progressbar.ProgressBar
}

func validateRun(c *cli.Context) error {

	prometheusEnabled := false
	var header []string

	if c.Int("concurrency") == 0 {
		return errConcurrencyLevel
	}

	if c.Int("requests") == 0 {
		return errRequestNo
	}

	if isValidURL(c.String("url")) == false {
		return errNotValidURL
	}

	if c.String("prompushgwurl") != "" {
		prometheusEnabled = true
	}

	if c.String("header") != "" {
		length := 0
		length, header = splitHeader(c.String("header"))
		if length != 2 {
			return errNotValidHeader
		}
	}

	cass := &cassowary{
		fileMode:         false,
		baseURL:          c.String("url"),
		concurrencyLevel: c.Int("concurrency"),
		requests:         c.Int("requests"),
		requestHeader:    header,
		promExport:       prometheusEnabled,
		promURL:          c.String("prompushgwurl"),
		exportMetrics:    c.Bool("json-metrics"),
		// could have a single --json-output=file.json
		// as it is not even documented (yet)
		// but don't break existing usage.
		// However, the author should make any necessary changes (even breaking ones)
		// as soon as possible before huge amount of users.
		exportMetricsFile: c.String("json-metrics-file"),
	}

	//fmt.Printf("%+v\n", cass)
	return cass.coordinate()
}

func validateRunFile(c *cli.Context) error {

	prometheusEnabled := false
	var header []string

	if c.Int("concurrency") == 0 {
		return errConcurrencyLevel
	}

	if isValidURL(c.String("url")) == false {
		return errNotValidURL
	}

	if c.String("prompushgwurl") != "" {
		prometheusEnabled = true
	}

	if c.String("header") != "" {
		length := 0
		length, header = splitHeader(c.String("header"))
		if length != 2 {
			return errNotValidHeader
		}
	}

	cass := &cassowary{
		fileMode:          true,
		inputFile:         c.String("file"),
		baseURL:           c.String("url"),
		concurrencyLevel:  c.Int("concurrency"),
		requestHeader:     header,
		promExport:        prometheusEnabled,
		promURL:           c.String("prompushgwurl"),
		exportMetrics:     c.Bool("json-metrics"),
		exportMetricsFile: c.String("json-metrics-file"),
	}

	return cass.coordinate()
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
			Usage: "start load test in spread mode",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:     "u, url",
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
					Usage: "add Arbitrary header line, eg. 'Host: www.example.com'",
				},
				cli.BoolFlag{
					Name:  "F, json-metrics",
					Usage: "outputs metrics to a json file by setting flag to true",
				},
				cli.StringFlag{
					Name:  "json-metrics-file",
					Usage: "outputs metrics to a custom json filepath, if json-metrics is set to true",
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
					Usage: "add Arbitrary header line, eg. 'Host: www.example.com'",
				},
				cli.BoolFlag{
					Name:  "F, json-metrics",
					Usage: "outputs metrics to a json file by setting flag to true",
				},
				cli.StringFlag{
					Name:  "json-metrics-file",
					Usage: "outputs metrics to a custom json filepath, if json-metrics is set to true",
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
