import { describe, expect, it } from 'vitest'

import { safeSeatMessages } from '@i18n/messages/index'
import { routing } from '@i18n/routing'

import {
	categoryIds,
	evidenceFollowUps,
	evidenceIds,
	evidenceSources,
	flowStepIds,
	limitIds,
	principleIds,
	safeSeatFaqIds,
	sourceIds,
	sourceLinks,
	whyIds,
} from './safe-seat-content.data'

describe('website/platform/safe-seat/safe-seat-content.data', () => {
	it.each(routing.locales)(
		'names every catalog entry by its id in %s',
		(locale) => {
			const copy = safeSeatMessages[locale]

			expect(Object.keys(copy.why.items)).toEqual([...whyIds])
			expect(Object.keys(copy.evidence.items)).toEqual([...evidenceIds])
			expect(Object.keys(copy.reserved.categories)).toEqual([...categoryIds])
			expect(Object.keys(copy.flow.steps)).toEqual([...flowStepIds])
			expect(Object.keys(copy.principles.items)).toEqual([...principleIds])
			expect(Object.keys(copy.limits.items)).toEqual([...limitIds])
			expect(Object.keys(copy.sources.items)).toEqual([...sourceIds])
			expect(Object.keys(copy.faq.items)).toEqual([...safeSeatFaqIds])
		},
	)

	it('links every source to an absolute https record', () => {
		for (const id of sourceIds) {
			expect(new URL(sourceLinks[id]).protocol).toBe('https:')
		}
	})

	it('links every evidence card to a listed source', () => {
		for (const id of evidenceIds) {
			expect(sourceIds).toContain(evidenceSources[id])
		}
	})

	it('links the METR card to its 2026 follow-up beside its 2025 study', () => {
		expect(evidenceSources.metr).toBe('metr-2025')
		expect(evidenceFollowUps).toEqual({ metr: 'metr-2026' })
	})

	it.each(routing.locales)('never abbreviates the name in %s', (locale) => {
		const text = JSON.stringify(safeSeatMessages[locale])

		expect(text).not.toMatch(/\bSS\b/)
		expect(text).toContain('Safe Seat')
	})
})
