package ec2

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/service/ec2"
	zaia_auth "github.com/youyo/zaia/auth"
)

func fetchInstanceStatus(ec2Service *ec2.EC2, ec2InstanceID string) (resp *ec2.DescribeInstanceStatusOutput, err error) {
	params := &ec2.DescribeInstanceStatusInput{
		InstanceIds: []*string{&ec2InstanceID},
	}
	ctx, cancelFn := context.WithTimeout(
		context.Background(),
		config.RequestTimeout,
	)
	defer cancelFn()
	resp, err = ec2Service.DescribeInstanceStatusWithContext(ctx, params)
	return
}

func buildMaintenanceMessage(resp *ec2.DescribeInstanceStatusOutput, noMaintenanceMessage string) (message string) {
	message = noMaintenanceMessage
	if len(resp.InstanceStatuses) > 0 {
		if len(resp.InstanceStatuses[0].Events) > 0 {
			message = fmt.Sprintf("Code: %s, Description: %s, NotAfter: %s, NotBefore: %s",
				*resp.InstanceStatuses[0].Events[0].Code,
				*resp.InstanceStatuses[0].Events[0].Description,
				*resp.InstanceStatuses[0].Events[0].NotAfter,
				*resp.InstanceStatuses[0].Events[0].NotBefore,
			)
		}
	}
	return
}

func Maintenance(args []string) (message string, err error) {
	ec2InstanceID := args[0]
	noMaintenanceMessage := args[1]
	arn := args[2]
	region := args[3]
	sess, config := zaia_auth.Auth(arn, region)
	ec2Service := ec2.New(sess, config)
	resp, err := fetchInstanceStatus(ec2Service, ec2InstanceID)
	if err != nil {
		return
	}
	message = buildMaintenanceMessage(resp, noMaintenanceMessage)
	return
}
