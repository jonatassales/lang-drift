import { renderToString } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import { Skeleton } from './skeleton'

describe('web/ui/skeleton', () => {
	it('is a decorative block that rests for reduced motion', () => {
		const html = renderToString(<Skeleton className="h-4 w-24" />)

		expect(html).toContain('aria-hidden="true"')
		expect(html).toContain('motion-reduce:animate-none')
		expect(html).toContain('h-4 w-24')
	})
})
