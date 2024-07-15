package activities

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/surajsub/temporal-opentofu-eks/models"
	"github.com/surajsub/temporal-opentofu-eks/utils"
	"go.temporal.io/sdk/activity"
)

func RTInitActivity(ctx context.Context, prov string) (string, error) {
	provisioner, engine, dir := utils.GetProvisioner(provisioner, prov)
	output, err := provisioner.Init(dir + utils.RT_DIR)
	if err != nil {
		return "", err
	}
	activity.GetLogger(ctx).Info("Calling Init Activity with engine ", engine)
	return output, nil
}

func RTApplyActivity(ctx context.Context, prov, vpcId, igwId, natId, privateSubnetId, publicSubnetId string) (string, error) {
	// output, err := utils.RunTFApplyCommand(utils.SUBNET_TF_DIRECTORY)
	activity.GetLogger(ctx).Info("The vpc input to the Route Table Activity is %s", vpcId)
	provisioner, engine, dir := utils.GetProvisioner(provisioner, prov)
	activity.GetLogger(ctx).Info("Calling Init Activity with engine ", engine)
	output, err := provisioner.Apply(dir+utils.RT_DIR,
		"-var", fmt.Sprintf("vpc_id=%s", vpcId),
		"-var", fmt.Sprintf("igw_id=%s", igwId),
		"-var", fmt.Sprintf("nat_id=%s", natId),
		"-var", fmt.Sprintf("private_subnet_id=%s", privateSubnetId),
		"-var", fmt.Sprintf("public_subnet_id=%s", publicSubnetId))

	if err != nil {
		return "", err
	}
	return output, nil
}

func RTOutputActivity(ctx context.Context, prov string) (map[string]string, error) {
	provisioner, engine, dir := utils.GetProvisioner(provisioner, prov)
	outputValues, err := provisioner.Output(dir + utils.RT_DIR)
	if err != nil {
		return nil, err
	}
	activity.GetLogger(ctx).Info("Calling Output Activity with engine ", engine)

	var tfOutput map[string]models.RTCommonOutput
	if err := json.Unmarshal([]byte(outputValues), &tfOutput); err != nil {
		return nil, fmt.Errorf("error unmarshaling  output: %v", err)
	}

	rtOutput := map[string]string{
		"rt_public_id":  tfOutput[utils.RTPUBLICID].Value,
		"rt_private_id": tfOutput[utils.RTPRIVATEID].Value,
	}

	return rtOutput, nil
}
