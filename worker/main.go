package main

import (
	"log"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"

	app "insulinReminder"
)

// Temporal Worker
func main() {
	var err error
	temporal, err := client.NewClient(client.Options{})
	if err != nil {
		log.Fatalln("unable to create Temporal client", err)
	}

	wkr := worker.New(temporal, "insulinFlowWorker", worker.Options{})

	wkr.RegisterWorkflowWithOptions(app.InsulinWorkflow, workflow.RegisterOptions{Name: "InsulinWorkflow"})
	wkr.RegisterActivity(app.InsulinSMSSendActivity)

	err = wkr.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}

}
