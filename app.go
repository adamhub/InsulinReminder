package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/workflow"

	openapi "github.com/twilio/twilio-go/rest/api/v2010"

	"github.com/twilio/twilio-go"
)

type ErrorResponse struct {
	Message string
}

type Message struct {
	Type   string `json:"insulinType"`
	Amount int    `json:"insulinAmount"`
	From   string `json:"phoneNumber"`
}

type Patient struct {
	Name   string
	Number string
}

func TwilioClient() *twilio.RestClient {

	accountSid := os.Getenv("TWILIO_ACCOUNT_SID")
	authToken := os.Getenv("TWILIO_AUTH_TOKEN")
	return twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})
}

// Main workflow
func InsulinWorkflow(ctx workflow.Context, patient Patient, watchers []string) error {
	var result string
	insulinTaken := false
	logger := workflow.GetLogger(ctx)
	logger.Info("Starting InsulinWorkflow")
	aoptions := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, aoptions)
	err := workflow.ExecuteActivity(ctx, InsulinSMSSendActivity, patient.Number).Get(ctx, &result)
	if err != nil {
		return err
	}
	logger.Info("InsulinSMSSendActivity completed.", "result", result)

	// wait for signal or 10 mins to pass
	// signal handler for SMS response
	var signal Message
	signalChan := workflow.GetSignalChannel(ctx, "CHANNEL_YOLO")
	selector := workflow.NewSelector(ctx)
	selector.AddReceive(signalChan, func(channel workflow.ReceiveChannel, more bool) {
		channel.Receive(ctx, &signal)
		// check if insulin taken amount is acceptable
		if signal.Amount > 0 && signal.Amount < 20 {
			insulinTaken = true
			logger.Info("Insulin has been taken. Sending notification")
			// trigger parent notification of insulin taken
			mes := fmt.Sprintf("%s took %d units of insulin just now.", patient.Name, signal.Amount)
			if insulinTaken {
				SendSMS(TwilioClient(), mes, watchers)
			}

		}

	})

	selector.AddFuture(workflow.NewTimer(ctx, time.Second*10), func(f workflow.Future) {
		// trigger next step after time shown above
		if !insulinTaken {
			// since sms hasn't came in with insulin amount, trigger parent notification
			logger.Info("Insulin not taken. Sending notification")
			// trigger parent notification of insulin taken
			mes := fmt.Sprintf("Alert: %s hasn't taken Insulin yet", patient.Name)
			SendSMS(TwilioClient(), mes, watchers)
		}

	})

	// wait for the timer or sms
	selector.Select(ctx)

	return nil
}

// Send SMS and wait for SMS message to be returned
func InsulinSMSSendActivity(ctx context.Context, to string) (string, error) {

	// Sens SMS alert to take insulin
	fmt.Println("patient_number:", to)
	mes := "Time to take Insulin"
	SendSMS(TwilioClient(), mes, []string{to})

	return "twillio send success", nil
}

// Trigger a signal from incoming SMS
func SMSPOSTHandler(ctx context.Context, temporal client.Client, runID, wfID string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/x-www-form-urlencoded")

		err := req.ParseForm()
		if err != nil {
			log.Fatal(err)
		}

		amount, _ := strconv.Atoi(req.PostForm.Get("Body"))
		if err != nil {
			log.Fatal(err)
		}

		signal := Message{
			Amount: amount,
		}
		err = temporal.SignalWorkflow(context.Background(), wfID, "", "CHANNEL_YOLO", signal)
		if err != nil {
			log.Fatalln("Error sending the Signal foo", err)
			return
		}

	}
}

// Send SMS via twilio
func SendSMS(tclient *twilio.RestClient, message string, to []string) error {
	//from := os.Getenv("TWILIO_FROM_PHONE_NUMBER")
	testing := false
	rrom := "+" + os.Getenv("TWILIO_FROM_PHONE_NUMBER")
	params := &openapi.CreateMessageParams{}
	params.SetFrom(rrom)
	params.SetBody(message)

	// send sms
	for i := range to {
		log.Println("SMS sending...")
		params.SetTo(to[i])

		if !testing {
			resp, err := tclient.Api.CreateMessage(params)
			if err != nil {
				fmt.Println(err.Error())
				return err
			} else {
				response, _ := json.Marshal(*resp)
				fmt.Println("Response: " + string(response))
			}
		} else {
			log.Println("Mock SMS to:", *params.To)
			log.Println("Mock SMS message:", *params.Body)
		}
	}

	return nil
}
