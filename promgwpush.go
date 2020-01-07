package main

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	promTCPConnMean = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "tcp_connect_mean",
		Help: "Duration of TCP Connection mean",
	})
	promTCPConnMedian = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "tcp_connect_median",
		Help: "Duration of TCP Connection median",
	})
	promTCPConn95P = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "tcp_connect_95p",
		Help: "Duration of TCP Connection 95th percentile",
	})
	promServerProcessingMean = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "server_processing_mean",
		Help: "Duration of server processing in last load test",
	})
	promServerProcessingMedian = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "server_processing_median",
		Help: "Duration of server processing in last load test",
	})
	promServerProcessing95p = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "server_processing_95p",
		Help: "Duration of server processing in last load test",
	})
	promContentTransferMean = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "content_transfer_mean",
		Help: "Duration of content transfer in last load test",
	})
	promContentTransferMedian = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "content_transfer_median",
		Help: "Duration of content transfer in last load test",
	})
	promContentTransfer95p = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "content_transfer_95p",
		Help: "Duration of content transfer in last load test",
	})
	promTotalRequests = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "total_requests",
		Help: "Total requests in last load test",
	})
	promFailedRequests = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "failed_requests",
		Help: "Failed requests in last load test",
	})
	promRequestPerSecond = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "requests_per_second",
		Help: "Requests per second in last load test",
	})
)

func (c *cassowary) pushPrometheusMetrics() error {
	return nil
}
