// Package judge asks the planner's questions of a Claude model through the
// official Anthropic SDK. The credential is the SDK's to resolve, from the
// environment or a saved profile; it never passes through ex. The calls go
// through the judge's own HTTP client, which follows no redirect, sends
// nothing to a base URL that is not HTTPS (loopback aside) and reads no
// response past 4 MiB.
package judge

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"

	"github.com/jonatassales/tao/lib/cli/src/application/plan"
)

// Anthropic is the Claude adapter. Options are for tests: a local server
// and a test key. The judge's HTTP client is applied after them, so no
// option lifts its guards.
type Anthropic struct {
	Options []option.RequestOption
}

// Ask sends one question and returns the text and the usage. A credential
// the API refuses, or a base URL that may not carry it, makes the model
// unavailable; a refusal stop reason is an answer that says nothing.
func (adapter Anthropic) Ask(ctx context.Context, question plan.Question) (plan.Answer, error) {
	if environmentBaseURLRefused() {
		return plan.Answer{}, fmt.Errorf("%w: %w", plan.ErrModelUnavailable, errInsecureBaseURL)
	}

	client := anthropic.NewClient(adapter.options()...)

	message, err := client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.Model(question.Model),
		MaxTokens: int64(question.MaxTokens),
		System:    []anthropic.TextBlockParam{{Text: question.System}},
		Messages:  []anthropic.MessageParam{anthropic.NewUserMessage(anthropic.NewTextBlock(question.Prompt))},
	})
	if err != nil {
		return plan.Answer{}, classify(question, err)
	}

	answer := plan.Answer{
		InputTokens:  message.Usage.InputTokens,
		OutputTokens: message.Usage.OutputTokens,
		Refused:      message.StopReason == anthropic.StopReasonRefusal,
	}

	var text strings.Builder
	for _, block := range message.Content {
		if textBlock, isText := block.AsAny().(anthropic.TextBlock); isText {
			text.WriteString(textBlock.Text)
		}
	}
	answer.Text = text.String()

	return answer, nil
}

// options are two retries, the adapter's own options, and the judge's HTTP
// client last.
func (adapter Anthropic) options() []option.RequestOption {
	options := append([]option.RequestOption{option.WithMaxRetries(2)}, adapter.Options...)

	return append(options, option.WithHTTPClient(httpClient()))
}

// classify names the failure without the provider's words or the URL it
// called: an error status is its code and request id, and a failure on the
// way is a fixed phrase. A refused credential or base URL makes the model
// unavailable; anything else is the call's own failure.
func classify(question plan.Question, err error) error {
	var apiError *anthropic.Error
	var oversized *http.MaxBytesError
	var unreached *url.Error

	switch {
	case errors.As(err, &apiError) && credentialRefused(apiError):
		return fmt.Errorf("%w: the API refused the credential for %s · %s", plan.ErrModelUnavailable, question.Model, status(apiError))
	case errors.As(err, &apiError):
		return fmt.Errorf("ask %s about %s: the API answered %s", question.Model, question.Kind, status(apiError))
	case errors.Is(err, errInsecureBaseURL):
		return fmt.Errorf("%w: %w", plan.ErrModelUnavailable, errInsecureBaseURL)
	case errors.Is(err, errRedirect):
		return fmt.Errorf("ask %s about %s: %w", question.Model, question.Kind, errRedirect)
	case errors.As(err, &oversized):
		return fmt.Errorf("ask %s about %s: %w", question.Model, question.Kind, errOversized)
	case errors.As(err, &unreached):
		return fmt.Errorf("ask %s about %s: the API connection failed", question.Model, question.Kind)
	}

	return fmt.Errorf("ask %s about %s: %w", question.Model, question.Kind, err)
}

// credentialRefused reports whether the API refused the credential itself.
func credentialRefused(apiError *anthropic.Error) bool {
	return apiError.StatusCode == http.StatusUnauthorized || apiError.StatusCode == http.StatusForbidden
}

// status is an error status as the session may print and keep it: its code
// and request id, never the provider's body.
func status(apiError *anthropic.Error) string {
	if apiError.RequestID == "" {
		return fmt.Sprintf("HTTP %d · no request id", apiError.StatusCode)
	}

	return fmt.Sprintf("HTTP %d · request %s", apiError.StatusCode, apiError.RequestID)
}
