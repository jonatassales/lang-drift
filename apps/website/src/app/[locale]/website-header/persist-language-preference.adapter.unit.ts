// @vitest-environment jsdom
import { afterEach, describe, expect, it, vi } from 'vitest'

import { persistLanguagePreference } from './persist-language-preference.adapter'

afterEach(() => {
	vi.restoreAllMocks()
	window.localStorage.clear()
})

describe('website/website-header/persist-language-preference.adapter', () => {
	it('keeps the language and answers true', () => {
		const kept = persistLanguagePreference('pt-BR')

		expect(kept).toBe(true)
		expect(window.localStorage.getItem('langdrift-website-locale')).toBe(
			'pt-BR',
		)
	})

	it('answers false when the browser refuses the write', () => {
		vi.spyOn(Storage.prototype, 'setItem').mockImplementation(() => {
			throw new Error('quota exceeded')
		})

		const kept = persistLanguagePreference('ja')

		expect(kept).toBe(false)
	})
})
