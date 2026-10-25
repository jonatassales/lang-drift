'use client'

import { ProductVisionCurve } from '@langdrift/web/ui/product-vision-curve'
import { useAtom, useAtomValue } from 'jotai'
import { useLocale, useTranslations } from 'next-intl'
import { useEffect, useMemo, useRef } from 'react'

import { heroSelectionAtom, selectedPointAtom } from '@state'

import { useHeroPoints } from './hero-chart/use-hero-points.hook'

export function HeroChart() {
	const t = useTranslations('demo')
	const locale = useLocale()
	const points = useHeroPoints()
	const [selectedPoint, setSelectedPoint] = useAtom(selectedPointAtom)
	const selection = useAtomValue(heroSelectionAtom)
	const previousSelection = useRef(selection)
	const frameRef = useRef<HTMLSpanElement>(null)
	const passage = useRef(false)
	const selectedEventId = (points[selectedPoint]?.event ?? points.at(-1)?.event)
		?.id
	const percent = new Intl.NumberFormat(locale, { style: 'percent' })
	const number = useMemo(() => new Intl.NumberFormat(locale), [locale])
	const delta = new Intl.NumberFormat(locale, { signDisplay: 'exceptZero' })

	useEffect(() => {
		const previous = previousSelection.current
		previousSelection.current = selection
		const frame = frameRef.current
		const motionPreference = window.matchMedia(
			'(prefers-reduced-motion: reduce)',
		)
		if (
			!frame ||
			selection.source !== 'graph' ||
			selection.pointIndex === previous.pointIndex ||
			motionPreference.matches
		) {
			return
		}

		// Alternate finite CSS animations to replace an in-flight passage without
		// remounting the chart, measuring layout or scheduling a frame loop.
		function stopInterruptedPassage() {
			if (frame && motionPreference.matches) delete frame.dataset.passage
		}
		frame.addEventListener('animationcancel', stopInterruptedPassage)
		passage.current = !passage.current
		frame.dataset.passage = passage.current ? 'a' : 'b'

		return () => {
			frame.removeEventListener('animationcancel', stopInterruptedPassage)
			delete frame.dataset.passage
		}
	}, [selection])

	return (
		<div className="hero-chart-card">
			<span
				aria-hidden="true"
				className="hero-chart-chroma"
				onAnimationEnd={(event) => {
					delete event.currentTarget.dataset.passage
				}}
				ref={frameRef}
			/>
			<div className="hero-chart-heading">
				<div className="hero-chart-score">
					<span>{t('chart.title')}</span>
					<strong>{percent.format(0.73)}</strong>
				</div>
				<div className="hero-chart-movement">
					<span>{t('chart.movement')}</span>
					<strong>
						{t('chart.movementDetail', { baseline: percent.format(0.91) })}
					</strong>
				</div>
			</div>
			<ProductVisionCurve
				data={points}
				formatNumber={number.format}
				labels={{
					classifications: {
						baseline: t('classifications.baseline'),
						intentional: t('classifications.intentional'),
						review: t('classifications.review'),
						unexplained: t('classifications.unexplained'),
					},
					describeEvent: (event) =>
						t('chart.event', {
							title: event.title,
							delta: delta.format(event.delta),
							classification: t(`classifications.${event.classification}`),
							date: event.date,
						}),
					product: t('chart.product'),
					summary: t('chart.summary', {
						from: percent.format(0.91),
						to: percent.format(0.73),
					}),
					why: t('chart.why'),
				}}
				onSelectEvent={(event) => {
					const index = points.findIndex(
						(point) => point.event?.id === event.id,
					)
					if (index >= 0) setSelectedPoint(index)
				}}
				pointPadding={16}
				selectedEventId={selectedEventId}
			/>
		</div>
	)
}
