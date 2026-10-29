package plan

import (
	"fmt"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/jonatassales/tao/packages/contracts/configuration"
	"github.com/jonatassales/tao/packages/ex/src/application"
	"github.com/jonatassales/tao/packages/ex/src/domain/collection"
	"github.com/jonatassales/tao/packages/ex/src/domain/finops"
	"github.com/jonatassales/tao/packages/ex/src/domain/harness"
	"github.com/jonatassales/tao/packages/ex/src/domain/microruntime"
	"github.com/jonatassales/tao/packages/ex/src/domain/project"
	"github.com/jonatassales/tao/packages/ex/src/domain/run"
)

// The results a check may have.
const (
	passed = "passed"
	failed = "failed"
)

// noManifest names a component path without any manifest the rules know.
const noManifest = "no manifest"

var criterionID = regexp.MustCompile(`\bFIN-\d{2}\b`)

// resolved is one component, as named and as its manifests describe it.
type resolved struct {
	input        ComponentInput
	manifestName string
	technologies []run.Technology
	kind         string
	proposed     harness.Resources
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
	seen := map[int]bool{}
	for _, component := range components {
		ports = append(ports, strconv.Itoa(component.input.Port))
		seen[component.input.Port] = true
	}
	portsCheck := run.Check{ID: "ports-distinct", Result: passed, Detail: text(strings.Join(ports, " · "))}
	if len(seen) != len(components) {
		portsCheck.Result = failed
		portsCheck.Detail = text("ports repeat: " + strings.Join(ports, " · "))
	}

	images := make([]string, 0, len(components))
	imagesCheck := run.Check{ID: "images-resolved", Result: passed}
	for _, component := range components {
		images = append(images, component.input.Image)
		if !configuration.IsImageReference(component.input.Image) {
			imagesCheck.Result = failed
		}
	}
	imagesCheck.Detail = text(strings.Join(images, " · "))

	recipes := make([]string, 0, len(components))
	recipeCheck := run.Check{ID: "recipe-allowed", Result: passed}
	for _, component := range components {
		if err := component.input.Build.Validate(); err != nil {
			recipeCheck.Result = failed
			recipes = append(recipes, err.Error())

			continue
		}
		recipes = append(recipes, recipeName(component.input.Build))
	}
	recipeCheck.Detail = text(strings.Join(recipes, " · "))

	ceilingCheck := run.Check{ID: "within-ceiling", Result: passed, Detail: text(fmt.Sprintf("%sm of %sm · %s MiB of %s MiB", finops.Grouped(bounded.CPUMillis), finops.Grouped(ceiling.CPUMillis), finops.Grouped(bounded.MemoryMiB), finops.Grouped(ceiling.MemoryMiB)))}
	if bounded.CPUMillis > ceiling.CPUMillis || bounded.MemoryMiB > ceiling.MemoryMiB {
		ceilingCheck.Result = failed
	}

	return []run.Check{portsCheck, imagesCheck, recipeCheck, ceilingCheck}
}

func recipeName(build configuration.Build) string {
	if strings.Contains(build.Executable, "/") {
		return "committed script " + build.Executable
	}

	return strings.TrimSpace(build.Executable + " " + strings.Join(build.Args, " "))
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

// completeness is the check that every component names what the fence
// needs: an image, a port, a readiness path and a build recipe. A missing
// one blocks the plan with the field named.
func completeness(components []resolved) run.Check {
	var missing []string
	for _, component := range components {
		if fields := missingFields(component.input); len(fields) > 0 {
			missing = append(missing, component.input.Path+": "+strings.Join(fields, ", "))
		}
	}

	if len(missing) == 0 {
		return run.Check{ID: "components-complete", Result: passed, Detail: text("image, port, readiness and build named for every component")}
	}

	return run.Check{ID: "components-complete", Result: failed, Detail: text("missing " + strings.Join(missing, " · "))}
}

func missingFields(input ComponentInput) []string {
	var fields []string
	if input.Image == "" {
		fields = append(fields, "image")
	}
	if input.Port == 0 {
		fields = append(fields, "port")
	}
	if input.ReadinessPath == "" {
		fields = append(fields, "readiness")
	}
	if input.Build.Executable == "" {
		fields = append(fields, "build")
	}

	return fields
}
