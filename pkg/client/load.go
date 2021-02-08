package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptrace"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/schollz/progressbar"
)

type durationMetrics struct {
	DNSLookup        float64
	TCPConn          float64
	TLSHandshake     float64
	ServerProcessing float64
	ContentTransfer  float64
	StatusCode       int
}

func (c *Cassowary) runLoadTest(outPutChan chan<- durationMetrics, workerChan chan string) {
	for URLitem := range workerChan {

		var request *http.Request
		var err error

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

		if len(c.RequestHeader) == 2 {
			request.Header.Add(c.RequestHeader[0], c.RequestHeader[1])
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
			_, err = io.Copy(ioutil.Discard, resp.Body)
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

		outPutChan <- out
	}
}

// Coordinate bootstraps the load test based on values in Cassowary struct
func (c *Cassowary) Coordinate() (ResultMetrics, error) {
	var dnsDur []float64
	var tcpDur []float64
	var tlsDur []float64
	var serverDur []float64
	var transferDur []float64
	var statusCodes []int

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
			DisableCompression:  false,
			DisableKeepAlives:   c.DisableKeepAlive,
		},
	}

	if c.FileMode {
		if c.Requests > len(c.URLPaths) {
			c.URLPaths = generateSuffixes(c.URLPaths, c.Requests)
		}
		c.Requests = len(c.URLPaths)
	}

	c.Bar = progressbar.New(c.Requests)

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
		ticker := time.NewTicker(time.Duration(nextTick) * time.Millisecond)
		if nextTick == 0 {
			log.Fatalf("The combination of %v requests and %v(s) duration is invalid. Try raising the duration to a greater value", c.Requests, c.Duration)
		}
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

	for item := range channel {
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
	return outPut, nil
}
