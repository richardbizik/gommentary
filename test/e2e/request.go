package e2e

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

type AddHeaderTransport struct {
	T http.RoundTripper
}

func (adt *AddHeaderTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJqb2huZG9lIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.Ks_BdfH4CWilyzLNk8S2gDARFhuxIauLa8PwhdEQhEo")
	req.Header.Add("Content-Type", "application/json")
	return adt.T.RoundTrip(req)
}

func NewAddHeaderTransport(T http.RoundTripper) *AddHeaderTransport {
	if T == nil {
		T = &http.Transport{}
	}
	return &AddHeaderTransport{T}
}

func Request(t testing.TB, method string, path string, requestBody interface{}) (string, int) {
	c := http.Client{
		Transport: NewAddHeaderTransport(nil),
	}
	var reader io.Reader
	if requestBody != nil {
		json, err := json.Marshal(requestBody)
		if err != nil {
			t.Fatalf("Could not serialize %T to json", requestBody)
		}
		reader = bytes.NewBuffer(json)
	}
	req, err := http.NewRequest(method, path, reader)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	res, err := c.Do(req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			t.Fatalf("Could not read response body %v", err)
		}
	}(res.Body)
	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("Could not read response body %v", err)
	}
	return string(bytes), res.StatusCode
}
