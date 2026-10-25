import { renderToString } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import type { VisionPoint } from './product-vision-curve'
import { ProductVisionCurve } from './product-vision-curve'

const data: VisionPoint[] = [
	{
		label: 'Aug',
		value: 73,
		event: {
			actors: [
				{
					initials: 'AN',
					name: 'Ana',
				},
			],
			classification: 'review',
			date: 'Aug 20',
			delta: -3,
			id: 'exports',
			reason: 'Rationale is incomplete.',
			title: 'Export behavior changed',
		},
	},
]

describe('Product Vision server-rendered content', () => {
	it('keeps event text and chart dimensions in HTML before loading the renderer', () => {
		const html = renderToString(<ProductVisionCurve data={data} />)
		expect(html).toContain('Export behavior changed')
		expect(html).toContain('Rationale is incomplete.')
		expect(html).toContain('Under Review')
		expect(html).toContain('height:250px')
		expect(html).not.toContain('recharts-wrapper')
	})

	it('preserves compact dimensions without introducing the event detail panel', () => {
		const html = renderToString(<ProductVisionCurve compact data={data} />)
		expect(html).toContain('height:128px')
		expect(html).not.toContain('product-vision-event-detail')
	})
})
