package judge_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/anthropics/anthropic-sdk-go/option"

	"github.com/jonatassales/tao/lib/cli/src/application/plan"
	"github.com/jonatassales/tao/lib/cli/src/infrastructure/judge"
)

func server(t *testing.T, status int, body string, seen *map[string]any) *httptest.Server {
	t.Helper()

	return httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		data, _ := io.ReadAll(request.Body)
		if seen != nil {
			_ = json.Unmarshal(data, seen)
		}
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(status)
		_, _ = io.WriteString(writer, body)
	}))
}

var question = plan.Question{Kind: plan.SourcesQuestion, Model: "claude-sonnet-5-5", System: "the system", Prompt: "the material", MaxTokens: 8000}

func TestAnthropicReturnsTheTextAndTheUsage(t *testing.T) {
	t.Parallel()

	var seen map[string]any
	api := server(t, 200, `{"id":"msg_1","type":"message","role":"assistant","model":"claude-sonnet-5-5","content":[{"type":"text","text":"{\"sources\":[]}"}],"stop_reason":"end_turn","usage":{"input_tokens":10000,"output_tokens":1000}}`, &seen)
	defer api.Close()

	adapter := judge.Anthropic{Options: []option.RequestOption{option.WithBaseURL(api.URL), option.WithAPIKey("test-key"), option.WithMaxRetries(0)}}

	answer, err := adapter.Ask(context.Background(), question)
	if err != nil {
		t.Fatal(err)
	}

	if answer.Text != `{"sources":[]}` || answer.InputTokens != 10000 || answer.OutputTokens != 1000 || answer.Refused {
		t.Fatalf("answer = %+v", answer)
	}

	if seen["model"] != "claude-sonnet-5-5" || seen["max_tokens"] != 8000.0 || !strings.Contains(string(mustJSON(t, seen["system"])), "the system") {
		t.Fatalf("request = %v", seen)
	}
}

func TestAnthropicReportsARefusal(t *testing.T) {
	t.Parallel()

	api := server(t, 200, `{"id":"msg_2","type":"message","role":"assistant","model":"claude-sonnet-5-5","content":[],"stop_reason":"refusal","usage":{"input_tokens":500,"output_tokens":0}}`, nil)
	defer api.Close()

	adapter := judge.Anthropic{Options: []option.RequestOption{option.WithBaseURL(api.URL), option.WithAPIKey("test-key"), option.WithMaxRetries(0)}}

	answer, err := adapter.Ask(context.Background(), question)
	if err != nil || !answer.Refused || answer.InputTokens != 500 {
		t.Fatalf("answer = %+v, %v", answer, err)
	}
}

func TestAnthropicTreatsARefusedCredentialAsUnavailable(t *testing.T) {
	t.Parallel()

	api := server(t, 401, `{"type":"error","error":{"type":"authentication_error","message":"invalid x-api-key"}}`, nil)
	defer api.Close()

	adapter := judge.Anthropic{Options: []option.RequestOption{option.WithBaseURL(api.URL), option.WithAPIKey("bad-key"), option.WithMaxRetries(0)}}

	if _, err := adapter.Ask(context.Background(), question); !errors.Is(err, plan.ErrModelUnavailable) {
		t.Fatalf("error = %v, want ErrModelUnavailable", err)
	}
}

func TestAnthropicReportsAServerFailure(t *testing.T) {
	t.Parallel()

	api := server(t, 500, `{"type":"error","error":{"type":"api_error","message":"boom"}}`, nil)
	defer api.Close()

	adapter := judge.Anthropic{Options: []option.RequestOption{option.WithBaseURL(api.URL), option.WithAPIKey("test-key"), option.WithMaxRetries(0)}}

	_, err := adapter.Ask(context.Background(), question)
	if err == nil || errors.Is(err, plan.ErrModelUnavailable) {
		t.Fatalf("error = %v, want a call failure", err)
	}
}

// messageBody is one answer of the Messages API.
const messageBody = `{"id":"msg_1","type":"message","role":"assistant","model":"claude-sonnet-5-5","content":[{"type":"text","text":"{\"sources\":[]}"}],"stop_reason":"end_turn","usage":{"input_tokens":10000,"output_tokens":1000}}`

// counting is a server that counts the calls it receives and answers each
// with the status and the body.
func counting(calls *atomic.Int32, status int, body string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		calls.Add(1)
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(status)
		_, _ = io.WriteString(writer, body)
	}))
}

// unspecified is the server's address at http://0.0.0.0: not a loopback
// host, yet one the operating system delivers locally, so a call that got
// past the judge's guard would reach the server.
func unspecified(t *testing.T, api *httptest.Server) string {
	t.Helper()

	address, isTCP := api.Listener.Addr().(*net.TCPAddr)
	if !isTCP {
		t.Fatalf("the test server listens on %v", api.Listener.Addr())
	}

	return "http://0.0.0.0:" + strconv.Itoa(address.Port)
}

// The SDK sends the key as X-Api-Key, a header Go carries to the host a
// redirect names; the judge refuses the redirect instead of following it.
func TestAnthropicNeverFollowsARedirect(t *testing.T) {
	t.Parallel()

	var secondCalls atomic.Int32
	second := counting(&secondCalls, http.StatusOK, messageBody)
	defer second.Close()

	var keyedCalls atomic.Int32
	first := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Header.Get("X-Api-Key") == "test-key" {
			keyedCalls.Add(1)
		}
		http.Redirect(writer, request, second.URL+"/v1/messages", http.StatusTemporaryRedirect)
	}))
	defer first.Close()

	adapter := judge.Anthropic{Options: []option.RequestOption{option.WithBaseURL(first.URL), option.WithAPIKey("test-key"), option.WithMaxRetries(0)}}

	_, err := adapter.Ask(context.Background(), question)
	if err == nil {
		t.Fatal("a redirected call was answered")
	}

	if keyedCalls.Load() != 1 || secondCalls.Load() != 0 {
		t.Fatalf("the first host received %d calls with the key and the second %d calls; the redirect must not be followed", keyedCalls.Load(), secondCalls.Load())
	}

	if !strings.Contains(err.Error(), "redirect") || strings.Contains(err.Error(), first.URL) || strings.Contains(err.Error(), second.URL) {
		t.Fatalf("error = %v, want the refused redirect without a URL", err)
	}
}

// A response is read up to 4 MiB: a well-formed answer past the bound fails
// instead of being decoded.
func TestAnthropicBoundsTheResponseItReads(t *testing.T) {
	t.Parallel()

	oversized := strings.Replace(messageBody, `{\"sources\":[]}`, strings.Repeat("a", 4<<20), 1)
	api := server(t, http.StatusOK, oversized, nil)
	defer api.Close()

	adapter := judge.Anthropic{Options: []option.RequestOption{option.WithBaseURL(api.URL), option.WithAPIKey("test-key"), option.WithMaxRetries(0)}}

	answer, err := adapter.Ask(context.Background(), question)
	if err == nil {
		t.Fatalf("an answer over 4 MiB was accepted with %d bytes of text", len(answer.Text))
	}

	if !strings.Contains(err.Error(), "4 MiB") {
		t.Fatalf("error = %v, want the 4 MiB bound", err)
	}
}

// An error status reaches the planner as its code and request id, never as
// the provider's body or the URL that was called.
func TestAnthropicReportsAnErrorStatusByItsCodeAndRequestID(t *testing.T) {
	t.Parallel()

	for _, status := range []int{http.StatusBadRequest, http.StatusUnauthorized, http.StatusForbidden, http.StatusTooManyRequests, http.StatusInternalServerError, 529} {
		t.Run(strconv.Itoa(status), func(t *testing.T) {
			t.Parallel()

			api := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				writer.Header().Set("Content-Type", "application/json")
				writer.Header().Set("Request-Id", "req_011judge")
				writer.WriteHeader(status)
				_, _ = io.WriteString(writer, `{"type":"error","error":{"type":"api_error","message":"provider body that echoes the prompt"}}`)
			}))
			defer api.Close()

			adapter := judge.Anthropic{Options: []option.RequestOption{option.WithBaseURL(api.URL), option.WithAPIKey("test-key"), option.WithMaxRetries(0)}}

			_, err := adapter.Ask(context.Background(), question)
			if err == nil {
				t.Fatal("an error status was answered")
			}

			text := err.Error()
			if !strings.Contains(text, "HTTP "+strconv.Itoa(status)) || !strings.Contains(text, "req_011judge") {
				t.Fatalf("error = %q, want the status and the request id", text)
			}

			for _, leak := range []string{"provider body", "api_error", api.URL, api.Listener.Addr().String(), "test-key"} {
				if strings.Contains(text, leak) {
					t.Fatalf("error %q carries %q", text, leak)
				}
			}
		})
	}
}

// A base URL that is neither HTTPS nor HTTP on loopback is refused before
// any call leaves, wherever it was configured.
func TestAnthropicRefusesAPlainHTTPBaseURLBeforeAnyCall(t *testing.T) {
	t.Parallel()

	var calls atomic.Int32
	api := counting(&calls, http.StatusOK, messageBody)
	defer api.Close()

	adapter := judge.Anthropic{Options: []option.RequestOption{option.WithBaseURL(unspecified(t, api)), option.WithAPIKey("test-key"), option.WithMaxRetries(0)}}

	_, err := adapter.Ask(context.Background(), question)
	if !errors.Is(err, plan.ErrModelUnavailable) {
		t.Fatalf("error = %v, want ErrModelUnavailable", err)
	}

	if calls.Load() != 0 {
		t.Fatalf("the server received %d calls; a plain HTTP base URL must send none", calls.Load())
	}
}

// ANTHROPIC_BASE_URL is held to the same rule before the client is built,
// and a malformed one is refused without being quoted.
func TestAnthropicChecksTheEnvironmentBaseURL(t *testing.T) {
	var calls atomic.Int32
	api := counting(&calls, http.StatusOK, messageBody)
	defer api.Close()

	adapter := judge.Anthropic{Options: []option.RequestOption{option.WithAPIKey("test-key"), option.WithMaxRetries(0)}}

	for _, refused := range []struct{ name, baseURL string }{
		{name: "plain HTTP off loopback", baseURL: unspecified(t, api)},
		{name: "malformed with userinfo", baseURL: "http://leaked-user:leaked-pass@127.0.0.1:1:x"},
	} {
		t.Run(refused.name, func(t *testing.T) {
			t.Setenv("ANTHROPIC_BASE_URL", refused.baseURL)
			before := calls.Load()

			_, err := adapter.Ask(context.Background(), question)
			if !errors.Is(err, plan.ErrModelUnavailable) {
				t.Fatalf("error = %v, want ErrModelUnavailable", err)
			}

			if calls.Load() != before || strings.Contains(err.Error(), refused.baseURL) || strings.Contains(err.Error(), "leaked") {
				t.Fatalf("the refusal sent %d calls or quoted the URL: %v", calls.Load()-before, err)
			}
		})
	}

	t.Run("loopback", func(t *testing.T) {
		t.Setenv("ANTHROPIC_BASE_URL", api.URL)

		answer, err := adapter.Ask(context.Background(), question)
		if err != nil || answer.Text != `{"sources":[]}` {
			t.Fatalf("a loopback base URL gave %+v, %v", answer, err)
		}
	})
}

func mustJSON(t *testing.T, value any) []byte {
	t.Helper()

	data, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}

	return data
}
