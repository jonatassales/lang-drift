package judge_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/jonatassales/tao/lib/cli/src/application/connected"
	"github.com/jonatassales/tao/lib/cli/src/application/plan"
	judgeinfra "github.com/jonatassales/tao/lib/cli/src/infrastructure/judge"
)

var question = plan.Question{
	Kind:          plan.SuggestionQuestion,
	Model:         "claude-sonnet-5-5",
	System:        "system text",
	Prompt:        "prompt text",
	MaxTokens:     8000,
	PromptVersion: "judge-1",
}

const judged = `{"text":"{}","usage":{"inputTokens":1840,"outputTokens":96},"refused":false,"reason":"","cost":{"cents":1,"provenance":"console","illustrative":true},"judgeVersion":"console-judge-1"}`

// A question goes out once to the judge route, as a body the contract's
// request schema accepts, with the recorded secret as a bearer; the judgment
// comes back with the console's usage, price and version.
func TestConsoleAsksTheJudgeRouteWithTheContract(t *testing.T) {
	t.Parallel()

	server, calls := route(t, func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPost || request.URL.Path != "/api/judge" {
			t.Errorf("asked %s %s", request.Method, request.URL.Path)
		}

		bearer := request.Header.Get("Authorization") == "Bearer "+consoleSecret
		if !bearer || request.Header.Get("Content-Type") != "application/json" {
			t.Errorf("headers = %v", request.Header)
		}

		body, err := io.ReadAll(request.Body)
		if err != nil {
			t.Fatal(err)
		}

		if violations := requestSchemaViolations(t, body); len(violations) > 0 {
			t.Errorf("the request breaks request.schema.json: %v", violations)
		}

		reply(t, writer, http.StatusOK, judged)
	})

	answer, outcome := judgeinfra.Console{}.Ask(context.Background(), connectionTo(server), question)

	if outcome.Status != plan.AskAnswered || outcome.Kind != plan.SuggestionQuestion {
		t.Fatalf("outcome = %+v", outcome)
	}

	want := plan.Answer{
		Text:         "{}",
		InputTokens:  1840,
		OutputTokens: 96,
		Cents:        1,
		Provenance:   "console",
		Illustrative: true,
		JudgeVersion: "console-judge-1",
	}
	if answer != want {
		t.Errorf("answer = %+v, want %+v", answer, want)
	}

	if calls.Load() != 1 {
		t.Errorf("the route was asked %d times", calls.Load())
	}
}

// Each answer of the route maps to the planner's outcome with its cause, and
// none is retried: a paid call is sent once. A server error, or a body
// outside the contract, leaves what the call cost unknown.
func TestConsoleMapsTheRoutesAnswers(t *testing.T) {
	t.Parallel()

	for _, scenario := range []struct {
		name       string
		status     int
		body       string
		want       plan.AskStatus
		reason     string
		unmeasured bool
	}{
		{"the secret refused", http.StatusUnauthorized, `{"code":"UNAUTHORIZED"}`, plan.AskUnauthorized, "refused RUNTIME_SECRET", false},
		{"a model not offered", http.StatusBadRequest, `{"code":"MODEL_NOT_OFFERED","offered":["claude-opus-5-5","claude-haiku-5-5"]}`, plan.AskNotOffered, "does not offer claude-sonnet-5-5; it offers claude-opus-5-5, claude-haiku-5-5", false},
		{"a question outside the contract", http.StatusBadRequest, `{"code":"INVALID_QUESTION","field":"maxTokens"}`, plan.AskMalformed, "INVALID_QUESTION maxTokens", false},
		{"an unreadable refusal", http.StatusBadRequest, `not json`, plan.AskMalformed, "unreadable", false},
		{"the platform's key missing", http.StatusServiceUnavailable, `{"code":"UNAVAILABLE","reason":"platform-key-missing"}`, plan.AskUnavailable, "unavailable · platform-key-missing", false},
		{"a body too large for the route", http.StatusRequestEntityTooLarge, `{"code":"PAYLOAD_TOO_LARGE"}`, plan.AskMalformed, "HTTP 413", false},
		{"another media type", http.StatusUnsupportedMediaType, `{"code":"UNSUPPORTED_MEDIA_TYPE"}`, plan.AskMalformed, "HTTP 415", false},
		{"another method", http.StatusMethodNotAllowed, `{"code":"METHOD_NOT_ALLOWED"}`, plan.AskMalformed, "HTTP 405", false},
		{"a route that is not there", http.StatusNotFound, ``, plan.AskUnavailable, "HTTP 404", false},
		{"a server error", http.StatusInternalServerError, ``, plan.AskUnavailable, "HTTP 500", true},
		{"the model's refusal", http.StatusOK, `{"text":"","usage":{"inputTokens":10,"outputTokens":2},"refused":true,"reason":"the model stopped with refusal, not end_turn","cost":{"cents":1,"provenance":"console","illustrative":true},"judgeVersion":"console-judge-1"}`, plan.AskRefused, "the model stopped with refusal", false},
		{"a field outside the contract", http.StatusOK, `{"text":"{}","extra":true}`, plan.AskMalformed, "outside the contract", true},
		{"no judge version", http.StatusOK, `{"text":"{}","usage":{"inputTokens":1,"outputTokens":1},"refused":false,"reason":"","cost":{"cents":1,"provenance":"console","illustrative":true},"judgeVersion":""}`, plan.AskMalformed, "judgeVersion is empty", true},
		{"a cost that is not the console's", http.StatusOK, `{"text":"{}","usage":{"inputTokens":1,"outputTokens":1},"refused":false,"reason":"","cost":{"cents":1,"provenance":"provider","illustrative":true},"judgeVersion":"console-judge-1"}`, plan.AskMalformed, "not console", true},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			t.Parallel()

			server, calls := route(t, func(writer http.ResponseWriter, _ *http.Request) {
				reply(t, writer, scenario.status, scenario.body)
			})

			answer, outcome := judgeinfra.Console{}.Ask(context.Background(), connectionTo(server), question)

			if outcome.Status != scenario.want || !strings.Contains(outcome.Reason, scenario.reason) {
				t.Errorf("outcome = %+v, want %s naming %q", outcome, scenario.want, scenario.reason)
			}
			if answer.Unmeasured != scenario.unmeasured {
				t.Errorf("unmeasured = %v, want %v", answer.Unmeasured, scenario.unmeasured)
			}
			if strings.Contains(outcome.Reason, consoleSecret) {
				t.Errorf("the reason carries the secret: %q", outcome.Reason)
			}
			if calls.Load() != 1 {
				t.Errorf("the route was asked %d times; a paid call is sent once", calls.Load())
			}
		})
	}
}

// The model's refusal still reports what the call used and cost, so the
// planner accounts for it.
func TestConsoleKeepsTheCostOfARefusal(t *testing.T) {
	t.Parallel()

	server, _ := route(t, func(writer http.ResponseWriter, _ *http.Request) {
		reply(t, writer, http.StatusOK, `{"text":"","usage":{"inputTokens":10,"outputTokens":2},"refused":true,"reason":"no","cost":{"cents":1,"provenance":"console","illustrative":true},"judgeVersion":"console-judge-1"}`)
	})

	answer, outcome := judgeinfra.Console{}.Ask(context.Background(), connectionTo(server), question)

	if outcome.Status != plan.AskRefused || answer.Cents != 1 || answer.InputTokens != 10 || answer.Text != "" {
		t.Errorf("answer = %+v, outcome = %+v", answer, outcome)
	}
}

// A console that does not answer within the bound times out, once, and what
// the call cost is unknown.
func TestConsoleTimesOut(t *testing.T) {
	t.Parallel()

	server, calls := route(t, func(_ http.ResponseWriter, request *http.Request) {
		select {
		case <-time.After(2 * time.Second):
		case <-request.Context().Done():
		}
	})

	answer, outcome := judgeinfra.Console{Timeout: 100 * time.Millisecond}.Ask(context.Background(), connectionTo(server), question)

	if outcome.Status != plan.AskTimedOut || !answer.Unmeasured || calls.Load() != 1 {
		t.Errorf("outcome = %+v, unmeasured %v, after %d calls", outcome, answer.Unmeasured, calls.Load())
	}
}

// A URL that would carry the secret in clear text, or with credentials of its
// own, is refused before anything is sent.
func TestConsoleRefusesAnUnsafeURLBeforeSending(t *testing.T) {
	t.Parallel()

	for _, raw := range []string{"http://console.example", "https://user:pass@console.example", "https://console.example?x=1", "https://console.example/#frag", "https://leaked:pass@console.example:port"} {
		t.Run(raw, func(t *testing.T) {
			t.Parallel()

			connection := plan.Connection{
				Recorded: true,
				URL:      raw,
				Secret:   consoleSecret,
			}

			answer, outcome := judgeinfra.Console{}.Ask(context.Background(), connection, question)

			if outcome.Status != plan.AskUnavailable || answer.Unmeasured || strings.Contains(outcome.Reason, "pass") {
				t.Errorf("outcome = %+v", outcome)
			}
		})
	}
}

// A redirect is refused: the secret never reaches the host it names.
func TestConsoleRefusesARedirect(t *testing.T) {
	t.Parallel()

	target, reached := route(t, func(http.ResponseWriter, *http.Request) {})
	server, _ := route(t, func(writer http.ResponseWriter, request *http.Request) {
		http.Redirect(writer, request, target.URL+"/api/judge", http.StatusTemporaryRedirect)
	})

	_, outcome := judgeinfra.Console{}.Ask(context.Background(), connectionTo(server), question)

	if outcome.Status != plan.AskUnavailable || reached.Load() != 0 {
		t.Errorf("outcome = %+v; the redirect target was reached %d times", outcome, reached.Load())
	}
}

// The connection check accepts the secret on a 204 and refuses it on a 401,
// naming the console.
func TestConsoleChecksTheConnection(t *testing.T) {
	t.Parallel()

	server, _ := route(t, func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodGet || request.URL.Path != "/api/judge" {
			t.Errorf("checked with %s %s", request.Method, request.URL.Path)
		}

		if request.Header.Get("Authorization") != "Bearer "+consoleSecret {
			writer.WriteHeader(http.StatusUnauthorized)

			return
		}

		writer.WriteHeader(http.StatusNoContent)
	})

	console := judgeinfra.Console{}

	if err := console.Check(context.Background(), server.URL, consoleSecret); err != nil {
		t.Errorf("the secret was refused: %v", err)
	}

	err := console.Check(context.Background(), server.URL, "wrong")
	refusedNamingTheConsole := errors.Is(err, connected.ErrAccessRevoked) && strings.Contains(err.Error(), server.URL)
	if !refusedNamingTheConsole {
		t.Errorf("a wrong secret gave %v", err)
	}
}

// Any answer to the connection check but 204 or 401 refuses the connection,
// naming the console and the status.
func TestConsoleCheckRefusesAnotherStatus(t *testing.T) {
	t.Parallel()

	server, _ := route(t, func(writer http.ResponseWriter, _ *http.Request) {
		writer.WriteHeader(http.StatusNotFound)
	})

	err := judgeinfra.Console{}.Check(context.Background(), server.URL, consoleSecret)
	namesTheStatus := err != nil && strings.Contains(err.Error(), "HTTP 404") && strings.Contains(err.Error(), server.URL)
	if !namesTheStatus {
		t.Errorf("a 404 gave %v", err)
	}
}
