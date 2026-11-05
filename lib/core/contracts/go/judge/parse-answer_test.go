package judge_test

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jonatassales/tao/lib/core/contracts/go/judge"
)

// TestParseAnswer reads every sample into its kind with every field
// populated, and refuses one answer per rule of the contract, naming the
// field that broke it.
func TestParseAnswer(t *testing.T) {
	t.Parallel()

	samples := map[judge.Kind]string{}
	for _, kind := range []judge.Kind{judge.KindSources, judge.KindEpisodes, judge.KindSuggestion, judge.KindReservedPart} {
		data, err := os.ReadFile(filepath.Join("schema", string(kind)+".answer.sample.json"))
		if err != nil {
			t.Fatal(err)
		}

		samples[kind] = string(data)
	}

	t.Run("the sources sample parses with every field populated", func(t *testing.T) {
		t.Parallel()

		answer, err := judge.ParseAnswer(judge.KindSources, samples[judge.KindSources])
		if err != nil {
			t.Fatal(err)
		}

		sources := answer.Sources
		if answer.Kind != judge.KindSources || sources == nil || len(sources.Sources) != 2 || len(sources.Facts) != 1 || len(sources.Criteria) != 1 || len(sources.Figures) != 1 {
			t.Fatalf("answer = %+v", answer)
		}

		first := sources.Sources[0]
		if first.Path == "" || first.Verdict != "kept" || first.Reason == "" || first.Confidence != 0.98 || first.Role != "ticket" || first.Summary == "" {
			t.Errorf("source = %+v", first)
		}
		if figure := sources.Figures[0]; figure.Topic != "resource" || figure.CPUMillis != 2000 || figure.MemoryMiB != 3072 || figure.Fact == "" || figure.Source == "" {
			t.Errorf("figure = %+v", figure)
		}
		if fact := sources.Facts[0]; fact.Topic != "component" || fact.Text == "" || fact.Source == "" {
			t.Errorf("fact = %+v", fact)
		}
		if criterion := sources.Criteria[0]; criterion.ID != "FIN-02" || criterion.Applies == "" {
			t.Errorf("criterion = %+v", criterion)
		}
	})

	t.Run("the episodes sample parses with every field populated", func(t *testing.T) {
		t.Parallel()

		answer, err := judge.ParseAnswer(judge.KindEpisodes, samples[judge.KindEpisodes])
		if err != nil {
			t.Fatal(err)
		}

		if answer.Episodes == nil || len(answer.Episodes.Episodes) != 1 || answer.Episodes.Episodes[0].Run == "" || answer.Episodes.Episodes[0].Reason == "" {
			t.Fatalf("answer = %+v", answer)
		}
	})

	t.Run("the suggestion sample parses with every field populated", func(t *testing.T) {
		t.Parallel()

		answer, err := judge.ParseAnswer(judge.KindSuggestion, samples[judge.KindSuggestion])
		if err != nil {
			t.Fatal(err)
		}

		suggestion := answer.Suggestion
		if suggestion == nil || suggestion.Model != "claude-sonnet-5-5" || suggestion.Effort != "medium" || suggestion.EffortReason == "" || suggestion.ModelReason == "" || len(suggestion.Reasons) != 2 {
			t.Fatalf("answer = %+v", answer)
		}
	})

	t.Run("the reserved-part sample parses with every field populated", func(t *testing.T) {
		t.Parallel()

		answer, err := judge.ParseAnswer(judge.KindReservedPart, samples[judge.KindReservedPart])
		if err != nil {
			t.Fatal(err)
		}

		part := answer.ReservedPart
		if part == nil || len(part.Paths) != 2 || part.Goal == "" || len(part.Acceptance) != 2 || part.Category == nil || *part.Category != "core-domain" || part.PlannedLines != 60 || part.TotalPlannedLines != 420 || part.ExcludedKind != nil || part.Reason == "" || part.Confidence != 0.8 {
			t.Fatalf("answer = %+v", answer)
		}
	})

	t.Run("a reserved part with no paths is no part, with its reason", func(t *testing.T) {
		t.Parallel()

		answer, err := judge.ParseAnswer(judge.KindReservedPart, `{"paths":[],"goal":"","acceptance":[],"category":null,"plannedLines":0,"totalPlannedLines":120,"excludedKind":"configuration","reason":"the ticket only wires configuration","confidence":0.7}`)
		if err != nil {
			t.Fatal(err)
		}

		if part := answer.ReservedPart; len(part.Paths) != 0 || part.ExcludedKind == nil || *part.ExcludedKind != "configuration" {
			t.Fatalf("answer = %+v", answer)
		}
	})

	refusals := []struct {
		name  string
		kind  judge.Kind
		text  string
		field string
		rule  string
	}{
		{"text that is not an object", judge.KindEpisodes, `[]`, "answer", "must be a JSON object"},
		{"a missing field", judge.KindEpisodes, `{}`, "episodes", "is missing"},
		{"a field outside the schema", judge.KindEpisodes, `{"episodes":[],"extra":1}`, "extra", "is not a field of the schema"},
		{"a value of the wrong type", judge.KindSuggestion, `{"model":1,"effort":"low","effortReason":"a","modelReason":"b","reasons":[]}`, "model", "must be a string, not a number"},
		{"an answer that is not an object", judge.KindEpisodes, `[]`, "answer", "must be a JSON object"},
		{"a list that is not an array", judge.KindEpisodes, `{"episodes":{}}`, "episodes", "must be an array"},
		{"an item that is not an object", judge.KindEpisodes, `{"episodes":[1]}`, "episodes[0]", "must be a JSON object"},
		{"a missing field of an item", judge.KindEpisodes, `{"episodes":[{"run":"r"}]}`, "episodes[0].reason", "is missing"},
		{"an empty reason", judge.KindEpisodes, `{"episodes":[{"run":"r","reason":""}]}`, "episodes[0].reason", "must not be empty"},
		{"a verdict outside the set", judge.KindSources, `{"sources":[{"path":"p","verdict":"maybe","reason":"r","confidence":0.5,"role":"","summary":""}],"facts":[],"criteria":[],"figures":[]}`, "sources[0].verdict", "must be one of kept, discarded"},
		{"a confidence above one", judge.KindSources, `{"sources":[{"path":"p","verdict":"kept","reason":"r","confidence":1.5,"role":"","summary":""}],"facts":[],"criteria":[],"figures":[]}`, "sources[0].confidence", "must be between 0 and 1"},
		{"a fact topic outside the set", judge.KindSources, `{"sources":[],"facts":[{"topic":"weather","text":"t","source":"s"}],"criteria":[],"figures":[]}`, "facts[0].topic", "must be one of"},
		{"a figure that is not about resources", judge.KindSources, `{"sources":[],"facts":[],"criteria":[],"figures":[{"topic":"budget","fact":"f","source":"s","cpuMillis":1,"memoryMiB":1}]}`, "figures[0].topic", "must be one of resource"},
		{"a negative figure", judge.KindSources, `{"sources":[],"facts":[],"criteria":[],"figures":[{"topic":"resource","fact":"f","source":"s","cpuMillis":-1,"memoryMiB":1}]}`, "figures[0].cpuMillis", "must not be negative"},
		{"an empty criterion id", judge.KindSources, `{"sources":[],"facts":[],"criteria":[{"id":"","applies":"a"}],"figures":[]}`, "criteria[0].id", "must not be empty"},
		{"an effort outside the set", judge.KindSuggestion, `{"model":"m","effort":"max","effortReason":"a","modelReason":"b","reasons":[]}`, "effort", "must be one of low, medium, high"},
		{"an empty model", judge.KindSuggestion, `{"model":"","effort":"low","effortReason":"a","modelReason":"b","reasons":[]}`, "model", "must not be empty"},
		{"too many reasons", judge.KindSuggestion, `{"model":"m","effort":"low","effortReason":"a","modelReason":"b","reasons":["1","2","3","4"]}`, "reasons", "at most 3"},
		{"an empty reason among the reasons", judge.KindSuggestion, `{"model":"m","effort":"low","effortReason":"a","modelReason":"b","reasons":[""]}`, "reasons[0]", "must not be empty"},
		{"a category outside the set", judge.KindReservedPart, `{"paths":["a"],"goal":"g","acceptance":["c"],"category":"easy","plannedLines":1,"totalPlannedLines":2,"excludedKind":null,"reason":"r","confidence":0.5}`, "category", "must be one of core-domain, hard-part, new-stack"},
		{"an excluded kind outside the set", judge.KindReservedPart, `{"paths":[],"goal":"","acceptance":[],"category":null,"plannedLines":0,"totalPlannedLines":2,"excludedKind":"tests","reason":"r","confidence":0.5}`, "excludedKind", "must be one of"},
		{"planned lines above the total", judge.KindReservedPart, `{"paths":["a"],"goal":"g","acceptance":["c"],"category":"hard-part","plannedLines":3,"totalPlannedLines":2,"excludedKind":null,"reason":"r","confidence":0.5}`, "plannedLines", "must not exceed totalPlannedLines"},
		{"paths without a goal", judge.KindReservedPart, `{"paths":["a"],"goal":"","acceptance":["c"],"category":"hard-part","plannedLines":1,"totalPlannedLines":2,"excludedKind":null,"reason":"r","confidence":0.5}`, "goal", "must not be empty"},
		{"paths without acceptance", judge.KindReservedPart, `{"paths":["a"],"goal":"g","acceptance":[],"category":"hard-part","plannedLines":1,"totalPlannedLines":2,"excludedKind":null,"reason":"r","confidence":0.5}`, "acceptance", "must name at least one criterion"},
		{"paths without a category", judge.KindReservedPart, `{"paths":["a"],"goal":"g","acceptance":["c"],"category":null,"plannedLines":1,"totalPlannedLines":2,"excludedKind":null,"reason":"r","confidence":0.5}`, "category", "must be set when paths are proposed"},
		{"negative planned lines", judge.KindReservedPart, `{"paths":[],"goal":"","acceptance":[],"category":null,"plannedLines":-1,"totalPlannedLines":2,"excludedKind":null,"reason":"r","confidence":0.5}`, "plannedLines", "must not be negative"},
		{"a negative total", judge.KindReservedPart, `{"paths":[],"goal":"","acceptance":[],"category":null,"plannedLines":0,"totalPlannedLines":-2,"excludedKind":null,"reason":"r","confidence":0.5}`, "totalPlannedLines", "must not be negative"},
		{"an empty acceptance criterion", judge.KindReservedPart, `{"paths":["a"],"goal":"g","acceptance":[""],"category":"hard-part","plannedLines":1,"totalPlannedLines":2,"excludedKind":null,"reason":"r","confidence":0.5}`, "acceptance[0]", "must not be empty"},
		{"an empty path", judge.KindReservedPart, `{"paths":[""],"goal":"g","acceptance":["c"],"category":"hard-part","plannedLines":1,"totalPlannedLines":2,"excludedKind":null,"reason":"r","confidence":0.5}`, "paths[0]", "must not be empty"},
	}

	for _, refusal := range refusals {
		t.Run("refuses "+refusal.name, func(t *testing.T) {
			t.Parallel()

			_, err := judge.ParseAnswer(refusal.kind, refusal.text)

			var failure judge.FieldError
			if !errors.As(err, &failure) {
				t.Fatalf("err = %v, want a FieldError", err)
			}
			if failure.Field != refusal.field || !strings.Contains(failure.Rule, refusal.rule) {
				t.Errorf("refused %s: %s, want %s: %s", failure.Field, failure.Rule, refusal.field, refusal.rule)
			}
		})
	}

	t.Run("refuses a kind it does not know", func(t *testing.T) {
		t.Parallel()

		if _, err := judge.ParseAnswer(judge.Kind("weather"), `{}`); err == nil || !strings.Contains(err.Error(), "weather") {
			t.Fatalf("err = %v", err)
		}
	})
}
