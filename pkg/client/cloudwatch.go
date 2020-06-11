package client

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/cloudwatch/cloudwatchiface"
)

// PutCloudwatchMetrics exports metrics to AWS Cloudwatch
func (c *Cassowary) PutCloudwatchMetrics(svc cloudwatchiface.CloudWatchAPI, metrics ResultMetrics) (*cloudwatch.PutMetricDataOutput, error) {
	resp, err := svc.PutMetricData(&cloudwatch.PutMetricDataInput{
		Namespace: aws.String("Cassowary/Metrics"),
		MetricData: []*cloudwatch.MetricDatum{
			{
				MetricName: aws.String("tcp_connect_mean"),
				Unit:       aws.String("Milliseconds"),
				Value:      aws.Float64(metrics.TCPStats.TCPMean),
				Dimensions: []*cloudwatch.Dimension{
					{
						Name:  aws.String("Site"),
						Value: aws.String(c.BaseURL),
					},
				},
			},
			{
				MetricName: aws.String("tcp_connect_median"),
				Unit:       aws.String("Milliseconds"),
				Value:      aws.Float64(metrics.TCPStats.TCPMedian),
				Dimensions: []*cloudwatch.Dimension{
					{
						Name:  aws.String("Site"),
						Value: aws.String(c.BaseURL),
					},
				},
			},
			{
				MetricName: aws.String("tcp_connect_95p"),
				Unit:       aws.String("Milliseconds"),
				Value:      aws.Float64(metrics.TCPStats.TCP95p),
				Dimensions: []*cloudwatch.Dimension{
					{
						Name:  aws.String("Site"),
						Value: aws.String(c.BaseURL),
					},
				},
			},
			{
				MetricName: aws.String("server_processing_mean"),
				Unit:       aws.String("Milliseconds"),
				Value:      aws.Float64(metrics.ProcessingStats.ServerProcessingMean),
				Dimensions: []*cloudwatch.Dimension{
					{
						Name:  aws.String("Site"),
						Value: aws.String(c.BaseURL),
					},
				},
			},
			{
				MetricName: aws.String("server_processing_median"),
				Unit:       aws.String("Milliseconds"),
				Value:      aws.Float64(metrics.ProcessingStats.ServerProcessingMedian),
				Dimensions: []*cloudwatch.Dimension{
					{
						Name:  aws.String("Site"),
						Value: aws.String(c.BaseURL),
					},
				},
			},
			{
				MetricName: aws.String("server_processing_95p"),
				Unit:       aws.String("Milliseconds"),
				Value:      aws.Float64(metrics.ProcessingStats.ServerProcessing95p),
				Dimensions: []*cloudwatch.Dimension{
					{
						Name:  aws.String("Site"),
						Value: aws.String(c.BaseURL),
					},
				},
			},
			{
				MetricName: aws.String("content_transfer_mean"),
				Unit:       aws.String("Milliseconds"),
				Value:      aws.Float64(metrics.ContentStats.ContentTransferMean),
				Dimensions: []*cloudwatch.Dimension{
					{
						Name:  aws.String("Site"),
						Value: aws.String(c.BaseURL),
					},
				},
			},
			{
				MetricName: aws.String("content_transfer_median"),
				Unit:       aws.String("Milliseconds"),
				Value:      aws.Float64(metrics.ContentStats.ContentTransferMedian),
				Dimensions: []*cloudwatch.Dimension{
					{
						Name:  aws.String("Site"),
						Value: aws.String(c.BaseURL),
					},
				},
			},
			{
				MetricName: aws.String("content_transfer_95p"),
				Unit:       aws.String("Milliseconds"),
				Value:      aws.Float64(metrics.ContentStats.ContentTransfer95p),
				Dimensions: []*cloudwatch.Dimension{
					{
						Name:  aws.String("Site"),
						Value: aws.String(c.BaseURL),
					},
				},
			},
			{
				MetricName: aws.String("total_requests"),
				Unit:       aws.String("Count"),
				Value:      aws.Float64(float64(metrics.TotalRequests)),
				Dimensions: []*cloudwatch.Dimension{
					{
						Name:  aws.String("Site"),
						Value: aws.String(c.BaseURL),
					},
				},
			},
			{
				MetricName: aws.String("failed_requests"),
				Unit:       aws.String("Count"),
				Value:      aws.Float64(float64(metrics.FailedRequests)),
				Dimensions: []*cloudwatch.Dimension{
					{
						Name:  aws.String("Site"),
						Value: aws.String(c.BaseURL),
					},
				},
			},
			{
				MetricName: aws.String("requests_per_second"),
				Unit:       aws.String("Count/Second"),
				Value:      aws.Float64(metrics.RequestsPerSecond),
				Dimensions: []*cloudwatch.Dimension{
					{
						Name:  aws.String("Site"),
						Value: aws.String(c.BaseURL),
					},
				},
			},
		},
	})

	if err != nil {
		return nil, err
	}

	return resp, nil
}
