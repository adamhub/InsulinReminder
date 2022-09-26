# Insulin Reminder

[![Go Report Card](https://goreportcard.com/badge/github.com/joho/godotenv)](https://goreportcard.com/report/github.com/joho/godotenv)

A SMS based, Temploral.io based workflow which reminders a patient to take insulin and either escalates an alert to separate watchers should they not take it, or reports what they have taken. The missed alert happens after 15 minutes.

Screen capture animation here:

# Setup

You will need Go and Docker Compose installed. 

In the root of the repo directory, create a .env file. Here's an example:

```
TWILIO_ACCOUNT_SID=1289h23r987h237y327623498
TWILIO_AUTH_TOKEN=asjhasouihwefoief980308238
TIWLIO_FROM_NUMBER=15559998888
WATCHER1=5559998888
WATCHER2=5559998888
PATIENT=5559998888
PATIENT_NAME=Nathan
```

The timer is easily configured in app.go under `CronSchedule`.

```go
workflowOptions := client.StartWorkflowOptions{
    ID:        "ifworker-" + str2,
    TaskQueue: "insulinFlowWorker",
    // for immediate start, remove cron
    CronSchedule: "0 10,22 * * *",
}
```

# Run

Clone the repo. Then get all the dependencies:
> go mod init

Start Temporal
> docker-compose up

Start the worker
> go run woker/main.go

Run the workflow
> go run starter/main.go

You can see the running process in the Temploral web ui:
[http://localhost:8080]