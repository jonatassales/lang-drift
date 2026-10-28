package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"testing"

	"github.com/jonatassales/tao/packages/ex/src/domain/connection"
)

// The .langdrift sample of ticket 0069 lives in the fixtures, byte for byte,
// beside a frozen copy of the ticket its run planned and of every source the
// ticket's links reached. These tests pin the fixtures to each other: the
// hashes the sample records must be the hashes of the frozen sources, so a
// test that reproduces the run can expect the sample's figures.
const (
	sampleFixture  = "testdata/langdrift/sample"
	projectFixture = "testdata/langdrift/project"
	samplePlan     = sampleFixture + "/microruntimes/runs/0001-console-microruntime/plan.json"
	sampleLap      = sampleFixture + "/microruntimes/runs/0001-console-microruntime/laps/0002.lap.md"
	sampleRuntime  = sampleFixture + "/microruntimes/runtimes/console-local/runtime.md"
)

var sampleFiles = []string{
	".gitignore",
	"README.md",
	"config.json",
	"finops/budgets.json",
	"finops/forecasts.json",
	"finops/gates.json",
	"finops/ledger/0001-console-microruntime.ledger.json",
	"finops/prices/0001-illustrative.prices.json",
	"harness/envelope.json",
	"harness/guardrails/0001-finops.gr.md",
	"harness/workflow.json",
	"microruntimes/AGENTS.md",
	"microruntimes/runs/0001-console-microruntime/0001-console-microruntime.run.md",
	"microruntimes/runs/0001-console-microruntime/laps/0001.lap.md",
	"microruntimes/runs/0001-console-microruntime/laps/0002.lap.md",
	"microruntimes/runs/0001-console-microruntime/plan.json",
	"microruntimes/runs/0001-console-microruntime/run.json",
	"microruntimes/runtimes/console-local/runtime.md",
}

type samplePlanFile struct {
	Collection struct {
		Context struct {
			SHA256 string `json:"sha256"`
		} `json:"context"`
	} `json:"collection"`
	Sources []struct {
		Path   string  `json:"path"`
		Bytes  *int64  `json:"bytes"`
		SHA256 *string `json:"sha256"`
		State  string  `json:"state"`
	} `json:"sources"`
	Runtime struct {
		Path   string `json:"path"`
		SHA256 string `json:"sha256"`
	} `json:"runtime"`
}

type sampleLapFence struct {
	Runtime struct {
		MarkdownSHA256 string `json:"markdownSha256"`
	} `json:"runtime"`
}

func TestSampleFixtureHoldsTheWholeSample(t *testing.T) {
	var found []string

	err := filepath.WalkDir(sampleFixture, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			return nil
		}
		relative, err := filepath.Rel(sampleFixture, path)
		if err != nil {
			return err
		}
		found = append(found, filepath.ToSlash(relative))

		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	sort.Strings(found)

	if !slices.Equal(found, sampleFiles) {
		t.Fatalf("sample fixture files:\n%s\nwant:\n%s", strings.Join(found, "\n"), strings.Join(sampleFiles, "\n"))
	}
}

func TestFrozenSourcesMatchTheSamplePlan(t *testing.T) {
	plan := readSamplePlan(t)

	var contextLines []string

	for _, source := range plan.Sources {
		if source.State != "read" {
			t.Fatalf("source %s is %q; the sample reads every source", source.Path, source.State)
		}
		if source.Bytes == nil || source.SHA256 == nil {
			t.Fatalf("source %s records no size or hash", source.Path)
		}

		content, err := os.ReadFile(filepath.Join(projectFixture, filepath.FromSlash(source.Path)))
		if err != nil {
			t.Fatalf("frozen source: %v", err)
		}

		if int64(len(content)) != *source.Bytes {
			t.Errorf("%s: %d bytes frozen, %d recorded", source.Path, len(content), *source.Bytes)
		}
		if digest := sha256Hex(content); digest != *source.SHA256 {
			t.Errorf("%s: frozen %s, recorded %s", source.Path, digest, *source.SHA256)
		}

		contextLines = append(contextLines, source.Path+"\t"+*source.SHA256+"\n")
	}

	sort.Strings(contextLines)

	if digest := sha256Hex([]byte(strings.Join(contextLines, ""))); digest != plan.Collection.Context.SHA256 {
		t.Errorf("context hash %s, recorded %s", digest, plan.Collection.Context.SHA256)
	}
}

func TestSampleRuntimeIsTheOneTheRunRecorded(t *testing.T) {
	plan := readSamplePlan(t)

	markdown, err := os.ReadFile(sampleRuntime)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := connection.Parse(string(markdown)); err != nil {
		t.Fatalf("the sample runtime.md does not parse: %v", err)
	}

	digest := sha256Hex(markdown)
	if digest != plan.Runtime.SHA256 {
		t.Errorf("runtime.md hashes to %s, the plan recorded %s", digest, plan.Runtime.SHA256)
	}

	lap := readSampleLapFence(t)
	if lap.Runtime.MarkdownSHA256 != digest {
		t.Errorf("lap 0002 recorded %s for runtime.md, the file hashes to %s", lap.Runtime.MarkdownSHA256, digest)
	}
}

func readSamplePlan(t *testing.T) samplePlanFile {
	t.Helper()

	data, err := os.ReadFile(samplePlan)
	if err != nil {
		t.Fatal(err)
	}

	var plan samplePlanFile
	if err := json.Unmarshal(data, &plan); err != nil {
		t.Fatalf("decode the sample plan: %v", err)
	}

	return plan
}

func readSampleLapFence(t *testing.T) sampleLapFence {
	t.Helper()

	data, err := os.ReadFile(sampleLap)
	if err != nil {
		t.Fatal(err)
	}

	const opening = "```langdrift-lap\n"

	_, after, found := strings.Cut(string(data), opening)
	if !found {
		t.Fatal("the lap has no langdrift-lap fence")
	}

	body, _, closed := strings.Cut(after, "\n```")
	if !closed {
		t.Fatal("the lap's fence is not closed")
	}

	var fence sampleLapFence
	if err := json.Unmarshal([]byte(body), &fence); err != nil {
		t.Fatalf("decode the lap fence: %v", err)
	}

	return fence
}

func sha256Hex(content []byte) string {
	digest := sha256.Sum256(content)

	return hex.EncodeToString(digest[:])
}
