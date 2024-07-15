package resources

import (
	"github.com/surajsub/temporal-opentofu-eks/activities"
	"github.com/surajsub/temporal-opentofu-eks/utils"
	"go.temporal.io/sdk/workflow"
	"time"
)

func SubnetWorkflow(ctx workflow.Context, provengine, vpcid string) (map[string]string, error) {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
	}
	workflow.GetLogger(ctx).Info("The provision engine is %s", provengine)
	workflow.GetLogger(ctx).Info("THE INPUT TO THE SUBNET IS %s", vpcid)
	ctx = workflow.WithActivityOptions(ctx, ao)

	err := workflow.ExecuteActivity(ctx, activities.SubnetInitActivity, provengine).Get(ctx, nil)
	if err != nil {
		return nil, err
	}

	err = workflow.ExecuteActivity(ctx, activities.SubnetApplyActivity, provengine, vpcid).Get(ctx, nil)
	if err != nil {
		return nil, err
	}

	var subnetOutput map[string]string
	err = workflow.ExecuteActivity(ctx, activities.SubnetOutputActivity, provengine).Get(ctx, &subnetOutput)
	if err != nil {
		return nil, err
	}

	workflow.GetLogger(ctx).Info(utils.SubnetWorkflow, "Subnet  Value is ", subnetOutput["public_subnet_id"])
	return subnetOutput, nil
}
