import {
	type StubFrameResult,
	stubFrameValidator,
} from './stub-frame-validator.ts'

export type { StubFrameResult }

/** Validate every field consumed by widgets before trusting a transport frame. */
export function parseStubFrame(
	raw: unknown,
	previous: number,
	now = Date.now(),
): StubFrameResult | null {
	return stubFrameValidator.parse(raw, previous, now)
}
