package main

import (
	"testing"
)

func TestGetCommandName(t *testing.T) {
	userCammnd := "logs example-service 12"
	result,_ := getCommandName(userCammnd)
	expected := "logs"
	if result != expected {
		t.Errorf("getCommandName was incorrect, got: %v, want: %v.", result, expected)
	}
}