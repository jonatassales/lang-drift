package records_test

import (
	"strings"
	"testing"

	"github.com/jonatassales/tao/lib/cli/src/infrastructure/records"
)

type pair struct {
	Input  int `json:"input"`
	Output int `json:"output"`
}

type sample struct {
	Name   string   `json:"name"`
	Tags   []string `json:"tags"`
	Empty  []string `json:"empty"`
	Tokens pair     `json:"tokens"`
	Note   *string  `json:"note"`
	Rule   string   `json:"rule"`
}

func TestFormat_FileProfileKeepsBiomesLayout(t *testing.T) {
	t.Parallel()

	value := sample{Name: "demo", Tags: []string{"gpu", "disk"}, Empty: []string{}, Tokens: pair{Input: 3700, Output: 240}, Rule: "sorted path<TAB>sha256 lines & more"}

	got, err := records.Format(value, records.FileProfile)
	if err != nil {
		t.Fatal(err)
	}

	want := "{\n\t\"name\": \"demo\",\n\t\"tags\": [\"gpu\", \"disk\"],\n\t\"empty\": [],\n\t\"tokens\": {\n\t\t\"input\": 3700,\n\t\t\"output\": 240\n\t},\n\t\"note\": null,\n\t\"rule\": \"sorted path<TAB>sha256 lines & more\"\n}\n"
	if string(got) != want {
		t.Fatalf("Format() =\n%s\nwant\n%s", got, want)
	}
}

func TestFormat_FenceProfileCollapsesSmallScalarObjects(t *testing.T) {
	t.Parallel()

	value := sample{Name: "demo", Tags: []string{"nextjs", "go"}, Empty: []string{}, Tokens: pair{Input: 3700, Output: 240}}

	got, err := records.Format(value, records.FenceProfile)
	if err != nil {
		t.Fatal(err)
	}

	want := "{\n  \"name\": \"demo\",\n  \"tags\": [\"nextjs\", \"go\"],\n  \"empty\": [],\n  \"tokens\": { \"input\": 3700, \"output\": 240 },\n  \"note\": null,\n  \"rule\": \"\"\n}\n"
	if string(got) != want {
		t.Fatalf("Format() =\n%s\nwant\n%s", got, want)
	}
}

func TestFormat_BreaksAnArrayThatWouldNotFit(t *testing.T) {
	t.Parallel()

	criteria := struct {
		Criteria []string `json:"criteria"`
	}{Criteria: []string{"same-component-path", "same-stack", "measured-resources", "measured-tokens"}}

	got, err := records.Format(criteria, records.FileProfile)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(got), "\"criteria\": [\n\t\t\"same-component-path\",\n") {
		t.Fatalf("a wide array stayed on one line:\n%s", got)
	}
}

func TestFormat_CountsTheTrailingCommaAndTheTab(t *testing.T) {
	t.Parallel()

	// With two tabs of indent (four columns), the key, the brackets, the quotes
	// and the trailing comma, this array is 81 columns on one line, so it breaks.
	nested := struct {
		Inner struct {
			Values []string `json:"values"`
			After  int      `json:"after"`
		} `json:"inner"`
	}{}
	nested.Inner.Values = []string{strings.Repeat("a", 30), strings.Repeat("b", 28)}

	got, err := records.Format(nested, records.FileProfile)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(got), "\"values\": [\n") {
		t.Fatalf("an array one column too wide stayed on one line:\n%s", got)
	}

	// One character shorter, 80 columns, it fits.
	nested.Inner.Values = []string{strings.Repeat("a", 30), strings.Repeat("b", 27)}

	got, err = records.Format(nested, records.FileProfile)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(got), "\"values\": [\"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaa\", \"bbbbbbbbbbbbbbbbbbbbbbbbbbb\"],\n") {
		t.Fatalf("an array that fits was broken:\n%s", got)
	}
}
