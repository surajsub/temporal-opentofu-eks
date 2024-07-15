package resources

import (
	"github.com/surajsub/temporal-opentofu-eks/activities"
	"github.com/surajsub/temporal-opentofu-eks/utils"
	"go.temporal.io/sdk/workflow"
)

func IGWWorkflow(ctx workflow.Context, provengine, vpcID string) (map[string]string, error) {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: utils.IgwTimeOut,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	err := workflow.ExecuteActivity(ctx, activities.IGWInitActivity).Get(ctx, nil)
	if err != nil {
		return nil, err
	}

	err = workflow.ExecuteActivity(ctx, activities.IGWApplyActivity, provengine, vpcID).Get(ctx, nil)
	if err != nil {
		return nil, err
	}

	var iGWOutput map[string]string
	err = workflow.ExecuteActivity(ctx, activities.IGWOutputActivity).Get(ctx, &iGWOutput)
	if err != nil {
		return nil, err
	}

	return iGWOutput, nil
}
