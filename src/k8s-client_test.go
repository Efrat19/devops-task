package main

import (
	"testing"
	"time"
)

func TestGetTimeSince(t *testing.T) {
	yesterday := time.Now().Add(-24*time.Hour)
	expected := time.Hour*24
	timeSince := getTimeSince(yesterday).Round(time.Minute)
	if timeSince != expected {
		t.Errorf("getTimeSince was incorrect, got: %v, want: %v.", timeSince, expected)
	}
}