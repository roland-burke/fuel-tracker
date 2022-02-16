package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func init() {
	initLogger()
	initDb()
}

func TestSendResponseWithMessageAndStatus(t *testing.T) {
	recorder := httptest.NewRecorder()

	sendResponseWithMessageAndStatus(recorder, 200, "Test")

	if recorder.Code != 200 || recorder.Body.String() != "Test" {
		t.Errorf("Status code and message was: %d, %s, expected: %d, %s", recorder.Code, recorder.Body.String(), 200, "Test")
	}
}

func TestCheckCredentials(t *testing.T) {
	recorder := httptest.NewRecorder()
	recorder.Header().Add("username", "john")
	recorder.Header().Add("password", "john")

	req, err := http.NewRequest("GET", "/unimportant", nil)

	if err != nil {
		t.Errorf("Failed to create http Request")
	}

	req.Header.Set("username", "john")
	req.Header.Set("password", "john")

	result := checkCredentials(req)

	if !result {
		t.Errorf("Credentials check failed")
	}
}
