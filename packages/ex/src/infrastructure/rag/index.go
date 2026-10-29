// Package rag keeps the project's episode index under the private state,
// rebuilt from every runtime's memory fence, and answers retrieval with it.
package rag

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jonatassales/tao/packages/ex/src/domain/microruntime"
	"github.com/jonatassales/tao/packages/ex/src/domain/project"
	"github.com/jonatassales/tao/packages/ex/src/domain/rag"
	"github.com/jonatassales/tao/packages/ex/src/domain/run"
	"github.com/jonatassales/tao/packages/ex/src/infrastructure/files"
)

// The index lives under the private state; the runtimes under microruntimes.
const (
	indexPath         = ".langdrift/agent/state/rag/index.json"
	runtimesDirectory = ".langdrift/microruntimes/runtimes"
	indexVersion      = 1
)

// File is the index as written.
type File struct {
	Version  int           `json:"version"`
	Embedder EmbedderInfo  `json:"embedder"`
	Episodes []rag.Episode `json:"episodes"`
}

// EmbedderInfo names the embedder the vectors came from.
type EmbedderInfo struct {
	ID         string `json:"id"`
	Dimensions int    `json:"dimensions"`
}

// Index is the file index.
type Index struct{}

// Episodes loads the index, rebuilding it when it is missing or was built by
// another embedder.
func (index Index) Episodes(owner project.Project, embedder rag.Embedder) ([]rag.Episode, error) {
	path, err := files.Within(owner.Root, indexPath)
	if err != nil {
		return nil, err
	}

	var file File

	err = files.ReadJSON(path, &file)
	if err == nil && file.Version == indexVersion && file.Embedder.ID == embedder.ID() {
		return file.Episodes, nil
	}
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return nil, err
	}

	return index.Rebuild(owner, embedder)
}

// Rebuild reads every runtime's memory fence and writes the index anew.
func (Index) Rebuild(owner project.Project, embedder rag.Embedder) ([]rag.Episode, error) {
	episodes, err := episodesOf(owner, embedder)
	if err != nil {
		return nil, err
	}

	path, err := files.Within(owner.Root, indexPath)
	if err != nil {
		return nil, err
	}

	file := File{Version: indexVersion, Embedder: EmbedderInfo{ID: embedder.ID(), Dimensions: rag.Dimensions}, Episodes: episodes}

	return episodes, files.WriteJSON(path, file)
}

func episodesOf(owner project.Project, embedder rag.Embedder) ([]rag.Episode, error) {
	directory, err := files.Within(owner.Root, runtimesDirectory)
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(directory)
	if errors.Is(err, fs.ErrNotExist) {
		return []rag.Episode{}, nil
	}
	if err != nil {
		return nil, err
	}

	episodes := []rag.Episode{}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		runtime, err := microruntime.ParseIdentifier(entry.Name())
		if err != nil {
			continue
		}

		episode, found, err := episodeOf(owner, runtime)
		if err != nil {
			return nil, err
		}
		if !found {
			continue
		}

		episode.Vector = embedder.Embed(episode.Text())
		episodes = append(episodes, episode)
	}

	sort.SliceStable(episodes, func(left, right int) bool { return episodes[left].Runtime < episodes[right].Runtime })

	return episodes, nil
}

// episodeOf reads one runtime's memory fence, and its title from the plan.
func episodeOf(owner project.Project, runtime microruntime.Identifier) (rag.Episode, bool, error) {
	folder := filepath.FromSlash(microruntime.ConfigurationDirectory(runtime).String())

	memoryPath, err := files.Within(owner.Root, filepath.Join(folder, "run.md"))
	if err != nil {
		return rag.Episode{}, false, err
	}

	memory, err := files.ReadBounded(memoryPath, 1<<20)
	if errors.Is(err, fs.ErrNotExist) {
		return rag.Episode{}, false, nil
	}
	if err != nil {
		return rag.Episode{}, false, err
	}

	if !strings.Contains(string(memory), "```langdrift-run\n") {
		return rag.Episode{}, false, nil
	}

	var fence run.RunFence
	if err := decodeFence(string(memory), "langdrift-run", &fence); err != nil {
		return rag.Episode{}, false, fmt.Errorf("the memory of %s has an unreadable langdrift-run fence: %w", runtime, err)
	}

	episode := rag.Episode{
		Runtime:           fence.Run,
		Stacks:            nonNil(fence.Stacks),
		ComponentPaths:    nonNil(fence.ComponentPaths),
		MeasuredResources: fence.MeasuredResources,
		MeasuredTokens:    fence.MeasuredTokens,
		Facts:             nonNil(fence.Facts),
		Title:             heading(string(memory)),
	}
	if fence.Ticket != nil {
		episode.Ticket = *fence.Ticket
	}
	if fence.Outcome != nil {
		episode.Outcome = *fence.Outcome
	}
	if fence.Observed != nil {
		if fence.Observed.CPUMillisPeak != nil {
			episode.CPUMillisPeak = *fence.Observed.CPUMillisPeak
		}
		if fence.Observed.MemoryMiBPeak != nil {
			episode.MemoryMiBPeak = *fence.Observed.MemoryMiBPeak
		}
	}

	runtimePath, err := files.Within(owner.Root, filepath.Join(folder, "runtime.md"))
	if err != nil {
		return rag.Episode{}, false, err
	}
	if markdown, err := files.ReadBounded(runtimePath, 1<<20); err == nil {
		var plan run.Plan
		if err := decodeFence(string(markdown), "langdrift-plan", &plan); err == nil && plan.Ticket.Title != nil {
			episode.Title = *plan.Ticket.Title
		}
	}

	return episode, true, nil
}

// decodeFence reads one JSON fence of a Markdown document.
func decodeFence(markdown, language string, value any) error {
	_, after, found := strings.Cut(markdown, "```"+language+"\n")
	if !found {
		return errors.New("no " + language + " fence")
	}

	body, _, closed := strings.Cut(after, "\n```")
	if !closed {
		return errors.New("the " + language + " fence is not closed")
	}

	return json.Unmarshal([]byte(body), value)
}

// heading is the first heading of a Markdown document, without its marker.
func heading(markdown string) string {
	for _, line := range strings.Split(markdown, "\n") {
		if strings.HasPrefix(line, "# ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "# "))
		}
	}

	return ""
}

func nonNil(items []string) []string {
	if items == nil {
		return []string{}
	}

	return items
}
