import { describe, expect, it } from 'vitest'

import { cliMessages } from '@i18n/messages/index'
import { routing } from '@i18n/routing'

import { cliFlowStepIds, cliGovernIds, cliLinks } from './cli-content.data'

describe('website/platform/cli/cli-content.data', () => {
	it.each(routing.locales)(
		'names every catalog entry by its id in %s',
		(locale) => {
			const copy = cliMessages[locale]

			expect(Object.keys(copy.flow.steps)).toEqual([...cliFlowStepIds])
			expect(Object.keys(copy.governs.items)).toEqual([...cliGovernIds])
			expect(copy.hero.command).toBe('npx @langdrift/cli init')
		},
	)

	it('links every destination to an absolute https address', () => {
		for (const link of Object.values(cliLinks)) {
			expect(new URL(link).protocol).toBe('https:')
		}
	})
})
