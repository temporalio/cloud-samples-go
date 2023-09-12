package main

import (
	"fmt"
	"github.com/temporalio/cloud-operations-workflows/client/temporalcloud"
	"github.com/temporalio/cloud-operations-workflows/client/workflowclient"
	"github.com/temporalio/cloud-operations-workflows/workflows"
	"go.temporal.io/sdk/worker"
)

const (
	controlPlaneHostPort = "saas-api.tmprl.cloud:443"
	//apiKeyValue          = "tmprl_myCSpZXG5EyGYOhCCriejuCf814zFRNz_Np9b0JC6knytsMz3261G1VT85x4vIJ5fBPdhuKK4SuTbK0XruUBgobydixE7icGB"
	// demo3 (in the ps13i account)
	//apiKeyValue = "tmprl_46L87rLmDTmve2mqycurdcJE9DEDNF7j_VBZGpFXcL3WD4UWVtQ6wxz9nRmRuKNMs6BSFBtlkNF1BgSxvZAGL3Wi5lXQGWwj2"
	// demo4 (in the temporal-dev account)
	apiKeyValue = "tmprl_C6h4NKZXqgImQTSnJX9NVMIzAiawtBI5_j1pjzO1DG31JOdEZdl99k3GLGfkle94p0U3qJpmOArOzP41aI74Cupytzvh57aoF"
)

func main() {
	c := workflowclient.MustGetClient()
	defer c.Close()

	w := worker.New(c, "demo", worker.Options{})

	conn, err := temporalcloud.NewConnectionWithAPIKey(
		controlPlaneHostPort,
		false,
		apiKeyValue,
	)
	if err != nil {
		panic(fmt.Errorf("failed to create cloud api connection: %+v", err))
	}
	workflows.Register(w, workflows.NewWorkflows(), workflows.NewActivities(conn))
	err = w.Run(worker.InterruptCh())
	if err != nil {
		panic(fmt.Errorf("failed to run worker: %+v", err))
	}
}
