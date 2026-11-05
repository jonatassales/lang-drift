// @vitest-environment jsdom
import { cleanup, render, screen } from '@testing-library/react'
import { afterEach, describe, expect, it } from 'vitest'

import { headerMessages } from '@i18n/messages/index'
import { routing } from '@i18n/routing'

import { PlatformFaqSection } from './platform-faq-section'

afterEach(() => {
	cleanup()
})

describe('website/lib/components/home/platform-faq-section', () => {
	it.each(routing.locales)(
		'answers every product in its own tab in %s',
		(locale) => {
			const names = headerMessages[locale]

			render(<PlatformFaqSection locale={locale} />)

			const tabs = screen.getAllByRole('tab').map((tab) => tab.textContent)

			expect(tabs).toEqual([names.microruntimes, names.finops, names.safeSeat])
		},
	)
})
