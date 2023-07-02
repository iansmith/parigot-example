package main

import (
	"testing"

	"github.com/iansmith/parigot-example/hello-world/g/greeting/v1"
)

func TestBounds(t *testing.T) {
	svc := &myService{}

	req := &greeting.FetchGreetingRequest{
		Tongue: greeting.Tongue_English,
	}
	resp, err := svc.fetchGreeting(req)
	if err != greeting.GreetErr_NoError || resp == "" {
		t.Errorf("failed to get a response for english: %s, %s",
			greeting.Tongue_name[int32(greeting.Tongue_English)],
			greeting.GreetErr_name[int32(err)])
	}

	// out of bounds request
	req.Tongue = greeting.Tongue_Unspecified
	_, err = svc.fetchGreeting(req)
	if err == greeting.GreetErr_NoError {
		t.Errorf("expected error when doing fetchGreeting() with unspecified language")
	}
}