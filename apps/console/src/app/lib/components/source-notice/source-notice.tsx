import { RegionRetry } from '@components/route-state/region-retry'
import type { SourceEnvelope } from '@data/read-envelopes.types'
import { readConsole } from '@data/server'

import { StubBanner } from './stub-banner'

export async function SourceNotice() {
	try {
		const { source } = await readConsole<SourceEnvelope>('meta')

		return <StubBanner source={source} />
	} catch {
		return <RegionRetry />
	}
}
