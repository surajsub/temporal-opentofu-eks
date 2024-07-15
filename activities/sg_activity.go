package activities

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/surajsub/temporal-opentofu-eks/models"
	"github.com/surajsub/temporal-opentofu-eks/utils"
	"go.temporal.io/sdk/activity"
)

func SGInitActivity(ctx context.Context, prov string) (string, error) {
	provisioner, engine, dir := utils.GetProvisioner(provisioner, prov)
	output, err := provisioner.Init(dir + utils.SG_DIR)
	if err != nil {
		return "", err
	}
	activity.GetLogger(ctx).Info("Calling Init Activity with engine ", engine)
	return output, nil
}

func SGApplyActivity(ctx context.Context, prov, vpcID string, vpcCdirBlock string) (string, error) {
	templog := utils.GetTemporalZap()

	fmt.Printf("the vpc cdir is set to %s\n and the vpcid is %s", vpcCdirBlock, vpcID)
	provisioner, engine, dir := utils.GetProvisioner(provisioner, prov)
	activity.GetLogger(ctx).Info("Calling Init Activity with engine ", engine)
	//output, err := provisioner.Apply(dir+utils.SG_DIR, prov, vpcID)
	activity.GetLogger(ctx).Info("Calling SG Apply Activity with engine ", engine)
	output, err := provisioner.Apply(dir+utils.SG_DIR,
		"-var", fmt.Sprintf("vpc_id=%s", vpcID))

	if err != nil {
		templog.Error(utils.SG_APPLY, "Failed to perform the Apply for Security Group")
		return "", err
	}
	return output, nil

}

func SGOutputActivity(ctx context.Context, prov string) (map[string]string, error) {
	templog := utils.GetTemporalZap()
	provisioner, engine, dir := utils.GetProvisioner(provisioner, prov)
	outputValues, err := provisioner.Output(dir + utils.SG_DIR)
	if err != nil {
		return nil, err
	}
	activity.GetLogger(ctx).Info("Calling Output Activity with engine ", engine)

	var tfOutput map[string]models.SGCommonOutput
	if err := json.Unmarshal([]byte(outputValues), &tfOutput); err != nil {
		return nil, fmt.Errorf("error unmarshaling terraform output: %v", err)
	}

	sgOutput := map[string]string{
		"sg_id":  tfOutput[utils.SGID].Value,
		"sg_arn": tfOutput[utils.SGARN].Value,
	}

	templog.Info(utils.SG_APPLY, "Security Group ID", sgOutput[utils.SGID], "Security Group ARN ", sgOutput[utils.SGARN])
	// fmt.Println("SG ACTIVITY :  the value is %s", sgOutput["sg_id"])

	return sgOutput, nil
}
