// A nonce for one response: 16 random bytes, base64.
export function createNonce() {
	const bytes = crypto.getRandomValues(new Uint8Array(16))

	return btoa(String.fromCharCode(...bytes))
}
