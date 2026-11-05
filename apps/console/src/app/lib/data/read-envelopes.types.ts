import type {
	Episode,
	Execution,
	LifecycleEvent,
	Product,
	ProductSecret,
	Recommendation,
	Runtime,
	RuntimeLoadTestResult,
	User,
} from './types'

// What the console API wraps each read and each command's reply in: one
// object named after what it carries.

export interface SessionEnvelope {
	user: User | null
}

export interface SourceEnvelope {
	source: string
}

export interface ProductEnvelope {
	product: Product
}

export interface ProductsEnvelope {
	products: Product[]
}

export interface ProductSecretsEnvelope {
	secrets: ProductSecret[]
}

export interface IssuedSecretEnvelope {
	secret: string
}

export interface RuntimeEnvelope {
	runtime: Runtime
}

export interface RuntimeExecutionsEnvelope {
	executions: Execution[]
	runtime: Runtime
}

export interface SavedRuntimeEnvelope {
	runtime: Runtime
	secret?: string
}

export interface LifecycleEventsEnvelope {
	events: LifecycleEvent[]
}

export interface EpisodesEnvelope {
	episodes: Episode[]
}

export interface RecommendationEnvelope {
	recommendation: Recommendation
}

export interface LoadTestEnvelope {
	loadTest: RuntimeLoadTestResult
}
