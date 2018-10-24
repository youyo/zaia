package cmd

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/fujiwara/go-zabbix-get/zabbix"
	"github.com/spf13/cobra"
	zaia_cache "github.com/youyo/zaia/cache"
	zaia_cloudwatch "github.com/youyo/zaia/cloudwatch"
	zaia_ec2 "github.com/youyo/zaia/ec2"
	zaia_rds_cluster "github.com/youyo/zaia/rds/cluster"
	zaia_rds_instance "github.com/youyo/zaia/rds/instance"
	zaia_ses_account "github.com/youyo/zaia/ses/account"
)

func RootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "zaia",
		Short: "zaia",
		Run: func(cmd *cobra.Command, args []string) {
			err := runZabbixAgent("0.0.0.0:10050")
			log.Fatal(err)
		},
	}
	cobra.OnInitialize(initConfig)
	zaia_cache.InitializeCacheDb()
	return cmd
}

func runZabbixAgent(listenIp string) error {
	return zabbix.RunAgent(listenIp, func(key string) (string, error) {
		switch {
		case itemKeyIs(`agent.ping`, key):
			return "1", nil
		case itemKeyIs(`aws-integration.cloudwatch.get-metrics\[.*\]`, key):
			args := extractFromArgs([]byte(key))
			data, err := zaia_cloudwatch.GetMetrics(args)
			return data, err
		case itemKeyIs(`aws-integration.ec2.discovery\[.*\]`, key):
			args := extractFromArgs([]byte(key))
			data, err := zaia_ec2.Discovery(args)
			return data, err
		case itemKeyIs(`aws-integration.ec2.maintenance\[.*\]`, key):
			args := extractFromArgs([]byte(key))
			data, err := zaia_ec2.Maintenance(args)
			return data, err
		case itemKeyIs(`aws-integration.rds.instance.discovery\[.*\]`, key):
			args := extractFromArgs([]byte(key))
			data, err := zaia_rds_instance.Discovery(args)
			return data, err
		case itemKeyIs(`aws-integration.rds.instance.maintenance\[.*\]`, key):
			args := extractFromArgs([]byte(key))
			data, err := zaia_rds_instance.Maintenance(args)
			return data, err
		case itemKeyIs(`aws-integration.rds.cluster.discovery\[.*\]`, key):
			args := extractFromArgs([]byte(key))
			data, err := zaia_rds_cluster.Discovery(args)
			return data, err
		case itemKeyIs(`aws-integration.rds.cluster.maintenance\[.*\]`, key):
			args := extractFromArgs([]byte(key))
			data, err := zaia_rds_cluster.Maintenance(args)
			return data, err
		case itemKeyIs(`aws-integration.ses.account.discovery\[.*\]`, key):
			args := extractFromArgs([]byte(key))
			data, err := zaia_ses_account.Discovery(args)
			return data, err
		case itemKeyIs(`aws-integration.ses.account.get-send-quota\[.*\]`, key):
			args := extractFromArgs([]byte(key))
			data, err := zaia_ses_account.GetSendQuota(args)
			return data, err
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
	cmd := RootCmd()
	cmd.SetOutput(os.Stdout)
	if err := cmd.Execute(); err != nil {
		cmd.SetOutput(os.Stderr)
		cmd.Println(err)
		os.Exit(1)
	}
}

func init() {}

func initConfig() {}
