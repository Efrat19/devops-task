package main

import (
	"testing"
)

func TestGetCommandFirstArg(t *testing.T) {
	userCammnd := "logs example-service 12"
	result,_ := getCommandFirstArg(userCammnd)
	expected := "logs"
	if result != expected {
		t.Errorf("getCommandFirstArg was incorrect, got: %v, want: %v.", result, expected)
	}
}