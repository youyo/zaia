package ec2

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	zaia_auth "github.com/youyo/zaia/auth"
	"github.com/youyo/zaia/zaia"
)

type (
	DiscoveryItem struct {
		InstanceID      string `json:"{#INSTANCE_ID}"`
		InstanceName    string `json:"{#INSTANCE_NAME}"`
		InstanceRole    string `json:"{#INSTANCE_ROLE}"`
		ZabbixHostGroup string `json:"{#ZABBIX_HOST_GROUP}"`
	}
	DiscoveryItems []DiscoveryItem
	DiscoveryData  struct {
		Data DiscoveryItems `json:"data"`
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
		config.RequestTimeout,
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

func buildDiscoveryData(resp *ec2.DescribeInstancesOutput, zabbixHostGroup string) (ec2DiscoveryData DiscoveryData, err error) {
	var ec2DiscoveryItems DiscoveryItems
	for _, v := range resp.Reservations {
		for _, i := range v.Instances {
			instanceName := fetchInstanceName(i)
			instanceRole := fetchInstanceRole(i)
			ec2DiscoveryItems = append(ec2DiscoveryItems, DiscoveryItem{
				InstanceID:      *i.InstanceId,
				InstanceName:    instanceName,
				InstanceRole:    instanceRole,
				ZabbixHostGroup: zabbixHostGroup,
			})
		}
	}
	ec2DiscoveryData = DiscoveryData{ec2DiscoveryItems}
	return
}

func Discovery(args []string) (data string, err error) {
	zabbixHostGroup := args[0]
	arn := args[1]
	region := args[2]
	sess, config := zaia_auth.Auth(arn, region)
	ec2Service := ec2.New(sess, config)
	resp, err := fetchRunningInstances(ec2Service)
	if err != nil {
		return
	}
	ec2DiscoveryData, err := buildDiscoveryData(resp, zabbixHostGroup)
	if err != nil {
		return
	}
	data, err = zaia.Jsonize(ec2DiscoveryData)
	return
}
