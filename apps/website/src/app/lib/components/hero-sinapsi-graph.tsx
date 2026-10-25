'use client'

import { useAtomValue, useSetAtom } from 'jotai'
import { useLocale, useTranslations } from 'next-intl'
import { useEffect, useId, useMemo, useRef, useState } from 'react'
import type { SinapsiElement } from 'sinapsi'
import type {} from 'sinapsi/react-types'

import {
	getHeroGraph,
	getHeroPointFromNodeClick,
	getHeroPoints,
} from '@data/hero-demo-data'
import { selectHeroGraphPointAtom, themeAtom } from '@state'

export function HeroSinapsiGraph() {
	const t = useTranslations('demo')
	const locale = useLocale()
	const theme = useAtomValue(themeAtom)
	const [reducedMotion, setReducedMotion] = useState(true)
	const setSelectedPoint = useSetAtom(selectHeroGraphPointAtom)
	const elementRef = useRef<SinapsiElement>(null)
	const [status, setStatus] = useState<'loading' | 'ready' | 'failed'>(
		'loading',
	)
	const descriptionId = useId()
	const headingId = useId()
	const graph = useMemo(() => getHeroGraph(t, locale), [t, locale])
	const points = useMemo(() => getHeroPoints(t, locale), [t, locale])

	useEffect(() => {
		const preference = window.matchMedia('(prefers-reduced-motion: reduce)')
		const update = () => setReducedMotion(preference.matches)
		update()
		preference.addEventListener('change', update)
		return () => preference.removeEventListener('change', update)
	}, [])

	useEffect(() => {
		const element = elementRef.current
		if (!element) return

		function selectMovement(event: Event) {
			if (!(event instanceof CustomEvent)) return
			const index = getHeroPointFromNodeClick(event.detail, graph, points)
			if (index !== undefined) setSelectedPoint(index)
		}

		element.addEventListener('sinapsi-node-click', selectMovement)

		return () =>
			element.removeEventListener('sinapsi-node-click', selectMovement)
	}, [graph, points, setSelectedPoint])

	useEffect(() => {
		let cancelled = false
		setStatus('loading')
		const timeout = window.setTimeout(() => {
			if (!cancelled) setStatus('failed')
		}, 15000)

		async function mountGraph() {
			try {
				// Browser registration is idempotent, including remounts and HMR.
				if (!customElements.get('sinaps-i')) await import('sinapsi/browser')
				await customElements.whenDefined('sinaps-i')
				if (cancelled || !elementRef.current) return
				elementRef.current.nodes = graph
				setStatus('ready')
			} catch {
				if (!cancelled) setStatus('failed')
			} finally {
				window.clearTimeout(timeout)
			}
		}

		void mountGraph()

		return () => {
			cancelled = true
			window.clearTimeout(timeout)
		}
	}, [graph])

	return (
		<section
			aria-describedby={descriptionId}
			aria-labelledby={headingId}
			className="hero-graph-region"
		>
			<div className="hero-graph-heading">
				<h2 id={headingId}>{t('graph.title')}</h2>
				<p>{t('graph.description')}</p>
			</div>
			<p className="ld-visually-hidden" id={descriptionId}>
				{t('graph.instructions')} {t('graph.summary')}
			</p>
			<div className="hero-graph-stage" data-status={status}>
				<sinaps-i
					className="hero-sinapsi"
					close-label={t('graph.closeLabel')}
					color-muted={theme === 'dark' ? '#A1A1AA' : '#6B6B75'}
					color-primary="#F97316"
					color-text={theme === 'dark' ? '#F5F5F5' : '#171717'}
					move={reducedMotion ? 'idle' : 'rotate'}
					ref={elementRef}
					speed={0.2}
				/>
				{status !== 'ready' ? (
					<output className="hero-graph-fallback">
						{t(status === 'failed' ? 'graph.fallback' : 'graph.loading')}
					</output>
				) : null}
			</div>
		</section>
	)
}
