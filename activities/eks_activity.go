package activities

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/surajsub/temporal-opentofu-eks/models"
	"github.com/surajsub/temporal-opentofu-eks/utils"
	"go.temporal.io/sdk/activity"
)

func EKSInitActivity(ctx context.Context, prov string) (string, error) {
	provisioner, engine, dir := utils.GetProvisioner(provisioner, prov)
	output, err := provisioner.Init(dir + utils.EKS_DIR)
	if err != nil {
		return "", err
	}
	activity.GetLogger(ctx).Info("Calling Init Activity with engine ", engine)
	return output, nil
}

func EKSApplyActivity(ctx context.Context, prov, vpcid, privateSubnetID, publicSubnetID string) (string, error) {
	activity.GetLogger(ctx).Info("The vpc input to the subnet is %s", vpcid)
	provisioner, engine, dir := utils.GetProvisioner(provisioner, prov)
	output, err := provisioner.Apply(dir+utils.EKS_DIR,
		"-var", fmt.Sprintf("private_subnet_id=%s", privateSubnetID),
		"-var", fmt.Sprintf("public_subnet_id=%s", publicSubnetID))
	if err != nil {
		return "", err
	}
	activity.GetLogger(ctx).Info("Calling Apply Activity with engine ", engine)
	return output, nil
}

func EKSOutputActivity(ctx context.Context, prov string) (map[string]string, error) {
	provisioner, engine, dir := utils.GetProvisioner(provisioner, prov)
	outputValues, err := provisioner.Output(dir + utils.EKS_DIR)
	if err != nil {
		return nil, err
	}
	activity.GetLogger(ctx).Info("Calling Output Activity with engine ", engine)

	var tfOutput map[string]models.EKSCommonOutput
	if err := json.Unmarshal([]byte(outputValues), &tfOutput); err != nil {
		return nil, fmt.Errorf("error unmarshaling terraform output: %w", err)
	}

	eksOutput := map[string]string{
		"eks_id":       tfOutput[utils.EKS_ID].Value,
		"eks_arn":      tfOutput[utils.EKS_ARN].Value,
		"eks_endpoint": tfOutput[utils.EKS_ENDPOINT].Value,
	}

	return eksOutput, nil
}
