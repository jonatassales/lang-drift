import type { RuntimeWidgetSnapshot } from '@data/types'

import { isLoadTestResult } from './load-test-result.ts'

type RecordValue = Record<string, unknown>
function record(value: unknown): value is RecordValue {
	return typeof value === 'object' && value !== null && !Array.isArray(value)
}
function text(value: unknown): value is string {
	return typeof value === 'string' && value.length > 0 && value.length <= 1024
}
function date(value: unknown): value is string {
	return text(value) && Number.isFinite(Date.parse(value))
}
function number(value: unknown): value is number {
	return typeof value === 'number' && Number.isFinite(value) && value >= 0
}
function nullable(value: unknown, check: (value: unknown) => boolean): boolean {
	return value == null || check(value)
}
function fields(value: RecordValue, names: string[]): boolean {
	return names.every((name) => text(value[name]))
}
function component(value: unknown): boolean {
	if (!record(value) || !fields(value, ['path', 'name', 'kind'])) return false

	const validStacks =
		value.stacks === undefined ||
		(Array.isArray(value.stacks) &&
			value.stacks.length <= 4 &&
			value.stacks.every((stack) => text(stack) && stack.length <= 32))

	return (
		nullable(value.image, text) && nullable(value.port, number) && validStacks
	)
}
function widget(value: unknown): value is RuntimeWidgetSnapshot {
	if (!record(value) || value.source !== 'stub' || !record(value.runtime))
		return false
	const runtime = value.runtime
	if (runtime.loadTest !== undefined && !isLoadTestResult(runtime.loadTest))
		return false
	const config = runtime.configuration
	if (
		!fields(runtime, ['id', 'productId', 'nativeRuntimeId', 'name']) ||
		!number(runtime.revision) ||
		!record(config)
	)
		return false
	if (
		!number(config.cpu) ||
		config.cpu <= 0 ||
		!number(config.memoryMiB) ||
		config.memoryMiB <= 0 ||
		!text(config.model) ||
		!['low', 'medium', 'high'].includes(String(config.effort)) ||
		(config.llmBudgetUsd !== undefined &&
			(!number(config.llmBudgetUsd) || config.llmBudgetUsd <= 0))
	)
		return false
	if (
		!Array.isArray(runtime.components) ||
		runtime.components.length > 16 ||
		!runtime.components.every(component)
	)
		return false
	if (
		!['running', 'stopped', 'deleted'].includes(String(runtime.desiredState)) ||
		!['purple', 'green', 'yellow', 'orange', 'red'].includes(
			String(value.tone),
		) ||
		typeof value.partial !== 'boolean'
	)
		return false
	if (
		!Array.isArray(value.environments) ||
		value.environments.length > 2 ||
		!Array.isArray(value.executions) ||
		value.executions.length > 2 ||
		!record(value.startedAt)
	)
		return false
	const environments = value.environments
	if (
		!environments.every(
			(item) =>
				record(item) &&
				fields(item, ['id', 'name']) &&
				item.runtimeId === runtime.id &&
				['connected', 'unavailable'].includes(String(item.connection)) &&
				date(item.lastSeenAt) &&
				typeof item.cleanupPending === 'boolean',
		)
	)
		return false
	if (
		!value.executions.every(
			(item) =>
				record(item) &&
				text(item.id) &&
				item.runtimeId === runtime.id &&
				environments.some(
					(environment) => environment.id === item.environmentId,
				) &&
				[
					'starting',
					'running',
					'stopping',
					'stopped',
					'failed',
					'unknown',
				].includes(String(item.state)) &&
				date(item.updatedAt) &&
				nullable(item.observedAt, date) &&
				nullable(item.measuredFrom, date) &&
				(item.cpuUsed === null || number(item.cpuUsed)) &&
				(item.memoryUsedMiB === null || number(item.memoryUsedMiB)) &&
				nullable(item.llmCostUsd, number) &&
				nullable(item.downtimeStartedAt, date) &&
				(item.downtimeMilliseconds === null ||
					number(item.downtimeMilliseconds)) &&
				['observed', 'estimated', 'unknown'].includes(
					String(item.downtimeQuality),
				),
		)
	)
		return false
	const executions = value.executions
	if (
		!Object.entries(value.startedAt).every(
			([id, timestamp]) =>
				date(timestamp) && executions.some((item) => item.id === id),
		)
	)
		return false
	const components = runtime.components

	return (
		value.componentLinks === undefined ||
		(Array.isArray(value.componentLinks) &&
			value.componentLinks.length <= 64 &&
			value.componentLinks.every(
				(link) =>
					record(link) &&
					link.from !== link.to &&
					components.some((item) => item.path === link.from) &&
					components.some((item) => item.path === link.to),
			))
	)
}

/** Validate every field consumed by widgets before trusting a transport frame. */
export function parseStubFrame(
	raw: unknown,
	previous: number,
	now = Date.now(),
): {
	sequence: number
	widgets: RuntimeWidgetSnapshot[]
} | null {
	if (typeof raw !== 'string' || raw.length > 262144) return null
	try {
		const frame: unknown = JSON.parse(raw)
		if (
			!record(frame) ||
			frame.source !== 'stub' ||
			frame.kind !== 'runtime-widgets' ||
			!number(frame.sequence) ||
			!Number.isSafeInteger(frame.sequence) ||
			frame.sequence <= previous ||
			!date(frame.generatedAt) ||
			Math.abs(now - Date.parse(frame.generatedAt)) > 5000 ||
			!Array.isArray(frame.widgets) ||
			frame.widgets.length > 16 ||
			!frame.widgets.every(widget)
		)
			return null
		const ids = frame.widgets.map((item) => item.runtime.id)
		if (new Set(ids).size !== ids.length) return null

		return {
			sequence: frame.sequence,
			widgets: frame.widgets,
		}
	} catch {
		return null
	}
}
