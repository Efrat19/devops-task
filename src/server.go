package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"os"
	"github.com/apex/log"
)

func main() {
	port := os.Getenv("PORT")
	http.HandleFunc("/k-bot", slashCommandHandler)
	http.HandleFunc("/healthz", healthzHandler)
	http.Handle("/metrics", Adapt(promhttp.Handler(), logAccess()))
	http.ListenAndServe(fmt.Sprintf(":%s",port), nil)
	log.Infof("Server listening on port %s",port)
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("Receiving /healthz request")
	w.WriteHeader(http.StatusOK)
	return
}

