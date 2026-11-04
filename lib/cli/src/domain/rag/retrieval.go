package rag

import (
	"fmt"
	"math"
	"sort"
	"strings"
)

// Episode is one earlier run as the index holds it: what it was, what it
// measured, and its vector.
type Episode struct {
	Runtime           string    `json:"runtime"`
	Ticket            string    `json:"ticket"`
	Title             string    `json:"title"`
	Stacks            []string  `json:"stacks"`
	ComponentPaths    []string  `json:"componentPaths"`
	MeasuredResources bool      `json:"measuredResources"`
	MeasuredTokens    bool      `json:"measuredTokens"`
	Outcome           string    `json:"outcome"`
	Facts             []string  `json:"facts"`
	CPUMillisPeak     int       `json:"cpuMillisPeak"`
	MemoryMiBPeak     int       `json:"memoryMiBPeak"`
	Vector            []float32 `json:"vector"`
}

// Text is what an episode is embedded from.
func (episode Episode) Text() string {
	return strings.Join(append([]string{episode.Title, episode.Ticket, strings.Join(episode.Stacks, " "), strings.Join(episode.ComponentPaths, " ")}, episode.Facts...), "\n")
}

// The retrieval's settings, as the plan records them.
const (
	Method        = "vector+structured"
	TopK          = 6
	Threshold     = 0.55
	KeepAtMatches = 2
)

// The criteria, as the plan names them.
const (
	SameComponentPath = "same-component-path"
	SameStack         = "same-stack"
	MeasuredResources = "measured-resources"
	MeasuredTokens    = "measured-tokens"
)

// Criteria lists the criteria in the order the plan records them.
var Criteria = []string{SameComponentPath, SameStack, MeasuredResources, MeasuredTokens}

// Query is what the plan is about: the stacks and component paths it
// names, and the vector of its ticket.
type Query struct {
	Stacks         []string
	ComponentPaths []string
	Vector         []float32
}

// Candidate is one episode the retrieval considered.
type Candidate struct {
	Episode    Episode
	Similarity float64
	Criteria   []string
	Score      float64
	Kept       bool
	Reason     string
}

// Retrieve takes the top episodes by cosine similarity, scores the criteria
// each matches and keeps the ones above the threshold or matching enough
// criteria. The reason says why, until a judge says better.
func Retrieve(query Query, episodes []Episode) []Candidate {
	candidates := make([]Candidate, 0, len(episodes))

	for _, episode := range episodes {
		candidate := Candidate{Episode: episode, Similarity: round(Cosine(query.Vector, episode.Vector))}
		candidate.Criteria = matches(query, episode)
		candidate.Score = round(float64(len(candidate.Criteria)) / float64(len(Criteria)))
		candidate.Kept = candidate.Similarity >= Threshold || len(candidate.Criteria) >= KeepAtMatches
		candidate.Reason = reason(candidate)

		candidates = append(candidates, candidate)
	}

	sort.SliceStable(candidates, func(left, right int) bool { return candidates[left].Similarity > candidates[right].Similarity })

	if len(candidates) > TopK {
		candidates = candidates[:TopK]
	}

	return candidates
}

func matches(query Query, episode Episode) []string {
	var matched []string

	if shares(query.ComponentPaths, episode.ComponentPaths) {
		matched = append(matched, SameComponentPath)
	}
	if shares(query.Stacks, episode.Stacks) {
		matched = append(matched, SameStack)
	}
	if episode.MeasuredResources {
		matched = append(matched, MeasuredResources)
	}
	if episode.MeasuredTokens {
		matched = append(matched, MeasuredTokens)
	}

	if matched == nil {
		return []string{}
	}

	return matched
}

func shares(left, right []string) bool {
	for _, item := range left {
		for _, other := range right {
			if item == other {
				return true
			}
		}
	}

	return false
}

func reason(candidate Candidate) string {
	similarity := fmt.Sprintf("similarity %.2f", candidate.Similarity)

	switch {
	case candidate.Kept && len(candidate.Criteria) > 0:
		return similarity + " · " + strings.Join(candidate.Criteria, ", ")
	case candidate.Kept:
		return similarity + " · above the threshold"
	case len(candidate.Criteria) > 0:
		return similarity + fmt.Sprintf(" below %.2f · %s", Threshold, strings.Join(candidate.Criteria, ", "))
	}

	return similarity + fmt.Sprintf(" below %.2f · no criterion matched", Threshold)
}

func round(value float64) float64 {
	return math.Round(value*100) / 100
}
