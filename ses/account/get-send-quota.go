package ses_account

import (
	"strconv"

	"github.com/aws/aws-sdk-go/service/ses"
	zaia_auth "github.com/youyo/zaia/auth"
)

func GetSendQuota(args []string) (value string, err error) {
	metricName := args[0]
	arn := args[1]
	region := args[2]

	sess, config := zaia_auth.Auth(arn, region)

	sesService := ses.New(sess, config)
	input := &ses.GetSendQuotaInput{}
	result, err := sesService.GetSendQuota(input)
	if err != nil {
		return
	}

	switch metricName {
	case "Max24HourSend":
		value = strconv.FormatFloat(*result.Max24HourSend, 'f', 4, 64)
	case "MaxSendRate":
		value = strconv.FormatFloat(*result.MaxSendRate, 'f', 4, 64)
	case "SentLast24Hours":
		value = strconv.FormatFloat(*result.SentLast24Hours, 'f', 4, 64)
	}

	return
}
