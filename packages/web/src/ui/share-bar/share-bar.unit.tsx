import { renderToString } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import { ShareBar } from './share-bar'

describe('ShareBar', () => {
	it('fills its share of the track and never more than all of it', () => {
		expect(renderToString(<ShareBar share={0.25} />)).toContain('width:25%')
		expect(renderToString(<ShareBar share={1.6} />)).toContain('width:100%')
	})

	it('is decorative', () => {
		expect(renderToString(<ShareBar share={0.5} />)).toContain(
			'aria-hidden="true"',
		)
	})
})
