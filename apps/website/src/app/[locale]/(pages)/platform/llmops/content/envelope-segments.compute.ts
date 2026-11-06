import type { Envelope, EnvelopeSegmentId } from './llmops-content.types'

export function envelopeSegments(
	funding: Envelope,
): Record<EnvelopeSegmentId, number> {
	const { authorized, protected: protectedAmount, reserved, spent } = funding

	return {
		spent,
		reserved,
		protected: protectedAmount,
		available: authorized - spent - reserved - protectedAmount,
	}
}
