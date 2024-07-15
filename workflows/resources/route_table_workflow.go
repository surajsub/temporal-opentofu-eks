package resources

import (
	"github.com/surajsub/temporal-opentofu-eks/activities"
	"github.com/surajsub/temporal-opentofu-eks/utils"
	"go.temporal.io/sdk/workflow"
	"time"
)

func RouteTableWorkflow(ctx workflow.Context, prov, vpc_id, internet_gateway_id, nat_gateway_id, private_subnet_id, public_subnet_id string) (map[string]string, error) {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	templog := workflow.GetLogger(ctx)
	templog.Info(utils.RtWorkflow, "Internet Gateway ID ", internet_gateway_id)

	err := workflow.ExecuteActivity(ctx, activities.RTInitActivity).Get(ctx, nil)
	if err != nil {
		return nil, err
	}

	err = workflow.ExecuteActivity(ctx, activities.RTApplyActivity, prov, vpc_id, internet_gateway_id, nat_gateway_id, private_subnet_id, public_subnet_id).Get(ctx, nil)
	if err != nil {
		return nil, err
	}

	var RTOutput map[string]string
	err = workflow.ExecuteActivity(ctx, activities.RTOutputActivity).Get(ctx, &RTOutput)
	if err != nil {
		return nil, err
	}

	return RTOutput, nil
}
