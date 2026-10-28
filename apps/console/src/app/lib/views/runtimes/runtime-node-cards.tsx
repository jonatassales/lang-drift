'use client'
import { type RefObject, useEffect, useRef } from 'react'
import type { SinapsiElement } from 'sinapsi'

import type { Component } from '@data/types'

import { RuntimeStackIcons } from './runtime-stack-icons'

interface RuntimeNodeCardsProps {
	graph: RefObject<SinapsiElement | null>
	components: Component[]
	ready: boolean
}
export function RuntimeNodeCards(props: RuntimeNodeCardsProps) {
	const { graph, components, ready } = props

	const cards = useRef(new Map<string, HTMLDivElement>())
	useEffect(() => {
		const element = graph.current
		if (!element || !ready) return
		let frame = 0
		let visible = false
		const place = () => {
			frame = 0
			if (!visible || window.document.visibilityState !== 'visible') return
			const bounds = element.getBoundingClientRect()
			const positions = element.getNodePositions()
			const columns = Math.max(1, Math.floor(bounds.width / 112))
			const rows = Math.max(3, Math.ceil(components.length / columns))
			const occupied = new Set<number>()
			const slots = Array.from({ length: columns * rows }, (_, slot) => ({
				slot,
				x: ((slot % columns) * bounds.width) / columns + 6,
				y:
					Math.floor(slot / columns) *
						((bounds.height - 60) / Math.max(1, rows - 1)) +
					12,
			}))
			for (const point of positions) {
				const card = cards.current.get(point.id)
				if (!card) continue
				slots.sort(
					(left, right) =>
						Math.hypot(
							left.x - point.x + bounds.left,
							left.y - point.y + bounds.top,
						) -
						Math.hypot(
							right.x - point.x + bounds.left,
							right.y - point.y + bounds.top,
						),
				)
				const position = slots.find((slot) => !occupied.has(slot.slot))
				if (!position) continue
				occupied.add(position.slot)
				// Viewport coordinates come from Sinapsi's public geometry API.
				card.style.transform = `translate(${position.x}px, ${position.y}px)`
				card.style.visibility = 'visible'
			}
			if (!element.paused) frame = requestAnimationFrame(place)
		}
		const update = () => {
			if (!frame) frame = requestAnimationFrame(place)
		}
		const observer = new IntersectionObserver(([entry]) => {
			visible = entry.isIntersecting
			update()
		})
		observer.observe(element)
		const resize = new ResizeObserver(update)
		resize.observe(element)
		const mutations = new MutationObserver(update)
		mutations.observe(element, {
			attributes: true,
			attributeFilter: ['paused', 'nodes'],
		})
		window.document.addEventListener('visibilitychange', update)

		return () => {
			cancelAnimationFrame(frame)
			observer.disconnect()
			resize.disconnect()
			mutations.disconnect()
			window.document.removeEventListener('visibilitychange', update)
		}
	}, [graph, components, ready])

	return (
		<div
			className="pointer-events-none absolute inset-0 overflow-hidden"
			aria-hidden
		>
			{components.map((component) => (
				<div
					key={component.path}
					data-component-card={component.path}
					ref={(element) => {
						if (element) cards.current.set(component.path, element)
						else cards.current.delete(component.path)
					}}
					className="invisible absolute left-0 top-0 w-24 rounded-md border border-border bg-card/95 px-2 py-1.5 text-foreground shadow-sm"
				>
					<p className="truncate text-[10px] font-medium">{component.name}</p>
					<p className="mt-1 flex items-center gap-1 text-muted-foreground">
						<RuntimeStackIcons stacks={component.stacks} size={14} />
						{!component.stacks?.length ? (
							<span className="truncate text-[9px] text-muted-foreground">
								{component.kind}
							</span>
						) : null}
					</p>
				</div>
			))}
		</div>
	)
}
