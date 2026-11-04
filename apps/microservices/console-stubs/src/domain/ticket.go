package domain

import (
	"fmt"
	"strconv"
	"time"
)

// Ticket phases, in order. A ticket's phase is the last one it has reached;
// a section exists once its phase is reached and may hold unknown fields,
// as null, while the phase is current. Unknown is never zero.
var TicketPhases = []string{"ticket", "analysis", "plan", "workflow"}

// The seeded project the stub tickets belong to; a project created at
// runtime has none.
const TicketProductID = "seed-microruntimes"

const ticketPageSize = 20

type Resources struct {
	CPUMillis int `json:"cpuMillis"`
	MemoryMiB int `json:"memoryMiB"`
}

type TicketArtifacts struct {
	PRD       string `json:"prd"`
	DesignDoc string `json:"designDoc"`
	Spec      string `json:"spec"`
	Guardrail string `json:"guardrail"`
}

// A model with its list prices in cents per million tokens, as read on
// the date the workflow records.
type TicketModel struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	ShortName   string `json:"shortName"`
	InputCents  int    `json:"inputCents"`
	OutputCents int    `json:"outputCents"`
}

type TicketEpisode struct {
	ID         string  `json:"id"`
	Slug       string  `json:"slug"`
	Similarity float64 `json:"similarity"`
	Verdict    string  `json:"verdict"`
	Reason     *string `json:"reason"`
}

// What the analysis read from each linked document, and what it suggests.
type TicketAnalysis struct {
	DocumentsRead      []string        `json:"documentsRead"`
	UseCase            *string         `json:"useCase"`
	Bounded            *Resources      `json:"bounded"`
	BudgetCents        *int            `json:"budgetCents"`
	SpecComponents     *int            `json:"specComponents"`
	AcceptanceCriteria *int            `json:"acceptanceCriteria"`
	Criteria           []string        `json:"criteria"`
	Threshold          float64         `json:"threshold"`
	Episodes           []TicketEpisode `json:"episodes"`
	MeasuredResources  *Resources      `json:"measuredResources"`
	Model              *TicketModel    `json:"model"`
	Effort             *string         `json:"effort"`
}

type TicketComponent struct {
	Path       string `json:"path"`
	Technology string `json:"technology"`
	Kind       string `json:"kind"`
	CPUMillis  int    `json:"cpuMillis"`
	MemoryMiB  int    `json:"memoryMiB"`
	Port       int    `json:"port"`
}

type TicketPlan struct {
	Bounded          Resources         `json:"bounded"`
	Components       []TicketComponent `json:"components"`
	Suggested        *Resources        `json:"suggested"`
	ValidationErrors *int              `json:"validationErrors"`
	Revision         int               `json:"revision"`
	Authorization    string            `json:"authorization"`
}

type TicketStage struct {
	ID       string  `json:"id"`
	Skills   *string `json:"skills"`
	Human    bool    `json:"human"`
	Progress string  `json:"progress"`
}

type TicketPullRequest struct {
	Number       int    `json:"number"`
	Commit       string `json:"commit"`
	ChecksPassed int    `json:"checksPassed"`
	ChecksTotal  int    `json:"checksTotal"`
}

type TicketWorkflow struct {
	Model         TicketModel        `json:"model"`
	Effort        string             `json:"effort"`
	Stages        []TicketStage      `json:"stages"`
	BudgetCents   int                `json:"budgetCents"`
	EstimateCents int                `json:"estimateCents"`
	ReserveCents  int                `json:"reserveCents"`
	InputTokensK  int                `json:"inputTokensK"`
	OutputTokensK int                `json:"outputTokensK"`
	ObservedCents *int               `json:"observedCents"`
	ReleasedCents *int               `json:"releasedCents"`
	PullRequest   *TicketPullRequest `json:"pullRequest"`
	Outcome       string             `json:"outcome"`
	PriceDate     string             `json:"priceDate"`
}

type Ticket struct {
	ID        string          `json:"id"`
	ProductID string          `json:"productId"`
	Number    string          `json:"number"`
	Slug      string          `json:"slug"`
	Title     string          `json:"title"`
	Status    string          `json:"status"`
	Phase     string          `json:"phase"`
	CreatedAt string          `json:"createdAt"`
	Artifacts TicketArtifacts `json:"artifacts"`
	Analysis  *TicketAnalysis `json:"analysis"`
	Plan      *TicketPlan     `json:"plan"`
	Workflow  *TicketWorkflow `json:"workflow"`
}

func phaseIndex(phase string) int {
	for index, candidate := range TicketPhases {
		if candidate == phase {
			return index
		}
	}
	return -1
}

// workTemplate is one kind of work the generated tickets are made of.
type workTemplate struct {
	slug       string
	title      string
	useCase    string
	components []TicketComponent
}

var workTemplates = []workTemplate{
	{"early-access-runtime", "Run the early-access sign-up use case as a microruntime", "early-access sign-up", []TicketComponent{
		{"services/signup-api", "Go", "api", 250, 256, 8080},
		{"services/notify-api", "Node.js", "api", 300, 512, 8081},
		{"platform/nats", "NATS", "broker", 200, 256, 4222},
		{"platform/redis", "Redis", "cache", 200, 256, 6379},
		{"apps/signup", "Next.js", "web", 500, 1024, 3000},
	}},
	{"checkout-retries", "Reduce checkout retries without a slower checkout", "checkout retries", []TicketComponent{
		{"services/checkout-api", "Go", "api", 400, 512, 8080},
		{"platform/redis", "Redis", "cache", 200, 256, 6379},
		{"apps/storefront", "Next.js", "web", 500, 1024, 3000},
	}},
	{"invoice-export", "Export invoices as a nightly job", "invoice export", []TicketComponent{
		{"services/billing-api", "Node.js", "api", 300, 512, 8081},
		{"platform/nats", "NATS", "broker", 200, 256, 4222},
	}},
	{"tenant-migration", "Migrate tenants to the new schema in batches", "tenant migration", []TicketComponent{
		{"services/tenants-api", "Go", "api", 500, 768, 8080},
		{"platform/nats", "NATS", "broker", 200, 256, 4222},
		{"platform/redis", "Redis", "cache", 200, 256, 6379},
	}},
	{"onboarding-copy", "Refresh the onboarding copy in five locales", "onboarding copy", []TicketComponent{
		{"apps/onboarding", "Next.js", "web", 400, 768, 3000},
	}},
	{"search-reindex", "Reindex search nightly within the reservation", "search reindex", []TicketComponent{
		{"services/search-api", "Go", "api", 600, 1024, 8080},
		{"platform/nats", "NATS", "broker", 200, 256, 4222},
	}},
	{"usage-webhooks", "Send usage webhooks with retries and a dead letter", "usage webhooks", []TicketComponent{
		{"services/webhooks-api", "Node.js", "api", 300, 512, 8081},
		{"platform/nats", "NATS", "broker", 200, 256, 4222},
		{"platform/redis", "Redis", "cache", 200, 256, 6379},
	}},
	{"pricing-page", "Ship the pricing page behind a flag", "pricing page", []TicketComponent{
		{"apps/website", "Next.js", "web", 500, 1024, 3000},
		{"services/plans-api", "Go", "api", 250, 256, 8080},
	}},
}

var ticketModels = []TicketModel{
	{"claude-sonnet-5-5", "Claude Sonnet 5.5", "Sonnet 5.5", 200, 1000},
	{"claude-haiku-5-5", "Claude Haiku 5.5", "Haiku 5.5", 10, 50},
	{"claude-opus-5-5", "Claude Opus 5.5", "Opus 5.5", 400, 2000},
}

// The date the list prices above were read.
const priceDate = "2026-10-06"

var ticketEfforts = []string{"medium", "low", "high"}

var workflowStageIDs = []string{"planning", "spec", "tickets", "implement", "review", "merge"}

// The harness skills each stage runs; review and merge run none.
var stageSkills = map[string]string{"planning": "grilling", "spec": "to-spec", "tickets": "to-tickets", "implement": "implement · tdd · code-review"}

// The workflow outcomes the generated tickets cycle through, and the stage
// each one has reached.
var ticketOutcomes = []string{"readyForReview", "running", "paused", "readyToMerge"}

var outcomeStage = map[string]int{"running": 3, "paused": 0, "readyForReview": 4, "readyToMerge": 5}

const ticketCount = 120

// The early-access ticket carries the website's figures, so the console
// and the site describe it the same way.
const earlyAccessNumber = 31

func isEarlyAccess(number int) bool {
	return number == earlyAccessNumber
}

func intPointer(value int) *int {
	return &value
}

func stringPointer(value string) *string {
	return &value
}

func sumResources(components []TicketComponent) Resources {
	total := Resources{}
	for _, component := range components {
		total.CPUMillis += component.CPUMillis
		total.MemoryMiB += component.MemoryMiB
	}
	return total
}

// The early-access ticket's episodes are the website's; the others reuse
// their shape with similarities shifted by the ticket number.
func ticketEpisodes(number int) []TicketEpisode {
	shift := float64(number%5) * 0.01
	if isEarlyAccess(number) {
		shift = 0
	}
	return []TicketEpisode{
		{"ep-0412", "checkout-preview-sizing", 0.81 - shift, "kept", nil},
		{"ep-0388", "notify-api-rollout", 0.74 - shift, "kept", nil},
		{"ep-0351", "signup-redesign-preview", 0.66 - shift, "kept", nil},
		{"ep-0297", "docs-search-index", 0.52 - shift, "dropped", stringPointer("docsOnly")},
		{"ep-0240", "batch-export-worker", 0.47 - shift, "dropped", stringPointer("otherRuntime")},
		{"ep-0188", "storefront-cache-warmup", 0.42 - shift, "dropped", stringPointer("noBroker")},
	}
}

// ticketPhase spreads the generated tickets across the phases, the newest
// earliest in the work. The early-access ticket has reached the workflow.
func ticketPhase(number int) string {
	if isEarlyAccess(number) {
		return "workflow"
	}
	switch {
	case number > ticketCount-8:
		return TicketPhases[(ticketCount-number)%2]
	case number > ticketCount-16:
		return TicketPhases[1+(number%2)]
	default:
		return TicketPhases[number%4]
	}
}

func generateTicket(number int, startedAt time.Time) Ticket {
	work := workTemplates[number%len(workTemplates)]
	slug := fmt.Sprintf("%s-%d", work.slug, number)
	if isEarlyAccess(number) {
		work = workTemplates[0]
		slug = work.slug
	}
	phase := ticketPhase(number)
	reached := phaseIndex(phase)
	padded := fmt.Sprintf("%04d", number)
	ticket := Ticket{
		ID:        "ticket-" + padded,
		ProductID: TicketProductID,
		Number:    padded,
		Slug:      slug,
		Title:     work.title,
		Status:    "needs-triage",
		Phase:     phase,
		CreatedAt: startedAt.Add(-time.Duration(ticketCount-number) * 3 * time.Hour).UTC().Format(time.RFC3339),
		Artifacts: TicketArtifacts{
			PRD:       fmt.Sprintf("%04d-PRD-%s.prd.md", 1+number%7, work.slug),
			DesignDoc: fmt.Sprintf("%04d-DD-%s.dd.md", 3+number%9, work.slug),
			Spec:      fmt.Sprintf("%04d-SPEC-%s.spec.md", 9+number%11, work.slug),
			Guardrail: "0001-finops.gr.md",
		},
	}
	if reached < 1 {
		return ticket
	}

	suggested := sumResources(work.components)
	bounded := Resources{
		CPUMillis: ((suggested.CPUMillis*4/3)/500 + 1) * 500,
		MemoryMiB: ((suggested.MemoryMiB*4/3)/512 + 1) * 512,
	}
	model := ticketModels[number%len(ticketModels)]
	effort := ticketEfforts[number%len(ticketEfforts)]
	if isEarlyAccess(number) {
		ticket.Artifacts = TicketArtifacts{
			PRD:       "0002-PRD-early-access.prd.md",
			DesignDoc: "0005-DD-early-access-runtime.dd.md",
			Spec:      "0011-SPEC-early-access-runtime.spec.md",
			Guardrail: "0001-finops.gr.md",
		}
		bounded = Resources{CPUMillis: 2000, MemoryMiB: 3072}
		model = ticketModels[0]
		effort = "medium"
	}
	ticket.Analysis = generateAnalysis(number, work, bounded, model, effort, reached == 1)
	ticket.Status = "in-analysis"
	if reached < 2 {
		return ticket
	}

	ticket.Plan = generatePlan(number, work, bounded, suggested, reached == 2)
	ticket.Status = "planned"
	if reached < 3 {
		return ticket
	}

	ticket.Plan.Authorization = "authorized"
	ticket.Workflow = generateWorkflow(number, suggested, model, effort)
	ticket.Status = "in-workflow"
	return ticket
}

// An analysis still in progress has read only part of the documents: what it
// has not read, judged or suggested yet is null.
func generateAnalysis(number int, work workTemplate, bounded Resources, model TicketModel, effort string, current bool) *TicketAnalysis {
	analysis := &TicketAnalysis{
		DocumentsRead:      []string{"prd", "designDoc", "spec", "guardrail"},
		UseCase:            stringPointer(work.useCase),
		Bounded:            &bounded,
		BudgetCents:        intPointer(2500),
		SpecComponents:     intPointer(len(work.components)),
		AcceptanceCriteria: intPointer(2 + number%3),
		Criteria:           []string{"FIN-02", "FIN-03", "FIN-04"},
		Threshold:          0.55,
		Episodes:           ticketEpisodes(number),
		Model:              &model,
		Effort:             stringPointer(effort),
	}
	if !current || number%2 == 1 {
		return analysis
	}
	analysis.DocumentsRead = []string{"prd", "designDoc"}
	analysis.SpecComponents = nil
	analysis.AcceptanceCriteria = nil
	analysis.Criteria = nil
	analysis.Episodes = nil
	analysis.Model = nil
	analysis.Effort = nil
	return analysis
}

// A plan still sizing has resolved only part of its components and has no
// totals or validation yet.
func generatePlan(number int, work workTemplate, bounded, suggested Resources, current bool) *TicketPlan {
	plan := &TicketPlan{
		Bounded:          bounded,
		Components:       work.components,
		Suggested:        &suggested,
		ValidationErrors: intPointer(0),
		Revision:         1,
		Authorization:    "pending",
	}
	if !current || number%2 == 0 {
		return plan
	}
	plan.Components = work.components[:len(work.components)/2+1]
	plan.Suggested = nil
	plan.ValidationErrors = nil
	return plan
}

// A workflow reserves its estimate plus a fifth before any paid call. A run
// the gate paused asked for more than the budget, so no stage ran and
// nothing was paid; a running one has no observed spend or pull request yet.
func generateWorkflow(number int, suggested Resources, model TicketModel, effort string) *TicketWorkflow {
	outcome := ticketOutcomes[(number/4)%len(ticketOutcomes)]
	inputK := 2000 + suggested.CPUMillis
	outputK := 150 + suggested.MemoryMiB/20
	estimate := 300 + suggested.CPUMillis/2
	if isEarlyAccess(number) {
		outcome = "readyForReview"
		inputK, outputK, estimate = 3500, 250, 950
	}
	if outcome == "paused" {
		estimate = 2200
	}
	reserve := estimate * 6 / 5
	workflow := &TicketWorkflow{
		Model:         model,
		Effort:        effort,
		BudgetCents:   2500,
		EstimateCents: estimate,
		ReserveCents:  reserve,
		InputTokensK:  inputK,
		OutputTokensK: outputK,
		Outcome:       outcome,
		PriceDate:     priceDate,
	}
	reachedStage := outcomeStage[outcome]
	for index, id := range workflowStageIDs {
		progress := "pending"
		switch {
		case outcome == "paused":
		case index < reachedStage:
			progress = "done"
		case index == reachedStage && outcome == "running":
			progress = "active"
		case index == reachedStage:
			progress = "waiting"
		}
		stage := TicketStage{
			ID:       id,
			Human:    id != "tickets" && id != "implement",
			Progress: progress,
		}
		if skills, ok := stageSkills[id]; ok {
			stage.Skills = stringPointer(skills)
		}
		workflow.Stages = append(workflow.Stages, stage)
	}
	switch outcome {
	case "paused":
		workflow.ObservedCents = intPointer(0)
	case "readyForReview", "readyToMerge":
		observed := estimate + estimate/30
		pullRequest := TicketPullRequest{
			Number:       number%40 + 1,
			Commit:       fmt.Sprintf("%07x", number*0x3f9c2e1%0xfffffff),
			ChecksPassed: 4,
			ChecksTotal:  4,
		}
		if isEarlyAccess(number) {
			observed = 980
			pullRequest.Number = 19
			pullRequest.Commit = "3f9c2e1"
		}
		workflow.ObservedCents = intPointer(observed)
		workflow.ReleasedCents = intPointer(reserve - observed)
		workflow.PullRequest = &pullRequest
	}
	return workflow
}

// Tickets are generated, newest first, for the seeded project only.
func (catalog Catalog) Tickets(product string) ([]Ticket, error) {
	if _, err := catalog.Product(product); err != nil {
		return nil, err
	}
	if product != TicketProductID {
		return []Ticket{}, nil
	}
	startedAt, err := time.Parse(time.RFC3339Nano, catalog.StartedAt)
	if err != nil {
		return nil, fmt.Errorf("catalog start time %q: %w", catalog.StartedAt, err)
	}
	tickets := make([]Ticket, 0, ticketCount)
	for number := ticketCount; number >= 1; number-- {
		tickets = append(tickets, generateTicket(number, startedAt))
	}
	return tickets, nil
}

// TicketPage lists the tickets that have reached a phase, twenty at a time.
// The cursor is the number of the last ticket returned; the tickets come
// newest first, so the ticket at index i has number ticketCount - i.
func (catalog Catalog) TicketPage(product, reached, cursor string) (map[string]any, error) {
	if reached == "" {
		reached = "ticket"
	}
	minimum := phaseIndex(reached)
	if minimum < 0 {
		return nil, ErrInvalid
	}
	before := ticketCount + 1
	if cursor != "" {
		value, err := strconv.Atoi(cursor)
		if err != nil || value < 1 || value > ticketCount || strconv.Itoa(value) != cursor {
			return nil, ErrInvalid
		}
		before = value
	}
	tickets, err := catalog.Tickets(product)
	if err != nil {
		return nil, err
	}
	items := []Ticket{}
	last := 0
	for index, ticket := range tickets {
		number := ticketCount - index
		if number >= before || phaseIndex(ticket.Phase) < minimum {
			continue
		}
		if len(items) == ticketPageSize {
			next := strconv.Itoa(last)
			return map[string]any{"source": "stub", "items": items, "nextCursor": &next}, nil
		}
		items = append(items, ticket)
		last = number
	}
	return map[string]any{"source": "stub", "items": items, "nextCursor": nil}, nil
}

// Ticket finds one generated ticket by id.
func (catalog Catalog) Ticket(id string) (Ticket, error) {
	tickets, err := catalog.Tickets(TicketProductID)
	if err != nil {
		return Ticket{}, err
	}
	for _, ticket := range tickets {
		if ticket.ID == id {
			return ticket, nil
		}
	}
	return Ticket{}, ErrNotFound
}
