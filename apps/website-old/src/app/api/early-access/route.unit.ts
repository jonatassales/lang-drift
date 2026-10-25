import { afterEach, expect, it, vi } from 'vitest'

import { POST } from './route'

afterEach(() => {
	vi.unstubAllGlobals()
})

it('answers unavailable without calling a provider while sign-up is paused', async () => {
	const transport = vi.fn()
	vi.stubGlobal('fetch', transport)

	const response = POST()

	expect(response.status).toBe(503)
	expect(response.headers.get('Retry-After')).toBe('30')
	expect(response.headers.get('Cache-Control')).toBe('no-store')
	expect(await response.json()).toEqual({
		ok: false,
		code: 'UNAVAILABLE',
	})
	expect(transport).not.toHaveBeenCalled()
})
