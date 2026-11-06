import { renderToString } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import {
	CiGatesIcon,
	CliIcon,
	FinopsIcon,
	LlmopsIcon,
	MicroruntimesIcon,
	SafeSeatIcon,
} from './index'

describe('web/ui/icons/platform/platform-icons', () => {
	it.each([
		['CiGatesIcon', CiGatesIcon],
		['CliIcon', CliIcon],
		['FinopsIcon', FinopsIcon],
		['LlmopsIcon', LlmopsIcon],
		['MicroruntimesIcon', MicroruntimesIcon],
		['SafeSeatIcon', SafeSeatIcon],
	])('%s renders a decorative monochrome svg', (_name, Icon) => {
		const html = renderToString(<Icon className="size-4" />)

		expect(html.startsWith('<svg')).toBe(true)
		expect(html).toContain('aria-hidden="true"')
		expect(html).toContain('size-4')
		expect(html).not.toContain('#')
	})
})
