# Insulin Reminder

[![Go Report Card](https://goreportcard.com/badge/github.com/joho/godotenv)](https://goreportcard.com/report/github.com/joho/godotenv)

A SMS based, Temploral.io based workflow which reminders a patient to take insulin and either escalates an alert to separate watchers should they not take it, or reports what they have taken. The missed alert happens after 15 minutes.

Screen capture animation here:

# Setup

You will need Go and Docker Compose installed. If you don't want to spin up Temporal from Docker, a lighter option is [Temporalite](https://github.com/temporalio/temporalite). Note that you have to create ~/.config/temporalite/db directories. Skip the Docker Compose directions below if you went this route.

In the root of the repo directory, create a .env file. Here's an example:

```conf
TWILIO_ACCOUNT_SID=1289h23r987h237y327623498
TWILIO_AUTH_TOKEN=asjhasouihwefoief980308238
TWILIO_FROM=15559998888
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


The code runs a server on port 4000 for recieving SMS messages via POST requests from Twilio. Make sure that is open on your server. If you are running locally you can use ngrok for testing:
> ngrok http 4000
The public test url that you recieve there, or alternatively your production server address, will need to be copied to your Twilio SMS webhook URL setting in your Twilio console. It's composed with /sms at the end. For example, if you setup Temporal on your live server, the POST request is open at this URL:
> https://prod-site.com:4000/sms
That is the URL you need to add to Twilio as a SMS receive webhook.


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
http://localhost:8080 or temporalite at http://localhost:8233


![image](https://user-images.githubusercontent.com/763917/192355793-8b4339c8-cfe8-4cb2-8609-e70f46172027.png)
