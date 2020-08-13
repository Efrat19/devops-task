package main

import (
	"bytes"
	"encoding/json"
	"os"
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


func logRequest(user string ,userID string, devops string) {
	webhook := os.Getenv("SLACK_WEBHOOK_URL")
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
	signingSecret := os.Getenv("SLACK_SIGNING_SECRET")
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

	countRequest(s.Command,s.UserID)
	switch s.Command {
	case "/k-bot pods":
		pods, err := getPodInfoList()
		if err != nil {
			panic(err.Error())
		}
		var buffer bytes.Buffer
		for _,pod := range *pods {
			buffer.WriteString(fmt.Sprintf("pod %s uptime %s version %s\n",pod.Name,pod.Uptime,pod.Version))
		}
		response := buffer.String()
		w.Write([]byte(response))
	case "/k-bot logs service tail":
		log := getServiceLog(20,"service")
		w.Write([]byte(log))
	default:
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
