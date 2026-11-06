import type {
	ComputeEnvelope,
	ComputeEnvelopeSegments,
} from './finops-content.types'

// Splits the team's compute envelope. A metered environment's observed cost
// is used, less its idle cost, which counts apart as idle. An unmetered
// environment holds its whole requested ceiling: its cost is unknown, never
// zero. Available is what remains of the authorized amount, never below zero.
export function computeEnvelope(
	envelope: ComputeEnvelope,
): ComputeEnvelopeSegments {
	const { authorized, environments } = envelope

	let used = 0
	let idle = 0
	let held = 0

	for (const environment of environments) {
		if (environment.observed === null) {
			held += environment.ceiling
			continue
		}

		used += environment.observed - environment.idle
		idle += environment.idle
	}

	return {
		used,
		idle,
		held,
		available: Math.max(0, authorized - used - idle - held),
	}
}
