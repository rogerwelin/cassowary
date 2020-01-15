package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptrace"
	"strconv"
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

func (c *cassowary) runLoadTest(outPutChan chan<- durationMetrics, workerChan chan string) {
	for URLitem := range workerChan {

		var request *http.Request
		var err error

		if c.fileMode {
			request, err = http.NewRequest("GET", c.baseURL+URLitem, nil)
			if err != nil {
				panic(err)
			}
		} else {
			request, err = http.NewRequest("GET", c.baseURL, nil)
			if err != nil {
				panic(err)
			}
		}

		if len(c.requestHeader) == 2 {
			request.Header.Add(c.requestHeader[0], c.requestHeader[1])
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
		resp, err := c.client.Do(request)
		if err != nil {
			panic(err)
		}
		if resp != nil {
			_, err = io.Copy(ioutil.Discard, resp.Body)
			if err != nil {
				fmt.Println("Failed to read HTTP response body", err)
			}
			resp.Body.Close()
		}
		c.bar.Add(1)

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

		if c.isTLS {
			out.TCPConn = float64(t2.Sub(t1) / time.Millisecond)
			out.TLSHandshake = float64(t6.Sub(t5) / time.Millisecond) // tls handshake
		} else {
			out.TCPConn = float64(t3.Sub(t1) / time.Millisecond)
		}

		outPutChan <- out
	}
}

func (c *cassowary) coordinate() error {
	var dnsDur []float64
	var tcpDur []float64
	var tlsDur []float64
	var serverDur []float64
	var transferDur []float64
	var statusCodes []int

	tls, err := isTLS(c.baseURL)
	if err != nil {
		return err
	}
	c.isTLS = tls

	var urlSuffixes []string

	c.client = &http.Client{
		Timeout: time.Second * 5,
		Transport: &http.Transport{
			MaxIdleConns:        300,
			MaxIdleConnsPerHost: 300,
			MaxConnsPerHost:     300,
			DisableCompression:  false,
		},
	}

	c.bar = progressbar.New(c.requests)

	if c.fileMode {
		urlSuffixes, err = readFile(c.inputFile)
		if err != nil {
			return err
		}
		c.requests = len(urlSuffixes)
		c.bar = progressbar.New(c.requests)
		//fmt.Println(urlSuffixes)
	}

	col := color.New(color.FgCyan).Add(color.Underline)
	col.Printf("\nStarting Load Test with %d requests using %d concurrent users\n\n", c.requests, c.concurrencyLevel)

	var wg sync.WaitGroup
	channel := make(chan durationMetrics, c.requests)
	workerChan := make(chan string)

	wg.Add(c.concurrencyLevel)
	start := time.Now()

	for i := 0; i < c.concurrencyLevel; i++ {
		go func() {
			c.runLoadTest(channel, workerChan)
			wg.Done()
		}()
	}

	if c.fileMode {
		for _, line := range urlSuffixes {
			workerChan <- line
		}
	} else {
		for i := 0; i < c.requests; i++ {
			workerChan <- "a"
		}
	}

	close(workerChan)
	wg.Wait()
	close(channel)

	end := time.Since(start)
	fmt.Println(end)

	for item := range channel {
		if item.DNSLookup != 0 {
			dnsDur = append(dnsDur, item.DNSLookup)
		}
		if item.TCPConn < 1000 {
			tcpDur = append(tcpDur, item.TCPConn)
		}
		if c.isTLS {
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
	reqS := requestsPerSecond(c.requests, end)

	// Failed Requests
	failedR := failedRequests(statusCodes)

	printf(summaryTable,
		color.CyanString(fmt.Sprintf("%.2f", tcpMean)),
		color.CyanString(fmt.Sprintf("%.2f", tcpMedian)),
		color.CyanString(tcp95),
		color.CyanString(fmt.Sprintf("%.2f", serverMean)),
		color.CyanString(fmt.Sprintf("%.2f", serverMedian)),
		color.CyanString(server95),
		color.CyanString(fmt.Sprintf("%.2f", transferMean)),
		color.CyanString(fmt.Sprintf("%.2f", transferMedian)),
		color.CyanString(transfer95),
		color.CyanString(strconv.Itoa(c.requests)),
		color.CyanString(strconv.Itoa(failedR)),
		color.CyanString(fmt.Sprintf("%.2f", dnsMedian)),
		color.CyanString(reqS),
	)

	if c.promExport == true {
		err := c.pushPrometheusMetrics(
			tcpMean,
			tcpMedian,
			stringToFloat(tcp95),
			serverMean,
			serverMedian,
			stringToFloat(server95),
			transferMean,
			transferMedian,
			stringToFloat(transfer95),
			float64(c.requests),
			float64(failedR),
			stringToFloat(reqS),
		)
		if err != nil {
			return err
		}
	}

	if c.exportMetrics {
		return c.outPutJSON(failedR, stringToFloat(reqS), tcpMean, tcpMedian, tcp95, serverMean, serverMedian, server95, transferMean, transferMedian, transfer95)
	}

	//fmt.Println(tcpDur)
	//fmt.Println(dnsDur)
	//fmt.Println(tlsDur)
	return nil
}
