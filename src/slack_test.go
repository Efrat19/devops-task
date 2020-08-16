package main

import (
	"testing"
)

func TestGetCommandFirstArgReturnsCorrectArg(t *testing.T) {
	fullCommand := "logs example-service 12"
	result,_ := getCommandFirstArg(fullCommand)
	expected := "logs"
	if result != expected {
		t.Errorf("TestGetCommandFirstArgReturnsCorrectArg was incorrect, got: %v, want: %v.", result, expected)
	}
}

func TestGetCommandFirstArgFailsWithError(t *testing.T) {
	fullCommand := ""
	_,err := getCommandFirstArg(fullCommand)
	if err == nil {
		t.Errorf("TestGetCommandFirstArgFailsWithError was incorrect, it was supposed to return error for command: %v", fullCommand)
	}
}