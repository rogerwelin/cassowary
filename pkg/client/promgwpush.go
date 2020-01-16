package client

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
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

// PushPrometheusMetrics exports metrics to a PushGateway
func (c *Cassowary) PushPrometheusMetrics(metrics ResultMetrics) error {
	promTCPConnMean.Set(metrics.TCPStats.TCPMean)
	promTCPConnMedian.Set(metrics.TCPStats.TCPMedian)
	promTCPConn95P.Set(metrics.TCPStats.TCP95p)
	promServerProcessingMean.Set(metrics.ProcessingStats.ServerProcessingMean)
	promServerProcessingMedian.Set(metrics.ProcessingStats.ServerProcessingMedian)
	promServerProcessing95p.Set(metrics.ProcessingStats.ServerProcessing95p)
	promContentTransferMean.Set(metrics.ContentStats.ContentTransferMean)
	promContentTransferMedian.Set(metrics.ContentStats.ContentTransferMedian)
	promContentTransfer95p.Set(metrics.ContentStats.ContentTransfer95p)
	promTotalRequests.Set(float64(metrics.TotalRequests))
	promFailedRequests.Set(float64(metrics.FailedRequests))
	promRequestPerSecond.Set(metrics.RequestsPerSecond)

	if err := push.New(c.PromURL, "cassowary_load_test").
		Collector(promTCPConnMean).
		Collector(promTCPConnMedian).
		Collector(promTCPConn95P).
		Collector(promServerProcessingMean).
		Collector(promServerProcessingMedian).
		Collector(promServerProcessing95p).
		Collector(promContentTransferMean).
		Collector(promContentTransferMedian).
		Collector(promContentTransfer95p).
		Collector(promTotalRequests).
		Collector(promFailedRequests).
		Collector(promRequestPerSecond).
		Grouping("url", c.BaseURL).
		Push(); err != nil {
		return err
	}
	return nil
}
