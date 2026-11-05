import { describe, expect, it } from 'vitest'

import { ciGatesMessages } from '@i18n/messages/index'
import { routing } from '@i18n/routing'

import { finopsGateIds, postureIds } from './ci-gates-content.data'

describe('website/platform/ci-gates/ci-gates-content.data', () => {
	it.each(routing.locales)(
		'names every catalog entry by its id in %s',
		(locale) => {
			const copy = ciGatesMessages[locale]

			expect(Object.keys(copy.posture.items)).toEqual([...postureIds])
			expect(Object.keys(copy.gates.items)).toEqual([...finopsGateIds])
			expect(copy.install.command).toBe('npx @langdrift/cli init --ci-gates')
		},
	)
})
