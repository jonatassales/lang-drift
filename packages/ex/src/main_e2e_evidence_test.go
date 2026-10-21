//go:build e2e

package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"
)

// excerptLimit bounds each recorded output.
const excerptLimit = 4096

// evidence is the sanitized record of one run. It never holds a kubeconfig,
// Terraform state, a plan, a variable file, or an ownership token.
type evidence struct {
	SchemaVersion int               `json:"schemaVersion"`
	RunID         string            `json:"runId"`
	StartedAt     time.Time         `json:"startedAt"`
	FinishedAt    time.Time         `json:"finishedAt"`
	Result        string            `json:"result"`
	Platform      map[string]string `json:"platform"`
	Tools         map[string]string `json:"tools"`
	Steps         []stepRecord      `json:"steps"`
	Cleanup       []stepRecord      `json:"cleanup"`
	LeftBehind    []leftBehind      `json:"leftBehind"`
}

// stepRecord is one command. ToolInvocations counts the native tools the
// CLI ran directly, not the tools those tools ran themselves.
type stepRecord struct {
	Name            string         `json:"name"`
	Command         []string       `json:"command"`
	ExitCode        int            `json:"exitCode"`
	DurationMS      int64          `json:"durationMs"`
	ToolInvocations map[string]int `json:"toolInvocations,omitempty"`
	Stdout          string         `json:"stdout,omitempty"`
	Stderr          string         `json:"stderr,omitempty"`
}

type leftBehind struct {
	Resource string `json:"resource"`
	Reason   string `json:"reason"`
}

func newEvidence(nonce string) evidence {
	return evidence{
		SchemaVersion: 1,
		RunID:         "ex-e2e-" + nonce,
		StartedAt:     time.Now().UTC(),
		Platform:      map[string]string{"goos": runtime.GOOS, "goarch": runtime.GOARCH, "go": runtime.Version()},
		Tools:         map[string]string{},
		Steps:         []stepRecord{},
		Cleanup:       []stepRecord{},
		LeftBehind:    []leftBehind{},
	}
}

// recordTools records the platform and native tool versions.
func (h *harness) recordTools() {
	ctx, cancel := context.WithTimeout(h.ctx, probeDeadline)
	defer cancel()

	probes := map[string][]string{
		"docker":    {"version", "--format", "{{.Server.Version}}"},
		"kind":      {"version"},
		"kubectl":   {"version", "--client", "-o", "json"},
		"terraform": {"version", "-json"},
	}

	for tool, arguments := range probes {
		result := h.native(ctx, tool, nil, arguments...)
		h.report.Tools[tool] = excerpt(sanitize(strings.TrimSpace(result.stdout), h.workdir), 512)
	}

	if host := h.native(ctx, "docker", nil, "info", "--format", "{{.OSType}}/{{.Architecture}} {{.KernelVersion}}"); host.exit == 0 {
		h.report.Platform["dockerHost"] = strings.TrimSpace(host.stdout)
	}
}

// record appends one sanitized step and the native tools it invoked.
func (h *harness) record(target *[]stepRecord, step string, arguments []string, result outcome, duration time.Duration, invocations map[string]int) {
	command := make([]string, len(arguments))
	for index, argument := range arguments {
		command[index] = sanitize(argument, h.workdir)
	}

	*target = append(*target, stepRecord{
		Name:            step,
		Command:         command,
		ExitCode:        result.exit,
		DurationMS:      duration.Milliseconds(),
		ToolInvocations: invocations,
		Stdout:          excerpt(sanitize(result.stdout, h.workdir), excerptLimit),
		Stderr:          excerpt(sanitize(result.stderr, h.workdir), excerptLimit),
	})
}

// countInvocations counts the native tool runs one CLI invocation logged.
func countInvocations(log string) map[string]int {
	data, err := os.ReadFile(log)
	if err != nil {
		return nil
	}

	counts := map[string]int{}
	for _, tool := range strings.Fields(string(data)) {
		counts[tool]++
	}

	return counts
}

// finish writes the evidence and fails the run if anything was left behind,
// keeping the workspace in that case so the operator can inspect it.
func (h *harness) finish() {
	h.report.FinishedAt = time.Now().UTC()

	if len(h.report.LeftBehind) > 0 {
		h.t.Errorf("resources left behind; workspace kept at %s", h.workdir)

		for _, item := range h.report.LeftBehind {
			h.t.Errorf("  %s: %s", item.Resource, item.Reason)
		}
	}

	h.report.Result = "passed"
	if h.t.Failed() {
		h.report.Result = "failed"
	}

	h.writeEvidence()

	if len(h.report.LeftBehind) == 0 {
		if err := os.RemoveAll(h.workdir); err != nil {
			h.t.Errorf("remove workspace: %v", err)
		}
	}
}

// writeEvidence writes report.json to .audits/evidence/runs/<run ID> at
// the repository root, the ignored destination of local run reports.
func (h *harness) writeEvidence() {
	root, err := repositoryRoot()
	if err != nil {
		h.t.Errorf("locate the evidence directory: %v", err)
		return
	}

	data, err := json.MarshalIndent(h.report, "", "  ")
	if err != nil {
		h.t.Errorf("encode evidence: %v", err)
		return
	}

	directory := filepath.Join(root, ".audits", "evidence", "runs", h.report.RunID)
	if err := os.MkdirAll(directory, 0o755); err != nil {
		h.t.Errorf("create the evidence directory: %v", err)
		return
	}

	path := filepath.Join(directory, "report.json")
	if err := os.WriteFile(path, append(data, '\n'), 0o644); err != nil {
		h.t.Errorf("write evidence: %v", err)
		return
	}

	h.t.Logf("evidence written to %s", path)
}

// repositoryRoot is the nearest parent of the test's directory that holds
// pnpm-workspace.yaml.
func repositoryRoot() (string, error) {
	directory, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		_, err := os.Stat(filepath.Join(directory, "pnpm-workspace.yaml"))
		switch {
		case err == nil:
			return directory, nil
		case !errors.Is(err, fs.ErrNotExist):
			return "", err
		}

		parent := filepath.Dir(directory)
		if parent == directory {
			return "", errors.New("no repository root holds pnpm-workspace.yaml above the test")
		}

		directory = parent
	}
}

// Values that must never reach the evidence or the test log.
var (
	pemBlock        = regexp.MustCompile(`(?s)-----BEGIN [A-Z ]+-----.*?-----END [A-Z ]+-----`)
	credentialValue = regexp.MustCompile(
		`(?i)((?:client-key-data|client-certificate-data|certificate-authority-data|token|password|secret)["']?\s*[:=]\s*["']?)[^\s"',}]+`,
	)
)

// sanitize replaces the run's private paths and redacts credential-like values.
func sanitize(text, workdir string) string {
	if workdir != "" {
		text = strings.ReplaceAll(text, workdir, "<workdir>")
	}

	if home, err := os.UserHomeDir(); err == nil && len(home) > 1 {
		text = strings.ReplaceAll(text, home, "<home>")
	}

	text = pemBlock.ReplaceAllString(text, "<redacted>")

	return credentialValue.ReplaceAllString(text, "${1}<redacted>")
}

// excerpt bounds a recorded output.
func excerpt(text string, limit int) string {
	if len(text) <= limit {
		return text
	}

	return text[:limit] + fmt.Sprintf("\n[truncated %d bytes]", len(text)-limit)
}
