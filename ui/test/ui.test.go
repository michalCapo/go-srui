package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"dasolutions.sk/test/ui"
)

func TestContext(body []ui.BodyItem) *ui.Context {
	// Create a sample request body
	requestBody, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}

	// Create a new HTTP request with the sample body
	req := httptest.NewRequest(http.MethodPost, "http://example.com", bytes.NewReader(requestBody))

	ctx := &ui.Context{
		Request: req,
	}

	return ctx
}

func Assert(t *testing.T, value, expected interface{}) {
	if value != expected {
		t.Errorf("Expected %v, got %v", expected, value)
	}
}

func Equal(t *testing.T, result bool) {
	if !result {
		t.Errorf("Expected true, got false")
	}
}
