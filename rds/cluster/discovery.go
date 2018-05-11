package rds_cluster

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go/service/rds"
	zaia_auth "github.com/youyo/zaia/auth"
)

type (
	DiscoveryItem struct {
		DBClusterIdentifier string `json:"{#DB_CLUSTER_IDENTIFIER}"`
		Engine              string `json:"{#ENGINE}"`
		ZabbixHostGroup     string `json:"{#ZABBIX_HOST_GROUP}"`
	}
	DiscoveryItems []DiscoveryItem
	DiscoveryData  struct {
		Data DiscoveryItems `json:"data"`
	}
)

func fetchDBClusters(svc *rds.RDS) (resp *rds.DescribeDBClustersOutput, err error) {
	params := &rds.DescribeDBClustersInput{}
	ctx, cancelFn := context.WithTimeout(
		context.Background(),
		config.RequestTimeout,
	)
	defer cancelFn()
	resp, err = svc.DescribeDBClustersWithContext(ctx, params)
	return
}

func buildDiscoveryData(resp *rds.DescribeDBClustersOutput, zabbixHostGroup string) (discoveryData DiscoveryData, err error) {
	var discoveryItems DiscoveryItems
	for _, v := range resp.DBClusters {
		discoveryItems = append(discoveryItems, DiscoveryItem{
			DBClusterIdentifier: *v.DBClusterIdentifier,
			Engine:              *v.Engine,
			ZabbixHostGroup:     zabbixHostGroup,
		})
	}
	discoveryData = DiscoveryData{discoveryItems}
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
	resp, err := fetchDBClusters(rdsService)
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
