import { describe, expect, it } from 'vitest'

import { objectKeyOf } from './object-key.compute.js'

describe('demo-requests/archive/object-key.compute', () => {
	it('names the object after the aggregate and the stream entry', () => {
		expect(
			objectKeyOf('701706be-7e70-4681-9608-826d3ad548b9', '1762178700000-0'),
		).toBe(
			'demo-request/701706be-7e70-4681-9608-826d3ad548b9/1762178700000-0.json',
		)
	})
})
