package manifest_test

import (
	"reflect"
	"testing"

	"github.com/jonatassales/tao/lib/cli/src/domain/harness"
	"github.com/jonatassales/tao/lib/cli/src/domain/manifest"
)

func TestClassifyReadsTheManifestsOfTheWebsiteComponents(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		manifests map[string][]byte
		want      manifest.Inspection
	}{
		{"a Go API", map[string][]byte{"go.mod": []byte("module example.com/signup-api\n\ngo 1.26\n"), "Dockerfile": []byte("FROM golang\nEXPOSE 8080\n")}, manifest.Inspection{Manifest: "go.mod", Technologies: []manifest.Technology{manifest.Go}, Kind: manifest.API, Port: 8080}},
		{"a Node.js API", map[string][]byte{"package.json": []byte(`{"scripts":{"start":"node server.js --port 8081"},"dependencies":{"express":"5"}}`)}, manifest.Inspection{Manifest: "package.json", Technologies: []manifest.Technology{manifest.NodeJS}, Kind: manifest.API, Port: 8081}},
		{"a NATS broker", map[string][]byte{"nats.conf": []byte("port: 4222\n")}, manifest.Inspection{Manifest: "nats.conf", Technologies: []manifest.Technology{manifest.NATS}, Kind: manifest.Broker, Port: 4222}},
		{"a Redis cache", map[string][]byte{"redis.conf": []byte("port 6379\n")}, manifest.Inspection{Manifest: "redis.conf", Technologies: []manifest.Technology{manifest.Redis}, Kind: manifest.Cache, Port: 6379}},
		{"a Next.js app", map[string][]byte{"package.json": []byte(`{"scripts":{"dev":"next dev -p 3000"},"dependencies":{"next":"16"}}`)}, manifest.Inspection{Manifest: "package.json", Technologies: []manifest.Technology{manifest.NodeJS, manifest.NextJS}, Kind: manifest.Web, Port: 3000}},
		{"nothing", map[string][]byte{}, manifest.Inspection{}},
	}

	for _, tt := range tests {
		if got := manifest.Classify(tt.manifests); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%s: Classify() = %+v, want %+v", tt.name, got, tt.want)
		}
	}
}

func TestProposedSharesFollowTheKind(t *testing.T) {
	t.Parallel()

	if got := manifest.Proposed(manifest.Web, nil); got != (harness.Resources{CPUMillis: 500, MemoryMiB: 1024}) {
		t.Errorf("web = %+v", got)
	}
	if got := manifest.Proposed(manifest.API, []manifest.Technology{manifest.NodeJS}); got != (harness.Resources{CPUMillis: 300, MemoryMiB: 512}) {
		t.Errorf("node api = %+v", got)
	}
	if got := manifest.Proposed(manifest.API, []manifest.Technology{manifest.Go}); got != (harness.Resources{CPUMillis: 250, MemoryMiB: 256}) {
		t.Errorf("go api = %+v", got)
	}
	if got := manifest.Proposed(manifest.Broker, nil); got != (harness.Resources{CPUMillis: 200, MemoryMiB: 256}) {
		t.Errorf("broker = %+v", got)
	}

	if manifest.Labels([]manifest.Technology{manifest.NextJS, manifest.Go}) != "Next.js+Go" || manifest.Labels(nil) != "Unknown" {
		t.Error("Labels() misjoins")
	}
}
