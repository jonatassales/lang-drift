// Production fails closed until a trusted proxy contract for the host is
// implemented and verified; client-supplied forwarding headers never count.
export function trustedOrigin(mode: string): string | null {
	if (mode !== 'production') return '127.0.0.1'

	return null
}
