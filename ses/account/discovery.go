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
	if region == "us-east-1" || region == "us-west-2" || region == "eu-west-1" {
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
