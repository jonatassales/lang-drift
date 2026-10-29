'use client'

import { useEffect, useEffectEvent, useRef } from 'react'

import { Button } from './button'

interface ListEndProps {
	/** Set after a failed request, so the list waits for a retry. */
	failed: boolean
	hasMore: boolean
	label: string
	loading: boolean
	onEnd: () => void
}

// The end of a paged list in the page's flow: a Load more button that also
// asks for the next page as it comes within reach of the viewport. It re-arms
// after every page, so a short page keeps filling the screen, and stays
// quiet after a failure until the list is retried.
export function ListEnd(props: ListEndProps) {
	const { failed, hasMore, label, loading, onEnd } = props

	const end = useRef<HTMLDivElement>(null)
	const reachEnd = useEffectEvent(onEnd)
	const armed = hasMore && !loading && !failed

	useEffect(() => {
		const element = end.current
		if (!armed || !element) return

		const observer = new IntersectionObserver(
			(entries) => {
				if (entries.some((entry) => entry.isIntersecting)) reachEnd()
			},
			{ rootMargin: '0px 0px 600px 0px' },
		)
		observer.observe(element)

		return () => observer.disconnect()
	}, [armed])

	if (!hasMore) return null

	return (
		<div className="mt-4" ref={end}>
			<Button disabled={loading} onClick={onEnd}>
				{label}
			</Button>
		</div>
	)
}
