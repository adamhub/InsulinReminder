# Insulin Reminder

A SMS based, Temploral.io based workflow which reminders a patient to take insulin and either escalates an alert to separate watchers should they not take it, or reports what they have taken. The missed alert happens after 15 minutes.

Screen capture animation here:


# Running

After you have installed Go, clone the repo. Then get all the dependencies:
> go mod init

The timer is easily configured in app.go under `CronSchedule`.

```
workflowOptions := client.StartWorkflowOptions{
    ID:        "ifworker-" + str2,
    TaskQueue: "insulinFlowWorker",
    // for immediate start, remove cron
    CronSchedule: "0 10,22 * * *",
}
```


Start Temporal
> docker-compose up

Start the worker
> go run woker/main.go

Run the workflow
> go run starter/main.go

