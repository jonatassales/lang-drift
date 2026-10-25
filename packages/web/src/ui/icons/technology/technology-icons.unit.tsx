import { renderToString } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import { GoIcon, NatsIcon, NextjsIcon, NodejsIcon, RedisIcon } from './index'

describe('technology icons', () => {
	it.each([
		['GoIcon', GoIcon],
		['NodejsIcon', NodejsIcon],
		['NatsIcon', NatsIcon],
		['RedisIcon', RedisIcon],
		['NextjsIcon', NextjsIcon],
	])('%s renders a decorative monochrome svg', (_name, Icon) => {
		const html = renderToString(<Icon className="size-4" />)

		expect(html.startsWith('<svg')).toBe(true)
		expect(html).toContain('aria-hidden="true"')
		expect(html).toContain('size-4')
		expect(html).not.toContain('#')
	})
})
