package session

import (
	"strconv"

	"github.com/jonatassales/tao/lib/cli/src/application/session"
	"github.com/jonatassales/tao/lib/cli/src/domain/finops"
)

// line is one entry of the catalog: the kind the website gives the line and
// its words, filled from the event's data.
type line struct {
	kind session.Kind
	text func(session.Data) string
}

// catalog maps a session and a line id to its words. The keys are the
// website's ids; the data keys each line reads are the ones named in its
// function.
var catalog = map[string]line{
	// The plan session.
	key(session.PlanSession, session.PlanScope): {session.Comment, func(session.Data) string {
		return Layout(Comment, Cells{Text: "plan only · nothing is built, started or paid for in this step"})
	}},
	key(session.PlanSession, session.ReadTicket): {session.Output, func(data session.Data) string {
		return Layout(Narrative, Cells{Text: "tool call: read " + text(data, "ticket")})
	}},
	key(session.PlanSession, session.Links): {session.Output, func(data session.Data) string {
		return Layout(Narrative, Cells{Text: text(data, "status") + " · links " + joined(data, "links")})
	}},
	key(session.PlanSession, session.Analysis): {session.Comment, func(session.Data) string {
		return Layout(Comment, Cells{Text: "analysis · reading the linked documents, then the project's episodes"})
	}},
	key(session.PlanSession, session.ReadDocument): {session.Output, func(data session.Data) string {
		return Layout(Narrative, Cells{Text: "tool call: read " + text(data, "name") + " · " + text(data, "summary")})
	}},
	key(session.PlanSession, session.Reads): {session.Success, func(data session.Data) string {
		refused := number(data, "refused")
		if refused == 0 {
			return Layout(Result, Cells{
				Label: "read",
				Value: strconv.Itoa(number(data, "count")) + " documents · nothing outside the ticket's links",
			})
		}

		return Layout(Result, Cells{
			Label: "read",
			Value: strconv.Itoa(number(data, "count")) + " documents · " + strconv.Itoa(refused) + " refused",
		})
	}},
	key(session.PlanSession, session.Guardrail): {session.Output, func(data session.Data) string {
		return Layout(Status, Cells{
			Label: "reads",
			Value: text(data, "path") + " · " + joined(data, "criteria"),
		})
	}},
	key(session.PlanSession, session.JudgeGuardrail): {session.Output, func(data session.Data) string {
		return Layout(Narrative, Cells{Text: "judging · " + joined(data, "judgements")})
	}},
	key(session.PlanSession, session.Bounds): {session.Output, func(data session.Data) string {
		return Layout(Status, Cells{
			Label: "bounded",
			Value: "cpuMillis " + strconv.Itoa(number(data, "cpuMillis")) + " · memoryMiB " + strconv.Itoa(number(data, "memoryMiB")) + " · set by " + text(data, "by") + " in the CLI · record " + text(data, "record"),
		})
	}},
	key(session.PlanSession, session.Discrepancy): {session.Warning, func(data session.Data) string {
		return Layout(Status, Cells{
			Label: "discrepancy",
			Value: text(data, "source") + " says " + text(data, "fact") + " · the input holds",
		})
	}},
	key(session.PlanSession, session.Budget): {session.Output, func(data session.Data) string {
		return Layout(Status, Cells{
			Label: "budget",
			Value: finops.Dollars(number(data, "perRunCents")) + " per run for model calls · set by " + text(data, "by") + " in the CLI",
		})
	}},
	key(session.PlanSession, session.Memory): {session.Output, func(session.Data) string {
		return Layout(Narrative, Cells{Text: "RAG over project episodes · a model and effort for this ticket"})
	}},
	key(session.PlanSession, session.Embed): {session.Output, func(session.Data) string {
		return Layout(Narrative, Cells{Text: "embedding the ticket intent · 1 query vector"})
	}},
	key(session.PlanSession, session.VectorSearch): {session.Output, func(data session.Data) string {
		return Layout(Narrative, Cells{Text: "vector search · project episodes · top " + strconv.Itoa(number(data, "topK")) + " by cosine similarity"})
	}},
	key(session.PlanSession, session.Candidates): {session.Output, func(data session.Data) string {
		count := number(data, "count")
		if count == 0 {
			return Layout(Narrative, Cells{Text: "0 candidates · no episode indexed yet"})
		}

		return Layout(Narrative, Cells{Text: strconv.Itoa(count) + " candidates · similarity " + text(data, "highest") + " to " + text(data, "lowest")})
	}},
	key(session.PlanSession, session.JudgeRelevance): {session.Output, func(data session.Data) string {
		return Layout(Narrative, Cells{Text: "judging · " + strconv.Itoa(number(data, "kept")) + " share the stack and scale · " + strconv.Itoa(number(data, "dropped")) + " differ: " + text(data, "reasons")})
	}},
	key(session.PlanSession, session.Drop): {session.Output, func(data session.Data) string {
		return Layout(Narrative, Cells{Text: "dropped " + strconv.Itoa(number(data, "dropped")) + " below " + text(data, "threshold") + " · kept " + strconv.Itoa(number(data, "kept"))})
	}},
	key(session.PlanSession, session.Episodes): {session.Output, func(data session.Data) string {
		return Layout(Narrative, Cells{Text: strconv.Itoa(number(data, "cited")) + " episodes cited · " + text(data, "measured")})
	}},
	key(session.PlanSession, session.ReasonResources): {session.Output, func(data session.Data) string {
		return Layout(Narrative, Cells{Text: "reasoning · " + text(data, "reason")})
	}},
	key(session.PlanSession, session.WeighEffort): {session.Output, func(data session.Data) string {
		return Layout(Narrative, Cells{Text: "weighing effort · " + text(data, "reason")})
	}},
	key(session.PlanSession, session.WeighModel): {session.Output, func(data session.Data) string {
		return Layout(Narrative, Cells{Text: "weighing model · " + text(data, "reason")})
	}},
	key(session.PlanSession, session.Suggest): {session.Success, func(data session.Data) string {
		return Layout(Result, Cells{
			Label: "suggests",
			Value: text(data, "model") + " · effort " + text(data, "effort"),
		})
	}},
	key(session.PlanSession, session.Planning): {session.Comment, func(session.Data) string {
		return Layout(Comment, Cells{Text: "plan · resolving the components, then sizing them inside the envelope"})
	}},
	key(session.PlanSession, session.Resolve): {session.Output, func(data session.Data) string {
		return Layout(Narrative, Cells{Text: "resolving " + strconv.Itoa(number(data, "count")) + " components named in " + text(data, "namedIn") + " · paths, images, ports"})
	}},
	key(session.PlanSession, session.Inspect): {session.Output, func(data session.Data) string {
		return Layout(Narrative, Cells{Text: "tool call: inspect " + text(data, "path") + " · " + text(data, "manifest") + " · :" + strconv.Itoa(number(data, "port"))})
	}},
	key(session.PlanSession, session.Component): {session.Success, func(data session.Data) string {
		return Layout(ComponentRow, Cells{
			Path:      text(data, "path"),
			StackRole: text(data, "technology") + " · " + text(data, "kind"),
			Share:     strconv.Itoa(number(data, "cpuMillis")) + "m · " + strconv.Itoa(number(data, "memoryMiB")) + " MiB",
		})
	}},
	key(session.PlanSession, session.CrossCheck): {session.Output, func(data session.Data) string {
		return Layout(Narrative, Cells{Text: "cross-checking · " + strconv.Itoa(number(data, "ports")) + " ports distinct · " + strconv.Itoa(number(data, "images")) + " images resolved · " + text(data, "collisions")})
	}},
	key(session.PlanSession, session.Size): {session.Success, func(data session.Data) string {
		return Layout(Result, Cells{
			Label: "size",
			Value: strconv.Itoa(number(data, "cpuMillis")) + "m of " + strconv.Itoa(number(data, "boundedCpuMillis")) + "m · " + strconv.Itoa(number(data, "memoryMiB")) + " of " + strconv.Itoa(number(data, "boundedMemoryMiB")) + " MiB · " + strconv.Itoa(number(data, "deployments")) + " one-replica Deployments · proposed shares",
		})
	}},
	key(session.PlanSession, session.Limits): {session.Warning, func(data session.Data) string {
		return Layout(Status, Cells{
			Label: "limits",
			Value: "gpu and disk are not runtime.md fields · " + text(data, "requested"),
		})
	}},
	key(session.PlanSession, session.Price): {session.Output, func(session.Data) string {
		return Layout(Status, Cells{
			Label: "estimate",
			Value: "runtime has no price source · cpu and memory stay in millicores and MiB",
		})
	}},
	key(session.PlanSession, session.Validate): {session.Output, func(data session.Data) string {
		return Layout(Narrative, Cells{Text: "tool call: validate runtime.md draft against the schema · " + strconv.Itoa(number(data, "errors")) + " errors"})
	}},
	key(session.PlanSession, session.Draft): {session.Success, func(data session.Data) string {
		return Layout(Result, Cells{
			Label: "plan",
			Value: "runtime.md draft, revision " + strconv.Itoa(number(data, "revision")) + " · proposed · authorize with dft agent run",
		})
	}},
	key(session.PlanSession, session.Handoff): {session.Comment, func(session.Data) string {
		return Layout(Comment, Cells{Text: "next: once the team authorizes, dft agent run applies the recorded intent"})
	}},
	key(session.PlanSession, session.PlanBlocked): {session.Error, func(data session.Data) string {
		return "✗ plan       runtime.md draft, revision " + strconv.Itoa(number(data, "revision")) + " · blocked · " + joined(data, "checks") + " failed"
	}},

	// The run session.
	key(session.RunSession, session.RunScope): {session.Comment, func(data session.Data) string {
		return Layout(Comment, Cells{Text: "run · " + text(data, "ticket") + " · the recorded intent, authorized by " + text(data, "by")})
	}},
	key(session.RunSession, session.Intent): {session.Output, func(data session.Data) string {
		return Layout(Labeled, Cells{
			Label: "intent",
			Value: text(data, "model") + " · effort " + text(data, "effort") + " · runtime.md revision " + strconv.Itoa(number(data, "revision")),
		})
	}},
	key(session.RunSession, session.Estimate): {session.Output, func(data session.Data) string {
		return Layout(Labeled, Cells{
			Label: "estimate",
			Value: finops.Dollars(number(data, "estimateCents")) + " · " + finops.Tokens(number(data, "inputThousands")) + " input · " + finops.Tokens(number(data, "outputThousands")) + " output · list prices · FIN-02",
		})
	}},
	key(session.RunSession, session.RunBudget): {session.Output, func(data session.Data) string {
		return Layout(Labeled, Cells{
			Label: "budget",
			Value: finops.Dollars(number(data, "perRunCents")) + " per run · set by " + text(data, "by") + " · FIN-03",
		})
	}},
	key(session.RunSession, session.GateOpened): {session.Success, func(data session.Data) string {
		return Layout(Labeled, Cells{
			Label: "gate",
			Value: "reserved " + finops.Dollars(number(data, "reservedCents")) + " of " + finops.Dollars(number(data, "perRunCents")) + " · opened",
		})
	}},
	key(session.RunSession, session.GatePaused): {session.Warning, func(data session.Data) string {
		bound, boundCents := text(data, "bound"), number(data, "boundCents")
		if bound == "" {
			bound, boundCents = "budget", number(data, "perRunCents")
		}

		return Layout(Labeled, Cells{
			Label: "gate",
			Value: "reservation " + finops.Dollars(number(data, "reservedCents")) + " exceeds the " + finops.Dollars(boundCents) + " " + bound + " · paused",
		})
	}},
	key(session.RunSession, session.GateBlocked): {session.Error, func(data session.Data) string {
		return Layout(Labeled, Cells{
			Label: "gate",
			Value: "blocked · " + text(data, "gate") + " · " + text(data, "detail"),
		})
	}},
	key(session.RunSession, session.GateExhausted): {session.Warning, func(data session.Data) string {
		return Layout(Labeled, Cells{
			Label: "gate",
			Value: "paused · FG-04 · " + text(data, "detail"),
		})
	}},
	key(session.RunSession, session.LapResumed): {session.Output, func(data session.Data) string {
		return Layout(Labeled, Cells{
			Label: "lap",
			Value: text(data, "lap") + " open · " + finops.Dollars(number(data, "reservedCents")) + " reserved · resumed",
		})
	}},
	key(session.RunSession, session.Paused): {session.Comment, func(session.Data) string {
		return Layout(Comment, Cells{Text: "needs a decision: raise the budget or lower the effort · nothing was paid for"})
	}},
	key(session.RunSession, session.Waiting): {session.Output, func(data session.Data) string {
		return Layout(Labeled, Cells{
			Label: "waiting",
			Value: text(data, "stage") + " · the owner confirms",
		})
	}},
	key(session.RunSession, session.Confirmed): {session.Success, func(data session.Data) string {
		return Layout(Result, Cells{
			Label: "confirmed",
			Value: text(data, "stage") + " · by " + text(data, "by"),
		})
	}},
	key(session.RunSession, session.Stage): {session.Output, func(data session.Data) string {
		return Layout(Labeled, Cells{
			Label: text(data, "skill"),
			Value: text(data, "summary"),
		})
	}},
	key(session.RunSession, session.Hook): {session.Output, func(data session.Data) string {
		return Layout(Labeled, Cells{
			Label: "hook",
			Value: text(data, "summary"),
		})
	}},
	key(session.RunSession, session.RuntimeReady): {session.Output, func(data session.Data) string {
		return Layout(Status, Cells{
			Label: "runtime",
			Value: text(data, "runtime") + " · revision " + strconv.Itoa(number(data, "revision")) + " · " + strconv.Itoa(number(data, "deployments")) + " Deployments ready · " + strconv.Itoa(number(data, "cpuMillis")) + "m · " + strconv.Itoa(number(data, "memoryMiB")) + " MiB",
		})
	}},
	key(session.RunSession, session.Completed): {session.Success, func(data session.Data) string {
		return Layout(Result, Cells{
			Label: "completed",
			Value: "implementation · " + finops.Dollars(number(data, "observedCents")) + " observed of " + finops.Dollars(number(data, "reservedCents")) + " reserved · usage record",
		})
	}},
	key(session.RunSession, session.Pushed): {session.Success, func(data session.Data) string {
		return Layout(Result, Cells{
			Label: "pushed",
			Value: text(data, "branch") + " · " + text(data, "commit"),
		})
	}},
	key(session.RunSession, session.NotPushed): {session.Warning, func(data session.Data) string {
		return Layout(Status, Cells{
			Label: "pushed",
			Value: "unknown · " + text(data, "reason"),
		})
	}},
	key(session.RunSession, session.PullRequest): {session.Success, func(data session.Data) string {
		return Layout(Result, Cells{
			Label: "pr #" + strconv.Itoa(number(data, "number")),
			Value: text(data, "title"),
		})
	}},
	key(session.RunSession, session.PullRequestUnknown): {session.Warning, func(data session.Data) string {
		return Layout(Status, Cells{
			Label: "pr",
			Value: "unknown · " + text(data, "reason"),
		})
	}},
	key(session.RunSession, session.Checks): {session.Output, func(data session.Data) string {
		return Layout(Status, Cells{
			Label: "checks",
			Value: joined(data, "names") + " · running",
		})
	}},
	key(session.RunSession, session.ChecksPassed): {session.Success, func(data session.Data) string {
		return Layout(Result, Cells{
			Label: "checks",
			Value: strconv.Itoa(number(data, "passed")) + " of " + strconv.Itoa(number(data, "total")) + " passed",
		})
	}},
	key(session.RunSession, session.ChecksUnknown): {session.Warning, func(data session.Data) string {
		return Layout(Status, Cells{
			Label: "checks",
			Value: "unknown · " + text(data, "reason"),
		})
	}},
	key(session.RunSession, session.StopCommand): {session.Command, func(data session.Data) string {
		return "dft agent stop " + text(data, "runtime")
	}},
	key(session.RunSession, session.LapAbandoned): {session.Success, func(data session.Data) string {
		if text(data, "outcome") != finops.LapAbandoned {
			return Layout(Result, Cells{
				Label: "stopped",
				Value: "lap " + text(data, "lap") + " · by " + text(data, "by") + " · outcome unknown, the reservation is held",
			})
		}

		if data["held"] == true {
			return Layout(Result, Cells{
				Label: "abandoned",
				Value: "lap " + text(data, "lap") + " · by " + text(data, "by") + " · a call's spend is unknown, the reservation is held",
			})
		}

		return Layout(Result, Cells{
			Label: "abandoned",
			Value: "lap " + text(data, "lap") + " · by " + text(data, "by") + " · " + finops.Dollars(number(data, "releasedCents")) + " released to the envelope",
		})
	}},
	key(session.RunSession, session.NextPlan): {session.Comment, func(data session.Data) string {
		return Layout(Comment, Cells{Text: "next: dft agent plan plans " + text(data, "runtime") + " again"})
	}},
	key(session.RunSession, session.RuntimeStopped): {session.Output, func(data session.Data) string {
		return Layout(Status, Cells{
			Label: "runtime",
			Value: "stopped · " + strconv.Itoa(number(data, "removed")) + " Deployments removed · " + text(data, "left"),
		})
	}},
	key(session.RunSession, session.ReadyForReview): {session.Success, func(data session.Data) string {
		return Layout(Result, Cells{
			Label: "ready",
			Value: pullRequest(data, "the branch") + " for review · " + finops.Dollars(number(data, "releasedCents")) + " released to the envelope",
		})
	}},
	key(session.RunSession, session.ReadyToMerge): {session.Success, func(data session.Data) string {
		return Layout(Result, Cells{
			Label: "ready",
			Value: pullRequest(data, "the branch") + " to merge · " + finops.Dollars(number(data, "releasedCents")) + " released to the envelope",
		})
	}},
	key(session.RunSession, session.NextReview): {session.Comment, func(data session.Data) string {
		return Layout(Comment, Cells{Text: "next: the owner reviews " + pullRequest(data, "the branch") + " and decides whether to merge"})
	}},
	key(session.RunSession, session.RuntimeUnknown): {session.Warning, func(data session.Data) string {
		return Layout(Status, Cells{
			Label: "runtime",
			Value: "unknown · " + text(data, "reason") + " · the reservation is held",
		})
	}},
	key(session.RunSession, session.ChecksFailed): {session.Error, func(data session.Data) string {
		return "✗ checks     " + strconv.Itoa(number(data, "passed")) + " of " + strconv.Itoa(number(data, "total")) + " passed · " + joined(data, "failed") + " failed"
	}},
	key(session.RunSession, session.NextMerge): {session.Comment, func(session.Data) string {
		return Layout(Comment, Cells{Text: "next: the merge stays the owner's"})
	}},
}

// websiteAliases are the website's per-document, per-component and per-stage
// ids, which the catalog renders through one generic line each.
var websiteAliases = map[string]string{
	key(session.PlanSession, "readPrd"):       key(session.PlanSession, session.ReadDocument),
	key(session.PlanSession, "readDesignDoc"): key(session.PlanSession, session.ReadDocument),
	key(session.PlanSession, "readSpec"):      key(session.PlanSession, session.ReadDocument),
	key(session.PlanSession, "inspectGo"):     key(session.PlanSession, session.Inspect),
	key(session.PlanSession, "inspectNode"):   key(session.PlanSession, session.Inspect),
	key(session.PlanSession, "inspectNats"):   key(session.PlanSession, session.Inspect),
	key(session.PlanSession, "inspectRedis"):  key(session.PlanSession, session.Inspect),
	key(session.PlanSession, "inspectNext"):   key(session.PlanSession, session.Inspect),
	key(session.PlanSession, "go"):            key(session.PlanSession, session.Component),
	key(session.PlanSession, "node"):          key(session.PlanSession, session.Component),
	key(session.PlanSession, "nats"):          key(session.PlanSession, session.Component),
	key(session.PlanSession, "redis"):         key(session.PlanSession, session.Component),
	key(session.PlanSession, "next"):          key(session.PlanSession, session.Component),
	key(session.RunSession, "grilling"):       key(session.RunSession, session.Stage),
	key(session.RunSession, "toSpec"):         key(session.RunSession, session.Stage),
	key(session.RunSession, "toTickets"):      key(session.RunSession, session.Stage),
	key(session.RunSession, "implement"):      key(session.RunSession, session.Stage),
	key(session.RunSession, "codeReview"):     key(session.RunSession, session.Stage),
}

func key(sessionName, id string) string { return sessionName + "." + id }

// Text renders the words of one event; an id the catalog does not know
// renders as itself, so a line is never lost.
func Text(event session.Event) string {
	entry, known := catalog[key(event.Session, event.ID)]
	if !known {
		return event.ID
	}

	return entry.text(event.Data)
}

// KindOf is the kind the catalog gives a line, or the event's own when the
// catalog does not know the id.
func KindOf(event session.Event) session.Kind {
	entry, known := catalog[key(event.Session, event.ID)]
	if !known {
		return event.Kind
	}

	return entry.kind
}

// Covers reports whether a session's line, by the website's id, has words.
func Covers(sessionName, id string) bool {
	if _, known := catalog[key(sessionName, id)]; known {
		return true
	}

	_, aliased := websiteAliases[key(sessionName, id)]

	return aliased
}

// pullRequest names a pull request by its number, or says what stands in
// for it when the run observed none.
func pullRequest(data session.Data, fallback string) string {
	if number(data, "number") == 0 {
		return fallback
	}

	return "#" + strconv.Itoa(number(data, "number"))
}
