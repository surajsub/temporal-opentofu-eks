package activities

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/surajsub/temporal-opentofu-eks/models"
	"github.com/surajsub/temporal-opentofu-eks/utils"
	"go.temporal.io/sdk/activity"
)

var provisioner utils.Provisioner

// This is the common vpc provisioner
func VPCInitActivity(ctx context.Context, prov string) (string, error) {
	provisioner, engine, dir := utils.GetProvisioner(provisioner, prov)
	output, err := provisioner.Init(dir + utils.VPC_DIR)
	if err != nil {
		return "", err
	}
	activity.GetLogger(ctx).Info("Calling Init Activity with engine ", engine)
	return output, nil
}

func VPCApplyActivity(ctx context.Context, prov string, vpc string) (string, error) {
	provisioner, engine, dir := utils.GetProvisioner(provisioner, prov)
	output, err := provisioner.Apply(dir+utils.VPC_DIR, "-var", fmt.Sprintf("cidr_block=%s", vpc))
	if err != nil {
		return "", err
	}
	activity.GetLogger(ctx).Info("Calling Init Activity with engine ", engine)
	return output, nil
}

func VPCOutputActivity(ctx context.Context, prov string) (map[string]string, error) {
	provisioner, engine, dir := utils.GetProvisioner(provisioner, prov)
	outputValues, err := provisioner.Output(dir + utils.VPC_DIR)
	if err != nil {
		return nil, err
	}
	activity.GetLogger(ctx).Info("Calling Init Activity with engine ", engine)

	var tfOutput map[string]models.VPCCommonOutput
	if err := json.Unmarshal([]byte(outputValues), &tfOutput); err != nil {
		return nil, fmt.Errorf("error unmarshaling terraform output: %w", err)
	}

	vpcOutput := map[string]string{
		"vpc_id":         tfOutput[utils.VPCID].Value,
		"vpc_cidr_block": tfOutput[utils.VPCCIDR].Value,
	}

	return vpcOutput, nil
}
