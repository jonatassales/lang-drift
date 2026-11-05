package configuration

import (
	"math"
	"testing"
)

func TestConfiguration_MarkdownRefusesABudgetThatIsNotFinite(t *testing.T) {
	config := validConfig()
	budget := math.NaN()
	config.LLMBudgetUSD = &budget

	defer func() {
		if recover() == nil {
			t.Fatal("Markdown encoded a budget that is not finite")
		}
	}()

	config.Markdown()
}
