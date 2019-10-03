package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/schollz/progressbar"
)

type durationMetrics struct {
	DnsLookup        int
	TCPConn          int
	TLSHandshake     int
	ServerProcessing int
	ContentTransfer  int
	StatusCode       int
	URL              string
}

func (c *cassowary) runLoadTest(cmdOutputChan chan<- durationMetrics, workerChan chan string) {

}

func (c *cassowary) coordinate() error {

	color := color.New(color.FgCyan).Add(color.Underline)
	color.Printf("\nStarting Load Test with %d concurrent users\n\n", c.concurrencyLevel)

	var urlSuffixes []string
	var err error

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
		fmt.Println(urlSuffixes)
	}

	var wg sync.WaitGroup
	channel := make(chan durationMetrics, c.requests)
	workerChan := make(chan string)

	return nil
}
