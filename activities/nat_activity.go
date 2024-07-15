package activities

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/surajsub/temporal-opentofu-eks/models"
	_ "github.com/surajsub/temporal-opentofu-eks/models"
	"github.com/surajsub/temporal-opentofu-eks/utils"
	"go.temporal.io/sdk/activity"
)

func NATInitActivity(ctx context.Context, prov string) (string, error) {
	provisioner, engine, dir := utils.GetProvisioner(provisioner, prov)
	output, err := provisioner.Init(dir + utils.NAT_DIR)
	if err != nil {
		return "", err
	}
	activity.GetLogger(ctx).Info("Calling Init Activity with engine ", engine)
	return output, nil
}

func NATApplyActivity(ctx context.Context, prov, publicSubnetID string) (string, error) {
	activity.GetLogger(ctx).Info("Starting the NAT Apply Activity  with input [public_subnet_id]  %s", publicSubnetID)
	provisioner, engine, dir := utils.GetProvisioner(provisioner, prov)
	output, err := provisioner.Apply(dir+utils.NAT_DIR, "-var", fmt.Sprintf("subnet_id=%s", publicSubnetID))
	if err != nil {
		return "", err
	}
	activity.GetLogger(ctx).Info("Calling Init Activity with engine ", engine)
	return output, nil
}

func NATOutputActivity(ctx context.Context, prov string) (map[string]string, error) {
	provisioner, engine, dir := utils.GetProvisioner(provisioner, prov)
	outputValues, err := provisioner.Output(dir + utils.NAT_DIR)
	if err != nil {
		return nil, err
	}
	activity.GetLogger(ctx).Info("Calling Init Activity with engine ", engine)

	var tfOutput map[string]models.NATCommonOutput
	if err := json.Unmarshal([]byte(outputValues), &tfOutput); err != nil { //nolint:typecheck
		return nil, fmt.Errorf("error unmarshaling terraform output: %v", err)
	}

	natOutput := map[string]string{
		"nat_id":            tfOutput[utils.NATID].Value,
		"nat_gateway_id":    tfOutput[utils.NATGATEWAYID].Value,
		"nat_allocation_id": tfOutput[utils.NATALLOCATIONID].Value,
	}

	return natOutput, nil
}
