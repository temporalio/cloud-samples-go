package main

import (
	"context"
	"github.com/pborman/uuid"
	"log"

	"github.com/temporalio/cloud-operations-workflows/client/workflowclient"
	"github.com/temporalio/cloud-operations-workflows/workflows"
	"go.temporal.io/sdk/client"
)

func main() {
	c := workflowclient.MustGetClient()
	defer c.Close()

	workflowOptions := client.StartWorkflowOptions{
		ID:        uuid.NewUUID().String(),
		TaskQueue: "demo",
	}

	input := &workflows.PeriodicReconcileUsersInput{
		FilePath: "/Users/liang/demo/users.json",
	}
	we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, workflows.PeriodicReconcileUsersWorkflowType, input)
	if err != nil {
		log.Fatalln("Unable to start workflow", err)
	}

	log.Println("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())

	// Synchronously wait for the workflow completion.
	var result string
	err = we.Get(context.Background(), &result)
	if err != nil {
		log.Fatalln("Unable get workflow result", err)
	}
	log.Println("Workflow result:", result)
}
