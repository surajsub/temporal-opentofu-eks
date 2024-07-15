package workflows

import (
	"github.com/surajsub/temporal-opentofu-eks/utils"
	"github.com/surajsub/temporal-opentofu-eks/workflows/resources" //nolint:typecheck
	"go.temporal.io/sdk/workflow"
	"log"
	"time"
)

//nolint:funlen
func ParentWorkflow(ctx workflow.Context, vpc, prov string) (map[string]interface{}, error) {
	cwo := workflow.ChildWorkflowOptions{
		WorkflowExecutionTimeout: time.Hour,
		WorkflowRunTimeout:       time.Minute * 30,
	}
	ctx = workflow.WithChildOptions(ctx, cwo)
	workflowID := workflow.GetInfo(ctx).OriginalRunID

	log.Printf("Printing the wortkflow id from the PARENT %s\n", workflowID)

	// Start VPC Workflow

	var vpcOutput map[string]string
	err := workflow.ExecuteChildWorkflow(ctx, resources.VPCWorkflow, prov, vpc).Get(ctx, &vpcOutput)
	if err != nil {
		return nil, err
	}
	workflow.GetLogger(ctx).Info("VPC created", "vpc_id", vpcOutput[utils.VPCID], "and the vpc cidr is ", vpcOutput["vpc_cidr_block"])

	// Start Subnet Workflow
	var subnetOutput map[string]string
	err = workflow.ExecuteChildWorkflow(ctx, resources.SubnetWorkflow, prov, vpcOutput["vpc_id"]).Get(ctx, &subnetOutput)
	if err != nil {
		return nil, err
	}
	workflow.GetLogger(ctx).Info("Subnet created", "private_subnet_id", subnetOutput[utils.PRIVATE_SUBNET_ID])
	workflow.GetLogger(ctx).Info("Public Subnet Created", "public_subnet_id", subnetOutput[utils.PUBLIC_SUBNET_ID])

	// Start IGW Workflow

	var igwOutput map[string]string
	err = workflow.ExecuteChildWorkflow(ctx, resources.IGWWorkflow, prov, vpcOutput["vpc_id"]).Get(ctx, &igwOutput)
	if err != nil {
		return nil, err
	}
	workflow.GetLogger(ctx).Info("IGW Created", "igw_id", igwOutput[utils.IGWID])

	// Start the NAT Workflow

	var natoutput map[string]string
	err = workflow.ExecuteChildWorkflow(ctx, resources.NATWorkflow, prov, subnetOutput["public_subnet_id"]).Get(ctx, &natoutput)
	if err != nil {
		workflow.GetLogger(ctx).Info("Failed to create the NAT ", "public subnet id", subnetOutput["public_subnet_id"])
	}

	// Start the Route Table and Association workflow

	var rtOutput map[string]string
	err = workflow.ExecuteChildWorkflow(ctx, resources.RouteTableWorkflow, prov, vpcOutput[utils.VPCID], igwOutput[utils.IGWID], natoutput[utils.NATGATEWAYID], subnetOutput[utils.PRIVATE_SUBNET_ID], subnetOutput[utils.PUBLIC_SUBNET_ID]).Get(ctx, &rtOutput)
	if err != nil {
		workflow.GetLogger(ctx).Info(utils.RtError, "igw_id", igwOutput[utils.IGWID])
	}

	// Start the Security Group Workflow

	var sgOutPut map[string]string
	err = workflow.ExecuteChildWorkflow(ctx, resources.SGWorkflow, prov, vpcOutput[utils.VPCID], vpcOutput["vpc_cidr_block"]).Get(ctx, &sgOutPut)
	if err != nil {
		log.Fatalln("Failed to execute the ", utils.SgWorkflow)
		return nil, err
	}
	workflow.GetLogger(ctx).Info("Security Group created", "sg_id", sgOutPut[utils.SGID])

	// Start the EKS Workflow

	var eksOutput map[string]string
	err = workflow.ExecuteChildWorkflow(ctx, resources.EKSWorkflow, prov, vpcOutput[utils.VPCID], subnetOutput[utils.PRIVATE_SUBNET_ID], subnetOutput[utils.PUBLIC_SUBNET_ID]).Get(ctx, &eksOutput)
	if err != nil {
		log.Fatalln("Failed to execute the EKS workflow")
		return nil, err
	}
	workflow.GetLogger(ctx).Info("EKS Sucessfully Created", "eks_id", eksOutput[utils.EKS_ID])

	// Aggregate results

	results := map[string]interface{}{
		"EKSWorkflow": eksOutput,
	}

	return results, nil

}
