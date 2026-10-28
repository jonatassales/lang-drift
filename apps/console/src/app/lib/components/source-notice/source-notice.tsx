import { RegionRetry } from '@components/route-state/region-retry'
import { readConsole } from '@data/server'

import { StubBanner } from './stub-banner'

export async function SourceNotice() {
	try {
		const { source } = await readConsole<{ source: string }>('meta')

		return <StubBanner source={source} />
	} catch {
		return <RegionRetry />
	}
}
