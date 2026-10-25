'use client'

import { cn } from '@langdrift/web/utilities'
import {
	type MotionValue,
	motion,
	motionValue,
	useReducedMotion,
} from 'motion/react'
import type { ReactElement, ReactNode } from 'react'
import { useEffect, useRef, useState } from 'react'

// Adapted from SmoothUI's Magnetic Field (smoothui.dev): elements attracted
// to, repelled from or orbiting the pointer, with a springy return to rest.

export type MagneticFieldMode = 'attract' | 'repel' | 'orbit'

export interface MagneticFieldItem {
	element: ReactElement
	id: string
}

export interface MagneticFieldProps {
	children?: ReactNode
	className?: string
	damping?: number
	items?: MagneticFieldItem[]
	mode?: MagneticFieldMode
	paused?: boolean
	radius?: number
	rotate?: boolean
	stiffness?: number
	strength?: number
}

interface ParticipantState {
	centerX: number
	centerY: number
	rotate: MotionValue<number>
	vr: number
	vx: number
	vy: number
	x: MotionValue<number>
	y: MotionValue<number>
}

const DEFAULT_RADIUS = 160
const DEFAULT_STRENGTH = 0.35
const DEFAULT_STIFFNESS = 180
const DEFAULT_DAMPING = 18
const MAX_ROTATE_DEG = 18
const FALLOFF_POWER = 2
const MAX_DELTA_SECONDS = 0.05
const ORBIT_RADIAL_FACTOR = 0.2
const ORBIT_TANGENTIAL_FACTOR = 0.4
const MAGNETIC_SELECTOR = '[data-magnetic]'
// Below this offset and speed a participant counts as at rest.
const REST_EPSILON = 0.01

const clamp = (value: number, min: number, max: number) =>
	Math.min(max, Math.max(min, value))

const computeFalloff = (distance: number, radius: number) => {
	if (distance >= radius) return 0

	return (1 - distance / radius) ** FALLOFF_POWER
}

const computeOffset = (
	dx: number,
	dy: number,
	distance: number,
	falloff: number,
	strength: number,
	mode: MagneticFieldMode,
): { x: number; y: number } => {
	if (distance === 0) {
		return {
			x: 0,
			y: 0,
		}
	}
	const directional = mode === 'repel' ? -Math.abs(strength) : strength
	if (mode === 'orbit') {
		const tangentX = -dy / distance
		const tangentY = dx / distance
		const pull = directional * falloff * distance * ORBIT_TANGENTIAL_FACTOR

		return {
			x: dx * directional * falloff * ORBIT_RADIAL_FACTOR + tangentX * pull,
			y: dy * directional * falloff * ORBIT_RADIAL_FACTOR + tangentY * pull,
		}
	}

	return {
		x: dx * directional * falloff,
		y: dy * directional * falloff,
	}
}

const createParticipant = (): ParticipantState => ({
	centerX: 0,
	centerY: 0,
	rotate: motionValue(0),
	vr: 0,
	vx: 0,
	vy: 0,
	x: motionValue(0),
	y: motionValue(0),
})

const resetParticipant = (state: ParticipantState) => {
	state.x.set(0)
	state.y.set(0)
	state.rotate.set(0)
	state.vx = 0
	state.vy = 0
	state.vr = 0
}

export function MagneticField(props: MagneticFieldProps) {
	const {
		children,
		className,
		damping = DEFAULT_DAMPING,
		items,
		mode = 'attract',
		paused = false,
		radius = DEFAULT_RADIUS,
		rotate = false,
		stiffness = DEFAULT_STIFFNESS,
		strength = DEFAULT_STRENGTH,
	} = props

	const shouldReduceMotion = useReducedMotion()
	const [isPointerCapable, setIsPointerCapable] = useState(false)
	const containerRef = useRef<HTMLDivElement>(null)
	const pointerRef = useRef<{ x: number; y: number } | null>(null)
	const participantsRef = useRef<Map<string, ParticipantState>>(new Map())
	const nodeMapRef = useRef<Map<string, HTMLElement>>(new Map())
	const elementsRef = useRef<Map<HTMLElement, ParticipantState>>(new Map())
	const rafRef = useRef<number | null>(null)
	const wakeRef = useRef<(() => void) | null>(null)
	const lastFrameRef = useRef<number | null>(null)

	useEffect(() => {
		const mediaQuery = window.matchMedia('(hover: hover) and (pointer: fine)')
		setIsPointerCapable(mediaQuery.matches)
		const handleChange = (event: MediaQueryListEvent) => {
			setIsPointerCapable(event.matches)
		}
		mediaQuery.addEventListener('change', handleChange)

		return () => mediaQuery.removeEventListener('change', handleChange)
	}, [])

	const isEffectDisabled = paused || shouldReduceMotion || !isPointerCapable

	const getOrCreateParticipant = (id: string): ParticipantState => {
		let state = participantsRef.current.get(id)
		if (!state) {
			state = createParticipant()
			participantsRef.current.set(id, state)
		}

		return state
	}

	// Prune stale participants when the `items` list changes shape.
	useEffect(() => {
		if (!items) {
			participantsRef.current.clear()
			nodeMapRef.current.clear()

			return
		}
		const ids = new Set(items.map((item) => item.id))
		for (const key of Array.from(participantsRef.current.keys())) {
			if (!ids.has(key)) {
				participantsRef.current.delete(key)
				nodeMapRef.current.delete(key)
			}
		}
	}, [items])

	// Measure rest positions for `items` participants once per layout change.
	useEffect(() => {
		if (!items) return
		const container = containerRef.current
		if (!container) return

		const measure = () => {
			const containerRect = container.getBoundingClientRect()
			for (const [id, node] of nodeMapRef.current) {
				const state = participantsRef.current.get(id)
				if (!state) continue
				const rect = node.getBoundingClientRect()
				state.centerX = rect.left + rect.width / 2 - containerRect.left
				state.centerY = rect.top + rect.height / 2 - containerRect.top
			}
		}

		measure()
		const resizeObserver = new ResizeObserver(measure)
		for (const node of nodeMapRef.current.values()) {
			resizeObserver.observe(node)
		}
		window.addEventListener('resize', measure)

		return () => {
			resizeObserver.disconnect()
			window.removeEventListener('resize', measure)
		}
	}, [items])

	// Discover and measure `data-magnetic` descendants of `children`.
	useEffect(() => {
		if (items) return
		const container = containerRef.current
		if (!container) return

		const syncElements = () => {
			const found = Array.from(
				container.querySelectorAll<HTMLElement>(MAGNETIC_SELECTOR),
			)
			const map = elementsRef.current
			for (const element of Array.from(map.keys())) {
				if (!found.includes(element)) {
					element.style.transform = ''
					element.style.willChange = ''
					map.delete(element)
				}
			}
			const containerRect = container.getBoundingClientRect()
			for (const element of found) {
				if (!map.has(element)) {
					element.style.willChange = 'transform'
					map.set(element, createParticipant())
				}
				const state = map.get(element)
				if (!state) continue
				const rect = element.getBoundingClientRect()
				state.centerX = rect.left + rect.width / 2 - containerRect.left
				state.centerY = rect.top + rect.height / 2 - containerRect.top
			}
		}

		syncElements()
		const mutationObserver = new MutationObserver(syncElements)
		mutationObserver.observe(container, {
			childList: true,
			subtree: true,
		})
		const resizeObserver = new ResizeObserver(syncElements)
		resizeObserver.observe(container)
		window.addEventListener('resize', syncElements)

		return () => {
			mutationObserver.disconnect()
			resizeObserver.disconnect()
			window.removeEventListener('resize', syncElements)
			for (const element of elementsRef.current.keys()) {
				element.style.transform = ''
				element.style.willChange = ''
			}
		}
	}, [items])

	// Track the pointer position relative to the container with one listener.
	useEffect(() => {
		if (isEffectDisabled) {
			pointerRef.current = null

			return
		}
		const container = containerRef.current
		if (!container) return

		const handlePointerMove = (event: PointerEvent) => {
			const rect = container.getBoundingClientRect()
			pointerRef.current = {
				x: event.clientX - rect.left,
				y: event.clientY - rect.top,
			}
			wakeRef.current?.()
		}
		const handlePointerLeave = () => {
			pointerRef.current = null
		}

		container.addEventListener('pointermove', handlePointerMove, {
			passive: true,
		})
		container.addEventListener('pointerleave', handlePointerLeave, {
			passive: true,
		})

		return () => {
			container.removeEventListener('pointermove', handlePointerMove)
			container.removeEventListener('pointerleave', handlePointerLeave)
		}
	}, [isEffectDisabled])

	// Single rAF loop driving a manual spring integrator per participant.
	useEffect(() => {
		if (isEffectDisabled) {
			for (const state of participantsRef.current.values()) {
				resetParticipant(state)
			}
			for (const [element, state] of elementsRef.current) {
				resetParticipant(state)
				element.style.transform = 'translate3d(0px, 0px, 0) rotate(0deg)'
			}

			return
		}

		const applyPhysics = (state: ParticipantState, dt: number) => {
			const pointer = pointerRef.current
			let targetX = 0
			let targetY = 0
			let targetRotate = 0

			if (pointer) {
				const dx = pointer.x - state.centerX
				const dy = pointer.y - state.centerY
				const distance = Math.hypot(dx, dy)
				const falloff = computeFalloff(distance, radius)
				if (falloff > 0) {
					const offset = computeOffset(
						dx,
						dy,
						distance,
						falloff,
						strength,
						mode,
					)
					targetX = offset.x
					targetY = offset.y
					if (rotate) {
						targetRotate = clamp(
							((Math.atan2(dy, dx) * 180) / Math.PI) * falloff,
							-MAX_ROTATE_DEG,
							MAX_ROTATE_DEG,
						)
					}
				}
			}

			const currentX = state.x.get()
			const currentY = state.y.get()
			const currentRotate = state.rotate.get()

			state.vx += ((targetX - currentX) * stiffness - state.vx * damping) * dt
			state.vy += ((targetY - currentY) * stiffness - state.vy * damping) * dt
			state.vr +=
				((targetRotate - currentRotate) * stiffness - state.vr * damping) * dt

			state.x.set(currentX + state.vx * dt)
			state.y.set(currentY + state.vy * dt)
			state.rotate.set(currentRotate + state.vr * dt)
		}

		const isAtRest = (state: ParticipantState) =>
			Math.abs(state.x.get()) < REST_EPSILON &&
			Math.abs(state.y.get()) < REST_EPSILON &&
			Math.abs(state.rotate.get()) < REST_EPSILON &&
			Math.abs(state.vx) < REST_EPSILON &&
			Math.abs(state.vy) < REST_EPSILON &&
			Math.abs(state.vr) < REST_EPSILON

		const isSettled = () => {
			if (pointerRef.current) return false
			const states = [
				...participantsRef.current.values(),
				...elementsRef.current.values(),
			]

			return states.every(isAtRest)
		}

		const settle = () => {
			for (const state of participantsRef.current.values()) {
				resetParticipant(state)
			}
			for (const [element, state] of elementsRef.current) {
				resetParticipant(state)
				element.style.transform = 'translate3d(0px, 0px, 0) rotate(0deg)'
			}
			rafRef.current = null
			lastFrameRef.current = null
		}

		const step = (time: number) => {
			const last = lastFrameRef.current ?? time
			const dt = clamp((time - last) / 1000, 0, MAX_DELTA_SECONDS)
			lastFrameRef.current = time

			for (const state of participantsRef.current.values()) {
				applyPhysics(state, dt)
			}
			for (const [element, state] of elementsRef.current) {
				applyPhysics(state, dt)
				const x = state.x.get()
				const y = state.y.get()
				const angle = state.rotate.get()
				element.style.transform = `translate3d(${x}px, ${y}px, 0) rotate(${angle}deg)`
			}

			if (isSettled()) {
				settle()

				return
			}
			rafRef.current = requestAnimationFrame(step)
		}

		// The loop sleeps once every participant is back at rest with no
		// pointer over the field, and the next pointer move wakes it.
		const wake = () => {
			if (rafRef.current !== null) return
			rafRef.current = requestAnimationFrame(step)
		}
		wakeRef.current = wake
		wake()

		return () => {
			if (rafRef.current !== null) cancelAnimationFrame(rafRef.current)
			rafRef.current = null
			lastFrameRef.current = null
			wakeRef.current = null
		}
	}, [damping, isEffectDisabled, mode, radius, rotate, stiffness, strength])

	const content = items
		? items.map((item) => {
				const state = getOrCreateParticipant(item.id)

				return (
					<motion.div
						className="inline-block will-change-transform"
						key={item.id}
						ref={(node: HTMLDivElement | null) => {
							if (node) {
								nodeMapRef.current.set(item.id, node)

								return
							}
							nodeMapRef.current.delete(item.id)
						}}
						style={{
							rotate: state.rotate,
							x: state.x,
							y: state.y,
						}}
					>
						{item.element}
					</motion.div>
				)
			})
		: children

	return (
		<div className={cn('relative', className)} ref={containerRef}>
			{content}
		</div>
	)
}
