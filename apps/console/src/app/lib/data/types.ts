/** Browser-safe public contracts; credentials occur only in one-time responses. */
export interface User {
	id: string
	email: string
}

export interface Repository {
	fullName: string
	url: string
	defaultBranch: string
	private: boolean
}

export interface Component {
	path: string
	name: string
	kind: string
	/** Explicit technology identifiers; absence does not imply a detected stack. */
	stacks?: string[]
	image?: string
	port?: number
	readinessPath?: string
	build?: {
		executable: string
		args: string[]
	}
}

export interface Product {
	id: string
	name: string
	repository?: Repository
	ref: string
	source?: {
		kind: 'local'
		path: string
	}
	components: Component[]
	createdAt: string
}

export interface ProductSecret {
	id: string
	name: string
	createdAt: string
	revokedAt: string | null
}

export interface HomeOverview {
	products: Product[]
	running: {
		runtime: Runtime
		productName: string
		executions: Execution[]
	}[]
	unknownCount: number
}

export interface RuntimeConfiguration {
	cpu: number
	memoryMiB: number
	model: string
	effort: 'low' | 'medium' | 'high'
	nativeRuntimeId?: string
	llmBudgetUsd?: number
}

export type DesiredState = 'running' | 'stopped' | 'deleted'

export interface Runtime {
	id: string
	nativeRuntimeId: string
	productId: string
	name: string
	components: Component[]
	configuration: RuntimeConfiguration
	revision: number
	markdown: string
	desiredState: DesiredState
	intentVersion: number
	credentialRevoked: boolean
	createdAt: string
	updatedAt: string
	loadTest?: RuntimeLoadTestResult
}

/** Observed local console-stubs HTTP timings; not workload or security-test evidence. */
export interface RuntimeLoadTestResult {
	target: 'console-stubs-loopback'
	requests: number
	concurrency: number
	succeeded: number
	failed: number
	durationMs: number
	p50Ms: number
	p95Ms: number
	maxMs: number
	startedAt: string
	completedAt: string
}

export interface Environment {
	id: string
	runtimeId: string
	name: string
	lastSeenAt: string
	connection: 'connected' | 'unavailable'
	cleanupPending: boolean
}

export type ExecutionState =
	| 'starting'
	| 'running'
	| 'stopping'
	| 'stopped'
	| 'failed'
	| 'unknown'

export interface MeasurementWindow {
	pod: string
	observedAt: string
	window: string
}

export interface Execution {
	id: string
	runtimeId: string
	environmentId: string
	state: ExecutionState
	revision: number
	intentVersion: number
	updatedAt: string
	llmCostUsd?: number | null
	cpuUsed: number | null
	memoryUsedMiB: number | null
	observedAt: string | null
	measuredFrom?: string
	measurementWindows?: MeasurementWindow[]
	downtimeStartedAt: string | null
	downtimeMilliseconds: number | null
	downtimeQuality: 'observed' | 'estimated' | 'unknown'
}

export interface LifecycleEvent {
	id: string
	runtimeId: string
	environmentId?: string
	executionId?: string
	kind: string
	cause: string
	actor: string
	occurredAt: string
	receivedAt: string
	revision?: number
	detail?: string
	operationId?: string
	requestedBy?: string
	sourceRevision?: string
	concreteHash?: string
	markdownHash?: string
}

export type RecommendationReasonCode =
	| 'review_cited_evidence'
	| 'resource_observations_missing'
	| 'provider_not_configured'
	| 'model_not_allowed'
	| 'provider_endpoint_invalid'
	| 'provider_request_failed'
	| 'provider_response_invalid'
	| 'stub_suggestion'

export type RecommendationLimitationCode =
	| 'future_workload_unverified'
	| 'cpu_observation_missing'
	| 'memory_observation_missing'
	| 'agent_quality_unmeasured'

export interface RecommendationEvidence {
	basis: 'observed_resource' | 'historical_configuration' | 'ticket_inference'
	episodeIds: string[]
}

export interface Recommendation {
	ticket?: {
		name: string
		hash: string
	}
	retrieval?: 'local_lexical_vectors_v1'
	status: 'suggested' | 'insufficient_evidence' | 'unavailable'
	values?: Partial<RuntimeConfiguration>
	episodeIds: string[]
	reason: string
	reasonCode: RecommendationReasonCode
	limitationCodes: RecommendationLimitationCode[]
	createdAt: string
	generationModel?: string
	generationEffort?: RuntimeConfiguration['effort']
	generationSettingsSource?: 'draft' | 'runtime' | 'operator'
	usage?: Record<string, number> | null
	cost?: null
	limits?: {
		outputTokens: number
		timeoutSeconds: number
		retries: number
	}
	limitations?: string[]
	evidence?: Partial<
		Record<'cpu' | 'memoryMiB' | 'model' | 'effort', RecommendationEvidence>
	>
}

export interface AgentObservation {
	llmCostUsd?: number | null
	cpuUsed: number | null
	memoryUsedMiB: number | null
	observedAt: string
	measuredFrom?: string
	measurementWindows?: MeasurementWindow[]
}

export interface Episode {
	id: string
	runtimeId: string
	productId: string
	environmentId?: string
	executionId?: string
	kind: string
	cause: string
	revision: number
	occurredAt: string
	receivedAt: string
	components: string[]
	configuration: RuntimeConfiguration
	observation?: AgentObservation
}

export type RuntimeWidgetTone = 'purple' | 'green' | 'yellow' | 'orange' | 'red'

export interface RuntimeWidgetSnapshot {
	source?: 'stub'
	componentLinks?: {
		from: string
		to: string
	}[]
	runtime: Runtime
	environments: Environment[]
	executions: Execution[]
	tone: RuntimeWidgetTone
	partial: boolean
	startedAt: Record<string, string>
}

export interface RuntimeDetailSnapshot {
	source?: 'stub'
	runtime: Runtime
	environments: Environment[]
	executions: Execution[]
	widget: RuntimeWidgetSnapshot
}

/** Cursor pages are scoped to the authenticated product; null marks the end. */
export interface Page<Item> {
	items: Item[]
	nextCursor: string | null
}

export interface ProductExecution {
	runtimeName: string
	environmentName: string | null
	execution: Execution
}

export interface ProductEvent {
	runtimeName: string
	environmentName: string | null
	event: LifecycleEvent
}

export interface ProductNavigation {
	products: Product[]
	runtimeProducts: Record<string, string>
}

export interface RuntimeListSnapshot {
	source?: 'stub'
	runtimes: Runtime[]
	widgets: RuntimeWidgetSnapshot[]
}

export interface FolderListing {
	path: string
	parent: string | null
	folders: {
		name: string
		path: string
	}[]
}

// A ticket's phases, in order. Its phase is the last one reached; a
// section exists once its phase is reached, and a null field inside a
// current phase is not known yet: unknown, never zero.
export type TicketPhase = 'ticket' | 'analysis' | 'plan' | 'workflow'

export interface Resources {
	cpuMillis: number
	memoryMiB: number
}

export type TicketDocumentId = 'prd' | 'designDoc' | 'spec' | 'guardrail'

// A model with its list prices in cents per million tokens.
export interface TicketModel {
	id: string
	name: string
	shortName: string
	inputCents: number
	outputCents: number
}

export type EpisodeDropReason = 'docsOnly' | 'otherRuntime' | 'noBroker'

export interface TicketEpisode {
	id: string
	slug: string
	similarity: number
	verdict: 'kept' | 'dropped'
	reason: EpisodeDropReason | null
}

export type TicketEffort = 'low' | 'medium' | 'high'

export interface TicketAnalysis {
	documentsRead: TicketDocumentId[]
	useCase: string | null
	bounded: Resources | null
	budgetCents: number | null
	specComponents: number | null
	acceptanceCriteria: number | null
	criteria: string[] | null
	threshold: number
	episodes: TicketEpisode[] | null
	measuredResources: Resources | null
	model: TicketModel | null
	effort: TicketEffort | null
}

export type ComponentKind = 'api' | 'broker' | 'cache' | 'web'

export type Technology = 'Go' | 'Node.js' | 'NATS' | 'Redis' | 'Next.js'

export interface TicketComponent extends Resources {
	path: string
	technology: Technology
	kind: ComponentKind
	port: number
}

export type ReservedPartCategory = 'core-domain' | 'hard-part' | 'new-stack'

export type ReservedPartChooser = 'judge' | 'lead'

/** The part of a ticket a person writes by hand under Safe Seat. */
export interface TicketReservedPart {
	paths: string[]
	goal: string
	category: ReservedPartCategory
	chosenBy: ReservedPartChooser
}

export interface TicketPlan {
	bounded: Resources
	components: TicketComponent[]
	suggested: Resources | null
	validationErrors: number | null
	revision: number
	authorization: 'pending' | 'authorized'
	/** Null until the plan has sized its components and reserved a part. */
	reservedPart: TicketReservedPart | null
}

export type WorkflowStageId =
	| 'planning'
	| 'spec'
	| 'tickets'
	| 'implement'
	| 'dev-ticket'
	| 'review'
	| 'merge'

/** Who does a stage's work: the agent, or a person writing by hand. */
export type StageDriver = 'agent' | 'person'

export interface TicketStage {
	id: WorkflowStageId
	/** The harness skills the stage runs; review and merge run none. */
	skills: string | null
	/** A person approves the stage's work. */
	human: boolean
	drivenBy: StageDriver
	progress: 'pending' | 'active' | 'waiting' | 'done'
}

export interface TicketPullRequest {
	number: number
	commit: string
	checksPassed: number
	checksTotal: number
}

export type WorkflowOutcome =
	| 'running'
	| 'paused'
	| 'readyForReview'
	| 'readyToMerge'

export interface TicketWorkflow {
	model: TicketModel
	effort: TicketEffort
	stages: TicketStage[]
	budgetCents: number
	estimateCents: number
	reserveCents: number
	inputTokensK: number
	outputTokensK: number
	observedCents: number | null
	releasedCents: number | null
	pullRequest: TicketPullRequest | null
	outcome: WorkflowOutcome
	/** The date the model's list prices were read, as YYYY-MM-DD. */
	priceDate: string
}

export interface Ticket {
	id: string
	productId: string
	number: string
	slug: string
	title: string
	status: string
	phase: TicketPhase
	createdAt: string
	artifacts: Record<TicketDocumentId, string>
	analysis: TicketAnalysis | null
	plan: TicketPlan | null
	workflow: TicketWorkflow | null
	/** Null until the ticket's workflow starts. */
	safeSeat: TicketSafeSeat | null
}

export type SafeSeatState =
	| 'met'
	| 'needs-override'
	| 'overridden'
	| 'not-applicable'
	| 'unknown'

/** A ticket's Safe Seat record. It never holds the person's share. */
export interface TicketSafeSeat {
	state: SafeSeatState
	overrideReason: string | null
	illustrative: boolean
}

export interface TicketDetailSnapshot {
	source: 'stub'
	ticket: Ticket
}
