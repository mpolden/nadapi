package api

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mpolden/nadapi/nad"
)

func httpGet(url string) (string, int, error) {
	res, err := http.Get(url)
	if err != nil {
		return "", 0, err
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", 0, err
	}
	return string(data), res.StatusCode, nil
}

func httpPatch(url string, body string) (string, int, error) {
	r, err := http.NewRequest(http.MethodPatch, url, strings.NewReader(body))
	if err != nil {
		return "", 0, err
	}
	res, err := http.DefaultClient.Do(r)
	if err != nil {
		return "", 0, err
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", 0, err
	}
	return string(data), res.StatusCode, nil
}

func testServer() *httptest.Server {
	client := nad.NewTestClient()
	api := New(client)
	return httptest.NewServer(api.Handler())
}

func TestGetRequests(t *testing.T) {
	server := testServer()
	defer server.Close()

	var tests = []struct {
		url      string
		response string
		status   int
	}{
		{"/not-found", "404 page not found\n", 404},
		{"/api/not-found", `{"status":404,"message":"Not found"}`, 404},
		{"/api/v1/state/", `{"status":400,"message":"Missing path parameter"}`, 400},
		{"/api/v1/state/foo", `{"status":400,"message":"Invalid variable: \"foo\""}`, 400},
		{"/api/v1/state/power", `{"power":false}`, 200},
		{"/api/v1/state/mute", `{"mute":false}`, 200},
		{"/api/v1/state/speakera", `{"speakerA":true}`, 200},
		{"/api/v1/state/speakerb", `{"speakerB":false}`, 200},
		{"/api/v1/state/model", `{"model":"C356"}`, 200},
		{"/api/v1/state/source", `{"source":"CD"}`, 200},
	}

	for _, tt := range tests {
		data, status, err := httpGet(server.URL + tt.url)
		if err != nil {
			t.Fatal(err)
		}
		if got := status; status != tt.status {
			t.Errorf("want %d for %q, got %d", tt.status, tt.url, got)
		}
		if got := string(data); got != tt.response {
			t.Errorf("want %q, got %q", tt.response, got)
		}
	}
}

func TestPatchRequests(t *testing.T) {
	server := testServer()
	defer server.Close()

	var tests = []struct {
		url      string
		body     string
		response string
		status   int
	}{
		{"/api/v1/state/volume", `{"value":"off"}`, `{"status":400,"message":"Invalid command: volume=off"}`, 400},
		{"/api/v1/state/power", `{"value":"+"}`, `{"status":400,"message":"Invalid command: power+"}`, 400},
		// Model? is considered invalid in this case as it does not modify state
		{"/api/v1/state/model", `{"value":"?"}`, `{"status":400,"message":"Invalid command: model?"}`, 400},
		{"/api/v1/state/power", `{"value":"on"}`, `{"power":true}`, 200},
		{"/api/v1/state/power", `{"value":"off"}`, `{"power":false}`, 200},
		{"/api/v1/state/volume", `{"value":"+"}`, `{"volume":"+"}`, 200},
		{"/api/v1/state/volume", `{"value":"-"}`, `{"volume":"-"}`, 200},
		{"/api/v1/state/source", `{"value":"DISC/MDC"}`, `{"source":"DISC/MDC"}`, 200},
	}
	for _, tt := range tests {
		data, status, err := httpPatch(server.URL+tt.url, tt.body)
		if err != nil {
			t.Fatal(err)
		}
		if got := status; status != tt.status {
			t.Errorf("want %d for %q, got %d", tt.status, tt.url, got)
		}
		if got := string(data); got != tt.response {
			t.Errorf("want %q, got %q", tt.response, got)
		}
	}
}
