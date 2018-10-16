package cloudwatch

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	zaia_auth "github.com/youyo/zaia/auth"
)

type (
	CloudWatchParameters struct {
		DimensionName  string
		DimensionValue string
		Namespace      string
		MetricName     string
		Statistics     string
	}
)

func buildRequestParams(dimensionName, dimensionValue, namespace, metricName string) (params *cloudwatch.GetMetricStatisticsInput, err error) {
	endTime, err := time.Parse(config.TimeLayout, time.Now().UTC().Format(config.TimeLayout))
	if err != nil {
		return
	}
	startTime := endTime.Add(-600 * time.Second)
	params = func() *cloudwatch.GetMetricStatisticsInput {
		if namespace == "AWS/SES" {
			return &cloudwatch.GetMetricStatisticsInput{
				Namespace:  aws.String(namespace),
				MetricName: aws.String(metricName),
				Period:     aws.Int64(60),
				EndTime:    &endTime,
				StartTime:  &startTime,
				Statistics: []*string{
					aws.String("Minimum"),
					aws.String("Maximum"),
					aws.String("Average"),
					aws.String("SampleCount"),
					aws.String("Sum"),
				},
			}
		}
		return &cloudwatch.GetMetricStatisticsInput{
			Dimensions: []*cloudwatch.Dimension{
				{
					Name:  aws.String(dimensionName),
					Value: aws.String(dimensionValue),
				},
			},
			Namespace:  aws.String(namespace),
			MetricName: aws.String(metricName),
			Period:     aws.Int64(60),
			EndTime:    &endTime,
			StartTime:  &startTime,
			Statistics: []*string{
				aws.String("Minimum"),
				aws.String("Maximum"),
				aws.String("Average"),
				aws.String("SampleCount"),
				aws.String("Sum"),
			},
		}
	}()
	return
}

func fetchCloudWatchMetrics(cloudWatchService *cloudwatch.CloudWatch, params *cloudwatch.GetMetricStatisticsInput) (resp *cloudwatch.GetMetricStatisticsOutput, err error) {
	ctx, cancelFn := context.WithTimeout(
		context.Background(),
		config.RequestTimeout,
	)
	defer cancelFn()
	resp, err = cloudWatchService.GetMetricStatisticsWithContext(ctx, params)
	return
}

func extractValues(resp *cloudwatch.GetMetricStatisticsOutput, statistics string) (value float64, err error) {
	if len(resp.Datapoints) > 0 {
		sort.Slice(resp.Datapoints, func(i, j int) bool {
			return resp.Datapoints[i].Timestamp.Unix() > resp.Datapoints[j].Timestamp.Unix()
		})
		switch statistics {
		case "Minimum":
			value = *resp.Datapoints[0].Minimum
		case "Maximum":
			value = *resp.Datapoints[0].Maximum
		case "Average":
			value = *resp.Datapoints[0].Average
		case "SampleCount":
			value = *resp.Datapoints[0].SampleCount
		case "Sum":
			value = *resp.Datapoints[0].Sum
		default:
			err = errors.New("Statistics is not match")
		}
	} else {
		err = errors.New("Datapoint has not values")
	}
	return
}

func GetMetrics(args []string) (value string, err error) {
	namespace := args[0]
	dimensionName := args[1]
	dimensionValue := args[2]
	metricName := args[3]
	statistics := args[4]
	arn := args[5]
	region := args[6]

	sess, config := zaia_auth.Auth(arn, region)
	cloudWatchService := cloudwatch.New(sess, config)
	params, err := buildRequestParams(
		dimensionName,
		dimensionValue,
		namespace,
		metricName,
	)
	resp, err := fetchCloudWatchMetrics(cloudWatchService, params)
	if err != nil {
		return
	}
	valueFloat64, err := extractValues(resp, statistics)
	if err != nil {
		return
	}
	value = fmt.Sprint(valueFloat64)
	return
}
