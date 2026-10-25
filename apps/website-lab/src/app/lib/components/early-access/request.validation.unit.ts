import { expect, it } from 'vitest'

import { parseSubmissionJSON } from './request.validation'

it('rejects ambiguous JSON before properties can be overwritten', () => {
	for (const text of [
		'{"email":"a","email":"b"}',
		'{"email":"a","\\u0065mail":"b"}',
		'{"a":{"x":1,"x":2}}',
		'[]',
		'null',
		'{}{}',
	]) {
		expect(() => parseSubmissionJSON(text)).toThrow()
	}
	expect(
		parseSubmissionJSON('{"email":"x@example.com", "locale":"en"}'),
	).toEqual({
		email: 'x@example.com',
		locale: 'en',
	})
})
