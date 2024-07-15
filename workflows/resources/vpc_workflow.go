package resources

import (
	"github.com/surajsub/temporal-opentofu-eks/activities"
	"github.com/surajsub/temporal-opentofu-eks/utils"
	"go.temporal.io/sdk/workflow"
	"time"
)

func VPCWorkflow(ctx workflow.Context, prov, vpc string) (map[string]string, error) {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
	}
	templog := workflow.GetLogger(ctx)
	templog.Info(utils.VpcWorkflow, "VPC Value is ", vpc)
	ctx = workflow.WithActivityOptions(ctx, ao)

	err := workflow.ExecuteActivity(ctx, activities.VPCInitActivity, prov).Get(ctx, nil)
	if err != nil {
		return nil, err
	}

	workflow.GetLogger(ctx).Info("The value of the vpc cdir block passed is %s", vpc)
	err = workflow.ExecuteActivity(ctx, activities.VPCApplyActivity, prov, vpc).Get(ctx, nil)
	if err != nil {
		return nil, err
	}

	var vpcOutput map[string]string
	err = workflow.ExecuteActivity(ctx, activities.VPCOutputActivity, prov).Get(ctx, &vpcOutput)
	if err != nil {
		return nil, err
	}

	templog.Info(utils.VpcWorkflow, "VPC Value is ", vpcOutput["vpc_id"])
	return vpcOutput, nil
}
