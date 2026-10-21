package connected

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	app "github.com/jonatassales/tao/packages/ex/src/application/connected"
)

type roundTrip func(*http.Request) (*http.Response, error)

func (f roundTrip) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }
func TestMCP_RejectsUnsafeEndpointsBeforeSendingCredentials(t *testing.T) {
	for _, raw := range []string{"http://example.com", "https://secret@example.com", "https://example.com?secret=x", "https://example.com#x", "file:///tmp/socket"} {
		t.Run(raw, func(t *testing.T) {
			calls := 0
			client := &MCP{Client: &http.Client{Transport: roundTrip(func(r *http.Request) (*http.Response, error) { calls++; return nil, nil })}}
			_, err := client.Sync(context.Background(), app.Connection{URL: raw, Credential: "secret"}, app.Sync{})
			if err == nil || calls != 0 {
				t.Fatalf("unsafe endpoint sent %d requests, error %v", calls, err)
			}
		})
	}
}
func TestMCP_StructuredToolResultAndCredentialIsolation(t *testing.T) {
	client := &MCP{Client: &http.Client{Transport: roundTrip(func(r *http.Request) (*http.Response, error) {
		if r.URL.String() != "https://console.example/api/mcp" || r.Header.Get("Authorization") != "Bearer private" {
			t.Fatal("wrong scoped request")
		}
		var input map[string]any
		_ = json.NewDecoder(r.Body).Decode(&input)
		var result any = map[string]any{"structuredContent": app.Intent{DesiredState: "stopped", Revision: 1}}
		if input["method"] == "initialize" {
			result = map[string]any{"protocolVersion": "2025-11-25"}
		}
		data, _ := json.Marshal(map[string]any{"jsonrpc": "2.0", "id": input["id"], "result": result})
		return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(string(data))), Header: make(http.Header), Request: r}, nil
	})}}
	result, err := client.Sync(context.Background(), app.Connection{URL: "https://console.example", Credential: "private"}, app.Sync{State: "stopped"})
	if err != nil || result.DesiredState != "stopped" {
		t.Fatalf("result=%+v err=%v", result, err)
	}
	if client.httpClient().CheckRedirect(&http.Request{}, nil) == nil {
		t.Fatal("redirect followed")
	}
}
