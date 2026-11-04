import type { PipeTransform } from '@nestjs/common'
import { BadRequestException } from '@nestjs/common'
import type { ZodType } from 'zod'

// Validates a request body with a Zod schema and answers 400 with the first
// problem of each field, the way the website's edge does.
export class ZodValidationPipe implements PipeTransform<unknown, unknown> {
	constructor(private readonly schema: ZodType) {}

	transform(value: unknown): unknown {
		const result = this.schema.safeParse(value)
		if (result.success) return result.data

		const fieldErrors: Record<string, string> = {}
		for (const issue of result.error.issues) {
			const field = String(issue.path[0] ?? '')
			if (!field || fieldErrors[field]) continue

			fieldErrors[field] = issue.code
		}

		throw new BadRequestException({
			code: 'VALIDATION_ERROR',
			fieldErrors,
		})
	}
}
