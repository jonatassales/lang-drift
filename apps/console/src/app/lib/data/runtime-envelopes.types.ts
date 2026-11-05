import type {
	Execution,
	Recommendation,
	Runtime,
	RuntimeLoadTestResult,
} from './types'

// What the console API wraps a runtime read and a runtime command's reply in.

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

export interface LoadTestEnvelope {
	loadTest: RuntimeLoadTestResult
}

export interface RecommendationEnvelope {
	recommendation: Recommendation
}
