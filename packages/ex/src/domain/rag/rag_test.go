package rag_test

import (
	"testing"

	"github.com/jonatassales/tao/packages/ex/src/domain/rag"
)

func TestLocalEmbedderIsDeterministicAndNormalized(t *testing.T) {
	t.Parallel()

	embedder := rag.Local{}
	text := "Run the console with extub as a microruntime: Next.js and Go on port 3001"

	first, second := embedder.Embed(text), embedder.Embed(text)
	if len(first) != rag.Dimensions {
		t.Fatalf("a vector has %d dimensions", len(first))
	}

	for index := range first {
		if first[index] != second[index] {
			t.Fatal("the same text embedded differently")
		}
	}

	if self := rag.Cosine(first, second); self < 0.9999 || self > 1.0001 {
		t.Fatalf("cosine of a text with itself = %f", self)
	}

	similar := embedder.Embed("Run the console preview with extub as a microruntime on Next.js")
	different := embedder.Embed("Rotate the storage keys of the billing export job nightly")

	if rag.Cosine(first, similar) <= rag.Cosine(first, different) {
		t.Fatalf("a similar text scored %f, a different one %f", rag.Cosine(first, similar), rag.Cosine(first, different))
	}

	if rag.Cosine(first, embedder.Embed("")) != 0 {
		t.Fatal("an empty text is not at zero")
	}
}

func TestRetrieveKeepsByThresholdOrCriteria(t *testing.T) {
	t.Parallel()

	embedder := rag.Local{}
	query := rag.Query{Stacks: []string{"nextjs", "go"}, ComponentPaths: []string{"."}, Vector: embedder.Embed("Run the console with extub as a microruntime")}

	episodes := []rag.Episode{
		{Runtime: "console-preview", Title: "Run the console with extub as a microruntime, second attempt", Stacks: []string{"nextjs", "go"}, ComponentPaths: []string{"."}, MeasuredResources: true, MeasuredTokens: true},
		{Runtime: "docs-search-index", Title: "Index the docs for search", Stacks: []string{"python"}, ComponentPaths: []string{"docs"}},
		{Runtime: "batch-export", Title: "Nightly batch export worker", Stacks: []string{"go"}, ComponentPaths: []string{"jobs/export"}, MeasuredResources: true},
	}
	for index := range episodes {
		episodes[index].Vector = embedder.Embed(episodes[index].Text())
	}

	candidates := rag.Retrieve(query, episodes)
	if len(candidates) != 3 || candidates[0].Episode.Runtime != "console-preview" {
		t.Fatalf("Retrieve() = %+v", candidates)
	}

	console := candidates[0]
	if !console.Kept || len(console.Criteria) != 4 || console.Score != 1 {
		t.Errorf("the matching episode = %+v", console)
	}

	for _, candidate := range candidates[1:] {
		if candidate.Episode.Runtime == "batch-export" && (!candidate.Kept || candidate.Score != 0.5) {
			t.Errorf("an episode matching two criteria was dropped: %+v", candidate)
		}
		if candidate.Episode.Runtime == "docs-search-index" && (candidate.Kept || candidate.Reason == "") {
			t.Errorf("an unrelated episode was kept: %+v", candidate)
		}
	}
}
