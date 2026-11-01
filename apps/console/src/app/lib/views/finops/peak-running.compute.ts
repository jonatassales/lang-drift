import type { Execution } from '@data/types'

import { isReported } from './is-reported.compute.ts'

// The highest observation among running executions, as the runtime widget
// reads it: never a sum across machines.
export function peakRunning(
	executions: Execution[],
	key: 'cpuUsed' | 'memoryUsedMiB',
): number | null {
	const observed = executions
		.filter((execution) => execution.state === 'running')
		.map((execution) => execution[key])
		.filter(isReported)
	if (observed.length === 0) return null

	return Math.max(...observed)
}
