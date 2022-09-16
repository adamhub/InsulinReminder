package main

import (
	"context"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"go.temporal.io/sdk/client"

	app "familyFlows"
)

func main() {

	// make temporal client
	temporal, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer temporal.Close()

	// Get psuedo-random number and trim it for the WorkflowID
	rand.Seed(time.Now().UnixNano())
	str1 := strconv.Itoa(rand.Int())
	str2 := string([]byte(str1)[:8])

	workflowOptions := client.StartWorkflowOptions{
		ID:           "ifworker-" + str2,
		TaskQueue:    "insulinFlowWorker",
		CronSchedule: "0 10,22 * * *",
	}

	we, err := temporal.ExecuteWorkflow(context.Background(), workflowOptions, app.InsulinWorkflow)
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}

	log.Println("Started workflow", "WorkflowID:", we.GetID(), "runid", we.GetRunID())

	// Recieve POST request

	// Setup POST request handler
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/sms", app.SMSPOSTHandler(context.Background(), temporal, we.GetRunID(), we.GetID()))

	var port string
	var found bool
	if port, found = os.LookupEnv("PORT"); !found {
		port = "4000"
	}

	go func() {
		log.Printf("Starting web server on port %s", port)
		panic(http.ListenAndServe(":"+port, serveMux))
	}()

	// Synchronously wait for the workflow completion.
	var result string
	err = we.Get(context.Background(), &result)
	if err != nil {
		log.Fatalln("Unable get workflow result", err)
	}
	log.Println("Workflow result:", result)

}
