'use client'

import { cn } from '@langdrift/web/utilities'
import { AnimatePresence, motion, useReducedMotion } from 'motion/react'
import type { ElementType, ReactNode } from 'react'
import { useMemo } from 'react'

// Adapted from SmoothUI's Text Morph (smoothui.dev): morphs one string into
// another by diffing tokens. Shared tokens keep their key and slide to their
// new position, removed tokens shrink out and added tokens pop in.

export type TextMorphMode = 'characters' | 'words'
export type TextMorphAlign = 'start' | 'center' | 'end'
export type TextMorphTag = 'span' | 'div' | 'p' | 'h1' | 'h2' | 'h3' | 'h4'

export interface TextMorphProps {
	/** Justifies the token row when the wrapper is wider than the text. */
	align?: TextMorphAlign
	/** Element rendered for the wrapper. */
	as?: TextMorphTag
	className?: string
	/** Spring duration for each token, in milliseconds. */
	duration?: number
	/** Splits the diff on graphemes or whole words. */
	mode?: TextMorphMode
	/** Per-token stagger, in milliseconds. Capped so long strings stay snappy. */
	stagger?: number
	/**
	 * Fixed content after the morphing string, plain text or an element of its
	 * own. It is not diffed: it glides as one piece when the string before it
	 * changes width.
	 */
	suffix?: ReactNode
	/** The current string. Changing it triggers the morph. */
	text: string
}

interface TextMorphToken {
	key: string
	value: string
}

const DEFAULT_DURATION = 250
const DEFAULT_STAGGER = 16
const MS = 1000
const POP_SCALE = 0.72
const POP_OFFSET = 6
const MAX_STAGGER_MS = 260
const EXIT_DURATION_S = 0.15
const EASE_OUT = [0.23, 1, 0.32, 1] as const

const ALIGN_CLASS: Record<TextMorphAlign, string> = {
	center: 'justify-center',
	end: 'justify-end',
	start: 'justify-start',
}

function segmentGraphemes(value: string): string[] {
	if (typeof Intl === 'undefined' || !('Segmenter' in Intl)) return [...value]
	const segmenter = new Intl.Segmenter(undefined, {
		granularity: 'grapheme',
	})

	return Array.from(segmenter.segment(value), (entry) => entry.segment)
}

function tokenizeText(text: string, mode: TextMorphMode): string[] {
	if (mode === 'characters') return segmentGraphemes(text)

	return text.split(/(\s+)/).filter((token) => token.length > 0)
}

function buildTokens(text: string, mode: TextMorphMode): TextMorphToken[] {
	const occurrences = new Map<string, number>()

	return tokenizeText(text, mode).map((value) => {
		const count = occurrences.get(value) ?? 0
		occurrences.set(value, count + 1)

		return {
			key: `${value}-${count}`,
			value,
		}
	})
}

// The token row is positioned on purpose: AnimatePresence in popLayout mode
// takes exiting tokens out of the flow with `position: absolute`, so without a
// positioned parent they would resolve against the page and fly off.
export function TextMorph(props: TextMorphProps) {
	const {
		align = 'start',
		as = 'span',
		className,
		duration = DEFAULT_DURATION,
		mode = 'characters',
		stagger = DEFAULT_STAGGER,
		suffix,
		text,
	} = props

	const shouldReduceMotion = useReducedMotion()
	const tokens = useMemo(() => buildTokens(text, mode), [text, mode])
	const Tag = as as ElementType

	if (shouldReduceMotion) {
		return (
			<Tag className={className}>
				{text}
				{suffix ? ' ' : null}
				{suffix}
			</Tag>
		)
	}

	return (
		<Tag className={className}>
			{/* The intact string is visually hidden text, so screen readers read
			the word while every animated token stays aria-hidden. */}
			<span className="sr-only">
				{text}
				{suffix ? ' ' : null}
				{suffix}
			</span>
			<motion.span
				className={cn('relative inline-flex flex-wrap', ALIGN_CLASS[align])}
				layout
				transition={{
					bounce: 0.1,
					duration: duration / MS,
					type: 'spring',
				}}
			>
				<AnimatePresence initial={false} mode="popLayout">
					{tokens.map((token, index) => (
						<motion.span
							animate={{
								opacity: 1,
								scale: 1,
								y: 0,
							}}
							aria-hidden="true"
							exit={{
								opacity: 0,
								scale: POP_SCALE,
								transition: {
									duration: EXIT_DURATION_S,
									ease: EASE_OUT,
								},
								y: -POP_OFFSET,
							}}
							initial={{
								opacity: 0,
								scale: POP_SCALE,
								y: POP_OFFSET,
							}}
							key={token.key}
							layout
							style={{
								display: 'inline-block',
								whiteSpace: 'pre',
							}}
							transition={{
								bounce: 0.1,
								delay: Math.min(index * stagger, MAX_STAGGER_MS) / MS,
								duration: duration / MS,
								type: 'spring',
							}}
						>
							{token.value}
						</motion.span>
					))}
					{suffix ? (
						<motion.span
							aria-hidden="true"
							key="suffix"
							layout="position"
							style={{
								display: 'inline-block',
								whiteSpace: 'pre',
							}}
							transition={{
								bounce: 0.1,
								duration: duration / MS,
								type: 'spring',
							}}
						>
							{' '}
							{suffix}
						</motion.span>
					) : null}
				</AnimatePresence>
			</motion.span>
		</Tag>
	)
}
