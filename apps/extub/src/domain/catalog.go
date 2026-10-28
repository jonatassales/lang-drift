// Package domain describes synthetic console fixtures, never host resources.
package domain

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"slices"
	"strings"
	"time"
)

var ErrNotFound = errors.New("not_found")
var ErrInvalid = errors.New("invalid_input")
var ErrConflict = errors.New("executions_not_stopped")
var ErrLimit = errors.New("seed_limit")
var ErrLoadTestBusy = errors.New("load_test_busy")

type Source struct {
	Kind string `json:"kind"`
	Path string `json:"path"`
}
type Component struct {
	Name          string          `json:"name"`
	Path          string          `json:"path"`
	Kind          string          `json:"kind"`
	Stacks        []string        `json:"stacks,omitempty"`
	Image         string          `json:"image,omitempty"`
	Port          int             `json:"port,omitempty"`
	ReadinessPath string          `json:"readinessPath,omitempty"`
	Build         json.RawMessage `json:"build,omitempty"`
}
type Configuration struct {
	CPU             float64  `json:"cpu"`
	MemoryMiB       int      `json:"memoryMiB"`
	Model           string   `json:"model"`
	Effort          string   `json:"effort"`
	LLMBudgetUSD    *float64 `json:"llmBudgetUsd,omitempty"`
	NativeRuntimeID string   `json:"nativeRuntimeId,omitempty"`
}
type Project struct {
	ID             string          `json:"id"`
	Name           string          `json:"name"`
	Source         Source          `json:"source"`
	Ref            string          `json:"ref"`
	Components     []Component     `json:"components"`
	ComponentLinks []ComponentLink `json:"componentLinks,omitempty"`
	CreatedAt      string          `json:"createdAt"`
}
type Runtime struct {
	ID                string          `json:"id"`
	NativeRuntimeID   string          `json:"nativeRuntimeId"`
	ProjectID         string          `json:"projectId"`
	Name              string          `json:"name"`
	Components        []Component     `json:"components"`
	Configuration     Configuration   `json:"configuration"`
	Revision          int             `json:"revision"`
	Markdown          string          `json:"markdown"`
	DesiredState      string          `json:"desiredState"`
	IntentVersion     int             `json:"intentVersion"`
	CredentialRevoked bool            `json:"credentialRevoked"`
	CreatedAt         string          `json:"createdAt"`
	UpdatedAt         string          `json:"updatedAt"`
	StartedAt         string          `json:"startedAt,omitempty"`
	LoadTest          *LoadTestResult `json:"loadTest,omitempty"`
}
type Secret struct {
	ID        string  `json:"id"`
	ProjectID string  `json:"projectId"`
	Name      string  `json:"name"`
	CreatedAt string  `json:"createdAt"`
	RevokedAt *string `json:"revokedAt"`
}
type Catalog struct {
	ActivitySequence int64             `json:"activitySequence"`
	ExecutionRecords []ExecutionRecord `json:"executionRecords"`
	EventRecords     []EventRecord     `json:"eventRecords"`
	Source           string            `json:"source"`
	Projects         []Project         `json:"projects"`
	Runtimes         []Runtime         `json:"runtimes"`
	Secrets          []Secret          `json:"secrets"`
	StartedAt        string            `json:"startedAt"`
}

// LoadTestResult measures a bounded HTTP burst against extub itself, never a workload.
type LoadTestResult struct {
	Target      string  `json:"target"`
	Requests    int     `json:"requests"`
	Concurrency int     `json:"concurrency"`
	Succeeded   int     `json:"succeeded"`
	Failed      int     `json:"failed"`
	DurationMs  float64 `json:"durationMs"`
	P50Ms       float64 `json:"p50Ms"`
	P95Ms       float64 `json:"p95Ms"`
	MaxMs       float64 `json:"maxMs"`
	StartedAt   string  `json:"startedAt"`
	CompletedAt string  `json:"completedAt"`
}

// Clone prevents callers or mutation callbacks from retaining store-owned slices.
func (catalog Catalog) Clone() Catalog {
	data, _ := json.Marshal(catalog)
	var copy Catalog
	_ = json.Unmarshal(data, &copy)
	return copy
}
func (catalog Catalog) Project(id string) (Project, error) {
	for _, item := range catalog.Projects {
		if item.ID == id {
			return item, nil
		}
	}
	return Project{}, ErrNotFound
}
func (catalog Catalog) RuntimeIndex(id string) int {
	return slices.IndexFunc(catalog.Runtimes, func(item Runtime) bool { return item.ID == id })
}

func ValidateRuntime(runtime Runtime) error {
	data, err := json.Marshal(runtime)
	config := runtime.Configuration
	if err != nil || len(data) > 8192 || strings.TrimSpace(runtime.Name) == "" || len(runtime.Name) > 100 || len(runtime.Components) == 0 || len(runtime.Components) > 16 || config.CPU <= 0 || config.CPU > 64 || math.IsNaN(config.CPU) || config.MemoryMiB <= 0 || config.MemoryMiB > 262144 || config.Model == "" || len(config.Model) > 200 || !slices.Contains([]string{"low", "medium", "high"}, config.Effort) {
		return ErrInvalid
	}
	if config.LLMBudgetUSD != nil && (*config.LLMBudgetUSD <= 0 || *config.LLMBudgetUSD > 1e6 || math.IsNaN(*config.LLMBudgetUSD)) {
		return ErrInvalid
	}
	for _, component := range runtime.Components {
		if component.Path == "" || component.Name == "" || component.Kind == "" || len(component.Stacks) > 4 {
			return ErrInvalid
		}
		for _, stack := range component.Stacks {
			if len(stack) == 0 || len(stack) > 32 {
				return ErrInvalid
			}
		}
	}
	return nil
}

type Environment struct {
	ID             string `json:"id"`
	RuntimeID      string `json:"runtimeId"`
	Name           string `json:"name"`
	LastSeenAt     string `json:"lastSeenAt"`
	Connection     string `json:"connection"`
	CleanupPending bool   `json:"cleanupPending"`
}
type Execution struct {
	ID                   string   `json:"id"`
	RuntimeID            string   `json:"runtimeId"`
	EnvironmentID        string   `json:"environmentId"`
	State                string   `json:"state"`
	Revision             int      `json:"revision"`
	IntentVersion        int      `json:"intentVersion"`
	UpdatedAt            string   `json:"updatedAt"`
	ObservedAt           *string  `json:"observedAt"`
	MeasuredFrom         *string  `json:"measuredFrom,omitempty"`
	CPUUsed              *float64 `json:"cpuUsed"`
	MemoryUsedMiB        *float64 `json:"memoryUsedMiB"`
	LLMCostUSD           *float64 `json:"llmCostUsd"`
	DowntimeStartedAt    *string  `json:"downtimeStartedAt"`
	DowntimeMilliseconds *int     `json:"downtimeMilliseconds"`
	DowntimeQuality      string   `json:"downtimeQuality"`
}
type ComponentLink struct {
	From string `json:"from"`
	To   string `json:"to"`
}
type Widget struct {
	Source         string            `json:"source"`
	Runtime        Runtime           `json:"runtime"`
	Environments   []Environment     `json:"environments"`
	Executions     []Execution       `json:"executions"`
	Tone           string            `json:"tone"`
	Partial        bool              `json:"partial"`
	StartedAt      map[string]string `json:"startedAt"`
	ComponentLinks []ComponentLink   `json:"componentLinks"`
}

// Widgets creates explicitly simulated observations; unknown values remain null.
func (catalog Catalog) Widgets(now time.Time, sequence uint64) []Widget {
	widgets := make([]Widget, 0, len(catalog.Runtimes))
	timestamp := now.UTC().Format(time.RFC3339Nano)
	for index, runtime := range catalog.Runtimes {
		active := runtime.DesiredState == "running"
		ratio := []float64{.25, .60, .80, .94}[index%4] + float64(sequence%10)*.003
		widget := Widget{Source: "stub", Runtime: runtime, Environments: []Environment{}, Executions: []Execution{}, Tone: "purple", Partial: runtime.Configuration.LLMBudgetUSD == nil, StartedAt: map[string]string{}, ComponentLinks: []ComponentLink{}}
		if active {
			widget.Tone = severity(ratio)
		}
		project, _ := catalog.Project(runtime.ProjectID)
		paths := make(map[string]bool, len(runtime.Components))
		for _, component := range runtime.Components {
			paths[component.Path] = true
		}
		// Only authored fixture dependencies whose endpoints are in this runtime.
		for _, link := range project.ComponentLinks {
			if paths[link.From] && paths[link.To] {
				widget.ComponentLinks = append(widget.ComponentLinks, link)
			}
		}
		environment := Environment{ID: runtime.ID + "-machine-0", RuntimeID: runtime.ID, Name: "Seed machine 1", LastSeenAt: timestamp, Connection: "connected"}
		widget.Environments = append(widget.Environments, environment)
		if !active {
			widgets = append(widgets, widget)
			continue
		}
		cpu, memory := runtime.Configuration.CPU*ratio, math.Round(float64(runtime.Configuration.MemoryMiB)*ratio)
		measured := now.Add(-5 * time.Second).UTC().Format(time.RFC3339Nano)
		zero := 0
		execution := Execution{
			ID: fmt.Sprintf("%s-execution-%d", runtime.ID, runtime.IntentVersion), RuntimeID: runtime.ID,
			EnvironmentID: environment.ID, State: "running", Revision: runtime.Revision,
			IntentVersion: runtime.IntentVersion, UpdatedAt: timestamp, ObservedAt: &timestamp,
			MeasuredFrom: &measured, CPUUsed: &cpu, MemoryUsedMiB: &memory,
			DowntimeMilliseconds: &zero, DowntimeQuality: "observed",
		}
		if runtime.Configuration.LLMBudgetUSD != nil {
			cost := *runtime.Configuration.LLMBudgetUSD * ratio
			execution.LLMCostUSD = &cost
		}
		started := runtime.StartedAt
		if started == "" {
			started = catalog.StartedAt
		}
		widget.StartedAt[execution.ID] = started
		widget.Executions = append(widget.Executions, execution)
		widgets = append(widgets, widget)
	}
	return widgets
}
func severity(ratio float64) string {
	if ratio < .5 {
		return "green"
	}
	if ratio < .75 {
		return "yellow"
	}
	if ratio < .9 {
		return "orange"
	}
	return "red"
}
