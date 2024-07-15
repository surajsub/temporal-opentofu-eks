package activities

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/surajsub/temporal-opentofu-eks/models"
	"github.com/surajsub/temporal-opentofu-eks/utils"
	"go.temporal.io/sdk/activity"
)

// This is the common subnet provisioner

func SubnetInitActivity(ctx context.Context, prov string) (string, error) {
	provisioner, engine, dir := utils.GetProvisioner(provisioner, prov)
	output, err := provisioner.Init(dir + utils.SUBNET_DIR)
	if err != nil {
		return "", err
	}
	activity.GetLogger(ctx).Info("Calling Init Activity with engine ", engine)
	return output, nil
}

func SubnetApplyActivity(ctx context.Context, prov string, vpcid string) (string, error) {
	activity.GetLogger(ctx).Info("The vpc input to the subnet is %s", vpcid)
	provisioner, engine, dir := utils.GetProvisioner(provisioner, prov)
	output, err := provisioner.Apply(dir+utils.SUBNET_DIR, "-var", fmt.Sprintf("vpc_id=%s", vpcid))
	if err != nil {
		return "", err
	}
	activity.GetLogger(ctx).Info("Calling Init Activity with engine ", engine)
	return output, nil
}

func SubnetOutputActivity(ctx context.Context, prov string) (map[string]string, error) {
	provisioner, engine, dir := utils.GetProvisioner(provisioner, prov)
	outputValues, err := provisioner.Output(dir + utils.SUBNET_DIR)
	if err != nil {
		return nil, err
	}
	activity.GetLogger(ctx).Info("Calling Init Activity with engine ", engine)

	var tfOutput map[string]models.SubnetCommonOutput
	if err := json.Unmarshal([]byte(outputValues), &tfOutput); err != nil {
		return nil, fmt.Errorf("error unmarshaling  output: %w", err)
	}

	subnetOutput := map[string]string{
		"private_subnet_id": tfOutput[utils.PRIVATE_SUBNET_ID].Value,
		"public_subnet_id":  tfOutput[utils.PUBLIC_SUBNET_ID].Value,
	}

	return subnetOutput, nil
}
