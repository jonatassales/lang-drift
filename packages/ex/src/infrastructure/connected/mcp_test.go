package connected

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"slices"
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

// Go's parse error quotes the URL, userinfo included; a malformed console URL
// is refused with neither, before any request.
func TestMCP_MalformedURLIsRefusedWithoutBeingQuoted(t *testing.T) {
	for _, raw := range []string{"https://leaked-user:leaked-pass@console.example:port", "https://leaked-user:leaked-pass@console example", "https://leaked-user:leaked-pass@console.example/\x7f"} {
		t.Run(raw, func(t *testing.T) {
			calls := 0
			client := &MCP{Client: &http.Client{Transport: roundTrip(func(r *http.Request) (*http.Response, error) { calls++; return nil, nil })}}

			_, err := client.Bind(context.Background(), raw, "secret", "env")
			if err == nil || calls != 0 {
				t.Fatalf("a malformed console URL sent %d requests, error %v", calls, err)
			}

			if strings.Contains(err.Error(), "leaked") || strings.Contains(err.Error(), "console.example") {
				t.Fatalf("the error quotes the console URL: %v", err)
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

func mcpResponse(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     http.Header{"Content-Type": []string{"application/json"}},
	}
}

func TestMCP_UnauthorizedMapsToRevokedAccess(t *testing.T) {
	client := &MCP{Client: &http.Client{Transport: roundTrip(func(r *http.Request) (*http.Response, error) {
		return mcpResponse(http.StatusUnauthorized, `{}`), nil
	})}}

	_, err := client.Sync(context.Background(), app.Connection{URL: "https://console.example/api/mcp", Credential: "private"}, app.Sync{})
	if !errors.Is(err, app.ErrAccessRevoked) {
		t.Fatalf("err = %v, want ErrAccessRevoked", err)
	}
	if err != nil && strings.Contains(err.Error(), "private") {
		t.Fatalf("error leaks the credential: %v", err)
	}
}

func TestMCP_RefusesOversizedResponses(t *testing.T) {
	huge := `{"jsonrpc":"2.0","id":1,"result":{"padding":"` + strings.Repeat("a", (4<<20)+1024) + `"}}`
	client := &MCP{Client: &http.Client{Transport: roundTrip(func(r *http.Request) (*http.Response, error) {
		return mcpResponse(http.StatusOK, huge), nil
	})}}

	_, err := client.Sync(context.Background(), app.Connection{URL: "https://console.example/api/mcp", Credential: "private"}, app.Sync{})
	if err == nil || !strings.Contains(err.Error(), "4 MiB") && !strings.Contains(err.Error(), "bounded") && !strings.Contains(err.Error(), "limit") {
		t.Fatalf("oversized response accepted: %v", err)
	}
}

func TestMCP_RefusesRedirects(t *testing.T) {
	var hosts []string
	client := &MCP{Client: &http.Client{Transport: roundTrip(func(r *http.Request) (*http.Response, error) {
		hosts = append(hosts, r.URL.Host)
		response := mcpResponse(http.StatusTemporaryRedirect, "")
		response.Header.Set("Location", "https://elsewhere.example/api/mcp")
		return response, nil
	})}}

	_, err := client.Sync(context.Background(), app.Connection{URL: "https://console.example/api/mcp", Credential: "private"}, app.Sync{})
	if err == nil {
		t.Fatal("a redirected call reported success")
	}
	if len(hosts) != 1 || hosts[0] != "console.example" {
		t.Fatalf("requests reached %v; the redirect target must never be contacted", hosts)
	}
}

func TestMCP_ToolMethodsReachTheWire(t *testing.T) {
	var tools []string
	client := &MCP{Client: &http.Client{Transport: roundTrip(func(r *http.Request) (*http.Response, error) {
		var input map[string]any
		_ = json.NewDecoder(r.Body).Decode(&input)
		if input["method"] == "initialize" {
			return mcpResponse(http.StatusOK, `{"jsonrpc":"2.0","id":`+jsonID(input)+`,"result":{"protocolVersion":"2025-11-25"}}`), nil
		}
		params, _ := input["params"].(map[string]any)
		name, _ := params["name"].(string)
		if name != "" {
			tools = append(tools, name)
		}
		return mcpResponse(http.StatusOK, `{"jsonrpc":"2.0","id":`+jsonID(input)+`,"result":{"structuredContent":{}}}`), nil
	})}}

	record := app.Connection{URL: "https://console.example", Credential: "private", EnvironmentID: "env"}
	if _, err := client.Bind(context.Background(), "https://console.example", "secret", "env"); err != nil {
		t.Fatal(err)
	}
	if _, err := client.Claim(context.Background(), record, app.Intent{Revision: 1, IntentVersion: 1}); err != nil {
		t.Fatal(err)
	}
	if _, err := client.ProjectSync(context.Background(), "https://console.example", "secret"); err != nil {
		t.Fatal(err)
	}
	if _, err := client.ProjectBind(context.Background(), "https://console.example", "secret", "runtime", "env"); err != nil {
		t.Fatal(err)
	}

	want := []string{"langdrift.bind", "langdrift.claim", "langdrift.project.sync", "langdrift.project.bind"}
	if !slices.Equal(tools, want) {
		t.Fatalf("tools = %v, want %v", tools, want)
	}
}

func jsonID(input map[string]any) string {
	data, _ := json.Marshal(input["id"])
	return string(data)
}
