package resources

import (
	"github.com/surajsub/temporal-opentofu-eks/activities"
	"go.temporal.io/sdk/workflow"
	"time"
)

func SGWorkflow(ctx workflow.Context, prov, vpcID, vpcdir string) (map[string]string, error) {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	err := workflow.ExecuteActivity(ctx, activities.SGInitActivity).Get(ctx, nil)
	if err != nil {
		return nil, err
	}

	err = workflow.ExecuteActivity(ctx, activities.SGApplyActivity, prov, vpcID, vpcdir).Get(ctx, nil)
	if err != nil {
		return nil, err
	}

	var SGOutput map[string]string
	err = workflow.ExecuteActivity(ctx, activities.SGOutputActivity).Get(ctx, &SGOutput)
	if err != nil {
		return nil, err
	}

	return SGOutput, nil
}
