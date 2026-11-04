// Package rag is the structured retrieval over the project's episodes: the
// earlier runs, embedded locally and compared by cosine similarity, then
// scored by the criteria a plan cares about. No model and no service: the
// embedding is deterministic, so the same text always lands at the same
// vector.
package rag

import (
	"hash/fnv"
	"math"
	"regexp"
	"strings"
)

// Dimensions is the width of a local vector.
const Dimensions = 1024

// LocalEmbedderID names the local embedder and its version; an index built
// by another embedder is rebuilt.
const LocalEmbedderID = "local-hashed-tfidf-1"

// Embedder turns text into a vector.
type Embedder interface {
	ID() string
	Embed(text string) []float32
}

// Local is the deterministic embedder: term frequencies of unigrams and
// bigrams, hashed into fixed dimensions, dampened and normalized, so a
// cosine between two texts rewards the words they share.
type Local struct{}

// ID names the embedder.
func (Local) ID() string { return LocalEmbedderID }

var tokenPattern = regexp.MustCompile(`[a-z0-9]+(?:[.-][a-z0-9]+)*`)

// Embed vectorizes text.
func (Local) Embed(text string) []float32 {
	tokens := Tokens(text)
	vector := make([]float32, Dimensions)

	for index, token := range tokens {
		bump(vector, token)
		if index+1 < len(tokens) {
			bump(vector, token+" "+tokens[index+1])
		}
	}

	for index, value := range vector {
		if value > 0 {
			vector[index] = float32(1 + math.Log(float64(value)))
		}
	}

	return Normalize(vector)
}

// Tokens are the lowercase words and dotted or dashed identifiers of a text.
func Tokens(text string) []string {
	return tokenPattern.FindAllString(strings.ToLower(text), -1)
}

func bump(vector []float32, term string) {
	hasher := fnv.New32a()
	_, _ = hasher.Write([]byte(term))
	slot := int(hasher.Sum32() % Dimensions)

	// The sign spreads collisions out instead of piling them up.
	if hasher.Sum32()&1 == 0 {
		vector[slot]++

		return
	}

	vector[slot]--
}

// Normalize scales a vector to unit length; a zero vector stays zero.
func Normalize(vector []float32) []float32 {
	var sum float64
	for _, value := range vector {
		sum += float64(value) * float64(value)
	}

	if sum == 0 {
		return vector
	}

	length := float32(math.Sqrt(sum))
	for index := range vector {
		vector[index] /= length
	}

	return vector
}

// Cosine is the similarity of two vectors, between -1 and 1 for unit
// vectors; it is zero when either is empty.
func Cosine(left, right []float32) float64 {
	if len(left) != len(right) || len(left) == 0 {
		return 0
	}

	var dot, leftSum, rightSum float64
	for index := range left {
		dot += float64(left[index]) * float64(right[index])
		leftSum += float64(left[index]) * float64(left[index])
		rightSum += float64(right[index]) * float64(right[index])
	}

	if leftSum == 0 || rightSum == 0 {
		return 0
	}

	return dot / (math.Sqrt(leftSum) * math.Sqrt(rightSum))
}
