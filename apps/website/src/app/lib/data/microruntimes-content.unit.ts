import { describe, expect, it } from 'vitest'

import { microruntimesMessages } from '@i18n/messages/index'

import {
	documentIds,
	documentReads,
	episodes,
	formatDescriptor,
	hasPlayed,
	lifecycleStageIds,
	phaseOf,
	planLineIndex,
	planPhaseIds,
	planPhases,
	planReveals,
	planSession,
	runtimeDescriptor,
	similarityThreshold,
	suggestedTotals,
	ticketPlan,
} from './microruntimes-content'

describe('the simulated planning session', () => {
	it('types exactly one command, first, and waits before every other line', () => {
		const commands = planSession.filter((line) => line.type === 'command')
		expect(commands).toHaveLength(1)
		expect(planSession[0].type).toBe('command')
		for (const line of planSession.slice(1)) {
			expect(line.delay).toBeGreaterThan(0)
		}
	})

	it('lasts between twenty-four and thirty-two seconds at the default typing speed', () => {
		const charactersPerSecond = 24
		const typed = microruntimesMessages.en.session.terminal.lines.plan
		const typingMs = (typed.length / charactersPerSecond) * 1000
		const waitingMs = planSession.reduce(
			(total, line) => total + (line.delay ?? 0),
			0,
		)
		expect(typingMs + waitingMs).toBeGreaterThanOrEqual(24_000)
		expect(typingMs + waitingMs).toBeLessThanOrEqual(32_000)
	})

	it('reveals every card field on a line of the session', () => {
		const ids = planSession.map((line) => line.id)
		for (const lineId of Object.values(planReveals)) {
			expect(ids).toContain(lineId)
		}
	})

	it('plays its three acts in order, each opened by a comment', () => {
		const opening = planPhaseIds.map((phase) => planPhases[phase])
		expect(opening).toEqual(['plan', 'analysis', 'planning'])
		for (const lineId of opening.slice(1)) {
			const line = planSession.find((candidate) => candidate.id === lineId)
			expect(line?.type).toBe('comment')
		}

		expect(phaseOf(0)).toBe('ticket')
		expect(phaseOf(planLineIndex('analysis'))).toBe('ticket')
		expect(phaseOf(planLineIndex('analysis') + 1)).toBe('analysis')
		expect(phaseOf(planLineIndex('planning') + 1)).toBe('plan')
		expect(phaseOf(planSession.length)).toBe('plan')
	})

	it('reads every linked document on a line of the analysis act', () => {
		for (const id of documentIds) {
			const index = planLineIndex(documentReads[id])
			expect(index).toBeGreaterThan(planLineIndex('analysis'))
			expect(index).toBeLessThan(planLineIndex('planning'))
		}
	})

	it('keeps the episodes at or above the threshold and drops the rest with a reason', () => {
		expect(episodes.map((episode) => episode.similarity)).toEqual(
			[...episodes].map((episode) => episode.similarity).sort((a, b) => b - a),
		)
		for (const episode of episodes) {
			const aboveThreshold = episode.similarity >= similarityThreshold
			expect(episode.verdict).toBe(aboveThreshold ? 'kept' : 'dropped')
			expect(episode.reason === undefined).toBe(aboveThreshold)
		}
		expect(
			episodes.filter((episode) => episode.verdict === 'kept'),
		).toHaveLength(3)
	})

	it('resolves five components whose shares fill the suggested totals inside the envelope', () => {
		const ids = planSession.map((line) => line.id)
		expect(ticketPlan.components).toHaveLength(5)
		for (const component of ticketPlan.components) {
			expect(ids).toContain(component.id)
		}

		const totals = suggestedTotals(ticketPlan.components)
		expect(totals).toEqual({
			cpuMillis: 1450,
			memoryMiB: 2304,
		})
		expect(totals.cpuMillis).toBeLessThanOrEqual(ticketPlan.bounded.cpuMillis)
		expect(totals.memoryMiB).toBeLessThanOrEqual(ticketPlan.bounded.memoryMiB)
	})

	it('counts a line as played once the terminal has completed it', () => {
		expect(hasPlayed(0, 'plan')).toBe(false)
		expect(hasPlayed(1, 'plan')).toBe(true)
		expect(hasPlayed(1, 'scope')).toBe(false)
		expect(hasPlayed(planSession.length, 'handoff')).toBe(true)
	})
})

describe('illustrative Microruntimes content', () => {
	it('keeps the six lifecycle stages in order', () => {
		expect(lifecycleStageIds).toEqual([
			'configure',
			'validate',
			'prepare',
			'execute',
			'inspect',
			'cleanup',
		])
	})

	it('keeps the runtime file within the limits the ex agent accepts', () => {
		const { components, cpuMillis, memoryMiB } = runtimeDescriptor
		expect(components.length).toBeGreaterThanOrEqual(1)
		expect(components.length).toBeLessThanOrEqual(16)
		expect(cpuMillis).toBeLessThanOrEqual(64_000)
		expect(cpuMillis).toBeGreaterThanOrEqual(components.length)
		expect(memoryMiB).toBeLessThanOrEqual(1_048_576)
		expect(['low', 'medium', 'high']).toContain(runtimeDescriptor.effort)
	})

	it('uses only committed scripts or allowlisted build tools', () => {
		const allowed = ['docker', 'go', 'npm', 'pnpm', 'make']
		for (const { build } of runtimeDescriptor.components) {
			const committedScript = build.executable.startsWith('./')
			expect(committedScript || allowed.includes(build.executable)).toBe(true)
		}
	})

	it('formats the descriptor as parseable JSON with short arrays inline', () => {
		const text = formatDescriptor(runtimeDescriptor)
		expect(JSON.parse(text)).toEqual(runtimeDescriptor)
		expect(text).toContain('"args": ["build", "-t", "storefront:local", "."]')
		expect(text).toContain('"args": []')
	})
})
