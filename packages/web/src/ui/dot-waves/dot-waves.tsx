'use client'

import { cn } from '@langdrift/web/utilities'
import { useEffect, useRef } from 'react'

import { type Dot, dashOf, layoutDots, type Point } from './dot-field'

export interface DotWavesProps {
	className?: string
}

const dashLength = 5
const dashWidth = 2
// How quickly the wave center catches up with the pointer, per frame.
const follow = 0.08

interface Field {
	context: CanvasRenderingContext2D
	dots: Dot[]
	height: number
	width: number
}

// The brand purple, lighter in the dark scheme so the dots read on black.
function purpleFor(canvas: HTMLCanvasElement) {
	const styles = getComputedStyle(canvas)
	const dark = document.documentElement.classList.contains('dark')
	const role = dark ? '--primary-hover' : '--primary'

	return styles.getPropertyValue(role).trim()
}

// With no pointer over the field, the center wanders slowly across it.
function wanderingCenter(field: Field, seconds: number): Point {
	return {
		left: field.width * (0.5 + 0.3 * Math.sin(seconds * 0.13)),
		top: field.height * (0.45 + 0.25 * Math.cos(seconds * 0.17)),
	}
}

function draw(field: Field, center: Point, seconds: number, color: string) {
	const { context, dots, height, width } = field
	context.clearRect(0, 0, width, height)
	context.strokeStyle = color
	context.lineWidth = dashWidth
	context.lineCap = 'round'
	for (const dot of dots) {
		const dash = dashOf(dot, center, seconds)
		const reach = dashLength / 2
		context.globalAlpha = dash.alpha
		context.beginPath()
		context.moveTo(
			dash.left - dash.toLeft * reach,
			dash.top - dash.toTop * reach,
		)
		context.lineTo(
			dash.left + dash.toLeft * reach,
			dash.top + dash.toTop * reach,
		)
		context.stroke()
	}
	context.globalAlpha = 1
}

// A field of small purple dashes in rings that travel outward from the
// pointer, the way Google Antigravity's hero moves. The canvas fills its
// positioned parent and ignores the pointer itself; it listens on the window.
// Under reduced motion it draws one still frame, and it pauses while it is off
// screen or the tab is hidden.
export function DotWaves(props: DotWavesProps) {
	const { className } = props

	const canvasRef = useRef<HTMLCanvasElement>(null)

	useEffect(() => {
		const canvas = canvasRef.current
		const context = canvas?.getContext('2d')
		if (!canvas || !context) return

		const reducedMotion = window.matchMedia(
			'(prefers-reduced-motion: reduce)',
		).matches
		const field: Field = {
			context,
			dots: [],
			height: 0,
			width: 0,
		}
		let color = purpleFor(canvas)
		let pointer: Point | null = null
		let center: Point = {
			left: 0,
			top: 0,
		}
		let visible = true
		let frame = 0
		const start = performance.now()

		function resize() {
			if (!canvas) return
			const ratio = window.devicePixelRatio || 1
			field.width = canvas.clientWidth
			field.height = canvas.clientHeight
			canvas.width = Math.round(field.width * ratio)
			canvas.height = Math.round(field.height * ratio)
			context?.setTransform(ratio, 0, 0, ratio, 0, 0)
			field.dots = layoutDots(field.width, field.height)
			center = wanderingCenter(field, 0)
			if (reducedMotion) draw(field, center, 0, color)
		}

		function tick(now: number) {
			frame = window.requestAnimationFrame(tick)
			if (!visible || document.hidden) return

			const seconds = (now - start) / 1000
			const target = pointer ?? wanderingCenter(field, seconds)
			center = {
				left: center.left + (target.left - center.left) * follow,
				top: center.top + (target.top - center.top) * follow,
			}
			draw(field, center, seconds, color)
		}

		function trackPointer(event: PointerEvent) {
			if (!canvas) return
			const bounds = canvas.getBoundingClientRect()
			const left = event.clientX - bounds.left
			const top = event.clientY - bounds.top
			const inside =
				left >= 0 && top >= 0 && left <= bounds.width && top <= bounds.height
			pointer = inside
				? {
						left,
						top,
					}
				: null
		}

		const sizeObserver = new ResizeObserver(resize)
		sizeObserver.observe(canvas)
		const viewObserver = new IntersectionObserver(([entry]) => {
			visible = entry?.isIntersecting ?? false
		})
		viewObserver.observe(canvas)
		// The theme toggles a class on <html>; the purple follows it.
		const schemeObserver = new MutationObserver(() => {
			color = purpleFor(canvas)
			if (reducedMotion) draw(field, center, 0, color)
		})
		schemeObserver.observe(document.documentElement, {
			attributeFilter: ['class'],
		})
		resize()
		if (!reducedMotion) {
			window.addEventListener('pointermove', trackPointer, { passive: true })
			frame = window.requestAnimationFrame(tick)
		}

		return () => {
			window.cancelAnimationFrame(frame)
			window.removeEventListener('pointermove', trackPointer)
			sizeObserver.disconnect()
			viewObserver.disconnect()
			schemeObserver.disconnect()
		}
	}, [])

	return (
		<div
			aria-hidden="true"
			className={cn('pointer-events-none absolute inset-0', className)}
		>
			<canvas className="block h-full w-full" ref={canvasRef} />
		</div>
	)
}
