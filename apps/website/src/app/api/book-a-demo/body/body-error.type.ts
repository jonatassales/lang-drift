// A request body the reader refuses, with the status the route answers.
export class BodyError extends Error {
	constructor(readonly status: number) {
		super('invalid_request')
	}
}
