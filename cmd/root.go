package cmd

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/fujiwara/go-zabbix-get/zabbix"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:     "zabbix-aws-integration-agent",
	Version: version,
	Short:   "zabbix-aws-integration-agent",
	//Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		err := runZabbixAgent("0.0.0.0:10050")
		log.Fatal(err)
	},
}

func runZabbixAgent(listenIp string) error {
	return zabbix.RunAgent(listenIp, func(key string) (string, error) {
		switch {
		case itemKeyIs(`aws-integration.ec2.discovery\[.*\]`, key):
			args := extractFromArgs([]byte(key))
			data, err := ec2Discovery(args)
			return data, err
		case itemKeyIs(`aws-integration.ec2.maintenance\[.*\]`, key):
			args := extractFromArgs([]byte(key))
			data, err := ec2Maintenance(args)
			return data, err
		case itemKeyIs(`aws-integration.cloudwatch.get-metrics\[.*\]`, key):
			args := extractFromArgs([]byte(key))
			data, err := cloudWatchGetMetrics(args)
			return data, err
		case itemKeyIs(`aws-integration.rds.discovery\[.*\]`, key):
			args := extractFromArgs([]byte(key))
			data, err := rdsDiscovery(args)
			return data, err
		case itemKeyIs(`agent.ping`, key):
			return "1", nil
		default:
			return "", fmt.Errorf("not supported")
		}
	})
}

func itemKeyIs(pattern, key string) bool {
	return regexp.MustCompile(pattern).Match([]byte(key))
}

func extractFromArgs(b []byte) []string {
	assigned := regexp.MustCompile(`.*\[(.*)\]`)
	group := assigned.FindSubmatch(b)
	args := strings.Split(string(group[1]), ",")
	return args
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	initializeCacheDb()
}

func initConfig() {}
