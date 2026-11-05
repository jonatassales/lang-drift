import { renderToString } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import { Brand } from './brand'
import { LogoMark } from './logo-mark'

describe('web/ui/brand/brand', () => {
	it('renders Brand with default size and wordmark', () => {
		const html = renderToString(<Brand />)

		expect(html).toContain('LangDrift')
		expect(html).toContain('viewBox="0 0 64 64"')
		expect(html).toContain('text-2xl/none')
	})

	it('renders compact Brand and LogoMark with custom size', () => {
		const html = renderToString(<Brand compact />)
		const logoHtml = renderToString(<LogoMark size={40} />)

		expect(html).toContain('text-lg/none')
		expect(html).not.toContain('style=')
		expect(logoHtml).toContain('width="40"')
		expect(logoHtml).toContain('height="40"')
	})
})
