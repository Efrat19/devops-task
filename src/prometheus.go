package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"os"
	"github.com/apex/log"
)

var requestLabels = []string{
	"command",
	"userID",}

var (
	counterVec = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "k-bot",
			Name:      "requests_total",
			Help:      "monitor requests sent to k-bot",
		},
		requestLabels,
	)
)

func countRequest(command string,userID string){
	log.Debugf("countRequest called with command %s and userID %s",command,userID)
	labels := []string{command,userID}
	counterVec.WithLabelValues(labels...).Add(1)
}

func init(){
	prometheus.MustRegister(counterVec)
	fmt.Fprintln(os.Stderr, "Registering counter vector")
}