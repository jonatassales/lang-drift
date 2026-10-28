package support

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Prompter asks questions on the diagnostics stream and reads the answers,
// so prompts never mix with a command's result. Every answer is required: an
// empty answer cancels the guidance before anything is written.
type Prompter struct {
	answers   *bufio.Reader
	questions io.Writer
}

// NewPrompter reads answers from input and writes questions to questions.
func NewPrompter(input io.Reader, questions io.Writer) *Prompter {
	return &Prompter{answers: bufio.NewReader(input), questions: questions}
}

// Say writes one line of guidance.
func (prompter *Prompter) Say(line string) error {
	_, err := fmt.Fprintln(prompter.questions, line)

	return err
}

// Text asks for a required answer.
func (prompter *Prompter) Text(ctx context.Context, label string) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}

	if _, err := fmt.Fprintf(prompter.questions, "%s: ", label); err != nil {
		return "", err
	}

	line, err := prompter.answers.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return "", err
	}

	answer := strings.TrimSpace(line)
	if answer == "" {
		return "", UsageError{Cause: fmt.Errorf("%s is required; guidance canceled without writing files", label)}
	}

	return answer, nil
}

// Integer asks for a required whole number.
func (prompter *Prompter) Integer(ctx context.Context, label string) (int, error) {
	answer, err := prompter.Text(ctx, label)
	if err != nil {
		return 0, err
	}

	value, err := strconv.Atoi(answer)
	if err != nil {
		return 0, UsageError{Cause: fmt.Errorf("%s must be a whole number", label)}
	}

	return value, nil
}

// List asks for a required comma-separated list.
func (prompter *Prompter) List(ctx context.Context, label string) ([]string, error) {
	answer, err := prompter.Text(ctx, label)
	if err != nil {
		return nil, err
	}

	var items []string

	for _, item := range strings.Split(answer, ",") {
		if item = strings.TrimSpace(item); item != "" {
			items = append(items, item)
		}
	}

	return items, nil
}
