// Where one event lives in the bucket: under its aggregate, named by the
// stream entry it came from, so the archive lists a request's history.
export function objectKeyOf(aggregateId: string, streamId: string): string {
	return `demo-request/${aggregateId}/${streamId}.json`
}
