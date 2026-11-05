// What a caught failure says, for the route's error log: the message of an
// Error, or the value itself.
export function errorReason(error: unknown): string {
	return error instanceof Error ? error.message : String(error)
}
