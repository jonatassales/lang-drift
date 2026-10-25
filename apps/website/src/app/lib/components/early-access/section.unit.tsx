// @vitest-environment jsdom
import { cleanup, render, screen } from '@testing-library/react'
import { afterEach, describe, expect, it, vi } from 'vitest'

vi.mock('@i18n/navigation', () => ({ useRouter: () => ({ push: vi.fn() }) }))
vi.mock('./form.action', () => ({ submitEarlyAccess: vi.fn() }))

import { earlyAccessMessages } from '@i18n/messages/index'
import type { WebsiteLocale } from '@i18n/routing'

import { EarlyAccessSection } from './section'

afterEach(() => {
	cleanup()
	vi.unstubAllEnvs()
})

describe.each(['true', undefined, 'invalid'])(
	'paused requests with flag %s',
	(flag) => {
		for (const locale of [
			'en',
			'pt-BR',
			'zh-Hant',
			'ja',
		] satisfies WebsiteLocale[]) {
			it.each(['landing', 'pricing'] as const)(
				`replaces the ${locale} %s form with a return-soon notice`,
				(source) => {
					vi.stubEnv('WEBSITE_PREVIEW_SIGNUP_DISABLED', flag)
					const { container } = render(
						<EarlyAccessSection locale={locale} source={source} />,
					)

					expect(
						screen.getByRole('heading', {
							name:
								source === 'landing'
									? earlyAccessMessages[locale].demoHeading
									: earlyAccessMessages[locale].pausedHeading,
						}),
					).toBeTruthy()
					expect(
						screen.getByText(earlyAccessMessages[locale].pausedHeading),
					).toBeTruthy()
					expect(
						screen.getByText(earlyAccessMessages[locale].pausedDescription),
					).toBeTruthy()
					expect(
						container.querySelector('form, input, button, script'),
					).toBeNull()
					expect(
						screen.queryByText(earlyAccessMessages[locale].helper),
					).toBeNull()
				},
			)
		}
	},
)

describe('active Early Access', () => {
	for (const locale of [
		'en',
		'pt-BR',
		'zh-Hant',
		'ja',
	] satisfies WebsiteLocale[]) {
		it(`renders the ${locale} email form`, () => {
			vi.stubEnv('WEBSITE_PREVIEW_SIGNUP_DISABLED', 'false')
			vi.stubEnv('NEXT_PUBLIC_TURNSTILE_SITE_KEY', '1x00000000000000000000BB')
			const { container } = render(
				<EarlyAccessSection locale={locale} source="landing" />,
			)
			expect(
				screen.getByRole('heading', {
					name: earlyAccessMessages[locale].homeHeading,
				}),
			).toBeTruthy()
			expect(
				screen.getByRole('textbox', {
					name: earlyAccessMessages[locale].emailLabel,
				}),
			).toBeTruthy()
			expect(
				container.querySelector('[name="cf-turnstile-response"]'),
			).toBeNull()
		})
	}
})
