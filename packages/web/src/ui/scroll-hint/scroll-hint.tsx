'use client'

import { useEffect, useRef, useState } from 'react'

// A blurred band along the bottom of the viewport that tells the visitor the
// page goes on below, for pages whose top holds moving content that never
// scrolls the page itself. Placed last in the content it covers, it watches
// a one-pixel marker at its own position and fades out once that end is on
// screen, so it never rests over what follows.
export function ScrollHint() {
	const markerRef = useRef<HTMLSpanElement>(null)
	const [atEnd, setAtEnd] = useState(false)

	useEffect(() => {
		const marker = markerRef.current
		if (!marker) return

		const observer = new IntersectionObserver(([entry]) =>
			setAtEnd(entry?.isIntersecting ?? false),
		)
		observer.observe(marker)

		return () => observer.disconnect()
	}, [])

	return (
		<>
			<span aria-hidden="true" className="block h-px" ref={markerRef} />
			<div
				aria-hidden="true"
				className="pointer-events-none fixed inset-x-0 bottom-0 z-40 h-8 bg-linear-to-t from-background/85 to-transparent backdrop-blur-sm mask-t-from-40% motion-safe:transition-opacity motion-safe:duration-300 data-[at-end=true]:opacity-0 sm:h-12"
				data-at-end={atEnd}
			/>
		</>
	)
}
