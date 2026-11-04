import type { DemoRequestReceived, DemoRequestSubmission } from 'streaming'
import {
	createEnvelope,
	demoRequestAggregateType,
	demoRequestEventNames,
	demoRequestReceivedSchema,
} from 'streaming'

// The first event of a demo request: the accepted submission, stamped with
// the time it arrived, under a fresh aggregate id.
export function demoRequestReceivedEnvelope(
	submission: DemoRequestSubmission,
	id: string,
	now: Date,
): DemoRequestReceived {
	const envelope = createEnvelope({
		type: demoRequestEventNames.received,
		version: 1,
		aggregate: {
			type: demoRequestAggregateType,
			id,
		},
		payload: {
			...submission,
			submittedAt: now.toISOString(),
		},
		now,
	})

	return demoRequestReceivedSchema.parse(envelope)
}
