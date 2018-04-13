package cmd

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// discovery

type (
	Ec2DiscoveryItem struct {
		InstanceID        string `json:"{#INSTANCE_ID}"`
		InstanceName      string `json:"{#INSTANCE_NAME}"`
		InstanceRole      string `json:"{#INSTANCE_ROLE}"`
		InstancePublicIp  string `json:"{#INSTANCE_PUBLIC_IP}"`
		InstancePrivateIp string `json:"{#INSTANCE_PRIVATE_IP}"`
		ZabbixHostGroup   string `json:"{#ZABBIX_HOST_GROUP}"`
	}
	Ec2DiscoveryItems []Ec2DiscoveryItem
	Ec2DiscoveryData  struct {
		Data Ec2DiscoveryItems `json:"data"`
	}
)

func fetchRunningInstances(ec2Service *ec2.EC2) (resp *ec2.DescribeInstancesOutput, err error) {
	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("instance-state-name"),
				Values: []*string{aws.String("running")},
			},
		},
	}
	ctx, cancelFn := context.WithTimeout(
		context.Background(),
		RequestTimeout,
	)
	defer cancelFn()
	resp, err = ec2Service.DescribeInstancesWithContext(ctx, params)
	return
}

func fetchInstanceName(ec2Instance *ec2.Instance) (instnaceName string) {
	for _, tag := range ec2Instance.Tags {
		if *tag.Key == "Name" {
			instnaceName = *tag.Value
		}
	}
	return
}

func fetchInstanceRole(ec2Instance *ec2.Instance) (instnaceRole string) {
	for _, tag := range ec2Instance.Tags {
		if *tag.Key == "Role" {
			instnaceRole = *tag.Value
		}
	}
	return
}

func buildEc2DiscoveryData(resp *ec2.DescribeInstancesOutput, zabbixHostGroup string) (ec2DiscoveryData Ec2DiscoveryData, err error) {
	var ec2DiscoveryItems Ec2DiscoveryItems
	for _, v := range resp.Reservations {
		for _, i := range v.Instances {
			instanceName := fetchInstanceName(i)
			instanceRole := fetchInstanceRole(i)
			ec2DiscoveryItems = append(ec2DiscoveryItems, Ec2DiscoveryItem{
				InstanceID:        *i.InstanceId,
				InstanceName:      instanceName,
				InstanceRole:      instanceRole,
				InstancePublicIp:  *i.PublicIpAddress,
				InstancePrivateIp: *i.PrivateIpAddress,
				ZabbixHostGroup:   zabbixHostGroup,
			})
		}
	}
	ec2DiscoveryData = Ec2DiscoveryData{ec2DiscoveryItems}
	return
}

func ec2Discovery(args []string) (data string, err error) {
	zabbixHostGroup := args[0]
	arn := args[1]
	region := args[2]
	sess, config := Auth(arn, region)
	ec2Service := ec2.New(sess, config)
	resp, err := fetchRunningInstances(ec2Service)
	if err != nil {
		return
	}
	ec2DiscoveryData, err := buildEc2DiscoveryData(resp, zabbixHostGroup)
	if err != nil {
		return
	}
	data, err = jsonize(ec2DiscoveryData)
	return
}

// maintenance

func fetchInstanceStatus(ec2Service *ec2.EC2, ec2InstanceID string) (resp *ec2.DescribeInstanceStatusOutput, err error) {
	params := &ec2.DescribeInstanceStatusInput{
		InstanceIds: []*string{&ec2InstanceID},
	}
	ctx, cancelFn := context.WithTimeout(
		context.Background(),
		RequestTimeout,
	)
	defer cancelFn()
	resp, err = ec2Service.DescribeInstanceStatusWithContext(ctx, params)
	return
}

func buildMaintenanceMessage(resp *ec2.DescribeInstanceStatusOutput, noMaintenanceMessage string) (message string) {
	message = noMaintenanceMessage
	if len(resp.InstanceStatuses[0].Events) > 0 {
		message = fmt.Sprintf("Code: %s, Description: %s, NotAfter: %s, NotBefore: %s",
			*resp.InstanceStatuses[0].Events[0].Code,
			*resp.InstanceStatuses[0].Events[0].Description,
			*resp.InstanceStatuses[0].Events[0].NotAfter,
			*resp.InstanceStatuses[0].Events[0].NotBefore,
		)
	}
	return
}

func ec2Maintenance(args []string) (message string, err error) {
	ec2InstanceID := args[0]
	noMaintenanceMessage := args[1]
	arn := args[2]
	region := args[3]
	sess, config := Auth(arn, region)
	ec2Service := ec2.New(sess, config)
	resp, err := fetchInstanceStatus(ec2Service, ec2InstanceID)
	if err != nil {
		return
	}
	message = buildMaintenanceMessage(resp, noMaintenanceMessage)
	return
}
