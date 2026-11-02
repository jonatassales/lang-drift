export function retrySeconds(response: Response): number {
	const seconds = Number(response.headers.get('Retry-After'))

	return Number.isSafeInteger(seconds) && seconds > 0
		? Math.min(seconds, 2_678_400)
		: 30
}
