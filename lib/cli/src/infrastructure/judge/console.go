package judge

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/jonatassales/tao/lib/cli/src/application/connected"
	"github.com/jonatassales/tao/lib/cli/src/application/plan"
	contract "github.com/jonatassales/tao/lib/core/contracts/go/judge"
)

// The bounds of the console's calls. A question waits longer than the
// route's own 120 seconds, so the route's answer arrives before the CLI gives
// up; the connection check answers at once.
const (
	askTimeout   = 130 * time.Second
	checkTimeout = 10 * time.Second
)

// Console is the planner's judge: the judge route of the console a project is
// connected to, which asks the platform's model with LangDrift's own key.
// Timeout bounds one question; zero takes askTimeout.
type Console struct {
	Timeout time.Duration
}

// Ask posts one question with the deployment secret as a bearer and reads the
// route's answer. A paid call is sent once and never retried.
func (console Console) Ask(ctx context.Context, connection plan.Connection, question plan.Question) (plan.Answer, plan.AskOutcome) {
	target, err := judgeEndpoint(connection.URL)
	if err != nil {
		return plan.Answer{}, outcome(question.Kind, plan.AskUnavailable, err.Error())
	}

	body, err := json.Marshal(contract.Request{
		Kind:          contract.Kind(question.Kind),
		Model:         question.Model,
		System:        question.System,
		Prompt:        question.Prompt,
		MaxTokens:     question.MaxTokens,
		PromptVersion: question.PromptVersion,
	})
	if err != nil {
		return plan.Answer{}, outcome(question.Kind, plan.AskMalformed, err.Error())
	}

	bounded, cancel := context.WithTimeout(ctx, console.timeout())
	defer cancel()

	request, err := http.NewRequestWithContext(bounded, http.MethodPost, target, bytes.NewReader(body))
	if err != nil {
		return plan.Answer{}, outcome(question.Kind, plan.AskUnavailable, err.Error())
	}

	request.Header.Set("Authorization", "Bearer "+connection.Secret)
	request.Header.Set("Content-Type", "application/json")

	response, err := httpClient().Do(request)
	if err != nil {
		return failedCall(question.Kind, bounded, err)
	}
	defer response.Body.Close()

	payload, err := io.ReadAll(response.Body)
	if err != nil {
		return failedCall(question.Kind, bounded, err)
	}

	return consoleAnswer(question, response.StatusCode, payload)
}

// Check is the connection check dft agent connect performs: the judge route
// answers 204 to the deployment's secret and 401 without it.
func (console Console) Check(ctx context.Context, consoleURL string, secret string) error {
	target, err := judgeEndpoint(consoleURL)
	if err != nil {
		return err
	}

	bounded, cancel := context.WithTimeout(ctx, checkTimeout)
	defer cancel()

	request, err := http.NewRequestWithContext(bounded, http.MethodGet, target, nil)
	if err != nil {
		return err
	}

	request.Header.Set("Authorization", "Bearer "+secret)

	response, err := httpClient().Do(request)
	if err != nil {
		return fmt.Errorf("the console at %s could not be reached: %w", consoleURL, err)
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusUnauthorized {
		return fmt.Errorf("%w: the console at %s refused RUNTIME_SECRET", connected.ErrAccessRevoked, consoleURL)
	}

	if response.StatusCode != http.StatusNoContent {
		return fmt.Errorf("the console at %s answered HTTP %d to the connection check", consoleURL, response.StatusCode)
	}

	return nil
}

func (console Console) timeout() time.Duration {
	if console.Timeout > 0 {
		return console.Timeout
	}

	return askTimeout
}
