package resources

import (
	"github.com/surajsub/temporal-opentofu-eks/activities"
	"go.temporal.io/sdk/workflow"
	"time"
)

func EKSWorkflow(ctx workflow.Context, prov, vpc_id, privateSubnetId, publicSubnetId string) (map[string]string, error) {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 60 * time.Minute,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	err := workflow.ExecuteActivity(ctx, activities.EKSInitActivity).Get(ctx, nil)
	if err != nil {
		return nil, err
	}

	err = workflow.ExecuteActivity(ctx, activities.EKSApplyActivity, prov, vpc_id, privateSubnetId, publicSubnetId).Get(ctx, nil)
	if err != nil {
		return nil, err
	}

	var RTOutput map[string]string
	err = workflow.ExecuteActivity(ctx, activities.EKSOutputActivity).Get(ctx, &RTOutput)
	if err != nil {
		return nil, err
	}

	return RTOutput, nil
}
