package resources

import (
	"github.com/surajsub/temporal-opentofu-eks/activities"
	"github.com/surajsub/temporal-opentofu-eks/utils"
	"go.temporal.io/sdk/workflow"
	"time"
)

func NATWorkflow(ctx workflow.Context, provengine, publicSubnetId string) (map[string]string, error) {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	templog := workflow.GetLogger(ctx)
	templog.Info(utils.NatWorkflow, "Public SubnetID ", publicSubnetId)

	err := workflow.ExecuteActivity(ctx, activities.NATInitActivity).Get(ctx, nil)
	if err != nil {
		return nil, err
	}

	err = workflow.ExecuteActivity(ctx, activities.NATApplyActivity, provengine, publicSubnetId).Get(ctx, nil)
	if err != nil {
		return nil, err
	}

	var NATOutput map[string]string
	err = workflow.ExecuteActivity(ctx, activities.NATOutputActivity).Get(ctx, &NATOutput)
	if err != nil {
		return nil, err
	}

	return NATOutput, nil
}
