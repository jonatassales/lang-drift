'use client'

import { usePathname } from 'next/navigation'
import { useEffect, useRef } from 'react'

// Pages render their own header and main, so Next sees the first element of a
// new page already in view and keeps the previous scroll position. Every
// pathname change therefore resets the scroll to the top, unless the URL
// targets a fragment. The initial load is left to the browser.
export function ScrollToTop() {
	const previous = useRef<string | null>(null)

	const pathname = usePathname()

	// Reset window scroll to top on navigation unless navigating to an anchor hash.
	useEffect(() => {
		const initialLoad = previous.current === null
		const samePath = previous.current === pathname

		previous.current = pathname

		if (initialLoad || samePath) return
		if (window.location.hash) return

		window.scrollTo({
			top: 0,
			left: 0,
			behavior: 'instant',
		})
	}, [pathname])

	return null
}
