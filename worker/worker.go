package worker

import (
	"github.com/surajsub/temporal-opentofu-eks/activities"
	"github.com/surajsub/temporal-opentofu-eks/logger"
	"github.com/surajsub/temporal-opentofu-eks/utils"
	tfworkflows "github.com/surajsub/temporal-opentofu-eks/workflows"
	"github.com/surajsub/temporal-opentofu-eks/workflows/resources"
	"go.uber.org/zap"

	"log"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func RunWorker() {

	tflogger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer tflogger.Sync() // flushes buffer, if any

	// Create a custom Temporal logg using Zap
	temporalLogger := logger.NewZapAdapter(tflogger)

	c, err := client.Dial(client.Options{
		Logger: temporalLogger,
	})
	if err != nil {
		log.Panic("Unable to create Temporal client", err)
	}
	defer c.Close()

	w := worker.New(c, utils.WORKFLOW_TASK_QUEUE, worker.Options{})

	// Register workflows
	w.RegisterWorkflow(tfworkflows.ParentWorkflow)
	w.RegisterWorkflow(resources.VPCWorkflow)
	w.RegisterWorkflow(resources.SubnetWorkflow)
	w.RegisterWorkflow(resources.IGWWorkflow)
	w.RegisterWorkflow(resources.NATWorkflow)
	w.RegisterWorkflow(resources.RouteTableWorkflow)
	w.RegisterWorkflow(resources.SGWorkflow)
	w.RegisterWorkflow(resources.EKSWorkflow)

	// Register activities
	w.RegisterActivity(activities.VPCInitActivity)
	w.RegisterActivity(activities.VPCApplyActivity)
	w.RegisterActivity(activities.VPCOutputActivity)

	// Register the Subnet work
	w.RegisterActivity(activities.SubnetInitActivity)
	w.RegisterActivity(activities.SubnetApplyActivity)
	w.RegisterActivity(activities.SubnetOutputActivity)

	// Register the IGW work
	w.RegisterActivity(activities.IGWInitActivity)
	w.RegisterActivity(activities.IGWApplyActivity)
	w.RegisterActivity(activities.IGWOutputActivity)

	// Register the NAT Work

	w.RegisterActivity(activities.NATInitActivity)
	w.RegisterActivity(activities.NATApplyActivity)
	w.RegisterActivity(activities.NATOutputActivity)

	// Register the RT Work

	w.RegisterActivity(activities.RTInitActivity)
	w.RegisterActivity(activities.RTApplyActivity)
	w.RegisterActivity(activities.RTOutputActivity)

	// Register the SG work
	w.RegisterActivity(activities.SGInitActivity)
	w.RegisterActivity(activities.SGApplyActivity)
	w.RegisterActivity(activities.SGOutputActivity)

	// Register the EKS work
	w.RegisterActivity(activities.EKSInitActivity)
	w.RegisterActivity(activities.EKSApplyActivity)
	w.RegisterActivity(activities.EKSOutputActivity)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Panic("Unable to start worker", err)
	}
}
