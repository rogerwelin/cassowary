package client

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)

var (
	promElapsedMin = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "elapsed_min",
		Help: "Duration of TCP Connection min",
	})
	promElapsedMax = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "elapsed_max",
		Help: "Duration of TCP Connection max",
	})
	promElapsedMean = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "elapsed_mean",
		Help: "Duration of TCP Connection mean",
	})
	promElapsedMedian = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "elapsed_median",
		Help: "Duration of TCP Connection median",
	})
	promElapsed95P = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "elapsed_95p",
		Help: "Duration of TCP Connection 95th percentile",
	})
	promElapsed99P = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "elapsed_99p",
		Help: "Duration of TCP Connection 99th percentile",
	})
	promTCPConnMin = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "tcp_connect_min",
		Help: "Duration of TCP Connection min",
	})
	promTCPConnMax = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "tcp_connect_max",
		Help: "Duration of TCP Connection max",
	})
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
	promTCPConn99P = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "tcp_connect_99p",
		Help: "Duration of TCP Connection 99th percentile",
	})
	promServerProcessingMin = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "server_processing_min",
		Help: "Duration of server processing in last load test",
	})
	promServerProcessingMax = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "server_processing_max",
		Help: "Duration of server processing in last load test",
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
	promServerProcessing99p = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "server_processing_99p",
		Help: "Duration of server processing in last load test",
	})
	promContentTransferMin = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "content_transfer_min",
		Help: "Duration of content transfer in last load test",
	})
	promContentTransferMax = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "content_transfer_max",
		Help: "Duration of content transfer in last load test",
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
	promContentTransfer99p = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "content_transfer_99p",
		Help: "Duration of content transfer in last load test",
	})
	promBodySizeMin = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "body_size_min",
		Help: "Duration of content transfer in last load test",
	})
	promBodySizeMax = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "body_size_max",
		Help: "Duration of content transfer in last load test",
	})
	promBodySizeMean = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "body_size_mean",
		Help: "Duration of content transfer in last load test",
	})
	promBodySizeMedian = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "body_size_median",
		Help: "Duration of content transfer in last load test",
	})
	promBodySize95p = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "body_size_95p",
		Help: "Duration of content transfer in last load test",
	})
	promBodySize99p = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "body_size_99p",
		Help: "Duration of content transfer in last load test",
	})
	promRespSizeMin = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "resp_size_min",
		Help: "Duration of content transfer in last load test",
	})
	promRespSizeMax = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "resp_size_max",
		Help: "Duration of content transfer in last load test",
	})
	promRespSizeMean = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "resp_size_mean",
		Help: "Duration of content transfer in last load test",
	})
	promRespSizeMedian = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "resp_size_median",
		Help: "Duration of content transfer in last load test",
	})
	promRespSize95p = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "resp_size_95p",
		Help: "Duration of content transfer in last load test",
	})
	promRespSize99p = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "resp_size_99p",
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
	promElapsedMin.Set(metrics.ElapsedStats.Min)
	promElapsedMax.Set(metrics.ElapsedStats.Max)
	promElapsedMean.Set(metrics.ElapsedStats.Mean)
	promElapsedMedian.Set(metrics.ElapsedStats.Median)
	promElapsed95P.Set(metrics.ElapsedStats.P95)
	promElapsed95P.Set(metrics.ElapsedStats.P99)
	promTCPConnMin.Set(metrics.TCPStats.Min)
	promTCPConnMax.Set(metrics.TCPStats.Max)
	promTCPConnMean.Set(metrics.TCPStats.Mean)
	promTCPConnMedian.Set(metrics.TCPStats.Median)
	promTCPConn95P.Set(metrics.TCPStats.P95)
	promTCPConn95P.Set(metrics.TCPStats.P99)
	promServerProcessingMin.Set(metrics.ProcessingStats.Min)
	promServerProcessingMax.Set(metrics.ProcessingStats.Max)
	promServerProcessingMean.Set(metrics.ProcessingStats.Mean)
	promServerProcessingMedian.Set(metrics.ProcessingStats.Median)
	promServerProcessing95p.Set(metrics.ProcessingStats.P95)
	promServerProcessing99p.Set(metrics.ProcessingStats.P95)
	promContentTransferMin.Set(metrics.ContentStats.Min)
	promContentTransferMax.Set(metrics.ContentStats.Max)
	promContentTransferMean.Set(metrics.ContentStats.Mean)
	promContentTransferMedian.Set(metrics.ContentStats.Median)
	promContentTransfer95p.Set(metrics.ContentStats.P95)
	promContentTransfer99p.Set(metrics.ContentStats.P99)
	promBodySizeMin.Set(metrics.BodySize.Min)
	promBodySizeMax.Set(metrics.BodySize.Max)
	promBodySizeMean.Set(metrics.BodySize.Mean)
	promBodySizeMedian.Set(metrics.BodySize.Median)
	promBodySize95p.Set(metrics.BodySize.P95)
	promBodySize99p.Set(metrics.BodySize.P99)
	promRespSizeMin.Set(metrics.RespSize.Min)
	promRespSizeMax.Set(metrics.RespSize.Max)
	promRespSizeMean.Set(metrics.RespSize.Mean)
	promRespSizeMedian.Set(metrics.RespSize.Median)
	promRespSize95p.Set(metrics.RespSize.P95)
	promRespSize99p.Set(metrics.RespSize.P99)
	promTotalRequests.Set(float64(metrics.TotalRequests))
	promFailedRequests.Set(float64(metrics.FailedRequests))
	promRequestPerSecond.Set(metrics.RequestsPerSecond)

	if err := push.New(c.PromURL, "cassowary_load_test").
		Collector(promTCPConnMin).
		Collector(promTCPConnMax).
		Collector(promTCPConnMean).
		Collector(promTCPConnMedian).
		Collector(promTCPConn95P).
		Collector(promTCPConn99P).
		Collector(promServerProcessingMin).
		Collector(promServerProcessingMax).
		Collector(promServerProcessingMean).
		Collector(promServerProcessingMedian).
		Collector(promServerProcessing95p).
		Collector(promServerProcessing99p).
		Collector(promContentTransferMin).
		Collector(promContentTransferMax).
		Collector(promContentTransferMean).
		Collector(promContentTransferMedian).
		Collector(promContentTransfer95p).
		Collector(promContentTransfer99p).
		Collector(promBodySizeMin).
		Collector(promBodySizeMax).
		Collector(promBodySizeMean).
		Collector(promBodySizeMedian).
		Collector(promBodySize95p).
		Collector(promBodySize99p).
		Collector(promRespSizeMin).
		Collector(promRespSizeMax).
		Collector(promRespSizeMean).
		Collector(promRespSizeMedian).
		Collector(promRespSize95p).
		Collector(promRespSize99p).
		Collector(promTotalRequests).
		Collector(promFailedRequests).
		Collector(promRequestPerSecond).
		Grouping("url", c.BaseURL).
		Push(); err != nil {
		return err
	}
	return nil
}
