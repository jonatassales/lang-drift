package connected

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	app "github.com/jonatassales/tao/lib/cli/src/application/connected"
	"github.com/jonatassales/tao/lib/core/contracts/go/protocol"
)

// MCP uses bounded Streamable HTTP JSON responses with no credential-bearing
// redirects. The endpoint must be HTTPS, except explicit loopback development.
type MCP struct {
	Client      *http.Client
	next        atomic.Uint64
	mutex       sync.Mutex
	initialized map[string]bool
}

func endpoint(raw string) (string, error) {
	u, err := url.Parse(raw)
	if err != nil {
		// The parse error quotes the URL, userinfo included.
		return "", errors.New("console URL is malformed")
	}
	ip := net.ParseIP(u.Hostname())
	loopback := u.Hostname() == "localhost" || (ip != nil && ip.IsLoopback())
	if u.User != nil || u.RawQuery != "" || u.Fragment != "" || u.Hostname() == "" || (u.Scheme != "https" && !(u.Scheme == "http" && loopback)) {
		return "", errors.New("console URL must be HTTPS (HTTP is allowed only on loopback), without credentials or query")
	}
	u.Path = strings.TrimRight(u.Path, "/")
	if !strings.HasSuffix(u.Path, "/api/mcp") {
		u.Path += "/api/mcp"
	}
	return u.String(), nil
}
func (client *MCP) httpClient() *http.Client {
	var out http.Client
	if client.Client != nil {
		out = *client.Client
	}
	out.Timeout = 4 * time.Second
	out.CheckRedirect = func(*http.Request, []*http.Request) error { return errors.New("console redirects are refused") }
	return &out
}
func (client *MCP) call(ctx context.Context, raw, credential, method string, params any, result any) error {
	target, err := endpoint(raw)
	if err != nil {
		return err
	}
	id := client.next.Add(1)
	body, err := json.Marshal(map[string]any{"jsonrpc": "2.0", "id": id, "method": method, "params": params})
	if err != nil {
		return err
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, target, bytes.NewReader(body))
	if err != nil {
		return err
	}
	request.Header.Set("Authorization", "Bearer "+credential)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json, text/event-stream")
	request.Header.Set("MCP-Protocol-Version", protocol.Version)
	response, err := client.httpClient().Do(request)
	if err != nil {
		return errors.New("console connection unavailable")
	}
	defer response.Body.Close()
	if response.StatusCode == http.StatusUnauthorized {
		return app.ErrAccessRevoked
	}
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("console returned HTTP %d", response.StatusCode)
	}
	data, err := io.ReadAll(io.LimitReader(response.Body, (4<<20)+1))
	if err != nil {
		return err
	}
	if len(data) > 4<<20 {
		return errors.New("console response exceeds 4 MiB")
	}
	var envelope struct {
		JSONRPC string          `json:"jsonrpc"`
		ID      uint64          `json:"id"`
		Result  json.RawMessage `json:"result"`
		Error   *struct {
			Code int `json:"code"`
		} `json:"error"`
	}
	if err = json.Unmarshal(data, &envelope); err != nil {
		return errors.New("console returned invalid JSON-RPC")
	}
	if envelope.ID != id || envelope.JSONRPC != "2.0" {
		return errors.New("console response identity mismatch")
	}
	if envelope.Error != nil {
		return fmt.Errorf("console operation rejected (%d)", envelope.Error.Code)
	}
	if result == nil {
		return nil
	}
	if method == "tools/call" {
		var tool struct {
			Structured json.RawMessage `json:"structuredContent"`
			Content    []struct {
				Type string `json:"type"`
				Text string `json:"text"`
			} `json:"content"`
			IsError bool `json:"isError"`
		}
		if err = json.Unmarshal(envelope.Result, &tool); err != nil {
			return err
		}
		if tool.IsError {
			return errors.New("console tool rejected the operation")
		}
		if len(tool.Structured) > 0 {
			return json.Unmarshal(tool.Structured, result)
		}
		for _, item := range tool.Content {
			if item.Type == "text" {
				return json.Unmarshal([]byte(item.Text), result)
			}
		}
		return errors.New("console tool returned no structured result")
	}
	return json.Unmarshal(envelope.Result, result)
}
func (client *MCP) tool(ctx context.Context, url, credential, name string, args, result any) error {
	if err := client.initialize(ctx, url, credential); err != nil {
		return err
	}
	return client.call(ctx, url, credential, "tools/call", map[string]any{"name": name, "arguments": args}, result)
}
func (client *MCP) Sync(ctx context.Context, record app.Connection, request app.Sync) (app.Intent, error) {
	var result app.Intent
	err := client.tool(ctx, record.URL, record.Credential, protocol.ToolSync, request, &result)
	return result, err
}
func (client *MCP) Claim(ctx context.Context, record app.Connection, intent app.Intent) (app.Claim, error) {
	var result app.Claim
	err := client.tool(ctx, record.URL, record.Credential, protocol.ToolClaim, map[string]any{"environmentId": record.EnvironmentID, "intentVersion": intent.IntentVersion, "revision": intent.Revision}, &result)
	return result, err
}

func (client *MCP) ProjectSync(ctx context.Context, url, secret string) (app.ProjectMembership, error) {
	var result app.ProjectMembership
	err := client.tool(ctx, url, secret, protocol.ToolProjectSync, map[string]any{}, &result)
	return result, err
}
func (client *MCP) ProjectBind(ctx context.Context, url, secret, runtime, environment string) (app.Binding, error) {
	var result app.Binding
	err := client.tool(ctx, url, secret, protocol.ToolProjectBind, map[string]any{"runtimeId": runtime, "environmentId": environment}, &result)
	return result, err
}

func (client *MCP) initialize(ctx context.Context, raw, credential string) error {
	client.mutex.Lock()
	defer client.mutex.Unlock()
	key := raw + "\x00" + credential
	if client.initialized[key] {
		return nil
	}
	var result struct {
		ProtocolVersion string `json:"protocolVersion"`
	}
	if err := client.call(ctx, raw, credential, "initialize", map[string]any{"protocolVersion": protocol.Version, "capabilities": map[string]any{}, "clientInfo": map[string]string{"name": "langdrift-project-agent", "version": "1"}}, &result); err != nil {
		return err
	}
	if result.ProtocolVersion != protocol.Version {
		return errors.New("console MCP protocol version is unsupported")
	}
	target, err := endpoint(raw)
	if err != nil {
		return err
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, target, strings.NewReader(`{"jsonrpc":"2.0","method":"notifications/initialized"}`))
	if err != nil {
		return err
	}
	request.Header.Set("Authorization", "Bearer "+credential)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json, text/event-stream")
	request.Header.Set("MCP-Protocol-Version", protocol.Version)
	response, err := client.httpClient().Do(request)
	if err != nil {
		return errors.New("console initialization unavailable")
	}
	response.Body.Close()
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return fmt.Errorf("console initialization returned HTTP %d", response.StatusCode)
	}
	if client.initialized == nil {
		client.initialized = map[string]bool{}
	}
	client.initialized[key] = true
	return nil
}
