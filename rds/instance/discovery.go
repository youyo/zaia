package rds_instance

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go/service/rds"
	zaia_auth "github.com/youyo/zaia/auth"
)

type (
	DiscoveryItem struct {
		DBInstanceIdentifier string `json:"{#DB_INSTANCE_IDENTIFIER}"`
		Engine               string `json:"{#ENGINE}"`
		ZabbixHostGroup      string `json:"{#ZABBIX_HOST_GROUP}"`
	}
	DiscoveryItems []DiscoveryItem
	DiscoveryData  struct {
		Data DiscoveryItems `json:"data"`
	}
)

func fetchRunningDBInstances(rdsService *rds.RDS) (resp *rds.DescribeDBInstancesOutput, err error) {
	params := &rds.DescribeDBInstancesInput{}
	ctx, cancelFn := context.WithTimeout(
		context.Background(),
		config.RequestTimeout,
	)
	defer cancelFn()
	resp, err = rdsService.DescribeDBInstancesWithContext(ctx, params)
	return
}

func buildDiscoveryData(resp *rds.DescribeDBInstancesOutput, zabbixHostGroup string) (rdsDiscoveryData DiscoveryData, err error) {
	var rdsDiscoveryItems DiscoveryItems
	for _, v := range resp.DBInstances {
		rdsDiscoveryItems = append(rdsDiscoveryItems, DiscoveryItem{
			DBInstanceIdentifier: *v.DBInstanceIdentifier,
			Engine:               *v.Engine,
			ZabbixHostGroup:      zabbixHostGroup,
		})
	}
	rdsDiscoveryData = DiscoveryData{rdsDiscoveryItems}
	return
}

func jsonize(data interface{}) (s string, err error) {
	b, err := json.Marshal(data)
	if err != nil {
		return
	}
	s = string(b)
	return
}

func Discovery(args []string) (data string, err error) {
	zabbixHostGroup := args[0]
	arn := args[1]
	region := args[2]
	sess, config := zaia_auth.Auth(arn, region)
	rdsService := rds.New(sess, config)
	resp, err := fetchRunningDBInstances(rdsService)
	if err != nil {
		return
	}
	discoveryData, err := buildDiscoveryData(resp, zabbixHostGroup)
	if err != nil {
		return
	}
	data, err = jsonize(discoveryData)
	return
}
