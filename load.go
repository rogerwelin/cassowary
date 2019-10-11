package main

import (
	"fmt"
	"io"
	"io/ioutil"
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
	for item := range workerChan {
		tt := newTransport(c.client.Transport)
		c.client.Transport = tt
		fmt.Println(item)
		fmt.Println(c)

		request, err := http.NewRequest("GET", c.baseURL, nil)
		if err != nil {
			panic(err)
		}

		trace := &httptrace.ClientTrace{
			DNSStart:             tt.DNSStart,
			DNSDone:              tt.DNSDone,
			ConnectStart:         tt.ConnectStart,
			ConnectDone:          tt.ConnectDone,
			GotConn:              tt.GotConn,
			GotFirstResponseByte: tt.GotFirstResponseByte,
		}

		request = request.WithContext(httptrace.WithClientTrace(request.Context(), trace))
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
		tt.current.end = time.Now()
		for _, trace := range tt.traces {
			out := durationMetrics{}

			out.DNSLookup = trace.dnsDone.Sub(trace.start).Seconds()
			out.TCPConn = trace.gotConn.Sub(trace.dnsDone).Seconds()
			out.ServerProcessing = trace.responseStart.Sub(trace.gotConn).Seconds()
			out.ContentTransfer = trace.end.Sub(trace.responseStart).Seconds()
			out.StatusCode = resp.StatusCode

			if trace.tls {
				out.TLSHandshake = trace.gotConn.Sub(trace.dnsDone).Seconds()
			}

			outPutChan <- out
		}
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

	col := color.New(color.FgCyan).Add(color.Underline)
	col.Printf("\nStarting Load Test with %d concurrent users\n\n", c.concurrencyLevel)

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
		fmt.Println("got here")
		fmt.Println(c.inputFile)
		fmt.Println("got here")
		urlSuffixes, err = readFile(c.inputFile)
		if err != nil {
			panic(err)
			return err
		}
		c.requests = len(urlSuffixes)
		fmt.Println(urlSuffixes)
	}
	fmt.Println("apa")
	fmt.Println(c)

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
		tcpDur = append(tcpDur, item.TCPConn)
		if c.isTLS {
			tlsDur = append(tlsDur, item.TLSHandshake)
		}
		serverDur = append(serverDur, item.ServerProcessing)
		transferDur = append(transferDur, item.ContentTransfer)
		statusCodes = append(statusCodes, item.StatusCode)
	}

	// DNS
	dnsMean := calcMean(dnsDur)
	dnsMedian := calcMedian(dnsDur)
	dns95 := calc95Percentile(dnsDur)

	// TCP
	tcpMean := calcMean(tcpDur)
	tcpMedian := calcMedian(tcpDur)
	tcp95 := calc95Percentile(tcpDur)

	// TLS
	var tlsMean float64
	var tlsMedian float64
	var tls95 string

	if c.isTLS {
		tlsMean = calcMean(tlsDur)
		tlsMedian = calcMedian(tlsDur)
		tls95 = calc95Percentile(tlsDur)
	}

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

	if c.isTLS {
		printf(summaryTLSTable,
			color.CyanString(fmt.Sprintf("%f", dnsMean)),
			color.CyanString(fmt.Sprintf("%f", dnsMedian)),
			color.CyanString(dns95),
			color.CyanString(fmt.Sprintf("%f", tcpMean)),
			color.CyanString(fmt.Sprintf("%f", tcpMedian)),
			color.CyanString(tcp95),
			color.CyanString(fmt.Sprintf("%f", tlsMean)),
			color.CyanString(fmt.Sprintf("%f", tlsMedian)),
			color.CyanString(tls95),
			color.CyanString(fmt.Sprintf("%f", serverMean)),
			color.CyanString(fmt.Sprintf("%f", serverMedian)),
			color.CyanString(server95),
			color.CyanString(fmt.Sprintf("%f", transferMean)),
			color.CyanString(fmt.Sprintf("%f", transferMedian)),
			color.CyanString(transfer95),
			color.CyanString(strconv.Itoa(c.requests)),
			color.CyanString(failedR),
			color.CyanString(reqS),
		)
		return nil
	}

	printf(summaryTable,
		color.CyanString(fmt.Sprintf("%f", dnsMean)),
		color.CyanString(fmt.Sprintf("%f", dnsMedian)),
		color.CyanString(dns95),
		color.CyanString(fmt.Sprintf("%f", tcpMean)),
		color.CyanString(fmt.Sprintf("%f", tcpMedian)),
		color.CyanString(tcp95),
		color.CyanString(fmt.Sprintf("%f", serverMean)),
		color.CyanString(fmt.Sprintf("%f", serverMedian)),
		color.CyanString(server95),
		color.CyanString(fmt.Sprintf("%f", transferMean)),
		color.CyanString(fmt.Sprintf("%f", transferMedian)),
		color.CyanString(transfer95),
		color.CyanString(strconv.Itoa(c.requests)),
		color.CyanString(failedR),
		color.CyanString(reqS),
	)
	return nil
}
