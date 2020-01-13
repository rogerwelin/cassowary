package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/fatih/color"
)

type jsonOutput struct {
	BaseURL           string                `json:"base_url"`
	TotalRequests     int                   `json:"total_requests"`
	FailedRequests    int                   `json:"failed_requests"`
	RequestsPerSecond float64               `json:"requests_per_second"`
	TCPStats          tcpStats              `json:"tcp_connect"`
	ProcessingStats   serverProcessingStats `json:"server_processing"`
	ContentStats      contentTransfer       `json:"content_transfer"`
}

type tcpStats struct {
	TCPMean   float64 `json:"mean"`
	TCPMedian float64 `json:"median"`
	TCP95p    int     `json:"95th_percentile"`
}

type serverProcessingStats struct {
	ServerProcessingMean   float64 `json:"mean"`
	ServerProcessingMedian float64 `json:"median"`
	ServerProcessing95p    int     `json:"95th_percentile"`
}

type contentTransfer struct {
	ContentTransferMean   float64 `json:"mean"`
	ContentTransferMedian float64 `json:"median"`
	ContentTransfer95p    int     `json:"95th_percentile"`
}

const (
	summaryTable = `` + "\n\n" +
		` TCP Connect.....................: Avg/mean=%sms ` + "\t" + `Median=%sms` + "\t" + `p(95)=%sms` + "\n" +
		` Server Processing...............: Avg/mean=%sms ` + "\t" + `Median=%sms` + "\t" + `p(95)=%sms` + "\n" +
		` Content Transfer................: Avg/mean=%sms ` + "\t" + `Median=%sms` + "\t" + `p(95)=%sms` + "\n" +
		`` + "\n" +
		`Summary: ` + "\n" +
		` Total Req.......................: %s` + "\n" +
		` Failed Req......................: %s` + "\n" +
		` DNS Lookup......................: %sms` + "\n" +
		` Req/s...........................: %s` + "\n\n"
)

func printf(format string, a ...interface{}) {
	fmt.Fprintf(color.Output, format, a...)
}

func (c *cassowary) outPutJSON(failedReq int, requestPerSec, tcpMean, tcpMed float64, tcp9p string, serverMean, serverMed float64, server95p string, contentMean, contentMed float64, content95p string) {
	tcp9P, _ := strconv.Atoi(tcp9p)
	server95P, _ := strconv.Atoi(server95p)
	content95P, _ := strconv.Atoi(content95p)
	output := jsonOutput{
		BaseURL:           c.baseURL,
		TotalRequests:     c.requests,
		FailedRequests:    failedReq,
		RequestsPerSecond: requestPerSec,
		TCPStats: tcpStats{
			TCPMean:   tcpMean,
			TCPMedian: tcpMed,
			TCP95p:    tcp9P,
		},
		ProcessingStats: serverProcessingStats{
			ServerProcessingMean:   serverMean,
			ServerProcessingMedian: serverMed,
			ServerProcessing95p:    server95P,
		},
		ContentStats: contentTransfer{
			ContentTransferMean:   contentMean,
			ContentTransferMedian: contentMed,
			ContentTransfer95p:    content95P,
		},
	}

	b, err := json.Marshal(output)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = ioutil.WriteFile("out.json", b, 0644)
	if err != nil {
		fmt.Println(err)
	}
}
