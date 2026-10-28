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

/** Observed local extub HTTP timings; not workload or security-test evidence. */
export interface RuntimeLoadTestResult {
	target: 'extub-loopback'
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
		Record<
			'cpu' | 'memoryMiB' | 'model' | 'effort',
			{
				basis:
					| 'observed_resource'
					| 'historical_configuration'
					| 'ticket_inference'
				episodeIds: string[]
			}
		>
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
export interface Page<T> {
	items: T[]
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
