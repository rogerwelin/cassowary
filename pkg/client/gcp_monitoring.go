package client

import (
	"context"
	"errors"
	"os"
	"time"

	monitoring "cloud.google.com/go/monitoring/apiv3"
	googlepb "github.com/golang/protobuf/ptypes/timestamp"
	metricpb "google.golang.org/genproto/googleapis/api/metric"
	monitoredrespb "google.golang.org/genproto/googleapis/api/monitoredres"
	monitoringpb "google.golang.org/genproto/googleapis/monitoring/v3"
)

// https://cloud.google.com/monitoring/docs/reference/libraries#client-libraries-install-go

// PutGCPMonitoringMetrics exports metrics to GCP Monitoring
func (c *Cassowary) PutGCPMonitoringMetrics(metrics ResultMetrics) error {

	ctx := context.Background()

	// creates a client
	client, err := monitoring.NewMetricClient(ctx)
	if err != nil {
		return err
	}

	gcpProjectID := os.Getenv("GCP_PROJECT")
	if gcpProjectID == "" {
		return errors.New("No GCP Project ID set")
	}

	tcpMean := &monitoringpb.Point{
		Interval: &monitoringpb.TimeInterval{
			EndTime: &googlepb.Timestamp{
				Seconds: time.Now().Unix(),
			},
		},
		Value: &monitoringpb.TypedValue{
			Value: &monitoringpb.TypedValue_DoubleValue{
				DoubleValue: metrics.TCPStats.TCPMedian,
			},
		},
	}

	// Writes time series data.
	if err := client.CreateTimeSeries(ctx, &monitoringpb.CreateTimeSeriesRequest{
		Name: monitoring.MetricProjectPath(gcpProjectID),
		TimeSeries: []*monitoringpb.TimeSeries{
			{
				Metric: &metricpb.Metric{
					Type: "custom.googleapis.com/cassowary/tcp_connect_mean",
					Labels: map[string]string{
						"site": c.BaseURL,
					},
				},
				Resource: &monitoredrespb.MonitoredResource{
					Type: "global",
					Labels: map[string]string{
						"project_id": gcpProjectID,
					},
				},
				Points: []*monitoringpb.Point{
					tcpMean,
				},
			},
			{
				Metric: &metricpb.Metric{
					Type: "custom.googleapis.com/cassowary/tcp_connect_median",
					Labels: map[string]string{
						"site": c.BaseURL,
					},
				},
				Resource: &monitoredrespb.MonitoredResource{
					Type: "global",
					Labels: map[string]string{
						"project_id": gcpProjectID,
					},
				},
				Points: []*monitoringpb.Point{
					tcpMean,
				},
			},
		},
	}); err != nil {
		return err
	}

	// Closes the client and flushes the data to Stackdriver.
	if err := client.Close(); err != nil {
		return err
	}

	return nil
}
