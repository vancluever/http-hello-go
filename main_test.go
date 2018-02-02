package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestClientAddress_noXFF(t *testing.T) {
	expected := "12.34.56.78"
	r := &http.Request{
		RemoteAddr: fmt.Sprintf("%s:8080", expected),
		Header:     http.Header{},
	}
	client := clientAddress(r)

	if expected != client {
		t.Fatalf("Expected %s, got %s", expected, client)
	}
}

func TestClientAddress_withXFF(t *testing.T) {
	expected := "90.12.34.56"
	proxy := "12.34.56.78"
	r := &http.Request{
		RemoteAddr: fmt.Sprintf("%s:8080", proxy),
		Header:     http.Header{},
	}
	r.Header.Add("X-Forwarded-For", fmt.Sprintf("%s, %s", expected, proxy))
	client := clientAddress(r)

	if expected != client {
		t.Fatalf("Expected %s, got %s", expected, client)
	}
}

func TestStartServer_version(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(handleVersion))
	defer ts.Close()
	resp, err := http.Get(ts.URL + "/version")
	if err != nil {
		t.Fatalf("Bad: %v", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Bad: %v", err)
	}
	defer resp.Body.Close()

	m := map[string]interface{}{}
	if err := json.Unmarshal(body, &m); err != nil {
		t.Fatalf("Bad: %v", err)
	}

	if m["version"] != version {
		t.Fatalf("Expected version to be %s, got %s", version, m["version"])
	}
}

func TestStartServer_default(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(handleDefault))
	defer ts.Close()
	resp, err := http.Get(ts.URL + "/")
	if err != nil {
		t.Fatalf("Bad: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Bad: %v", err)
	}

	if strings.Contains(string(body), fmt.Sprintf("version %s", version)) == false {
		t.Fatalf("Expected response body to contain version %s, but not found, body %s", version, body)
	}
}

func TestStartServer_default_404(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(handleDefault))
	defer ts.Close()
	resp, err := http.Get(ts.URL + "/bad")
	if err != nil {
		t.Fatalf("Bad: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Bad: %v", err)
	}

	if resp.StatusCode != 404 {
		t.Fatalf("Expected code to be 404, got %d", resp.StatusCode)
	}

	if strings.Contains(string(body), "Path /bad not found") == false {
		t.Fatalf("Expected Path /bad not found, but got %s", body)
	}
}
