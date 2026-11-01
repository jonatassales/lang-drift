import type { RuntimeWidgetSnapshot } from '@data/types'

export function runtimeState(snapshot: RuntimeWidgetSnapshot): string {
	const { runtime, executions } = snapshot
	const hasRunning = executions.some(
		(execution) => execution.state === 'running',
	)
	if (hasRunning) return 'running'

	const firstExecutionState = executions[0]?.state
	if (firstExecutionState) return firstExecutionState

	if (runtime.desiredState === 'running') return 'unknown'

	return runtime.desiredState
}
