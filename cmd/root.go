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
)

func RootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "zaia",
		Version: version,
		Short:   "zaia",
		//Long: ``,
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
		case itemKeyIs(`aws-integration.ec2.discovery\[.*\]`, key):
			args := extractFromArgs([]byte(key))
			data, err := zaia_ec2.Discovery(args)
			return data, err
		case itemKeyIs(`aws-integration.ec2.maintenance\[.*\]`, key):
			args := extractFromArgs([]byte(key))
			data, err := zaia_ec2.Maintenance(args)
			return data, err
		case itemKeyIs(`aws-integration.cloudwatch.get-metrics\[.*\]`, key):
			args := extractFromArgs([]byte(key))
			data, err := zaia_cloudwatch.GetMetrics(args)
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
