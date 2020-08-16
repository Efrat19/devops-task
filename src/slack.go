package main

import (
	"bytes"
	"os"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"github.com/nlopes/slack"
	"strconv"
	"strings"
	"github.com/apex/log"
)

const (
	LOGS_COMMAND = "logs"
	PODS_COMMAND = "pods"

	// logs command args
	SERVICE_ARG_INDEX = 1
	TAIL_ARG_INDEX = 2
	DEFAULT_TAIL_VALUE = 10
)
func slashCommandHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("Receiving /k-bot request")
	signingSecret := os.Getenv("SLACK_SIGNING_SECRET")
	verifier, err := slack.NewSecretsVerifier(r.Header, signingSecret)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	r.Body = ioutil.NopCloser(io.TeeReader(r.Body, &verifier))
	s, err := slack.SlashCommandParse(r)
	if err != nil {
		log.Errorf("On parsing: %v",err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = verifier.Ensure(); err != nil {
		log.Error("Invalid verifier")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	countRequest(s.Text,s.UserID)
	command,err := getCommandFirstArg(s.Text)
	if err != nil {
		response := "Available commands are:\nk-bot pods\nk-bot logs [service] [tail]"
		w.Write([]byte(response))
	} else {
		switch command {
		case PODS_COMMAND:
			response, err := getKbotPods()
			if err != nil {
				log.Error("Unable to getKbotPods")
				panic(err.Error())
			}
			w.Write([]byte(response))
		case LOGS_COMMAND:
			response,err := getKbotLogs(s.Text)
			if err != nil {
				log.Error("Unable to getKbotLogs")
				panic(err.Error())
			}
			w.Write([]byte(response))
		default:
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

}

func getCommandFirstArg(fullCommand string) (string,error) {
	log.Debugf("getCommandName called with userCommand %s",fullCommand)
	splittedCommand := strings.Split(fullCommand, " ")
	if len(splittedCommand) < 1 {
		return "",fmt.Errorf("No Command arg specified\n")
	}
	return splittedCommand[0],nil
}

func getKbotLogs(command string) (string,error) {
	splittedCommand := strings.Split(command, " ")
	if len(splittedCommand) < SERVICE_ARG_INDEX+1 {
		return "",fmt.Errorf("No service specified for logs\n")
	}
	tail := DEFAULT_TAIL_VALUE
	if len(splittedCommand) < TAIL_ARG_INDEX+1 {
		log.Warn("No tail specified, defaulting to 10\n")
	} else {
		tail, err := strconv.Atoi(splittedCommand[TAIL_ARG_INDEX])
		log.Infof("tail from command: %d",tail)
		if err != nil {
			log.Warn("No valid tail specified, defaulting to 10\n")
			tail = DEFAULT_TAIL_VALUE
		} else {
			log.Infof("tail is: %d",tail)
			return getServiceLog(int64(tail),splittedCommand[SERVICE_ARG_INDEX])
		}
	}
	return getServiceLog(int64(tail),splittedCommand[1])
}



func getKbotPods() (string,error) {
	pods, err := getPodInfoList()
	if err != nil {
		log.Error("Unable to getPodInfoList")
		return "",err
	}
	var buffer bytes.Buffer
	buffer.WriteString("```\n")
	for _,pod := range *pods {
		buffer.WriteString(fmt.Sprintf("pod %s uptime %s version %s\n",pod.Name,pod.Uptime,pod.Version))
	}
	buffer.WriteString("```\n")
	response := buffer.String()
	return response,nil
}
