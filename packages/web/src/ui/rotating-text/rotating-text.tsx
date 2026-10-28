'use client'

import { TextMorph } from '@langdrift/web/vendors/smoothui/text-morph'
import { useReducedMotion } from 'motion/react'
import { useEffect, useState } from 'react'

export interface RotatingTextProps {
	className?: string
	/** Time each word stays before the next one morphs in, in milliseconds. */
	interval?: number
	/** Fixed text after the word, gliding as the word changes width. */
	suffix?: string
	/** Classes for the suffix alone. */
	suffixClassName?: string
	words: string[]
}

const DEFAULT_INTERVAL = 2600

// Cycles through `words`, morphing each into the next. Under reduced motion it
// holds the first word, and it waits while the tab is hidden. Server and client both start on the first word, so
// hydration matches.
export function RotatingText(props: RotatingTextProps) {
	const {
		className,
		interval = DEFAULT_INTERVAL,
		suffix,
		suffixClassName,
		words,
	} = props

	const shouldReduceMotion = useReducedMotion()
	const [index, setIndex] = useState(0)

	useEffect(() => {
		const canRotate = !shouldReduceMotion && words.length > 1
		if (!canRotate) return
		// A hidden tab pauses the morph's animations, so the word waits too.
		const timer = window.setInterval(() => {
			if (document.hidden) return
			setIndex((current) => (current + 1) % words.length)
		}, interval)

		return () => window.clearInterval(timer)
	}, [interval, shouldReduceMotion, words.length])

	return (
		<TextMorph
			className={className}
			suffix={suffix}
			suffixClassName={suffixClassName}
			text={words[index] ?? ''}
		/>
	)
}
