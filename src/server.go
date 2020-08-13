package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	http.HandleFunc("/k-bot", slashCommandHandler)
	http.HandleFunc("/healthz", healthzHandler)
	http.Handle("/metrics", Adapt(promhttp.Handler(), logAccess()))
	http.ListenAndServe(fmt.Sprintf(":%s",port), nil)
	fmt.Println("[INFO] Server listening")
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[INFO] Receiving /healthz request")
	w.WriteHeader(http.StatusOK)
	return
}

