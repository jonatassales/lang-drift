'use client'

import { useTranslations } from 'next-intl'
import { useEffect, useMemo, useRef, useState } from 'react'
import type { SinapsiElement, SinapsiGraphDocument } from 'sinapsi'
import type {} from 'sinapsi/react-types'

import type { Component, RuntimeWidgetTone } from '@data/types'
import { Skeleton } from '@ui/skeleton'

import { RuntimeNodeCards } from './runtime-node-cards'

interface RuntimeGraphProps {
	components: Component[]
	active: boolean
	tone: RuntimeWidgetTone
	links?: {
		from: string
		to: string
	}[]
}

export function RuntimeGraph(props: RuntimeGraphProps) {
	const { components, active, links, tone } = props

	const t = useTranslations('ui')
	const graph = useRef<SinapsiElement>(null)
	const [ready, setReady] = useState(false)
	const [failed, setFailed] = useState(false)
	const tones = {
		purple: 'text-primary',
		green: 'text-status-success',
		yellow: 'text-status-warning',
		orange: 'text-drift-zero',
		red: 'text-status-danger',
	}
	const documentKey = JSON.stringify({
		graph: components.map((component) => ({
			id: component.path,
			name: component.name,
			payload: {
				kind: component.kind,
			},
			links: (links ?? [])
				.filter(
					(link) =>
						link.from === component.path &&
						components.some((target) => target.path === link.to),
				)
				.map((link) => ({
					id: link.to,
					name:
						components.find((target) => target.path === link.to)?.name ??
						link.to,
				})),
			presentation: {
				type: 'card',
				title: component.name,
				description: component.kind,
				reference: component.path,
			},
		})),
	})
	const document = useMemo<SinapsiGraphDocument>(
		() => JSON.parse(documentKey),
		[documentKey],
	)
	useEffect(() => {
		let disposed = false
		void import('sinapsi/browser')
			.then(() => {
				if (!disposed) setReady(true)
			})
			.catch(() => {
				if (!disposed) setFailed(true)
			})

		return () => {
			disposed = true
		}
	}, [])
	useEffect(() => {
		const element = graph.current
		if (!element || !ready) return
		element.nodes = document
		element.activeNodeIds = document.graph.map((node) => node.id)
		// Public palette properties cannot resolve CSS variables; read owned tokens.
		const updatePalette = () => {
			const style = getComputedStyle(element)
			element.palette = {
				primary: paletteHex(style.color),
				muted: paletteHex(style.color),
				text: paletteHex(getComputedStyle(window.document.body).color),
			}
		}
		updatePalette()
		const mutation = new MutationObserver(updatePalette)
		mutation.observe(element, {
			attributes: true,
			attributeFilter: ['class'],
		})
		mutation.observe(window.document.documentElement, {
			attributes: true,
			attributeFilter: ['class'],
		})
		const reduced = matchMedia('(prefers-reduced-motion: reduce)')
		let visible = false
		const pause = () => {
			element.paused =
				!active ||
				!visible ||
				reduced.matches ||
				window.document.visibilityState !== 'visible'
		}
		const observer = new IntersectionObserver(([entry]) => {
			visible = entry.isIntersecting
			pause()
		})
		observer.observe(element)
		reduced.addEventListener('change', pause)
		window.document.addEventListener('visibilitychange', pause)
		pause()

		return () => {
			observer.disconnect()
			mutation.disconnect()
			reduced.removeEventListener('change', pause)
			window.document.removeEventListener('visibilitychange', pause)
			element.paused = true
		}
	}, [ready, document, active])

	return (
		<div className="relative rounded-lg border border-border bg-background/40">
			<sinaps-i
				ref={graph}
				className={`block h-[32rem] w-full ${tones[tone]}`}
				move="rotate"
				speed={active ? 1.2 : 0.35}
				paused
				aria-label={t('componentGraph')}
			/>
			<RuntimeNodeCards graph={graph} components={components} ready={ready} />
			{failed ? (
				<p
					role="status"
					className="absolute inset-x-4 bottom-4 text-xs text-muted-foreground"
				>
					{t('graphUnavailable')}
				</p>
			) : !ready ? (
				<div
					role="status"
					className="absolute inset-4 grid grid-cols-3 content-center gap-8"
					aria-busy="true"
					aria-label={t('loading')}
					data-skeleton="graph-skeleton"
				>
					{['first', 'second', 'third', 'fourth', 'fifth', 'sixth'].map(
						(node) => (
							<Skeleton key={node} className="h-16 w-full" />
						),
					)}
				</div>
			) : null}
		</div>
	)
}

/** Sinapsi accepts hex tokens; canvas resolves the browser's computed color space. */
function paletteHex(color: string): string {
	const canvas = window.document.createElement('canvas')
	canvas.width = canvas.height = 1
	const context = canvas.getContext('2d')
	if (!context) return '#6d0ad6'
	context.fillStyle = color
	context.fillRect(0, 0, 1, 1)
	const rgba = context.getImageData(0, 0, 1, 1).data

	return `#${[rgba[0], rgba[1], rgba[2]].map((value) => value.toString(16).padStart(2, '0')).join('')}`
}
