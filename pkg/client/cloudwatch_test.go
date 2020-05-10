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
		TCPStats: tcpStats{
			TCPMean:   10.0,
			TCPMedian: 10.0,
			TCP95p:    10.0,
		},
		ProcessingStats: serverProcessingStats{
			ServerProcessingMean:   1.0,
			ServerProcessingMedian: 1.0,
			ServerProcessing95p:    1.0,
		},
	}
	mockSvc := &mockCloudWatchClient{}
	_, err := c.PutCloudwatchMetrics(mockSvc, metrics)

	if err != nil {
		t.Errorf("Wanted ok but got error: %v", err)
	}

}
