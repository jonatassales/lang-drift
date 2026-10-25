import { renderToString } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import { DetailList } from './detail-list'

describe('DetailList', () => {
	it('renders each item as a term and a detail', () => {
		const html = renderToString(
			<DetailList
				items={[
					{
						detail: '1,800 dollars',
						id: 'estimate',
						term: 'Estimate',
					},
					{
						detail: <a href="#evidence">Open</a>,
						id: 'evidence',
						term: 'Evidence',
					},
				]}
			/>,
		)
		expect(html).toContain('<dt')
		expect(html).toContain('Estimate')
		expect(html).toContain('1,800 dollars')
		expect(html).toContain('href="#evidence"')
	})

	it('colors a detail with its status tone', () => {
		const html = renderToString(
			<DetailList
				items={[
					{
						detail: 'Unknown',
						id: 'observed',
						term: 'Observed',
						tone: 'warning',
					},
				]}
			/>,
		)
		expect(html).toContain('text-status-warning')
	})
})
