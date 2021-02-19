package ses_account

import (
	"encoding/json"
)

type (
	DiscoveryItem struct {
		ZabbixHostId    string `json:"{#ZABBIX_HOST_ID}"`
		ZabbixHostGroup string `json:"{#ZABBIX_HOST_GROUP}"`
	}
	DiscoveryItems []DiscoveryItem
	DiscoveryData  struct {
		Data DiscoveryItems `json:"data"`
	}
)

func buildDiscoveryData(zabbixHostGroup, region string) (discoveryData DiscoveryData) {
	var discoveryItems DiscoveryItems

	regions := [18]string{
		"us-east-1", "us-east-2",
		"us-west-1", "us-west-2",
		"ap-south-1",
		"ap-southeast-1", "ap-southeast-2",
		"ap-northeast-1", "ap-northeast-2",
		"ca-central-1",
		"eu-central-1",
		"eu-west-1", "eu-west-2", "eu-west-3",
		"eu-north-1",
		"me-south-1",
		"sa-east-1",
		"us-gov-west-1",
	}

	for _, region := range regions {
		discoveryItems = append(discoveryItems, DiscoveryItem{
			ZabbixHostId:    zabbixHostGroup + "-" + region,
			ZabbixHostGroup: zabbixHostGroup + "-" + region,
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
	region := args[1]
	discoveryData := buildDiscoveryData(zabbixHostGroup, region)
	data, err = jsonize(discoveryData)
	return
}
