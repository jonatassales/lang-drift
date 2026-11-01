import type { RuntimeWidgetSnapshot } from '@data/types'

import { isLoadTestResult } from './load-test-result.ts'

export interface StubFrameResult {
	sequence: number
	widgets: RuntimeWidgetSnapshot[]
}

class StubFrameValidator {
	isRecord(value: unknown): value is Record<string, unknown> {
		return typeof value === 'object' && value !== null && !Array.isArray(value)
	}

	isText(value: unknown): value is string {
		return typeof value === 'string' && value.length > 0 && value.length <= 1024
	}

	isDate(value: unknown): value is string {
		return this.isText(value) && Number.isFinite(Date.parse(value))
	}

	isNumber(value: unknown): value is number {
		return typeof value === 'number' && Number.isFinite(value) && value >= 0
	}

	isNullable(value: unknown, check: (value: unknown) => boolean): boolean {
		return value == null || check(value)
	}

	hasFields(value: Record<string, unknown>, names: string[]): boolean {
		return names.every((name) => this.isText(value[name]))
	}

	isComponent(value: unknown): boolean {
		if (
			!this.isRecord(value) ||
			!this.hasFields(value, ['path', 'name', 'kind'])
		)
			return false

		const validStacks =
			value.stacks === undefined ||
			(Array.isArray(value.stacks) &&
				value.stacks.length <= 4 &&
				value.stacks.every((stack) => this.isText(stack) && stack.length <= 32))

		return (
			this.isNullable(value.image, this.isText.bind(this)) &&
			this.isNullable(value.port, this.isNumber.bind(this)) &&
			validStacks
		)
	}

	isWidget(value: unknown): value is RuntimeWidgetSnapshot {
		if (
			!this.isRecord(value) ||
			value.source !== 'stub' ||
			!this.isRecord(value.runtime)
		)
			return false
		const runtime = value.runtime
		if (runtime.loadTest !== undefined && !isLoadTestResult(runtime.loadTest))
			return false
		const config = runtime.configuration
		if (
			!this.hasFields(runtime, [
				'id',
				'productId',
				'nativeRuntimeId',
				'name',
			]) ||
			!this.isNumber(runtime.revision) ||
			!this.isRecord(config)
		)
			return false
		if (
			!this.isNumber(config.cpu) ||
			config.cpu <= 0 ||
			!this.isNumber(config.memoryMiB) ||
			config.memoryMiB <= 0 ||
			!this.isText(config.model) ||
			!['low', 'medium', 'high'].includes(String(config.effort)) ||
			(config.llmBudgetUsd !== undefined &&
				(!this.isNumber(config.llmBudgetUsd) || config.llmBudgetUsd <= 0))
		)
			return false
		if (
			!Array.isArray(runtime.components) ||
			runtime.components.length > 16 ||
			!runtime.components.every((item) => this.isComponent(item))
		)
			return false
		if (
			!['running', 'stopped', 'deleted'].includes(
				String(runtime.desiredState),
			) ||
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
			!this.isRecord(value.startedAt)
		)
			return false
		const environments = value.environments
		if (
			!environments.every(
				(item) =>
					this.isRecord(item) &&
					this.hasFields(item, ['id', 'name']) &&
					item.runtimeId === runtime.id &&
					['connected', 'unavailable'].includes(String(item.connection)) &&
					this.isDate(item.lastSeenAt) &&
					typeof item.cleanupPending === 'boolean',
			)
		)
			return false
		if (
			!value.executions.every(
				(item) =>
					this.isRecord(item) &&
					this.isText(item.id) &&
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
					this.isDate(item.updatedAt) &&
					this.isNullable(item.observedAt, this.isDate.bind(this)) &&
					this.isNullable(item.measuredFrom, this.isDate.bind(this)) &&
					(item.cpuUsed === null || this.isNumber(item.cpuUsed)) &&
					(item.memoryUsedMiB === null || this.isNumber(item.memoryUsedMiB)) &&
					this.isNullable(item.llmCostUsd, this.isNumber.bind(this)) &&
					this.isNullable(item.downtimeStartedAt, this.isDate.bind(this)) &&
					(item.downtimeMilliseconds === null ||
						this.isNumber(item.downtimeMilliseconds)) &&
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
					this.isDate(timestamp) && executions.some((item) => item.id === id),
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
						this.isRecord(link) &&
						link.from !== link.to &&
						components.some((item) => item.path === link.from) &&
						components.some((item) => item.path === link.to),
				))
		)
	}

	parse(
		raw: unknown,
		previous: number,
		now = Date.now(),
	): StubFrameResult | null {
		if (typeof raw !== 'string' || raw.length > 262144) return null
		try {
			const frame: unknown = JSON.parse(raw)
			if (
				!this.isRecord(frame) ||
				frame.source !== 'stub' ||
				frame.kind !== 'runtime-widgets' ||
				!this.isNumber(frame.sequence) ||
				!Number.isSafeInteger(frame.sequence) ||
				frame.sequence <= previous ||
				!this.isDate(frame.generatedAt) ||
				Math.abs(now - Date.parse(frame.generatedAt)) > 5000 ||
				!Array.isArray(frame.widgets) ||
				frame.widgets.length > 16 ||
				!frame.widgets.every((widget) => this.isWidget(widget))
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
}

export const stubFrameValidator = new StubFrameValidator()
