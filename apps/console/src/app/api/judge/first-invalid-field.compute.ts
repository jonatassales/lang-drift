import type { z } from 'zod'

// The field a question or an answer got wrong, as a dotted path: a key the
// contract does not know is named itself, and the object as a objectName, by the
// name the caller gives it, when no field is to blame.
export function firstInvalidField(
	error: z.ZodError,
	objectName: string,
): string {
	const issue = error.issues[0]

	if (!issue) return objectName
	if (issue.code === 'unrecognized_keys') return issue.keys.join(',')
	if (issue.path.length === 0) return objectName

	return issue.path.join('.')
}
