package main

import (
	"encoding/json"
	"strconv"
	"time"

	//"encoding/json"
	"fmt"
	"io"
	"io/ioutil"

	//"log"
	//"math"
	"net/http"
	//"github.com/joho/godotenv"
	"github.com/nlopes/slack"
)

func main() {

	http.HandleFunc("/devops-on-duty", slashCommandHandler)
	http.HandleFunc("/healthz", healthzHandler)
	fmt.Println("[INFO] Server listening")
	http.ListenAndServe(":8080", nil)
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[INFO] Receiving /healthz request")
	w.WriteHeader(http.StatusOK)
	return
}

func logRequest(user string ,userID string, devops string) {
	webhook := getEnv("SLACK_WEBHOOK_URL","xxx")
	fmt.Println("[INFO] Logging /devops-on-duty request")
	fmt.Printf("%s (%s) just issued the /devops-on-duty command\n",user,userID)
	attachment := slack.Attachment{
		Color:         "warning",
		Fallback:      fmt.Sprintf("Heads up for %s: %s (%s) just issued the /devops-on-duty command",devops,user,userID),
		//AuthorName:    "devops bot",
		//AuthorSubname: "github.com",
		//AuthorLink:    "https://github.com/nlopes/slack",
		//AuthorIcon:    "https://avatars2.githubusercontent.com/u/652790",
		Text:          fmt.Sprintf("Heads up for %s: %s just issued the /devops-on-duty command",devops,user),
		//Footer:        "slack api",
		//FooterIcon:    "https://platform.slack-edge.com/img/default_application_icon.png",
		Ts:            json.Number(strconv.FormatInt(time.Now().Unix(), 10)),
	}
	msg := slack.WebhookMessage{
		Attachments: []slack.Attachment{attachment},
	}

	err := slack.PostWebhook(webhook, &msg)
	if err != nil {
		fmt.Println(err)
	}
}

func slashCommandHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[INFO] Receiving /devops-on-duty request")
	signingSecret := getEnv("SLACK_SIGNING_SECRET","")
	verifier, err := slack.NewSecretsVerifier(r.Header, signingSecret)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	r.Body = ioutil.NopCloser(io.TeeReader(r.Body, &verifier))
	s, err := slack.SlashCommandParse(r)
	if err != nil {
		fmt.Printf("[ERROR] on parsing: %v",err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = verifier.Ensure(); err != nil {
		fmt.Printf("[ERROR] invalid verfier")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	switch s.Command {
	case "/devops-on-duty":
		err,onDuty := whoIsOnDutyNow()
		var response string
		if err != nil {
			fmt.Printf("[ERROR] Error finding DevOps on duty today %v",err)
			response = "Error finding DevOps on duty today"
			logRequest(s.UserName,s.UserID,"error")
		} else {
			response = fmt.Sprintf("%s is on duty today",onDuty)
			logRequest(s.UserName,s.UserID,onDuty)
		}
		w.Write([]byte(response))
	default:
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
