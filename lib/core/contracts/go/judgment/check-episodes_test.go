package judgment_test

import (
	"strings"
	"testing"

	"github.com/jonatassales/tao/lib/core/contracts/go/judgment"
)

// The simulation's six verdicts pass: three kept because they share the
// stack and scale, three dropped as docs only, another runtime and no
// broker.
func TestCheckEpisodesAcceptsTheSimulationsVerdicts(t *testing.T) {
	t.Parallel()

	if _, err := judgment.ParseAnswer(judgment.KindEpisodes, readSample(t, "episodes.answer.sample.json")); err != nil {
		t.Error(err)
	}
}

// A verdict the EPISODES rows could not show is refused naming its field: a
// malformed or repeated id, a verdict outside kept and dropped, a kept
// episode with a dropped reason, a dropped one with the kept reason, and a
// reason outside the set.
func TestCheckEpisodesRefusesAVerdictTheCardCannotShow(t *testing.T) {
	t.Parallel()

	for _, scenario := range []struct {
		name  string
		entry string
		field string
	}{
		{"a malformed id", `{"id":"412","verdict":"kept","reason":"shares the stack and scale"}`, "episodes[1].id"},
		{"a repeated id", `{"id":"ep-0412","verdict":"kept","reason":"shares the stack and scale"}`, "episodes[1].id"},
		{"another verdict", `{"id":"ep-0388","verdict":"maybe","reason":"shares the stack and scale"}`, "episodes[1].verdict"},
		{"a kept episode with a dropped reason", `{"id":"ep-0388","verdict":"kept","reason":"no broker"}`, "episodes[1].reason"},
		{"a dropped episode with the kept reason", `{"id":"ep-0388","verdict":"dropped","reason":"shares the stack and scale"}`, "episodes[1].reason"},
		{"a reason outside the set", `{"id":"ep-0388","verdict":"dropped","reason":"too old"}`, "episodes[1].reason"},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			t.Parallel()

			text := `{"episodes":[{"id":"ep-0412","verdict":"kept","reason":"shares the stack and scale"},` + scenario.entry + `]}`

			_, err := judgment.ParseAnswer(judgment.KindEpisodes, text)
			if err == nil || !strings.Contains(err.Error(), scenario.field+":") {
				t.Errorf("err = %v, want a refusal naming %s", err, scenario.field)
			}
		})
	}
}
