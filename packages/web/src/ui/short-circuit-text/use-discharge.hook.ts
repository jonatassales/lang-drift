import { useReducedMotion } from 'motion/react'
import { type RefObject, useEffect, useState } from 'react'

import { jolt } from './jolt.data'
import { lightningBolt } from './lightning-bolt.compute'
import { tremble } from './tremble.data'

export interface Discharge {
	left: string[]
	right: string[]
}

const quiet: Discharge = {
	left: [],
	right: [],
}

const trembleMs = 900
const firstArcMs = 600
const shortestPauseMs = 450
const pauseSpreadMs = 1300
const redrawMs = 70
const dischargeMs = 210
const mostBoltsPerSide = 3
const shortestReach = 0.5

// The bolts on each side of a word in a short circuit, while the word hums in
// place and, at irregular moments, jolts as they arc. The motion runs on the
// element through the Web Animations API and on timers, outside React's
// render, which is why it needs an effect. Server and client both start
// quiet, so hydration matches. Under reduced motion, or where the browser
// cannot animate an element, nothing moves and no bolt is drawn; a hidden tab
// skips the discharges.
export function useDischarge(
	wordRef: RefObject<HTMLElement | null>,
): Discharge {
	const [discharge, setDischarge] = useState<Discharge>(quiet)
	const shouldReduceMotion = useReducedMotion()

	useEffect(() => {
		const word = wordRef.current
		const canArc = !shouldReduceMotion && typeof word?.animate === 'function'

		if (!word || !canArc) return

		const timers = new Set<number>()
		const hum = word.animate(tremble, {
			duration: trembleMs,
			iterations: Number.POSITIVE_INFINITY,
		})
		const jerk = word.animate(jolt, {
			duration: dischargeMs,
			easing: 'linear',
		})

		// The jolt waits, idle, until the first discharge replays it.
		jerk.cancel()

		function later(callback: () => void, delay: number) {
			function fire() {
				timers.delete(timer)
				callback()
			}

			const timer = window.setTimeout(fire, delay)

			timers.add(timer)
		}

		function drawBolt() {
			return lightningBolt(
				Math.random,
				shortestReach + Math.random() * (1 - shortestReach),
			)
		}

		function bolts() {
			const count = 1 + Math.floor(Math.random() * mostBoltsPerSide)

			return Array.from({ length: count }, drawBolt)
		}

		function spark() {
			setDischarge({
				left: bolts(),
				right: bolts(),
			})
		}

		function settle() {
			setDischarge(quiet)
		}

		function arc() {
			later(arc, shortestPauseMs + Math.random() * pauseSpreadMs)

			if (document.hidden) return

			jerk.play()
			spark()
			later(spark, redrawMs)
			later(spark, redrawMs * 2)
			later(settle, dischargeMs)
		}

		function stop() {
			hum.cancel()
			jerk.cancel()
			for (const timer of timers) window.clearTimeout(timer)
			settle()
		}

		later(arc, firstArcMs)

		return stop
	}, [shouldReduceMotion, wordRef])

	return discharge
}
