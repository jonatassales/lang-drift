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

type TicketModel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
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
	ID       string `json:"id"`
	Human    bool   `json:"human"`
	Progress string `json:"progress"`
}

type TicketPullRequest struct {
	Number       int    `json:"number"`
	Commit       string `json:"commit"`
	ChecksPassed *int   `json:"checksPassed"`
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

// useCase is one kind of work the generated tickets are made of.
type useCase struct {
	slug       string
	title      string
	useCase    string
	components []TicketComponent
}

var useCases = []useCase{
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
	{"claude-sonnet-5-5", "Claude Sonnet 5.5"},
	{"claude-haiku-5-5", "Claude Haiku 5.5"},
	{"claude-opus-5-5", "Claude Opus 5.5"},
}

var ticketEfforts = []string{"medium", "low", "high"}

var workflowStageIDs = []string{"planning", "spec", "tickets", "implement", "review", "merge"}

// The workflow outcomes the generated tickets cycle through.
var ticketOutcomes = []string{"readyForReview", "running", "paused", "readyToMerge"}

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
// their shape with similarities scaled to the ticket.
func ticketEpisodes(seed int) []TicketEpisode {
	shift := float64(seed%5) * 0.01
	return []TicketEpisode{
		{"ep-0412", "checkout-preview-sizing", 0.81 - shift, "kept", nil},
		{"ep-0388", "signup-load-baseline", 0.74 - shift, "kept", nil},
		{"ep-0351", "notify-api-cold-start", 0.66 - shift, "kept", nil},
		{"ep-0297", "docs-portal-build", 0.52 - shift, "dropped", stringPointer("docsOnly")},
		{"ep-0240", "batch-runner-limits", 0.47 - shift, "dropped", stringPointer("otherRuntime")},
		{"ep-0188", "cron-cleanup", 0.42 - shift, "dropped", stringPointer("noBroker")},
	}
}

// ticketPhase spreads the generated tickets across the phases, newest
// earliest in the work. The early-access ticket has reached the workflow.
func ticketPhase(number int) string {
	if number == 31 {
		return "workflow"
	}
	switch {
	case number > 52:
		return TicketPhases[(60-number)%2]
	case number > 44:
		return TicketPhases[1+(number%2)]
	default:
		return TicketPhases[number%4]
	}
}

func generateTicket(number int, startedAt time.Time) Ticket {
	work := useCases[number%len(useCases)]
	if number == 31 {
		work = useCases[0]
	}
	phase := ticketPhase(number)
	reached := phaseIndex(phase)
	padded := fmt.Sprintf("%04d", number)
	slug := work.slug
	if number != 31 {
		slug = fmt.Sprintf("%s-%d", work.slug, number)
	}
	ticket := Ticket{
		ID:        "ticket-" + padded,
		ProductID: TicketProductID,
		Number:    padded,
		Slug:      slug,
		Title:     work.title,
		Status:    "ready-for-agent",
		Phase:     phase,
		CreatedAt: startedAt.Add(-time.Duration(60-number) * 3 * time.Hour).UTC().Format(time.RFC3339),
		Artifacts: TicketArtifacts{
			PRD:       fmt.Sprintf("%04d-PRD-%s.prd.md", 1+number%7, work.slug),
			DesignDoc: fmt.Sprintf("%04d-DD-%s.dd.md", 3+number%9, work.slug),
			Spec:      fmt.Sprintf("%04d-SPEC-%s.spec.md", 9+number%11, work.slug),
			Guardrail: "0001-finops.gr.md",
		},
	}
	if reached == 0 {
		ticket.Status = "needs-triage"
		return ticket
	}

	components := work.components
	suggested := sumResources(components)
	bounded := Resources{
		CPUMillis: ((suggested.CPUMillis*4/3)/500 + 1) * 500,
		MemoryMiB: ((suggested.MemoryMiB*4/3)/512 + 1) * 512,
	}
	if number == 31 {
		bounded = Resources{CPUMillis: 2000, MemoryMiB: 3072}
	}
	model := ticketModels[number%len(ticketModels)]
	effort := ticketEfforts[number%len(ticketEfforts)]
	if number == 31 {
		model = ticketModels[0]
		effort = "medium"
	}
	analysisCurrent := reached == 1
	analysis := &TicketAnalysis{
		DocumentsRead:      []string{"prd", "designDoc", "spec", "guardrail"},
		UseCase:            stringPointer(work.useCase),
		Bounded:            &bounded,
		BudgetCents:        intPointer(2500),
		SpecComponents:     intPointer(len(components)),
		AcceptanceCriteria: intPointer(2 + number%3),
		Criteria:           []string{"FIN-02", "FIN-03", "FIN-04"},
		Threshold:          0.55,
		Episodes:           ticketEpisodes(number),
		Model:              &model,
		Effort:             stringPointer(effort),
	}
	if analysisCurrent && number%2 == 0 {
		analysis.DocumentsRead = []string{"prd", "designDoc"}
		analysis.SpecComponents = nil
		analysis.AcceptanceCriteria = nil
		analysis.Episodes = []TicketEpisode{}
		analysis.Model = nil
		analysis.Effort = nil
	}
	ticket.Analysis = analysis
	ticket.Status = "in-analysis"
	if reached == 1 {
		return ticket
	}

	planCurrent := reached == 2
	plan := &TicketPlan{
		Bounded:          bounded,
		Components:       components,
		Suggested:        &suggested,
		ValidationErrors: intPointer(0),
		Revision:         1,
		Authorization:    "pending",
	}
	if planCurrent && number%2 == 1 {
		plan.Components = components[:len(components)/2+1]
		plan.Suggested = nil
		plan.ValidationErrors = nil
	}
	ticket.Plan = plan
	ticket.Status = "planned"
	if reached == 2 {
		return ticket
	}

	plan.Authorization = "authorized"
	outcome := ticketOutcomes[number%len(ticketOutcomes)]
	if number == 31 {
		outcome = "readyForReview"
	}
	inputK := 2000 + suggested.CPUMillis
	outputK := 150 + suggested.MemoryMiB/20
	estimate := 300 + suggested.CPUMillis/2
	if number == 31 {
		inputK, outputK, estimate = 3500, 250, 950
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
	}
	active := len(workflowStageIDs)
	switch outcome {
	case "running":
		active = 3
	case "paused":
		active = 1
	case "readyForReview":
		active = 4
	case "readyToMerge":
		active = 5
	}
	for index, id := range workflowStageIDs {
		progress := "pending"
		switch {
		case index < active:
			progress = "done"
		case index == active && outcome == "paused":
			progress = "waiting"
		case index == active && outcome == "readyForReview":
			progress = "waiting"
		case index == active && outcome == "readyToMerge":
			progress = "waiting"
		case index == active:
			progress = "active"
		}
		workflow.Stages = append(workflow.Stages, TicketStage{
			ID:       id,
			Human:    id != "tickets" && id != "implement",
			Progress: progress,
		})
	}
	if outcome != "running" && outcome != "paused" {
		observed := estimate + estimate/30
		if number == 31 {
			observed = 980
		}
		workflow.ObservedCents = intPointer(observed)
		workflow.ReleasedCents = intPointer(reserve - observed)
		workflow.PullRequest = &TicketPullRequest{
			Number:       number%40 + 1,
			Commit:       fmt.Sprintf("%07x", number*0x3f9c2e1%0xfffffff),
			ChecksPassed: intPointer(4),
			ChecksTotal:  4,
		}
		if number == 31 {
			workflow.PullRequest.Number = 19
			workflow.PullRequest.Commit = "3f9c2e1"
		}
	}
	if outcome == "paused" {
		workflow.ObservedCents = intPointer(0)
		workflow.ReleasedCents = intPointer(reserve)
	}
	ticket.Workflow = workflow
	ticket.Status = "in-workflow"
	return ticket
}

const ticketCount = 60

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
		startedAt = time.Time{}
	}
	tickets := make([]Ticket, 0, ticketCount)
	for number := ticketCount; number >= 1; number-- {
		tickets = append(tickets, generateTicket(number, startedAt))
	}
	return tickets, nil
}

// TicketPage lists the tickets that have reached a phase, twenty at a time.
// The cursor is the number of the last ticket returned.
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
	for _, ticket := range tickets {
		number, _ := strconv.Atoi(ticket.Number)
		if number >= before || phaseIndex(ticket.Phase) < minimum {
			continue
		}
		items = append(items, ticket)
		if len(items) == ticketPageSize+1 {
			break
		}
	}
	var next *string
	if len(items) > ticketPageSize {
		items = items[:ticketPageSize]
		last, _ := strconv.Atoi(items[ticketPageSize-1].Number)
		value := strconv.Itoa(last)
		next = &value
	}
	return map[string]any{"source": "stub", "items": items, "nextCursor": next}, nil
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
