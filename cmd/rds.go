package cmd

import (
	"context"

	"github.com/aws/aws-sdk-go/service/rds"
)

// discovery

type (
	RdsDiscoveryItem struct {
		DBInstanceArn        string `json:"{#DB_INSTANCE_ARN}"`
		DBInstanceIdentifier string `json:"{#DB_INSTANCE_IDENTIFIER}"`
		Engine               string `json:"{#ENGINE}"`
		ZabbixHostGroup      string `json:"{#ZABBIX_HOST_GROUP}"`
	}
	RdsDiscoveryItems []RdsDiscoveryItem
	RdsDiscoveryData  struct {
		Data RdsDiscoveryItems `json:"data"`
	}
)

func fetchRunningDBInstances(rdsService *rds.RDS) (resp *rds.DescribeDBInstancesOutput, err error) {
	params := &rds.DescribeDBInstancesInput{}
	ctx, cancelFn := context.WithTimeout(
		context.Background(),
		RequestTimeout,
	)
	defer cancelFn()
	resp, err = rdsService.DescribeDBInstancesWithContext(ctx, params)
	return
}

func buildRdsDiscoveryData(resp *rds.DescribeDBInstancesOutput, zabbixHostGroup string) (rdsDiscoveryData RdsDiscoveryData, err error) {
	var rdsDiscoveryItems RdsDiscoveryItems
	for _, v := range resp.DBInstances {
		rdsDiscoveryItems = append(rdsDiscoveryItems, RdsDiscoveryItem{
			DBInstanceArn:        *v.DBInstanceArn,
			DBInstanceIdentifier: *v.DBInstanceIdentifier,
			Engine:               *v.Engine,
			ZabbixHostGroup:      zabbixHostGroup,
		})
	}
	rdsDiscoveryData = RdsDiscoveryData{rdsDiscoveryItems}
	return
}

func rdsDiscovery(args []string) (data string, err error) {
	zabbixHostGroup := args[0]
	arn := args[1]
	region := args[2]
	sess, config := Auth(arn, region)
	rdsService := rds.New(sess, config)
	resp, err := fetchRunningDBInstances(rdsService)
	if err != nil {
		return
	}
	rdsDiscoveryData, err := buildRdsDiscoveryData(resp, zabbixHostGroup)
	if err != nil {
		return
	}
	data, err = jsonize(rdsDiscoveryData)
	return
}