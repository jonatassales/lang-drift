'use client'

import { defaultRangeExtractor, useVirtualizer } from '@tanstack/react-virtual'
import type { ReactNode } from 'react'
import { useEffect, useRef, useState } from 'react'

interface VirtualListProps<T> {
	items: T[]
	itemKey: (item: T) => string
	renderItem: (item: T) => ReactNode
	onEnd: () => void
	label: string
}

/** Measures variable-height rows and preserves the row owning keyboard focus. */
export function VirtualList<T>(props: VirtualListProps<T>) {
	const { items, itemKey, renderItem, onEnd, label } = props

	const viewport = useRef<HTMLElement>(null)
	const [focused, setFocused] = useState<number | null>(null)
	const virtualizer = useVirtualizer({
		count: items.length,
		getScrollElement: () => viewport.current,
		estimateSize: () => 180,
		getItemKey: (index) => itemKey(items[index]),
		overscan: 5,
		rangeExtractor: (range) => {
			const indices = defaultRangeExtractor(range)
			if (
				focused !== null &&
				focused < items.length &&
				!indices.includes(focused)
			)
				indices.push(focused)

			return indices.sort((left, right) => left - right)
		},
	})
	const rows = virtualizer.getVirtualItems()
	const last = rows.at(-1)?.index ?? -1
	useEffect(() => {
		if (last >= items.length - 3 && last >= 0) onEnd()
	}, [last, items.length, onEnd])

	return (
		<section
			ref={viewport}
			aria-label={label}
			// biome-ignore lint/a11y/noNoninteractiveTabindex: Keyboard users must be able to scroll the region.
			tabIndex={0}
			className="h-[65vh] min-h-80 overflow-auto rounded-xl border border-border bg-card"
		>
			<ul
				className="relative w-full"
				style={{
					height: virtualizer.getTotalSize(),
				}}
			>
				{rows.map((row) => (
					<li
						key={row.key}
						data-index={row.index}
						ref={virtualizer.measureElement}
						className="absolute top-0 left-0 w-full border-b border-border p-4"
						style={{
							transform: `translateY(${row.start}px)`,
						}}
						onFocusCapture={() => setFocused(row.index)}
						onBlurCapture={(event) => {
							if (!event.currentTarget.contains(event.relatedTarget))
								setFocused(null)
						}}
						aria-posinset={row.index + 1}
						aria-setsize={-1}
					>
						{renderItem(items[row.index])}
					</li>
				))}
			</ul>
		</section>
	)
}
