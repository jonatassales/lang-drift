import { renderToString } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import { Card } from './card'
import { CardContent } from './card-content'
import { CardHeader } from './card-header'
import { CardTitle } from './card-title'

describe('web/ui/card', () => {
	it('lays its parts out on a bordered surface', () => {
		const html = renderToString(
			<Card>
				<CardHeader>
					<CardTitle>Runtime</CardTitle>
				</CardHeader>
				<CardContent>Ready</CardContent>
			</Card>,
		)

		expect(html).toContain('data-slot="card"')
		expect(html).toContain('border-border')
		expect(html).toContain('data-slot="card-title"')
		expect(html).toContain('Ready')
	})

	it('lets the caller replace its layout', () => {
		const html = renderToString(<Card className="block p-5">Body</Card>)
		const classes = html.match(/class="([^"]*)"/)?.[1].split(' ')

		expect(classes).toContain('block')
		expect(classes).toContain('p-5')
		expect(classes).not.toContain('flex')
		expect(classes).not.toContain('py-6')
	})
})
