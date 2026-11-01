import type { Execution } from '@data/types'

/** Highest available observation per execution, never a sum across machines. */
export function peak(
	executions: Execution[],
	key: 'cpuUsed' | 'memoryUsedMiB' | 'llmCostUsd',
): number | null {
	const values = executions
		.filter((execution) => execution.state === 'running')
		.map((execution) => execution[key])
		.filter((value): value is number => value != null)

	return values.length ? Math.max(...values) : null
}
