package client

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/cloudwatch/cloudwatchiface"
)

type mockCloudWatchClient struct {
	cloudwatchiface.CloudWatchAPI
}

func (m *mockCloudWatchClient) PutMetricData(input *cloudwatch.PutMetricDataInput) (*cloudwatch.PutMetricDataOutput, error) {
	// mock response/functionality
	return nil, nil
}

func TestPutCloudwatchMetrics(t *testing.T) {
	c := Cassowary{}
	metrics := ResultMetrics{
		FailedRequests:    1,
		TotalRequests:     100,
		RequestsPerSecond: 100.10,
		TCPStats: stats{
			Min:    9.0,
			Max:    10.0,
			Mean:   10.0,
			Median: 10.0,
			P95:    10.0,
			P99:    10.0,
		},
		ProcessingStats: stats{
			Min:    1.0,
			Max:    1.0,
			Mean:   1.0,
			Median: 1.0,
			P95:    1.0,
			P99:    1.0,
		},
	}
	mockSvc := &mockCloudWatchClient{}
	_, err := c.PutCloudwatchMetrics(mockSvc, metrics)

	if err != nil {
		t.Errorf("Wanted ok but got error: %v", err)
	}

}
