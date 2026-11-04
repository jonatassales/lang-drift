package httpapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/jonatassales/tao/apps/microservices/console-stubs/src/application"
	"github.com/jonatassales/tao/apps/microservices/console-stubs/src/infrastructure/probe"
)

func TestLocalBurstMeasures50RequestsAndRejectsTargets(t *testing.T) {
	handler, store := fixture(t)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { handler.ServeHTTP(w, r) }))
	defer server.Close()
	address, err := url.Parse(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	handler.Service.LoadTests = &application.LoadTester{Probe: probe.New(address.Port())}
	endpoint := "/console/runtimes/seed-runtime-1/load-test"
	if code, _ := call(t, handler, "POST", endpoint, map[string]any{"url": "https://example.com", "requests": 50000}); code != 400 {
		t.Fatal("custom target/count accepted")
	}
	response, err := http.Post(server.URL+endpoint, "application/json", bytes.NewBufferString("{}"))
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	var payload struct {
		LoadTest struct {
			Target                                   string
			Requests, Concurrency, Succeeded, Failed int
		}
	}
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		t.Fatal(err)
	}
	if response.StatusCode != 200 || payload.LoadTest.Target != "console-stubs-loopback" || payload.LoadTest.Requests != 50 || payload.LoadTest.Concurrency != 50 || payload.LoadTest.Succeeded != 50 || payload.LoadTest.Failed != 0 {
		t.Fatalf("bad burst result: %+v", payload)
	}
	if store.Snapshot().Runtimes[0].LoadTest == nil || store.Snapshot().Runtimes[1].LoadTest != nil {
		t.Fatal("result ownership lost")
	}
	if code, _ := call(t, handler, "POST", endpoint, map[string]any{}); code != 429 {
		t.Fatal("cooldown bypassed")
	}
	if code, _ := call(t, handler, "POST", "/seeds/stopped", map[string]any{}); code != 404 {
		t.Fatal("global scenario endpoint survived")
	}
}
