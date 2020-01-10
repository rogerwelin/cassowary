package main

import (
	"fmt"

	"github.com/fatih/color"
)

type jsonOutput struct {
	BaseURL           string                `json:"base_url"`
	TotalRequests     string                `json:"total_requests"`
	FailedRequests    string                `json:"failed_requests"`
	RequestsPerSecond string                `json:"requests_per_second"`
	TCPStats          tcpStats              `json:"tcp_connect"`
	ProcessingStats   serverProcessingStats `json:"server_processing"`
	ContentStats      contentTransfer       `json:"content_transfer"`
}

type tcpStats struct {
	TCPMean   string `json:"mean"`
	TCPMedian string `json:"median"`
	TCP95p    string `json:"95th_percentile"`
}

type serverProcessingStats struct {
	ServerProcessingMean   string `json:"mean"`
	ServerProcessingMedian string `json:"median"`
	ServerProcessing95p    string `json:"95th_percentile"`
}

type contentTransfer struct {
	ContentTransferMean   string `json:"mean"`
	ContentTransferMedian string `json:"median"`
	ContentTransfer95p    string `json:"95th_percentile"`
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

func (c *cassowary) outPutJSON() {

}
