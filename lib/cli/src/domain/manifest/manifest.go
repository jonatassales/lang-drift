// Package manifest holds the deterministic rules that read a component from
// its manifests: which technology a file names, what kind of component that
// makes, which port it declares, and the share a kind is proposed to start
// from. Nothing here runs a build or asks a model.
package manifest

import (
	"encoding/json"
	"regexp"
	"strconv"
	"strings"

	"github.com/jonatassales/tao/lib/cli/src/domain/harness"
)

// Technology is one technology a manifest names.
type Technology struct {
	ID    string
	Label string
}

// The technologies the rules know.
var (
	Go      = Technology{ID: "go", Label: "Go"}
	NodeJS  = Technology{ID: "nodejs", Label: "Node.js"}
	NextJS  = Technology{ID: "nextjs", Label: "Next.js"}
	NATS    = Technology{ID: "nats", Label: "NATS"}
	Redis   = Technology{ID: "redis", Label: "Redis"}
	Python  = Technology{ID: "python", Label: "Python"}
	Docker  = Technology{ID: "docker", Label: "Docker"}
	Unknown = Technology{ID: "unknown", Label: "Unknown"}
)

// The kinds of component, by what the technology serves.
const (
	Web    = "web"
	API    = "api"
	Broker = "broker"
	Cache  = "cache"
	Worker = "worker"
)

// Names are the manifest files the rules read, in the order they are tried.
var Names = []string{"package.json", "go.mod", "nats.conf", "redis.conf", "pyproject.toml", "Dockerfile"}

// Inspection is what the manifests of one component path say.
type Inspection struct {
	// Manifest is the file that decided the technology, or empty.
	Manifest     string
	Technologies []Technology
	Kind         string
	// Port is the port a manifest declares, or zero.
	Port int
}

var (
	scriptPort  = regexp.MustCompile(`(?:-p|--port)[ =](\d{2,5})\b`)
	envPort     = regexp.MustCompile(`\bPORT=(\d{2,5})\b`)
	exposePort  = regexp.MustCompile(`(?m)^EXPOSE +(\d{2,5})\b`)
	natsPort    = regexp.MustCompile(`(?m)^\s*port\s*[:=]\s*(\d{2,5})\b`)
	redisPort   = regexp.MustCompile(`(?m)^port +(\d{2,5})\b`)
	kindRanking = map[string]int{Web: 4, Broker: 3, Cache: 2, API: 1, Worker: 0}
)

// Classify reads the manifests found at a component path, by name.
func Classify(manifests map[string][]byte) Inspection {
	var inspection Inspection

	if content, found := manifests["package.json"]; found {
		inspection.add("package.json", NodeJS, API)
		if dependsOn(content, "next") {
			inspection.add("package.json", NextJS, Web)
		}
		inspection.port(scriptPort, content)
		inspection.port(envPort, content)
	}

	if content, found := manifests["go.mod"]; found {
		inspection.add("go.mod", Go, API)
		inspection.port(envPort, content)
	}

	if content, found := manifests["nats.conf"]; found {
		inspection.add("nats.conf", NATS, Broker)
		inspection.port(natsPort, content)
		if inspection.Port == 0 {
			inspection.Port = 4222
		}
	}

	if content, found := manifests["redis.conf"]; found {
		inspection.add("redis.conf", Redis, Cache)
		inspection.port(redisPort, content)
		if inspection.Port == 0 {
			inspection.Port = 6379
		}
	}

	if _, found := manifests["pyproject.toml"]; found {
		inspection.add("pyproject.toml", Python, API)
	}

	if content, found := manifests["Dockerfile"]; found {
		inspection.port(exposePort, content)
		if len(inspection.Technologies) == 0 {
			inspection.add("Dockerfile", Docker, API)
		}
	}

	return inspection
}

func (inspection *Inspection) add(manifest string, technology Technology, kind string) {
	if inspection.Manifest == "" {
		inspection.Manifest = manifest
	}

	inspection.Technologies = append(inspection.Technologies, technology)

	if kindRanking[kind] > kindRanking[inspection.Kind] || inspection.Kind == "" {
		inspection.Kind = kind
	}
}

func (inspection *Inspection) port(pattern *regexp.Regexp, content []byte) {
	if inspection.Port != 0 {
		return
	}

	match := pattern.FindSubmatch(content)
	if match == nil {
		return
	}

	port, err := strconv.Atoi(string(match[1]))
	if err == nil && port > 0 && port < 65536 {
		inspection.Port = port
	}
}

// dependsOn reports whether a package.json lists a dependency.
func dependsOn(content []byte, name string) bool {
	var manifest struct {
		Dependencies    map[string]string `json:"dependencies"`
		DevDependencies map[string]string `json:"devDependencies"`
	}
	if err := json.Unmarshal(content, &manifest); err != nil {
		return strings.Contains(string(content), `"`+name+`"`)
	}

	_, direct := manifest.Dependencies[name]
	_, development := manifest.DevDependencies[name]

	return direct || development
}

// Labels joins the technologies' labels, as the lines show them.
func Labels(technologies []Technology) string {
	if len(technologies) == 0 {
		return Unknown.Label
	}

	labels := make([]string, 0, len(technologies))
	for _, technology := range technologies {
		labels = append(labels, technology.Label)
	}

	return strings.Join(labels, "+")
}

// Proposed is the share a component of a kind and technology starts from,
// before any measured run informs it.
func Proposed(kind string, technologies []Technology) harness.Resources {
	switch kind {
	case Web:
		return harness.Resources{CPUMillis: 500, MemoryMiB: 1024}
	case Broker, Cache:
		return harness.Resources{CPUMillis: 200, MemoryMiB: 256}
	case API:
		for _, technology := range technologies {
			if technology == NodeJS {
				return harness.Resources{CPUMillis: 300, MemoryMiB: 512}
			}
		}

		return harness.Resources{CPUMillis: 250, MemoryMiB: 256}
	}

	return harness.Resources{CPUMillis: 250, MemoryMiB: 256}
}
