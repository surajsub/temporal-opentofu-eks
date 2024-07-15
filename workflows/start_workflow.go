package workflows

import (
	"context"
	"github.com/surajsub/temporal-opentofu-eks/utils"
	"go.temporal.io/sdk/client"
	"log"
	"time"
)

func StartWorkflow(vpc string) {
	// Create Temporal client
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create Temporal client", err)
	}
	defer c.Close()

	// Define workflow options
	workflowOptions := client.StartWorkflowOptions{
		ID:                       "parent-workflow-opentofu", // Name of the workflow that will be visible in the Temporal UI
		TaskQueue:                utils.WORKFLOW_TASK_QUEUE,  // Queue Name - This can be made dynamic
		WorkflowExecutionTimeout: 60 * time.Minute,
	}

	var provengine = "opentofu"

	if vpc == "" {
		log.Panicln("failed to get the vpc cdir block")
	}
	log.Println("THE INPUT CDIR BLOCK IS ", vpc)

	// Start the workflow
	we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, ParentWorkflow, vpc, provengine)
	if err != nil {
		log.Panicln("Unable to execute workflow", err)
	}

	log.Printf("Started workflow with ID: %s and Run ID: %s", we.GetID(), we.GetRunID())

	var result interface{}
	err = we.Get(context.Background(), &result)
	if err != nil {
		log.Panic("Unable to get workflow result", err)
	}
	log.Printf("Workflow result: %v", result)
}
