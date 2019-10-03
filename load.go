package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/fatih/color"
)

var client = &http.Client{
	Timeout: time.Second * 5,
	Transport: &http.Transport{
		MaxIdleConns:        300,
		MaxIdleConnsPerHost: 300,
		MaxConnsPerHost:     300,
	},
}

func runLoadTest(client *http.Client, baseURL string) {

}

func (c *cassowary) coordinate() error {

	color := color.New(color.FgCyan).Add(color.Underline)
	fmt.Println()
	color.Printf("Starting Load Test with %d concurrent users\n\n", c.concurrencyLevel)

	var urlSuffixes []string
	var err error

	if c.fileMode {
		urlSuffixes, err = readFile(c.inputFile)
		if err != nil {
			return err
		}
		c.requests = len(urlSuffixes)
		fmt.Println(urlSuffixes)
	}

	fmt.Println(c.requests)

	return nil
}
