import { z } from 'zod'

import { demoRequestLocales } from './demo-request-locales.data.js'
import { demoRequestPositionIds } from './demo-request-positions.data.js'

// What the website forwards once its own checks accepted the form.
export const demoRequestSubmissionSchema = z
	.object({
		email: z.email().max(254),
		company: z.string().trim().min(1).max(120),
		position: z.enum(demoRequestPositionIds),
		locale: z.enum(demoRequestLocales),
	})
	.strict()
