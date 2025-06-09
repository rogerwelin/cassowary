package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptrace"
	"os"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/schollz/progressbar/v3"
)

type durationMetrics struct {
	DNSLookup        float64
	TCPConn          float64
	TLSHandshake     float64
	ServerProcessing float64
	ContentTransfer  float64
	StatusCode       int
	TotalDuration    float64
}

func (c *Cassowary) runLoadTest(outPutChan chan<- durationMetrics, workerChan chan string) {
	// Pre-allocate and reuse request objects for GET requests if possible
	var baseRequest *http.Request
	if c.HTTPMethod == "GET" && !c.FileMode {
		var err error
		baseRequest, err = http.NewRequest("GET", c.BaseURL, nil)
		if err != nil {
			log.Fatalf("%v", err)
		}

		// Set common headers once
		if len(c.RequestHeader)%2 == 0 {
			for idx := range c.RequestHeader {
				if idx%2 == 1 {
					continue
				}
				baseRequest.Header.Add(c.RequestHeader[idx], c.RequestHeader[idx+1])
			}
		}
	}

	for URLitem := range workerChan {
		var request *http.Request
		var err error

		if baseRequest != nil {
			// Clone the base request if we have one prepared
			request = baseRequest.Clone(context.Background())
		} else {
			if c.FileMode {
				request, err = http.NewRequest("GET", c.BaseURL+URLitem, nil)
				if err != nil {
					log.Fatalf("%v", err)
				}
			} else {
				switch c.HTTPMethod {
				case "POST":
					request, err = http.NewRequest("POST", c.BaseURL, bytes.NewBuffer(c.Data))
					request.Header.Set("Content-Type", "application/json")
					if err != nil {
						log.Fatalf("%v", err)
					}
				case "PUT":
					request, err = http.NewRequest("PUT", c.BaseURL, bytes.NewBuffer(c.Data))
					request.Header.Set("Content-Type", "application/json")
					if err != nil {
						log.Fatalf("%v", err)
					}
				case "PATCH":
					request, err = http.NewRequest("PATCH", c.BaseURL, bytes.NewBuffer(c.Data))
					request.Header.Set("Content-Type", "application/json")
					if err != nil {
						log.Fatalf("%v", err)
					}
				default:
					request, err = http.NewRequest("GET", c.BaseURL, nil)
					if err != nil {
						log.Fatalf("%v", err)
					}
				}
			}

			if len(c.RequestHeader)%2 == 0 {
				for idx := range c.RequestHeader {
					if idx%2 == 1 {
						continue
					}
					request.Header.Add(c.RequestHeader[idx], c.RequestHeader[idx+1])
				}
			}
		}

		var t0, t1, t2, t3, t4, t5, t6 time.Time

		trace := &httptrace.ClientTrace{
			DNSStart: func(_ httptrace.DNSStartInfo) { t0 = time.Now() },
			DNSDone:  func(_ httptrace.DNSDoneInfo) { t1 = time.Now() },
			ConnectStart: func(_, _ string) {
				if t1.IsZero() {
					// connecting directly to IP
					t1 = time.Now()
				}
			},
			ConnectDone: func(net, addr string, err error) {
				if err != nil {
					log.Fatalf("unable to connect to host %v: %v", addr, err)
				}
				t2 = time.Now()

			},
			GotConn:              func(_ httptrace.GotConnInfo) { t3 = time.Now() },
			GotFirstResponseByte: func() { t4 = time.Now() },
			TLSHandshakeStart:    func() { t5 = time.Now() },
			TLSHandshakeDone:     func(_ tls.ConnectionState, _ error) { t6 = time.Now() },
		}

		request = request.WithContext(httptrace.WithClientTrace(context.Background(), trace))
		resp, err := c.Client.Do(request)
		if err != nil {
			log.Fatalf("%v", err)
		}
		if resp != nil {
			_, err = io.Copy(io.Discard, resp.Body)
			if err != nil {
				fmt.Println("Failed to read HTTP response body", err)
			}
			resp.Body.Close()
		}

		if !c.DisableTerminalOutput {
			c.Bar.Add(1)
		}

		// Body fully read here
		t7 := time.Now()
		if t0.IsZero() {
			// we skipped DNS
			t0 = t1
		}

		out := durationMetrics{
			DNSLookup: float64(t1.Sub(t0) / time.Millisecond), // dns lookup
			//TCPConn:          float64(t3.Sub(t1) / time.Millisecond), // tcp connection
			ServerProcessing: float64(t4.Sub(t3) / time.Millisecond), // server processing
			ContentTransfer:  float64(t7.Sub(t4) / time.Millisecond), // content transfer
			StatusCode:       resp.StatusCode,
		}

		if c.IsTLS {
			out.TCPConn = float64(t2.Sub(t1) / time.Millisecond)
			out.TLSHandshake = float64(t6.Sub(t5) / time.Millisecond) // tls handshake
		} else {
			out.TCPConn = float64(t3.Sub(t1) / time.Millisecond)
		}

		out.TotalDuration = out.DNSLookup + out.TCPConn + out.ServerProcessing + out.ContentTransfer

		outPutChan <- out
	}
}

// Coordinate bootstraps the load test based on values in Cassowary struct
func (c *Cassowary) Coordinate() (ResultMetrics, error) {
	var (
		dnsDur      = make([]float64, 0, c.Requests)
		tcpDur      = make([]float64, 0, c.Requests)
		tlsDur      = make([]float64, 0, c.Requests)
		serverDur   = make([]float64, 0, c.Requests)
		transferDur = make([]float64, 0, c.Requests)
		statusCodes = make([]int, 0, c.Requests)
		totalDur    = make([]float64, 0, c.Requests)
	)

	tls, err := isTLS(c.BaseURL)
	if err != nil {
		return ResultMetrics{}, err
	}
	c.IsTLS = tls

	c.Client = &http.Client{
		Timeout: time.Second * time.Duration(c.Timeout),
		Transport: &http.Transport{
			TLSClientConfig:     c.TLSConfig,
			MaxIdleConnsPerHost: 10000,
			DisableCompression:  true,
			DisableKeepAlives:   c.DisableKeepAlive,
			Proxy:               http.ProxyFromEnvironment,
		},
	}

	if c.FileMode {
		if c.Requests > len(c.URLPaths) {
			c.URLPaths = generateSuffixes(c.URLPaths, c.Requests)
		}
		c.Requests = len(c.URLPaths)
	}

	c.Bar = progressbar.NewOptions(c.Requests,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))

	if !c.DisableTerminalOutput {
		col := color.New(color.FgCyan).Add(color.Underline)
		col.Printf("\nStarting Load Test with %d requests using %d concurrent users\n\n", c.Requests, c.ConcurrencyLevel)
	}

	var wg sync.WaitGroup
	channel := make(chan durationMetrics, c.Requests)
	workerChan := make(chan string)

	wg.Add(c.ConcurrencyLevel)
	start := time.Now()

	for i := 0; i < c.ConcurrencyLevel; i++ {
		go func() {
			c.runLoadTest(channel, workerChan)
			wg.Done()
		}()
	}

	if c.Duration > 0 {
		durationMS := c.Duration * 1000
		nextTick := durationMS / c.Requests

		// Check if nextTick is too small (less than 1ms)
		if nextTick < 1 {
			// If nextTick is less than 1 millisecond the duration is too short to handle the requested number of requests.
			// To avoid this, calculate the *minimum* valid duration (in seconds) needed to
			// accommodate c.Requests at a minimum interval of 1ms per request.
			// Use ceiling division to round up: (requests + 999) / 1000 ensures rounding up
			// when converting milliseconds to seconds.
			// To achieve ceiling division using integers, use the formula: (a + b - 1) / b.
			// In this case, a = c.Requests and b = 1000 (milliseconds per second), so:
			//   minDuration = (c.Requests + 1000 - 1) / 1000
			//               = (c.Requests + 999) / 1000
			// This ensures that any fractional second is rounded *up*, not down.
			minDuration := (c.Requests + 999) / 1000
			log.Fatalf("The combination of %v requests and %v(s) duration is invalid. The minimum duration required for %v requests is %v seconds. Try using -d %v or higher", c.Requests, c.Duration, c.Requests, minDuration, minDuration)
		}

		ticker := time.NewTicker(time.Duration(nextTick) * time.Millisecond)
		done := make(chan bool)
		iter := 0

		go func() {
			for {
				select {
				case <-done:
					return
				case _ = <-ticker.C:
					if c.FileMode {
						workerChan <- c.URLPaths[iter]
						iter++
					} else {
						workerChan <- "a"
					}
				}
			}
		}()

		time.Sleep(time.Duration(durationMS) * time.Millisecond)
		ticker.Stop()
		done <- true
	}

	if c.Duration == 0 && c.FileMode {
		for _, line := range c.URLPaths {
			workerChan <- line
		}
	} else if c.Duration == 0 && !c.FileMode {
		for i := 0; i < c.Requests; i++ {
			workerChan <- "a"
		}
	}

	close(workerChan)
	wg.Wait()
	close(channel)

	end := time.Since(start)
	if !c.DisableTerminalOutput {
		fmt.Println(end)
	}

	var w *csv.Writer

	if c.RawOutput {
		csvFile, err := os.Create("raw.csv")
		if err != nil {
			return ResultMetrics{}, err
		}
		w = csv.NewWriter(csvFile)
		headerInfo := structNames(&durationMetrics{})
		w.Write(headerInfo)
		defer csvFile.Close()
		defer w.Flush()
	}

	for item := range channel {
		if c.RawOutput {
			itemSlice := toSlice(item)
			w.Write(itemSlice)
		}
		if item.DNSLookup != 0 {
			dnsDur = append(dnsDur, item.DNSLookup)
		}
		if item.TCPConn < 1000 {
			tcpDur = append(tcpDur, item.TCPConn)
		}
		if c.IsTLS {
			tlsDur = append(tlsDur, item.TLSHandshake)
		}
		serverDur = append(serverDur, item.ServerProcessing)
		transferDur = append(transferDur, item.ContentTransfer)
		statusCodes = append(statusCodes, item.StatusCode)
		totalDur = append(totalDur, item.TotalDuration)
	}

	// DNS
	dnsMedian := calcMedian(dnsDur)

	// TCP
	tcpMean := calcMean(tcpDur)
	tcpMedian := calcMedian(tcpDur)
	tcp95 := calc95Percentile(tcpDur)

	// Server Processing
	serverMean := calcMean(serverDur)
	serverMedian := calcMedian(serverDur)
	server95 := calc95Percentile(serverDur)

	// Content Transfer
	transferMean := calcMean(transferDur)
	transferMedian := calcMedian(transferDur)
	transfer95 := calc95Percentile(transferDur)

	// Request per second
	reqS := requestsPerSecond(c.Requests, end)

	// Failed Requests
	failedR := failedRequests(statusCodes)

	outPut := ResultMetrics{
		BaseURL:           c.BaseURL,
		FailedRequests:    failedR,
		RequestsPerSecond: reqS,
		TotalRequests:     c.Requests,
		DNSMedian:         dnsMedian,
		TCPStats: tcpStats{
			TCPMean:   tcpMean,
			TCPMedian: tcpMedian,
			TCP95p:    stringToFloat(tcp95),
		},
		ProcessingStats: serverProcessingStats{
			ServerProcessingMean:   serverMean,
			ServerProcessingMedian: serverMedian,
			ServerProcessing95p:    stringToFloat(server95),
		},
		ContentStats: contentTransfer{
			ContentTransferMean:   transferMean,
			ContentTransferMedian: transferMedian,
			ContentTransfer95p:    stringToFloat(transfer95),
		},
	}

	// output histogram
	if c.Histogram {
		err := c.PlotHistogram(totalDur)
		if err != nil {
			fmt.Println(err)
		}
	}

	// output boxplot
	if c.Boxplot {
		err := c.PlotBoxplot(totalDur)
		if err != nil {
			fmt.Println(err)
		}
	}
	return outPut, nil
}
