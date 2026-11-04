import { z } from 'zod'

// The two emails a demo request produces: one to the visitor, one to the team.
export const demoRequestEmailKindSchema = z.enum([
	'confirmation',
	'lead-notification',
])
