package rds_cluster

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/service/rds"
	zaia_auth "github.com/youyo/zaia/auth"
)

func fetchMaintenances(svc *rds.RDS, dbClusterIdentifier string) (resp *rds.DescribePendingMaintenanceActionsOutput, err error) {
	f := &rds.Filter{}
	f.SetName("db-cluster-id")
	f.SetValues([]*string{&dbClusterIdentifier})
	params := &rds.DescribePendingMaintenanceActionsInput{
		Filters: []*rds.Filter{f},
	}
	ctx, cancelFn := context.WithTimeout(
		context.Background(),
		config.RequestTimeout,
	)
	defer cancelFn()
	resp, err = svc.DescribePendingMaintenanceActionsWithContext(ctx, params)
	return
}

func buildMaintenanceMessage(resp *rds.DescribePendingMaintenanceActionsOutput, noMaintenanceMessage string) (message string) {
	message = noMaintenanceMessage
	if len(resp.PendingMaintenanceActions) > 0 {
		action := resp.PendingMaintenanceActions[0].PendingMaintenanceActionDetails[0]
		message = fmt.Sprintf("Action: %s, Description: %s, AutoAppliedAfterDate: %s, CurrentApplyDate: %s, ForcedApplyDate: %s",
			*action.Action,
			*action.Description,
			action.AutoAppliedAfterDate,
			action.CurrentApplyDate,
			action.ForcedApplyDate,
		)
	}
	return
}

func Maintenance(args []string) (message string, err error) {
	dbClusterIdentifier := args[0]
	noMaintenanceMessage := args[1]
	arn := args[2]
	region := args[3]
	sess, config := zaia_auth.Auth(arn, region)
	svc := rds.New(sess, config)
	resp, err := fetchMaintenances(svc, dbClusterIdentifier)
	if err != nil {
		return
	}
	message = buildMaintenanceMessage(resp, noMaintenanceMessage)
	return
}
