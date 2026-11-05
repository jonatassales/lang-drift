import { BadRequestException } from '@nestjs/common'
import { describe, expect, it } from 'vitest'

import { demoRequestSubmissionSchema } from 'infra/streaming'

import { ZodValidationPipe } from './zod-validation.pipe.js'

interface RefusalBody {
	code: string
	fieldErrors: Record<string, string>
}

const pipe = new ZodValidationPipe(demoRequestSubmissionSchema)

describe('demo-requests/demo-requests/zod-validation.pipe', () => {
	it('passes a valid body through as the parsed value', () => {
		const value = pipe.transform({
			email: 'ana@acme.com',
			company: '  Acme ',
			position: 'cto',
			locale: 'pt-BR',
		})

		expect(value).toEqual({
			email: 'ana@acme.com',
			company: 'Acme',
			position: 'cto',
			locale: 'pt-BR',
		})
	})

	it('refuses an invalid body with 400 and the first error of each field', () => {
		let refusal: unknown
		try {
			pipe.transform({
				email: 'ana at acme',
				company: '',
				position: 'cto',
				locale: 'pt-BR',
			})
		} catch (error) {
			refusal = error
		}

		expect(refusal).toBeInstanceOf(BadRequestException)
		const body = (refusal as BadRequestException).getResponse() as RefusalBody
		expect(body.code).toBe('VALIDATION_ERROR')
		expect(Object.keys(body.fieldErrors).sort()).toEqual(['company', 'email'])
	})

	it('refuses a body that is not an object', () => {
		expect(() => pipe.transform('ana@acme.com')).toThrow(BadRequestException)
	})
})
