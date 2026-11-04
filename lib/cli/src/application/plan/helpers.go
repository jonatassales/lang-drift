package plan

import (
	"fmt"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/jonatassales/tao/lib/cli/src/application"
	"github.com/jonatassales/tao/lib/cli/src/domain/collection"
	"github.com/jonatassales/tao/lib/cli/src/domain/finops"
	"github.com/jonatassales/tao/lib/cli/src/domain/harness"
	"github.com/jonatassales/tao/lib/cli/src/domain/microruntime"
	"github.com/jonatassales/tao/lib/cli/src/domain/project"
	"github.com/jonatassales/tao/lib/cli/src/domain/run"
)

// The results a check may have.
const (
	passed = "passed"
	failed = "failed"
)

// noManifest names a component path without any manifest the rules know.
const noManifest = "no manifest"

var criterionID = regexp.MustCompile(`\bFIN-\d{2}\b`)

// resolved is one component, as named and as its manifests and its
// Dockerfile describe it.
type resolved struct {
	input        ComponentInput
	manifestName string
	technologies []run.Technology
	kind         string
	proposed     harness.Resources

	// dockerfile is the Dockerfile's path from the project root.
	dockerfile      string
	dockerfileFound bool
	// probe is the readiness, when it is known; readinessErr says why not.
	probe        *microruntime.Probe
	readinessErr error
	envErr       error
}

// ticketPathOf makes the ticket's path project-relative and checks it stays inside.
func ticketPathOf(owner project.Project, ticket string) (string, error) {
	relative := ticket
	if filepath.IsAbs(ticket) {
		inside, err := filepath.Rel(owner.Root, ticket)
		if err != nil || collection.Escapes(filepath.ToSlash(inside)) {
			return "", application.NotFound("the ticket %s is outside the project", ticket)
		}
		relative = inside
	}

	relative = filepath.ToSlash(filepath.Clean(relative))
	if collection.Escapes(relative) {
		return "", application.NotFound("the ticket %s is outside the project", ticket)
	}

	return relative, nil
}

func runtimeIDOf(override, slug string) (microruntime.Identifier, error) {
	text := override
	if text == "" {
		derived, err := collection.RuntimeID(slug)
		if err != nil {
			return microruntime.Identifier{}, invalid("%v; name one with --runtime", err)
		}
		text = derived
	}

	identifier, err := microruntime.ParseIdentifier(text)
	if err != nil {
		return microruntime.Identifier{}, invalid("runtime id %q: %v", text, err)
	}

	return identifier, nil
}

// crossChecks are the checks of the draft beyond the schema.
func crossChecks(components []resolved, bounded, ceiling harness.Resources) []run.Check {
	ports := make([]string, 0, len(components))
	var unnamed []string
	seen := map[int]bool{}
	for _, component := range components {
		if component.input.Port == 0 {
			unnamed = append(unnamed, component.input.Path)

			continue
		}

		ports = append(ports, strconv.Itoa(component.input.Port))
		seen[component.input.Port] = true
	}
	portsCheck := run.Check{ID: "ports-distinct", Result: passed, Detail: text(strings.Join(ports, " · "))}
	if len(seen) != len(ports) {
		portsCheck.Result = failed
		portsCheck.Detail = text("ports repeat: " + strings.Join(ports, " · "))
	}
	if len(unnamed) > 0 {
		portsCheck.Result = failed
		portsCheck.Detail = text("no port named or exposed: " + strings.Join(unnamed, " · "))
	}

	ceilingCheck := run.Check{ID: "within-ceiling", Result: passed, Detail: text(fmt.Sprintf("%sm of %sm · %s MiB of %s MiB", finops.Grouped(bounded.CPUMillis), finops.Grouped(ceiling.CPUMillis), finops.Grouped(bounded.MemoryMiB), finops.Grouped(ceiling.MemoryMiB)))}
	if bounded.CPUMillis > ceiling.CPUMillis || bounded.MemoryMiB > ceiling.MemoryMiB {
		ceilingCheck.Result = failed
	}

	return []run.Check{portsCheck, dockerfilesCheck(components), readinessCheck(components), envFileCheck(components), ceilingCheck}
}

func distinctPorts(components []resolved) int {
	ports := map[int]bool{}
	for _, component := range components {
		ports[component.input.Port] = true
	}

	return len(ports)
}

func collisionsOf(checks []run.Check) string {
	for _, check := range checks {
		if check.ID == "ports-distinct" && check.Result == failed {
			return "ports collide"
		}
	}

	return "no collisions"
}

func failingChecks(checks []run.Check) []string {
	var failing []string
	for _, check := range checks {
		if check.Result == failed {
			failing = append(failing, check.ID)
		}
	}

	return failing
}

// totalsOf sums the proposed shares against the envelope; headroom never
// goes below zero.
func totalsOf(components []resolved, bounded harness.Resources) run.Totals {
	var proposed harness.Resources
	for _, component := range components {
		proposed.CPUMillis += component.proposed.CPUMillis
		proposed.MemoryMiB += component.proposed.MemoryMiB
	}

	return run.Totals{
		Proposed: proposed,
		Bounded:  bounded,
		Headroom: harness.Resources{CPUMillis: max(0, bounded.CPUMillis-proposed.CPUMillis), MemoryMiB: max(0, bounded.MemoryMiB-proposed.MemoryMiB)},
	}
}

func isGuardrail(source collection.Source) bool {
	return source.State == collection.Read && (strings.HasSuffix(source.Path, ".gr.md") || strings.Contains(source.Path, "guardrail"))
}

// criteria are the criterion ids a guardrail names, in order, once each.
func criteria(content []byte) []string {
	var found []string
	seen := map[string]bool{}
	for _, id := range criterionID.FindAllString(string(content), -1) {
		if seen[id] {
			continue
		}
		seen[id] = true
		found = append(found, id)
	}

	return found
}

func linkNames(collected collection.Result) []string {
	var names []string
	for _, source := range collected.Sources {
		if source.Depth == 1 {
			names = append(names, path.Base(source.Path))
		}
	}

	return names
}

func componentName(componentPath string) string {
	if componentPath == "." {
		return "root"
	}

	return path.Base(componentPath)
}

func runtimePath(runtime microruntime.Identifier) string {
	return microruntime.ConfigurationDirectory(runtime).String() + "/runtime.md"
}

func distinct(items []string) []string {
	seen := map[string]bool{}
	var unique []string
	for _, item := range items {
		if seen[item] {
			continue
		}
		seen[item] = true
		unique = append(unique, item)
	}

	return unique
}

func text(value string) *string { return &value }

func optional(value string) *string {
	if value == "" {
		return nil
	}

	return &value
}

func number(value int) *int { return &value }

func size(value int64) *int64 { return &value }

func orDefault(value, fallback string) string {
	if value == "" {
		return fallback
	}

	return value
}

func nonNil[Item any](items []Item) []Item {
	if items == nil {
		return []Item{}
	}

	return items
}
