# Insulin Reminder

[![Go Report Card](https://goreportcard.com/badge/github.com/joho/godotenv)](https://goreportcard.com/report/github.com/joho/godotenv)

A SMS based, Temploral.io based workflow which reminders a patient to take insulin and either escalates an alert to separate watchers should they not take it, or reports what they have taken. The missed alert happens after 15 minutes.

Screen capture animation here:

# Setup

You will need Go and Docker Compose installed. 

In the root of the repo directory, create a .env file. Here's an example:

```conf
TWILIO_ACCOUNT_SID=1289h23r987h237y327623498
TWILIO_AUTH_TOKEN=asjhasouihwefoief980308238
TIWLIO_FROM_NUMBER=15559998888
WATCHER1=5559998888
WATCHER2=5559998888
PATIENT=5559998888
PATIENT_NAME=Nathan
CRONA=3
CRONB=15
```

The timer is made from CRONA and CRONB which currently constructs a twice daily timer implimented in `CronSchedule` in starter/main.go.

```go
CronSchedule: "0 " + cronA + "," + cronB + " * * *"
```

In a folder next to the main repo, create a docker-compose folder and clone the docker-compose files:
> git clone https://github.com/temporalio/docker-compose

# Run

Clone the repo. Then get all the dependencies:
```bash
cd InsulinReminder
go mod init
```

Start Temporal
```bash
cd docker-compose
docker-compose up
```

Start the worker
```bash
cd ../InsulinReminder
go run woker/main.go
```

Run the workflow
```bash
go run starter/main.go
```

You can see the running process in the Temploral web ui:
http://localhost:8080


![image](https://user-images.githubusercontent.com/763917/192355793-8b4339c8-cfe8-4cb2-8609-e70f46172027.png)
