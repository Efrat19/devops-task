package main

import (
	"github.com/apex/log/handlers/text"
	"os"
	"github.com/apex/log"
)

func main() {
	logLevel := getEnv("LOG_LEVEL","debug")
	log.SetLevelFromString(logLevel)
	log.SetHandler(text.New(os.Stderr))

	port := getEnv("PORT","1012")
	serve(port)
}


func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}