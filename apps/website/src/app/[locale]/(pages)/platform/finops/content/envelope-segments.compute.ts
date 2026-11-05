import type { Envelope, EnvelopeSegmentId } from './finops-content.data'

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
