package activities

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/surajsub/temporal-opentofu-eks/models"
	"github.com/surajsub/temporal-opentofu-eks/utils"
	"go.temporal.io/sdk/activity"
)

func IGWInitActivity(ctx context.Context, prov string) (string, error) {
	provisioner, engine, dir := utils.GetProvisioner(provisioner, prov)
	output, err := provisioner.Init(dir + utils.IGW_DIR)
	if err != nil {
		return "", err
	}
	activity.GetLogger(ctx).Info("Calling Init Activity with engine ", engine)
	return output, nil
}

func IGWApplyActivity(ctx context.Context, prov, vpcid string) (string, error) {
	activity.GetLogger(ctx).Info("The vpc input to the subnet is %s", vpcid)
	provisioner, engine, dir := utils.GetProvisioner(provisioner, prov)
	output, err := provisioner.Apply(dir+utils.IGW_DIR, "-var", fmt.Sprintf("vpc_id=%s", vpcid))
	if err != nil {
		return "", err
	}
	activity.GetLogger(ctx).Info("Calling Init Activity with engine ", engine)
	return output, nil
}

func IGWOutputActivity(ctx context.Context, prov string) (map[string]string, error) {
	provisioner, engine, dir := utils.GetProvisioner(provisioner, prov)
	outputValues, err := provisioner.Output(dir + utils.IGW_DIR)
	if err != nil {
		return nil, err
	}
	activity.GetLogger(ctx).Info("Calling Output Activity with engine ", engine)

	var tfOutput map[string]models.IGWCommonOutput
	if err := json.Unmarshal([]byte(outputValues), &tfOutput); err != nil {
		return nil, fmt.Errorf("error unmarshaling terraform output: %v", err)
	}

	igwOutput := map[string]string{
		"igw_id":  tfOutput["igw_id"].Value,
		"igw_arn": tfOutput["igw_arn"].Value,
	}

	return igwOutput, nil
}
