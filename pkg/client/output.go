package client

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/fatih/color"
)

const (
	summaryTable = "\n" +
		" Elapsed(ms)............: Min=%s\tMax=%s\tAvg=%s\tMedian=%s\tp(95)=%s\tp(99)=%s\n" +
		" TCP Connect(ms)........: Min=%s\tMax=%s\tAvg=%s\tMedian=%s\tp(95)=%s\tp(99)=%s\n" +
		" Server Processing(ms)..: Min=%s\tMax=%s\tAvg=%s\tMedian=%s\tp(95)=%s\tp(99)=%s\n" +
		" Content Transfer.......: Min=%s\tMax=%s\tAvg=%s\tMedian=%s\tp(95)=%s\tp(99)=%s\n" +
		" Body Size(bytes).......: Min=%s\tMax=%s\tAvg=%s\tMedian=%s\tp(95)=%s\tp(99)=%s\n" +
		" Response Size(bytes)...: Min=%s\tMax=%s\tAvg=%s\tMedian=%s\tp(95)=%s\tp(99)=%s\n" +
		"\n" +
		"Summary:\n" +
		" Total Req.......................: %s\n" +
		" Failed Req......................: %s\n" +
		" DNS Lookup......................: %sms\n" +
		" Req/s...........................: %s\n\n"
)

func printf(format string, a ...interface{}) {
	fmt.Fprintf(color.Output, format, a...)
}

func OutPutResults(metrics ResultMetrics) {
	if metrics.Name == TotalStr {
		printf(color.GreenString(fmt.Sprintf("\n%s:", metrics.Name)))
	} else {
		printf(color.GreenString(fmt.Sprintf("\nGroup %s:", metrics.Name)))
	}
	printf(summaryTable,
		color.CyanString(fmt.Sprintf("%.2f", metrics.ElapsedStats.Min)),
		color.CyanString(fmt.Sprintf("%.2f", metrics.ElapsedStats.Max)),
		color.CyanString(fmt.Sprintf("%.2f", metrics.ElapsedStats.Mean)),
		color.CyanString(fmt.Sprintf("%.2f", metrics.ElapsedStats.Median)),
		color.CyanString(fmt.Sprintf("%.2f", metrics.ElapsedStats.P95)),
		color.CyanString(fmt.Sprintf("%.2f", metrics.ElapsedStats.P99)),
		color.CyanString(fmt.Sprintf("%.2f", metrics.TCPStats.Min)),
		color.CyanString(fmt.Sprintf("%.2f", metrics.TCPStats.Max)),
		color.CyanString(fmt.Sprintf("%.2f", metrics.TCPStats.Mean)),
		color.CyanString(fmt.Sprintf("%.2f", metrics.TCPStats.Median)),
		color.CyanString(fmt.Sprintf("%.2f", metrics.TCPStats.P95)),
		color.CyanString(fmt.Sprintf("%.2f", metrics.TCPStats.P99)),
		color.CyanString(fmt.Sprintf("%.2f", metrics.ProcessingStats.Min)),
		color.CyanString(fmt.Sprintf("%.2f", metrics.ProcessingStats.Max)),
		color.CyanString(fmt.Sprintf("%.2f", metrics.ProcessingStats.Mean)),
		color.CyanString(fmt.Sprintf("%.2f", metrics.ProcessingStats.Median)),
		color.CyanString(fmt.Sprintf("%.2f", metrics.ProcessingStats.P95)),
		color.CyanString(fmt.Sprintf("%.2f", metrics.ProcessingStats.P99)),
		color.CyanString(fmt.Sprintf("%.2f", metrics.ContentStats.Min)),
		color.CyanString(fmt.Sprintf("%.2f", metrics.ContentStats.Max)),
		color.CyanString(fmt.Sprintf("%.2f", metrics.ContentStats.Mean)),
		color.CyanString(fmt.Sprintf("%.2f", metrics.ContentStats.Median)),
		color.CyanString(fmt.Sprintf("%.2f", metrics.ContentStats.P95)),
		color.CyanString(fmt.Sprintf("%.2f", metrics.ContentStats.P99)),
		color.CyanString(fmt.Sprintf("%.0f", metrics.BodySize.Min)),
		color.CyanString(fmt.Sprintf("%.0f", metrics.BodySize.Max)),
		color.CyanString(fmt.Sprintf("%.0f", metrics.BodySize.Mean)),
		color.CyanString(fmt.Sprintf("%.0f", metrics.BodySize.Median)),
		color.CyanString(fmt.Sprintf("%.0f", metrics.BodySize.P95)),
		color.CyanString(fmt.Sprintf("%.0f", metrics.BodySize.P99)),
		color.CyanString(fmt.Sprintf("%.0f", metrics.RespSize.Min)),
		color.CyanString(fmt.Sprintf("%.0f", metrics.RespSize.Max)),
		color.CyanString(fmt.Sprintf("%.0f", metrics.RespSize.Mean)),
		color.CyanString(fmt.Sprintf("%.0f", metrics.RespSize.Median)),
		color.CyanString(fmt.Sprintf("%.0f", metrics.RespSize.P95)),
		color.CyanString(fmt.Sprintf("%.0f", metrics.RespSize.P99)),
		color.CyanString(strconv.Itoa(metrics.TotalRequests)),
		color.CyanString(strconv.Itoa(metrics.FailedRequests)),
		color.CyanString(fmt.Sprintf("%.2f", metrics.DNSMedian)),
		color.CyanString(fmt.Sprintf("%.2f", metrics.RequestsPerSecond)),
	)
	if len(metrics.RespSuccess) > 0 {
		printf(color.GreenString("Success:\n"))
		for status, count := range metrics.RespSuccess {
			printf(color.GreenString(fmt.Sprintf("  %s: %d\n", status, count)))
		}
	}
	if len(metrics.RespFailed) > 0 {
		printf(color.RedString("Failed:\n"))
		for status, count := range metrics.RespFailed {
			printf(color.RedString(fmt.Sprintf("  %s: %d\n", status, count)))
		}
	}
}

func OutPutJSON(fileName string, metrics ResultMetrics, metricsGroup map[string]ResultMetrics) error {
	if fileName == "" {
		// default filename for json metrics output.
		fileName = "out.json"
	}
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	f.WriteString("[\n")
	enc := json.NewEncoder(f)
	err = enc.Encode(metrics)
	if err != nil {
		return err
	}
	for _, m := range metricsGroup {
		f.WriteString(",\n")
		err = enc.Encode(m)
	}
	f.WriteString("]\n")

	return err
}
