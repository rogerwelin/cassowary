package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/fatih/color"
	"github.com/rogerwelin/cassowary/pkg/client"
	"github.com/urfave/cli/v2"
)

var (
	version             = "dev"
	errConcurrencyLevel = errors.New("error: Concurrency level cannot be set to: 0")
	errRequestNo        = errors.New("error: No. of request cannot be set to: 0")
	errNotValidURL      = errors.New("error: Not a valid URL. Must have the following format: http{s}://{host}")
	errNotValidHeader   = errors.New("error: Not a valid header value. Did you forget : ?")
	errDurationValue    = errors.New("error: Duration cannot be set to 0 or negative")
	errSineAmplitude    = errors.New("error: Sine amplitude must be between 0.0 and 1.0")
	errSineFrequency    = errors.New("error: Sine frequency must be greater than 0")
	errSinePreset       = errors.New("error: Sine preset must be one of: gentle, normal, aggressive")
)

func outPutResults(metrics client.ResultMetrics) {
	printf(summaryTable,
		color.CyanString(fmt.Sprintf("%.2f", metrics.TCPStats.TCPMean)),
		color.CyanString(fmt.Sprintf("%.2f", metrics.TCPStats.TCPMedian)),
		color.CyanString(fmt.Sprintf("%.2f", metrics.TCPStats.TCP95p)),
		color.CyanString(fmt.Sprintf("%.2f", metrics.ProcessingStats.ServerProcessingMean)),
		color.CyanString(fmt.Sprintf("%.2f", metrics.ProcessingStats.ServerProcessingMedian)),
		color.CyanString(fmt.Sprintf("%.2f", metrics.ProcessingStats.ServerProcessing95p)),
		color.CyanString(fmt.Sprintf("%.2f", metrics.ContentStats.ContentTransferMean)),
		color.CyanString(fmt.Sprintf("%.2f", metrics.ContentStats.ContentTransferMedian)),
		color.CyanString(fmt.Sprintf("%.2f", metrics.ContentStats.ContentTransfer95p)),
		color.CyanString(strconv.Itoa(metrics.TotalRequests)),
		color.CyanString(strconv.Itoa(metrics.FailedRequests)),
		color.CyanString(fmt.Sprintf("%.2f", metrics.DNSMedian)),
		color.CyanString(fmt.Sprintf("%.2f", metrics.RequestsPerSecond)),
	)
}

func outPutJSON(fileName string, metrics client.ResultMetrics) error {
	if fileName == "" {
		// default filename for json metrics output.
		fileName = "out.json"
	}
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	return enc.Encode(metrics)
}

func applySinePreset(c *cli.Context, preset string) error {
	switch preset {
	case "gentle":
		// Gentle: Small variation, slow ramp-up/down
		c.Set("sine-frequency", "0.25") // Quarter cycle - just ramp up
		c.Set("sine-amplitude", "0.3")  // ±30% variation
		fmt.Fprintf(os.Stderr, "📈 Applied 'gentle' preset: frequency=0.25, amplitude=0.3 (gradual ramp-up to +30%)\n")
		
	case "normal":
		// Normal: Moderate variation, half cycle  
		c.Set("sine-frequency", "0.5")  // Half cycle - ramp up then down
		c.Set("sine-amplitude", "0.5")  // ±50% variation
		fmt.Fprintf(os.Stderr, "📈 Applied 'normal' preset: frequency=0.5, amplitude=0.5 (ramp up/down with ±50% variation)\n")
		
	case "aggressive":
		// Aggressive: High variation, fast oscillations
		c.Set("sine-frequency", "2.0")  // Two complete cycles
		c.Set("sine-amplitude", "0.7")  // ±70% variation
		fmt.Fprintf(os.Stderr, "📈 Applied 'aggressive' preset: frequency=2.0, amplitude=0.7 (rapid oscillations with ±70% variation)\n")
		
	default:
		return errSinePreset
	}
	
	return nil
}

func showSineCurveCompatibilityWarnings(c *cli.Context, duration int) {
	concurrency := c.Int("concurrency")
	requests := c.Int("requests")
	amplitude := c.Float64("sine-amplitude")
	frequency := c.Float64("sine-frequency")
	fileMode := c.String("file") != ""
	
	warnings := make([]string, 0)
	
	// High concurrency warning
	if concurrency > 1 {
		warnings = append(warnings, fmt.Sprintf("⚠️  High concurrency (%d users) may reduce sine wave timing precision", concurrency))
		warnings = append(warnings, "   Consider using -c 1 for most accurate sine wave patterns")
	}
	
	// Very short duration warning  
	if duration > 0 && duration < 10 {
		warnings = append(warnings, fmt.Sprintf("⚠️  Short duration (%ds) may not show clear sine wave pattern", duration))
		warnings = append(warnings, "   Consider using -d 30 or higher for better wave visibility")
	}
	
	// File mode warning
	if fileMode {
		warnings = append(warnings, "⚠️  File mode with sine curve uses URL paths in sequence")
		warnings = append(warnings, "   Each request in the sine pattern will cycle through your URL list")
	}
	
	// Extreme parameter combinations
	if amplitude >= 0.8 && frequency >= 3.0 {
		warnings = append(warnings, fmt.Sprintf("⚠️  High amplitude (%.1f) + high frequency (%.1f) creates dramatic load spikes", amplitude, frequency))
		warnings = append(warnings, "   This may cause rapid oscillation between very high and very low rates")
	}
	
	if amplitude >= 0.9 {
		warnings = append(warnings, fmt.Sprintf("⚠️  Very high amplitude (%.1f) creates extreme rate variations", amplitude))
		warnings = append(warnings, "   Rate will vary from near-zero to almost 2x baseline - ensure your system can handle this")
	}
	
	// Very high frequency warning
	if frequency >= 5.0 {
		warnings = append(warnings, fmt.Sprintf("⚠️  High frequency (%.1f) creates rapid load oscillations", frequency))
		if duration > 0 {
			cycleTime := float64(duration) / frequency
			warnings = append(warnings, fmt.Sprintf("   Each cycle lasts only %.1f seconds", cycleTime))
		}
	}
	
	// Very low amplitude warning
	if amplitude <= 0.1 {
		warnings = append(warnings, fmt.Sprintf("⚠️  Low amplitude (%.1f) creates minimal rate variation", amplitude))
		warnings = append(warnings, "   You may not notice much difference from constant load testing")
	}
	
	// Duration vs requests mismatch for very high rates
	if duration > 0 {
		baselineRate := float64(requests) / float64(duration)
		maxRate := baselineRate * (1.0 + amplitude)
		if maxRate > 100 {
			warnings = append(warnings, fmt.Sprintf("⚠️  Peak rate may reach %.0f req/s during sine wave peaks", maxRate))
			warnings = append(warnings, "   Ensure your target system can handle these high request rates")
		}
	}
	
	// Print warnings if any
	if len(warnings) > 0 {
		fmt.Fprintf(os.Stderr, "\n🌊 Sine Curve Compatibility Warnings:\n")
		for _, warning := range warnings {
			fmt.Fprintf(os.Stderr, "%s\n", warning)
		}
		fmt.Fprintf(os.Stderr, "\n💡 Use --sine-preview to visualize the pattern before running\n")
		fmt.Fprintf(os.Stderr, "⏱️  Waiting 3 seconds before continuing... (Ctrl+C to abort)\n\n")
		
		// Give user time to read warnings and cancel if needed
		time.Sleep(3 * time.Second)
	}
}

func runLoadTest(c *client.Cassowary) error {
	// Handle sine curve preview modes
	if c.SinePreviewLive {
		if !c.SineCurve {
			// Automatically enable sine curve if live preview is requested
			c.SineCurve = true
		}
		return c.ShowLiveSineVisualization()
	}
	
	if c.SinePreview {
		if !c.SineCurve {
			// Automatically enable sine curve if preview is requested
			c.SineCurve = true
		}
		c.ShowSineVisualization()
		return nil
	}

	metrics, err := c.Coordinate()
	if err != nil {
		return err
	}

	if !c.DisableTerminalOutput {
		outPutResults(metrics)
	}

	if c.ExportMetrics {
		return outPutJSON(c.ExportMetricsFile, metrics)
	}

	if c.PromExport {
		err := c.PushPrometheusMetrics(metrics)
		if err != nil {
			return err
		}
	}

	if c.Cloudwatch {
		session, err := session.NewSession()
		if err != nil {
			return err
		}

		svc := cloudwatch.New(session)
		_, err = c.PutCloudwatchMetrics(svc, metrics)
		if err != nil {
			return err
		}
	}

	return nil
}

func validateCLI(c *cli.Context) error {
	prometheusEnabled := false
	var header []string
	var httpMethod string
	var data []byte
	duration := 0
	var urlSuffixes []string
	fileMode := false

	if c.Int("concurrency") == 0 {
		return errConcurrencyLevel
	}

	if c.Int("requests") == 0 {
		return errRequestNo
	}

	if c.String("duration") != "" {
		var err error
		duration, err = strconv.Atoi(c.String("duration"))
		if err != nil {
			return err
		}
		if duration <= 0 {
			return errDurationValue
		}
	}

	// Validate sine curve parameters
	if c.Bool("sine-curve") || c.String("sine-preset") != "" {
		// Handle preset mode
		if preset := c.String("sine-preset"); preset != "" {
			if err := applySinePreset(c, preset); err != nil {
				return err
			}
			// Auto-enable sine curve when using preset
			c.Set("sine-curve", "true")
		}
		
		if c.Float64("sine-amplitude") < 0.0 || c.Float64("sine-amplitude") > 1.0 {
			return errSineAmplitude
		}
		if c.Float64("sine-frequency") <= 0 {
			return errSineFrequency
		}
		
		// Compatibility warnings for sine curve
		showSineCurveCompatibilityWarnings(c, duration)
	}

	if !client.IsValidURL(c.String("url")) {
		return errNotValidURL
	}

	if c.String("prompushgwurl") != "" {
		prometheusEnabled = true
	}

	if len(c.StringSlice("header")) != 0 {
		allHeaders := c.StringSlice("header")
		for _, hdr := range allHeaders {
			thisLen, thisHdr := client.SplitHeader(hdr)
			if thisLen != 2 {
				return errNotValidHeader
			}
			header = append(header, thisHdr...)
		}
	}

	if c.String("file") != "" {
		var err error
		urlSuffixes, err = readLocalRemoteFile(c.String("file"))
		if err != nil {
			return nil
		}
		fileMode = true
	}

	if c.String("postfile") != "" {
		httpMethod = "POST"
		fileData, err := readFile(c.String("postfile"))
		if err != nil {
			return err
		}
		data = fileData
	} else if c.String("putfile") != "" {
		httpMethod = "PUT"
		fileData, err := readFile(c.String("putfile"))
		if err != nil {
			return err
		}
		data = fileData
	} else if c.String("patchfile") != "" {
		httpMethod = "PATCH"
		fileData, err := readFile(c.String("patchfile"))
		if err != nil {
			return err
		}
		data = fileData
	} else {
		httpMethod = "GET"
	}

	tlsConfig := new(tls.Config)
	if c.Bool("insecure") {
		tlsConfig.InsecureSkipVerify = true
	}

	if c.String("renegotiation") != "" {
		switch option := c.String("renegotiation"); option {
		case "never":
			tlsConfig.Renegotiation = tls.RenegotiateNever
		case "once":
			tlsConfig.Renegotiation = tls.RenegotiateOnceAsClient
		case "freely":
			tlsConfig.Renegotiation = tls.RenegotiateFreelyAsClient
		default:
			return fmt.Errorf("invalid renegotiation option specified: %q", option)
		}
	}

	if c.String("ca") != "" {
		pemCerts, err := os.ReadFile(c.String("ca"))
		if err != nil {
			return err
		}
		ca := x509.NewCertPool()
		if !ca.AppendCertsFromPEM(pemCerts) {
			return fmt.Errorf("failed to read CA from PEM")
		}
		tlsConfig.RootCAs = ca
	}

	if c.String("cert") != "" && c.String("key") != "" {
		cert, err := tls.LoadX509KeyPair(c.String("cert"), c.String("key"))
		if err != nil {
			return err
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
	}

	cass := &client.Cassowary{
		FileMode:              fileMode,
		DisableTerminalOutput: c.Bool("silent"),
		BaseURL:               c.String("url"),
		ConcurrencyLevel:      c.Int("concurrency"),
		Requests:              c.Int("requests"),
		RequestHeader:         header,
		Duration:              duration,
		PromExport:            prometheusEnabled,
		TLSConfig:             tlsConfig,
		PromURL:               c.String("prompushgwurl"),
		Cloudwatch:            c.Bool("cloudwatch"),
		Boxplot:               c.Bool("boxplot"),
		Histogram:             c.Bool("histogram"),
		ExportMetrics:         c.Bool("json-metrics"),
		RawOutput:             c.Bool("raw-output"),
		ExportMetricsFile:     c.String("json-metrics-file"),
		DisableKeepAlive:      c.Bool("disable-keep-alive"),
		Timeout:               c.Int("timeout"),
		HTTPMethod:            httpMethod,
		URLPaths:              urlSuffixes,
		Data:                  data,
		SineCurve:             c.Bool("sine-curve"),
		SineFrequency:         c.Float64("sine-frequency"),
		SineAmplitude:         c.Float64("sine-amplitude"),
		SinePreview:           c.Bool("sine-preview"),
		SinePreviewLive:       c.Bool("sine-preview-live"),
	}

	return runLoadTest(cass)
}

func runCLI(args []string) {
	app := cli.NewApp()
	app.Name = "cassowary - 學名"
	setCustomCLITemplate(app)
	app.HelpName = "cassowary"
	app.UsageText = "cassowary [command] [command options] [arguments...]"
	app.EnableBashCompletion = true
	app.Usage = "Modern cross-platform HTTP load-testing tool"
	app.Version = version
	app.Commands = []*cli.Command{
		{
			Name:  "run",
			Usage: "start load-test",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "u",
					Aliases:  []string{"url"},
					Usage:    "the url (absoluteURI) to be used",
					Required: true,
				},
				&cli.IntFlag{
					Name:    "c",
					Aliases: []string{"concurrency"},
					Usage:   "number of concurrent users",
					Value:   1,
				},
				&cli.IntFlag{
					Name:    "n",
					Aliases: []string{"requests"},
					Usage:   "number of requests to perform",
					Value:   1,
				},
				&cli.StringFlag{
					Name:    "f",
					Aliases: []string{"file"},
					Usage:   "file-slurp mode: specify `FILE` path, local or www, containing the url suffixes",
				},
				&cli.StringFlag{
					Name:    "d",
					Aliases: []string{"duration"},
					Usage:   "set the duration in seconds of the load test (example: do 100 requests in a duration of 30s)",
				},
				&cli.IntFlag{
					Name:    "t",
					Aliases: []string{"timeout"},
					Usage:   "http client timeout",
					Value:   5,
				},
				&cli.StringFlag{
					Name:    "p",
					Aliases: []string{"prompushgwurl"},
					Usage:   "specify prometheus push gateway url to send metrics (optional)",
				},
				&cli.StringSliceFlag{
					Name:    "H",
					Aliases: []string{"header"},
					Usage:   "add arbitrary header, eg. 'Host: www.example.com'",
				},
				&cli.BoolFlag{
					Name:    "C",
					Aliases: []string{"cloudwatch"},
					Usage:   "enable to send metrics to AWS Cloudwatch",
				},
				&cli.BoolFlag{
					Name:    "R",
					Aliases: []string{"raw-output"},
					Usage:   "enable to export raw per-request metrics",
				},
				&cli.BoolFlag{
					Name:    "s",
					Aliases: []string{"silent"},
					Usage:   "Do not show progress and do not print results on terminal (useful for other output types)",
				},
				&cli.BoolFlag{
					Name:    "F",
					Aliases: []string{"json-metrics"},
					Usage:   "outputs metrics to a json file by setting flag to true",
				},
				&cli.BoolFlag{
					Name:    "b",
					Aliases: []string{"boxplot"},
					Usage:   "enable to generate a boxplot as png",
				},
				&cli.BoolFlag{
					Aliases: []string{"histogram"},
					Usage:   "enable to generate a histogram as png",
				},
				&cli.StringFlag{
					Name:  "postfile",
					Usage: "file containing data to POST (content type will default to application/json)",
				},
				&cli.StringFlag{
					Name:  "patchfile",
					Usage: "file containing data to PATCH (content type will default to application/json)",
				},
				&cli.StringFlag{
					Name:  "putfile",
					Usage: "file containing data to PUT (content type will default to application/json)",
				},
				&cli.StringFlag{
					Name:  "json-metrics-file",
					Usage: "outputs metrics to a custom json filepath, if json-metrics is set to true",
				},
				&cli.BoolFlag{
					Name:  "disable-keep-alive",
					Usage: "use this flag to disable http keep-alive",
				},
				&cli.BoolFlag{
					Name:  "insecure",
					Usage: "use this flag to skip ssl verification",
				},
				&cli.StringFlag{
					Name:  "renegotiation",
					Usage: "Allow client certificate renegotiation: never, once, freely",
					Value: "never",
				},
				&cli.StringFlag{
					Name:  "ca",
					Usage: "ca certificate to verify peer against",
				},
				&cli.StringFlag{
					Name:  "cert",
					Usage: "client authentication certificate",
				},
				&cli.StringFlag{
					Name:  "key",
					Usage: "client authentication key",
				},
				&cli.BoolFlag{
					Name:  "sine-curve",
					Usage: "enable sine curve request pattern (varies request rate over time)",
				},
				&cli.StringFlag{
					Name:  "sine-preset",
					Usage: "use predefined sine wave settings: gentle, normal, aggressive",
				},
				&cli.Float64Flag{
					Name:  "sine-frequency",
					Usage: "number of sine wave cycles during the test duration",
					Value: 0.5,
				},
				&cli.Float64Flag{
					Name:  "sine-amplitude",
					Usage: "amplitude factor for sine wave (0.0-1.0, where 1.0 means rate varies from 0 to 2x baseline)",
					Value: 0.5,
				},
				&cli.BoolFlag{
					Name:  "sine-preview",
					Usage: "show ASCII visualization of sine wave pattern (no HTTP requests sent)",
				},
				&cli.BoolFlag{
					Name:  "sine-preview-live",
					Usage: "show live animated ASCII visualization of sine wave pattern (no HTTP requests sent)",
				},
			},
			Action: validateCLI,
		},
	}

	err := app.Run(args)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}

func setCustomCLITemplate(c *cli.App) {
	whiteBold := color.New(color.Bold).SprintfFunc()
	whiteUnderline := color.New(color.Bold).Add(color.Underline).SprintfFunc()
	cyan := color.New(color.FgCyan).SprintFunc()

	c.CustomAppHelpTemplate = fmt.Sprintf(` %s:
		{{.Name}}{{if .Usage}} - {{.Usage}}{{end}}{{if .Description}}

	 DESCRIPTION:
		{{.Description | nindent 3 | trim}}{{end}}{{if len .Authors}}

	 AUTHOR{{with $length := len .Authors}}{{if ne 1 $length}}S{{end}}{{end}}:
		{{range $index, $author := .Authors}}{{if $index}}
		{{end}}{{$author}}{{end}}{{end}}{{if .VisibleCommands}}

 %s:{{range .VisibleCategories}}{{if .Name}}
	{{.Name}}:{{range .VisibleCommands}}
	  {{join .Names ", "}}{{"\t"}}{{.Usage}}{{end}}{{else}}{{range .VisibleCommands}}
	{{join .Names ", "}}{{"\t"}}{{.Usage}}{{end}}{{end}}{{end}}{{end}}{{if .VisibleFlags}}

 %s:
	{{range $index, $option := .VisibleFlags}}{{if $index}}
	{{end}}{{$option}}{{end}}{{end}}{{if .Copyright}}

 COPYRIGHT:
	{{.Copyright}}{{end}}

	%s
	Example running cassowary against a target with 100 requests using 10 concurrent users
	  %s
  `, whiteBold("NAME"),
		whiteBold("COMMANDS"),
		whiteBold("GLOBAL OPTIONS"),
		whiteUnderline("Example"),
		cyan("$ cassowary run -u http://www.example.com -c 10 -n 100"))
}
