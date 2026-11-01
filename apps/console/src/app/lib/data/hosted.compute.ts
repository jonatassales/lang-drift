export function hosted(token: string | undefined) {
	return Boolean(token?.trim())
}
