package main

import (
	"fmt"

	"github.com/caarlos0/env/v9"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	"github.com/temporalio/cloud-operations-workflows/client/temporalcloud"
	"github.com/temporalio/cloud-operations-workflows/internal/validator"
	"github.com/temporalio/cloud-operations-workflows/workflows"
)

type (
	config struct {
		TemporalCloudAPIAddress string `env:"TEMPORAL_CLOUD_API_ADDRESS" envDefault:"saas-api.tmprl.cloud:443" validate:"required"`
		AllowInsecure           bool   `env:"ALLOW_INSECURE"`
		TemporalCloudAPIKey     string `env:"TEMPORAL_CLOUD_API_KEY" validate:"required"`
	}
)

func main() {
	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		panic(fmt.Errorf("failed to parse env: %+v", err))
	}
	if err := validator.ValidateStruct(cfg); err != nil {
		panic(fmt.Errorf("invalid config: %+v", err))
	}
	c, err := newLocalTemporalClient()
	if err != nil {
		panic(fmt.Errorf("failed to create temporal client: %+v", err))
	}
	w, err := newWorker(c)
	if err != nil {
		panic(fmt.Errorf("failed to create temporal worker: %+v", err))
	}

	conn, err := temporalcloud.NewConnectionWithAPIKey(
		cfg.TemporalCloudAPIAddress,
		cfg.AllowInsecure,
		cfg.TemporalCloudAPIKey,
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

func newLocalTemporalClient() (client.Client, error) {
	return client.Dial(client.Options{})
}

func newWorker(client client.Client) (worker.Worker, error) {
	wo := worker.Options{
		MaxConcurrentActivityTaskPollers: 10,
		MaxConcurrentWorkflowTaskPollers: 10,
	}
	return worker.New(client, taskQueue, wo), nil
}
