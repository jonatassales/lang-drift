package judge_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/anthropics/anthropic-sdk-go/option"

	"github.com/jonatassales/tao/packages/ex/src/application/plan"
	"github.com/jonatassales/tao/packages/ex/src/infrastructure/judge"
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

func mustJSON(t *testing.T, value any) []byte {
	t.Helper()

	data, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}

	return data
}
