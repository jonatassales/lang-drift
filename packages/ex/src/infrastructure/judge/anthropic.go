// Package judge asks the planner's questions of a Claude model through the
// official Anthropic SDK. The credential is the SDK's to resolve, from the
// environment or a saved profile; it never passes through ex.
package judge

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"

	"github.com/jonatassales/tao/packages/ex/src/application/plan"
)

// Anthropic is the Claude adapter. Options are for tests: a local server
// and a test key.
type Anthropic struct {
	Options []option.RequestOption
}

// Ask sends one question and returns the text and the usage. A credential
// the API refuses makes the model unavailable; a refusal stop reason is an
// answer that says nothing.
func (adapter Anthropic) Ask(ctx context.Context, question plan.Question) (plan.Answer, error) {
	client := anthropic.NewClient(append([]option.RequestOption{option.WithMaxRetries(2)}, adapter.Options...)...)

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

// classify names the failure: a refused credential is the model being
// unavailable; anything else is the call's own failure.
func classify(question plan.Question, err error) error {
	var apiError *anthropic.Error
	if errors.As(err, &apiError) && (apiError.StatusCode == http.StatusUnauthorized || apiError.StatusCode == http.StatusForbidden) {
		return fmt.Errorf("%w: the API refused the credential for %s", plan.ErrModelUnavailable, question.Model)
	}

	return fmt.Errorf("ask %s about %s: %w", question.Model, question.Kind, err)
}
