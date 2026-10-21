package support_test

import (
	"context"
	"errors"
	"slices"
	"strings"
	"testing"

	"github.com/exphere/ex/src/presentation/cli/support"
)

func TestPrompter(t *testing.T) {
	t.Parallel()

	var questions strings.Builder
	prompter := support.NewPrompter(strings.NewReader("  team-payments  \n8080\n web, ,jobs/batch ,\nlast answer"), &questions)
	ctx := context.Background()

	if err := prompter.Say("Configure web"); err != nil {
		t.Fatal(err)
	}

	if answer, err := prompter.Text(ctx, "Owner reference"); err != nil || answer != "team-payments" {
		t.Errorf("Text() = %q, %v; want the trimmed answer", answer, err)
	}

	if port, err := prompter.Integer(ctx, "Application port"); err != nil || port != 8080 {
		t.Errorf("Integer() = %d, %v", port, err)
	}

	if items, err := prompter.List(ctx, "Components"); err != nil || !slices.Equal(items, []string{"web", "jobs/batch"}) {
		t.Errorf("List() = %q, %v; want trimmed items without empties", items, err)
	}

	if answer, err := prompter.Text(ctx, "Final"); err != nil || answer != "last answer" {
		t.Errorf("Text() at the end of input = %q, %v; want the unterminated line", answer, err)
	}

	if want := "Configure web\nOwner reference: Application port: Components: Final: "; questions.String() != want {
		t.Errorf("questions = %q, want %q", questions.String(), want)
	}
}

func TestPrompter_RefusesMissingAnswers(t *testing.T) {
	t.Parallel()

	var usage support.UsageError

	for name, input := range map[string]string{"blank answer": "   \n", "end of input": ""} {
		_, err := support.NewPrompter(strings.NewReader(input), &strings.Builder{}).Text(context.Background(), "Owner reference")
		if !errors.As(err, &usage) || !strings.Contains(err.Error(), "Owner reference is required") {
			t.Errorf("%s: Text() error = %v, want a usage error", name, err)
		}
	}

	if _, err := support.NewPrompter(strings.NewReader("eighty\n"), &strings.Builder{}).Integer(context.Background(), "Port"); !errors.As(err, &usage) {
		t.Errorf("Integer() error = %v, want a usage error", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	var questions strings.Builder
	if _, err := support.NewPrompter(strings.NewReader("answer\n"), &questions).Text(ctx, "Owner reference"); !errors.Is(err, context.Canceled) || questions.Len() != 0 {
		t.Errorf("Text() after cancellation = %v with questions %q; want context.Canceled before asking", err, questions.String())
	}
}
