import type { RuntimeWidgetSnapshot } from '@data/types'

import { stubFrameValidator } from './stub-frame-validator.ts'

/** Validate every field consumed by widgets before trusting a transport frame. */
export function parseStubFrame(
	raw: unknown,
	previous: number,
	now = Date.now(),
): {
	sequence: number
	widgets: RuntimeWidgetSnapshot[]
} | null {
	return stubFrameValidator.parse(raw, previous, now)
}
