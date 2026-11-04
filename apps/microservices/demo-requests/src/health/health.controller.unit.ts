import { describe, expect, it } from 'vitest'

import { HealthController } from './health.controller.js'

describe('demo-requests/health/health.controller', () => {
	it('answers that the service is up', () => {
		expect(new HealthController().check()).toEqual({
			status: 'ok',
		})
	})
})
