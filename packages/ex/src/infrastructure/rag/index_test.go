package rag_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jonatassales/tao/packages/ex/src/domain/project"
	domain "github.com/jonatassales/tao/packages/ex/src/domain/rag"
	"github.com/jonatassales/tao/packages/ex/src/infrastructure/rag"
)

const memory = "# Run console-local: the console preview ran\n\nNotes.\n\n```langdrift-run\n{\n  \"run\": \"console-local\",\n  \"ticket\": \"0070-console-microruntime\",\n  \"kind\": \"ticket\",\n  \"outcome\": \"completed\",\n  \"laps\": 2,\n  \"componentPaths\": [\".\"],\n  \"stacks\": [\"nextjs\", \"go\"],\n  \"measuredResources\": true,\n  \"measuredTokens\": true,\n  \"observed\": { \"cpuMillisPeak\": 410, \"memoryMiBPeak\": 690, \"tokensThousands\": { \"input\": 3700, \"output\": 240 } },\n  \"facts\": [\"prewarm-build-over-120s\"]\n}\n```\n"

func TestIndexRebuildsFromTheMemoryFences(t *testing.T) {
	t.Parallel()

	root, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	owner := project.Project{Root: root}

	folder := filepath.Join(root, ".langdrift", "microruntimes", "runtimes", "console-local")
	if err := os.MkdirAll(folder, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(folder, "run.md"), []byte(memory), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, ".langdrift", "microruntimes", "runtimes", "drafted"), 0o755); err != nil {
		t.Fatal(err)
	}

	episodes, err := (rag.Index{}).Episodes(owner, domain.Local{})
	if err != nil {
		t.Fatal(err)
	}

	if len(episodes) != 1 {
		t.Fatalf("%d episodes, want the one runtime with a memory", len(episodes))
	}

	episode := episodes[0]
	if episode.Runtime != "console-local" || episode.Ticket != "0070-console-microruntime" || !episode.MeasuredResources || episode.CPUMillisPeak != 410 || episode.MemoryMiBPeak != 690 || len(episode.Vector) != domain.Dimensions || episode.Title != "Run console-local: the console preview ran" {
		t.Fatalf("episode = %+v", episode)
	}

	indexPath := filepath.Join(root, ".langdrift", "agent", "state", "rag", "index.json")
	if info, err := os.Stat(indexPath); err != nil || info.Mode().Perm() != 0o600 {
		t.Fatalf("index file = %v, %v; want a private file", info, err)
	}

	// A second read serves the index; another embedder id rebuilds it.
	again, err := (rag.Index{}).Episodes(owner, domain.Local{})
	if err != nil || len(again) != 1 {
		t.Fatalf("second read = %d episodes, %v", len(again), err)
	}
}
