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
				MetricName: aws.String("elapsed_min"),
				Unit:       aws.String("Milliseconds"),
				Value:      aws.Float64(metrics.ElapsedStats.Min),
				Dimensions: []*cloudwatch.Dimension{
					{
						Name:  aws.String("Site"),
						Value: aws.String(c.BaseURL),
					},
				},
			},
			{
				MetricName: aws.String("elapsed_max"),
				Unit:       aws.String("Milliseconds"),
				Value:      aws.Float64(metrics.ElapsedStats.Max),
				Dimensions: []*cloudwatch.Dimension{
					{
						Name:  aws.String("Site"),
						Value: aws.String(c.BaseURL),
					},
				},
			},
			{
				MetricName: aws.String("elapsed_mean"),
				Unit:       aws.String("Milliseconds"),
				Value:      aws.Float64(metrics.ElapsedStats.Mean),
				Dimensions: []*cloudwatch.Dimension{
					{
						Name:  aws.String("Site"),
						Value: aws.String(c.BaseURL),
					},
				},
			},
			{
				MetricName: aws.String("elapsed_median"),
				Unit:       aws.String("Milliseconds"),
				Value:      aws.Float64(metrics.ElapsedStats.Median),
				Dimensions: []*cloudwatch.Dimension{
					{
						Name:  aws.String("Site"),
						Value: aws.String(c.BaseURL),
					},
				},
			},
			{
				MetricName: aws.String("elapsed_95p"),
				Unit:       aws.String("Milliseconds"),
				Value:      aws.Float64(metrics.ElapsedStats.P95),
				Dimensions: []*cloudwatch.Dimension{
					{
						Name:  aws.String("Site"),
						Value: aws.String(c.BaseURL),
					},
				},
			},
			{
				MetricName: aws.String("elapsed_99p"),
				Unit:       aws.String("Milliseconds"),
				Value:      aws.Float64(metrics.ElapsedStats.P99),
				Dimensions: []*cloudwatch.Dimension{
					{
						Name:  aws.String("Site"),
						Value: aws.String(c.BaseURL),
					},
				},
			},
			{
				MetricName: aws.String("tcp_connect_min"),
				Unit:       aws.String("Milliseconds"),
				Value:      aws.Float64(metrics.TCPStats.Min),
				Dimensions: []*cloudwatch.Dimension{
					{
						Name:  aws.String("Site"),
						Value: aws.String(c.BaseURL),
					},
				},
			},
			{
				MetricName: aws.String("tcp_connect_max"),
				Unit:       aws.String("Milliseconds"),
				Value:      aws.Float64(metrics.TCPStats.Max),
				Dimensions: []*cloudwatch.Dimension{
					{
						Name:  aws.String("Site"),
						Value: aws.String(c.BaseURL),
					},
				},
			},
			{
				MetricName: aws.String("tcp_connect_mean"),
				Unit:       aws.String("Milliseconds"),
				Value:      aws.Float64(metrics.TCPStats.Mean),
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
				Value:      aws.Float64(metrics.TCPStats.Median),
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
				Value:      aws.Float64(metrics.TCPStats.P95),
				Dimensions: []*cloudwatch.Dimension{
					{
						Name:  aws.String("Site"),
						Value: aws.String(c.BaseURL),
					},
				},
			},
			{
				MetricName: aws.String("tcp_connect_99p"),
				Unit:       aws.String("Milliseconds"),
				Value:      aws.Float64(metrics.TCPStats.P99),
				Dimensions: []*cloudwatch.Dimension{
					{
						Name:  aws.String("Site"),
						Value: aws.String(c.BaseURL),
					},
				},
			},
			{
				MetricName: aws.String("server_processing_min"),
				Unit:       aws.String("Milliseconds"),
				Value:      aws.Float64(metrics.ProcessingStats.Min),
				Dimensions: []*cloudwatch.Dimension{
					{
						Name:  aws.String("Site"),
						Value: aws.String(c.BaseURL),
					},
				},
			},
			{
				MetricName: aws.String("server_processing_max"),
				Unit:       aws.String("Milliseconds"),
				Value:      aws.Float64(metrics.ProcessingStats.Max),
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
				Value:      aws.Float64(metrics.ProcessingStats.Mean),
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
				Value:      aws.Float64(metrics.ProcessingStats.Median),
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
				Value:      aws.Float64(metrics.ProcessingStats.P95),
				Dimensions: []*cloudwatch.Dimension{
					{
						Name:  aws.String("Site"),
						Value: aws.String(c.BaseURL),
					},
				},
			},
			{
				MetricName: aws.String("server_processing_99p"),
				Unit:       aws.String("Milliseconds"),
				Value:      aws.Float64(metrics.ProcessingStats.P99),
				Dimensions: []*cloudwatch.Dimension{
					{
						Name:  aws.String("Site"),
						Value: aws.String(c.BaseURL),
					},
				},
			},
			{
				MetricName: aws.String("content_transfer_min"),
				Unit:       aws.String("Milliseconds"),
				Value:      aws.Float64(metrics.ContentStats.Min),
				Dimensions: []*cloudwatch.Dimension{
					{
						Name:  aws.String("Site"),
						Value: aws.String(c.BaseURL),
					},
				},
			},
			{
				MetricName: aws.String("content_transfer_max"),
				Unit:       aws.String("Milliseconds"),
				Value:      aws.Float64(metrics.ContentStats.Max),
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
				Value:      aws.Float64(metrics.ContentStats.Mean),
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
				Value:      aws.Float64(metrics.ContentStats.Median),
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
				Value:      aws.Float64(metrics.ContentStats.P95),
				Dimensions: []*cloudwatch.Dimension{
					{
						Name:  aws.String("Site"),
						Value: aws.String(c.BaseURL),
					},
				},
			},
			{
				MetricName: aws.String("content_transfer_99p"),
				Unit:       aws.String("Milliseconds"),
				Value:      aws.Float64(metrics.ContentStats.P99),
				Dimensions: []*cloudwatch.Dimension{
					{
						Name:  aws.String("Site"),
						Value: aws.String(c.BaseURL),
					},
				},
			},
			{
				MetricName: aws.String("body_size_min"),
				Unit:       aws.String("Bytes"),
				Value:      aws.Float64(metrics.ContentStats.Min),
				Dimensions: []*cloudwatch.Dimension{
					{
						Name:  aws.String("Site"),
						Value: aws.String(c.BaseURL),
					},
				},
			},
			{
				MetricName: aws.String("body_size_max"),
				Unit:       aws.String("Bytes"),
				Value:      aws.Float64(metrics.ContentStats.Max),
				Dimensions: []*cloudwatch.Dimension{
					{
						Name:  aws.String("Site"),
						Value: aws.String(c.BaseURL),
					},
				},
			},
			{
				MetricName: aws.String("body_size_mean"),
				Unit:       aws.String("Milliseconds"),
				Value:      aws.Float64(metrics.BodySize.Mean),
				Dimensions: []*cloudwatch.Dimension{
					{
						Name:  aws.String("Site"),
						Value: aws.String(c.BaseURL),
					},
				},
			},
			{
				MetricName: aws.String("body_size_median"),
				Unit:       aws.String("Bytes"),
				Value:      aws.Float64(metrics.BodySize.Median),
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
				Value:      aws.Float64(metrics.ContentStats.P95),
				Dimensions: []*cloudwatch.Dimension{
					{
						Name:  aws.String("Site"),
						Value: aws.String(c.BaseURL),
					},
				},
			},
			{
				MetricName: aws.String("body_size_99p"),
				Unit:       aws.String("Bytes"),
				Value:      aws.Float64(metrics.BodySize.P99),
				Dimensions: []*cloudwatch.Dimension{
					{
						Name:  aws.String("Site"),
						Value: aws.String(c.BaseURL),
					},
				},
			},
			{
				MetricName: aws.String("resp_size_min"),
				Unit:       aws.String("Bytes"),
				Value:      aws.Float64(metrics.ContentStats.Min),
				Dimensions: []*cloudwatch.Dimension{
					{
						Name:  aws.String("Site"),
						Value: aws.String(c.BaseURL),
					},
				},
			},
			{
				MetricName: aws.String("resp_size_max"),
				Unit:       aws.String("Bytes"),
				Value:      aws.Float64(metrics.ContentStats.Max),
				Dimensions: []*cloudwatch.Dimension{
					{
						Name:  aws.String("Site"),
						Value: aws.String(c.BaseURL),
					},
				},
			},
			{
				MetricName: aws.String("resp_size_mean"),
				Unit:       aws.String("Milliseconds"),
				Value:      aws.Float64(metrics.RespSize.Mean),
				Dimensions: []*cloudwatch.Dimension{
					{
						Name:  aws.String("Site"),
						Value: aws.String(c.BaseURL),
					},
				},
			},
			{
				MetricName: aws.String("resp_size_median"),
				Unit:       aws.String("Bytes"),
				Value:      aws.Float64(metrics.RespSize.Median),
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
				Value:      aws.Float64(metrics.ContentStats.P95),
				Dimensions: []*cloudwatch.Dimension{
					{
						Name:  aws.String("Site"),
						Value: aws.String(c.BaseURL),
					},
				},
			},
			{
				MetricName: aws.String("resp_size_99p"),
				Unit:       aws.String("Bytes"),
				Value:      aws.Float64(metrics.RespSize.P99),
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
