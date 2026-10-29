'use client'

import { useEffect } from 'react'

// Scrolls to the section the URL's fragment names once the page that holds
// it has rendered. Next gives up on a fragment whose target is missing while
// a route's loading state shows, so a link to a section would otherwise land
// at the top of the page.
export function useFragmentLanding(ready: boolean) {
	useEffect(() => {
		if (!ready) return

		const fragment = window.location.hash.slice(1)
		if (!fragment) return

		document.getElementById(fragment)?.scrollIntoView({ block: 'start' })
	}, [ready])
}
