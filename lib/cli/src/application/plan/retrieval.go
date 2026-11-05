package plan

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/jonatassales/tao/lib/cli/src/application"
	"github.com/jonatassales/tao/lib/cli/src/application/session"
	"github.com/jonatassales/tao/lib/cli/src/domain/collection"
	"github.com/jonatassales/tao/lib/cli/src/domain/harness"
	"github.com/jonatassales/tao/lib/cli/src/domain/rag"
	"github.com/jonatassales/tao/lib/cli/src/domain/run"
	contract "github.com/jonatassales/tao/lib/core/contracts/go/judge"
)

// retrieved is what the RAG found for a ticket.
type retrieved struct {
	candidates []rag.Candidate
	cited      []rag.Candidate
	// measured is the strongest observation among the cited episodes, or nil.
	measured *harness.Resources
	note     string
}

// noMeasuredNote is the sizing note when no cited episode measured anything.
const noMeasuredNote = "no cited episode measured cpu or memory · sizing follows the inputs, not history"

// retrieve embeds the ticket, searches the episodes, asks the judge why each
// candidate is or is not comparable, and prints lines 14 to 21.
func (planning *planning) retrieve(ctx context.Context) error {
	planning.emit(session.Memory, nil)

	embedder := planning.service.embedder()
	query := rag.Query{Vector: embedder.Embed(queryText(planning.ticket, planning.collected))}
	for _, component := range planning.request.Components {
		query.ComponentPaths = append(query.ComponentPaths, component.Path)
		for _, technology := range component.Technologies {
			query.Stacks = append(query.Stacks, technology.ID)
		}
	}
	planning.emit(session.Embed, nil)

	episodes, err := planning.service.Index.Episodes(planning.owner, embedder)
	if err != nil {
		return err
	}
	planning.emit(session.VectorSearch, session.Data{"topK": rag.TopK})

	retrieval := retrieved{candidates: rag.Retrieve(query, episodes), note: noMeasuredNote}
	threshold := fmt.Sprintf("%.2f", rag.Threshold)

	if len(retrieval.candidates) == 0 {
		planning.emit(session.Candidates, session.Data{"count": 0})
		planning.emit(session.Drop, session.Data{"dropped": 0, "threshold": threshold, "kept": 0})
		planning.emit(session.Episodes, session.Data{"cited": 0, "measured": "none with measured resources"})
		planning.emit(session.ReasonResources, session.Data{"reason": retrieval.note})
		planning.retrieval = retrieval

		return nil
	}

	judged, err := planning.judgeEpisodes(ctx, retrieval.candidates)
	if err != nil {
		return err
	}

	var droppedReasons []string
	for _, candidate := range retrieval.candidates {
		if candidate.Kept {
			retrieval.cited = append(retrieval.cited, candidate)

			continue
		}
		reason := "below " + threshold
		if judged[candidate.Episode.Runtime] {
			reason = candidate.Reason
		}
		droppedReasons = append(droppedReasons, reason)
	}

	planning.emit(session.Candidates, session.Data{"count": len(retrieval.candidates), "highest": fmt.Sprintf("%.2f", retrieval.candidates[0].Similarity), "lowest": fmt.Sprintf("%.2f", retrieval.candidates[len(retrieval.candidates)-1].Similarity)})
	planning.emit(session.JudgeRelevance, session.Data{"kept": len(retrieval.cited), "dropped": len(droppedReasons), "reasons": strings.Join(distinct(droppedReasons), ", ")})
	planning.emit(session.Drop, session.Data{"dropped": len(droppedReasons), "threshold": threshold, "kept": len(retrieval.cited)})

	measured := retrieval.measure()
	measuredWords := "none with measured resources"
	if measured > 0 {
		measuredWords = strconv.Itoa(measured) + " with measured resources"
		retrieval.note = strconv.Itoa(measured) + " cited episodes measured cpu and memory · sizing follows their peaks with headroom"
	}

	planning.emit(session.Episodes, session.Data{"cited": len(retrieval.cited), "measured": measuredWords})
	planning.emit(session.ReasonResources, session.Data{"reason": retrieval.note})
	planning.retrieval = retrieval

	return nil
}

// judgeEpisodes replaces each candidate's reason with the judge's and says
// which it judged. The verdict stays the threshold's and the criteria's.
func (planning *planning) judgeEpisodes(ctx context.Context, candidates []rag.Candidate) (map[string]bool, error) {
	judged := map[string]bool{}

	known := map[string]int{}
	var prompt strings.Builder
	prompt.WriteString("Ticket: " + planning.ticket.Title + "\n\n")
	for index, candidate := range candidates {
		known[candidate.Episode.Runtime] = index
		episode := candidate.Episode
		prompt.WriteString(material(episode.Runtime, "title: "+episode.Title+"\nstacks: "+strings.Join(episode.Stacks, ", ")+"\ncomponents: "+strings.Join(episode.ComponentPaths, ", ")+"\nmeasured resources: "+strconv.FormatBool(episode.MeasuredResources)+"\nsimilarity: "+fmt.Sprintf("%.2f", candidate.Similarity), 2000) + "\n\n")
	}
	prompt.WriteString(episodesPrompt)

	question := Question{
		Kind:   EpisodesQuestion,
		System: judgeSystem,
		Prompt: prompt.String(),
	}

	parsed, err := planning.judge.ask(ctx, question, func(answer contract.Answer) error {
		for _, episode := range answer.Episodes.Episodes {
			if _, found := known[episode.Run]; !found {
				return fmt.Errorf("it names %q, which was not retrieved", episode.Run)
			}
		}

		return nil
	})
	if err != nil {
		return judged, err
	}

	for _, episode := range parsed.Episodes.Episodes {
		if episode.Reason != "" {
			candidates[known[episode.Run]].Reason = episode.Reason
			judged[episode.Run] = true
		}
	}

	return judged, nil
}

// measure finds the strongest peak among the cited episodes that measured
// resources, and says how many did.
func (retrieval *retrieved) measure() int {
	measured := 0
	for _, candidate := range retrieval.cited {
		if !candidate.Episode.MeasuredResources {
			continue
		}
		measured++

		peak := harness.Resources{CPUMillis: candidate.Episode.CPUMillisPeak, MemoryMiB: candidate.Episode.MemoryMiBPeak}
		if retrieval.measured == nil {
			retrieval.measured = &peak

			continue
		}
		retrieval.measured.CPUMillis = max(retrieval.measured.CPUMillis, peak.CPUMillis)
		retrieval.measured.MemoryMiB = max(retrieval.measured.MemoryMiB, peak.MemoryMiB)
	}

	return measured
}

func (service Service) embedder() rag.Embedder {
	if service.Embedder != nil {
		return service.Embedder
	}

	return rag.Local{}
}

// queryText is what the ticket is embedded from: its title and its text.
func queryText(ticket collection.Ticket, collected collection.Result) string {
	content := collected.Sources[0].Content
	if len(content) > 4096 {
		content = content[:4096]
	}

	return ticket.Title + "\n" + string(content)
}

// share is the proposed share of one component when cited episodes measured
// resources: the strongest peak with half again as headroom, rounded up to
// 50 millicores and 128 MiB, split equally.
func (retrieval retrieved) share(count int) harness.Resources {
	cpu := roundUp(retrieval.measured.CPUMillis*3/2, 50)
	memory := roundUp(retrieval.measured.MemoryMiB*3/2, 128)

	return harness.Resources{CPUMillis: max(1, cpu/count), MemoryMiB: max(1, memory/count)}
}

func roundUp(value, step int) int {
	if value%step == 0 {
		return value
	}

	return value + step - value%step
}

func (retrieval retrieved) record() run.Retrieval {
	considered := make([]run.Considered, 0, len(retrieval.candidates))
	for _, candidate := range retrieval.candidates {
		verdict := "dropped"
		if candidate.Kept {
			verdict = "kept"
		}
		considered = append(considered, run.Considered{Run: candidate.Episode.Runtime, Similarity: candidate.Similarity, Criteria: nonNil(candidate.Criteria), Score: candidate.Score, Verdict: verdict, Reason: candidate.Reason})
	}

	note := "first run of this project; no earlier run to cite"
	if len(retrieval.candidates) > 0 {
		note = strconv.Itoa(len(retrieval.cited)) + " of " + strconv.Itoa(len(retrieval.candidates)) + " candidates cited"
	}

	return run.Retrieval{Method: rag.Method, Criteria: append([]string{}, rag.Criteria...), KeepAtMatches: rag.KeepAtMatches, Considered: considered, Note: note}
}

func (retrieval retrieved) sizing() run.Sizing {
	cited := []string{}
	for _, candidate := range retrieval.cited {
		if candidate.Episode.MeasuredResources {
			cited = append(cited, candidate.Episode.Runtime)
		}
	}

	basis := "input"
	if retrieval.measured != nil {
		basis = "measured-runs"
	}

	return run.Sizing{Basis: basis, MeasuredRunsCited: cited, Note: retrieval.note}
}

// ReindexRequest rebuilds the episode index.
type ReindexRequest struct {
	Project application.ProjectSelection
}

// ReindexResult says how many episodes the index holds.
type ReindexResult struct {
	Episodes int    `json:"episodes"`
	Embedder string `json:"embedder"`
	Status   string `json:"status"`
}

// Reindex rebuilds the index from every runtime's memory.
func (service Service) Reindex(_ context.Context, request ReindexRequest) (ReindexResult, error) {
	owner, err := service.Projects.Open(request.Project)
	if err != nil {
		return ReindexResult{}, err
	}

	embedder := service.embedder()
	episodes, err := service.Index.Rebuild(owner, embedder)
	if err != nil {
		return ReindexResult{}, err
	}

	return ReindexResult{Episodes: len(episodes), Embedder: embedder.ID(), Status: "rebuilt"}, nil
}
