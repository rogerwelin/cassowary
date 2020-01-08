package main

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

func (c *cassowary) pushPrometheusMetrics(t1, t2, t3, s1, s2, s3, tf1, tf2, tf3, r1, r2, r3 float64) error {
	promTCPConnMean.Set(t1)
	promTCPConnMedian.Set(t2)
	promTCPConn95P.Set(t3)
	promServerProcessingMean.Set(s1)
	promServerProcessingMedian.Set(s2)
	promServerProcessing95p.Set(s3)
	promContentTransferMean.Set(tf1)
	promContentTransferMedian.Set(tf2)
	promContentTransfer95p.Set(tf3)
	promTotalRequests.Set(r1)
	promFailedRequests.Set(r2)
	promRequestPerSecond.Set(r3)

	if err := push.New(c.promURL, "cassowary_load_test").
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
		Grouping("url", c.baseURL).
		Push(); err != nil {
		return err
	}
	return nil
}
