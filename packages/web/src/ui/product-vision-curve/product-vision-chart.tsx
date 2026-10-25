'use client'

import { driftColor } from '@langdrift/web/utilities'
import { useEffect, useId, useMemo, useState } from 'react'
import {
	CartesianGrid,
	Line,
	LineChart,
	ResponsiveContainer,
	XAxis,
	YAxis,
} from 'recharts'

import type {
	DriftClassification,
	ProductVisionCurveProps,
} from './product-vision-curve'

interface VisionChartProps
	extends Pick<
		ProductVisionCurveProps,
		'data' | 'labels' | 'onSelectEvent' | 'selectedEventId'
	> {
	compact: boolean
	formatNumber: NonNullable<ProductVisionCurveProps['formatNumber']>
	pointPadding: number
	eventClassifications: Record<DriftClassification, string>
}

const classificationColors: Record<DriftClassification, string> = {
	baseline: 'var(--muted-foreground)',
	intentional: '#3b82f6',
	review: '#a855f7',
	unexplained: '#dc2626',
}

const dataDomain: [string, string] = ['dataMin - 4', 'dataMax + 4']
const tickStyle = { fill: 'var(--muted-foreground)' }

export function ProductVisionChart(props: VisionChartProps) {
	const {
		compact,
		data,
		formatNumber,
		labels,
		onSelectEvent,
		pointPadding,
		selectedEventId,
		eventClassifications,
	} = props

	const [reduceMotion, setReduceMotion] = useState(false)
	const instanceId = useId()
	// Stable axis settings prevent Recharts from remounting focused event markers
	// when selection changes without changing the plotted data.
	const margin = useMemo(
		() => ({
			bottom: 4,
			left: 0,
			right: 24,
			top: compact ? 16 : 54,
		}),
		[compact],
	)
	const padding = useMemo(() => ({ left: pointPadding }), [pointPadding])
	useEffect(() => {
		const query = window.matchMedia('(prefers-reduced-motion: reduce)')
		const update = () => setReduceMotion(query.matches)
		update()
		query.addEventListener('change', update)

		return () => query.removeEventListener('change', update)
	}, [])
	const formatDelta = (value: number) =>
		`${value > 0 ? '+' : ''}${formatNumber(value)}`

	return (
		<ResponsiveContainer height="100%" width="100%">
			{/* Event markers own keyboard navigation and labels; avoid a second chart tab stop. */}
			<LineChart
				accessibilityLayer={false}
				className="product-vision-surface"
				data={data}
				margin={margin}
			>
				<CartesianGrid stroke="var(--border)" vertical={false} />
				<XAxis
					axisLine={false}
					dataKey="label"
					fontSize={compact ? 10 : 12}
					padding={padding}
					tickLine={false}
					tick={tickStyle}
				/>
				<YAxis
					axisLine={false}
					domain={dataDomain}
					fontSize={compact ? 10 : 12}
					tickCount={4}
					tick={tickStyle}
					tickFormatter={formatNumber}
					tickLine={false}
					width={40}
				/>
				<Line
					dataKey="value"
					dot={(props) => {
						const index = props.index ?? 0
						const point = data[index]
						const event = point?.event
						const selected = event?.id === selectedEventId
						const cx = Number(props.cx ?? 0)
						const cy = Number(props.cy ?? 0)
						const color = event
							? classificationColors[event.classification]
							: 'var(--foreground)'

						if (!event) {
							return (
								<circle
									cx={cx}
									cy={cy}
									fill="var(--foreground)"
									key={`vision-dot-${String(index)}`}
									r={3}
									stroke="var(--card)"
									strokeWidth={2}
								/>
							)
						}

						const actor = event.actors[0]
						const activate = () => onSelectEvent?.(event)
						const avatarClipId = `vision-avatar-${instanceId}-${event.id}`

						return (
							// biome-ignore lint/a11y/useSemanticElements: Recharts dot markers render inside SVG and cannot contain an HTML button.
							<g
								aria-label={
									labels?.describeEvent?.(event) ??
									`${event.title}, ${event.delta} Product Vision, ${eventClassifications[event.classification]}, ${event.date}`
								}
								className="product-vision-event-dot focus-visible:[&>circle:first-of-type]:stroke-[var(--foreground)] focus-visible:[&>circle:first-of-type]:stroke-4"
								aria-pressed={selected}
								key={event.id}
								onClick={activate}
								onFocus={activate}
								onKeyDown={(keyboardEvent) => {
									if (
										keyboardEvent.key === 'Enter' ||
										keyboardEvent.key === ' '
									) {
										keyboardEvent.preventDefault()
										activate()
									}
								}}
								role="button"
								tabIndex={0}
							>
								<circle
									cx={cx}
									cy={cy}
									fill={selected ? '#f97316' : color}
									r={selected ? 7 : 5}
									stroke="var(--card)"
									strokeWidth={3}
								/>
								{!compact ? (
									<>
										<circle
											cx={cx}
											cy={cy - 24}
											fill="var(--foreground)"
											r={12}
											stroke="var(--card)"
											strokeWidth={2}
										/>
										{actor?.src ? (
											<>
												<defs>
													<clipPath
														id={avatarClipId}
														clipPathUnits="userSpaceOnUse"
													>
														<circle cx={cx} cy={cy - 24} r={10} />
													</clipPath>
												</defs>
												<image
													clipPath={`url(#${avatarClipId})`}
													height="20"
													href={actor.src}
													preserveAspectRatio="xMidYMid slice"
													width="20"
													x={cx - 10}
													y={cy - 34}
												/>
											</>
										) : (
											<text
												fill="var(--card)"
												fontSize="7"
												fontWeight="700"
												textAnchor="middle"
												x={cx}
												y={cy - 21.5}
											>
												{actor?.initials ?? 'LD'}
											</text>
										)}
										<text
											fill={driftColor(event.delta)}
											fontSize="11"
											textAnchor="middle"
											x={cx}
											y={cy - 41}
										>
											{formatDelta(event.delta)}
										</text>
									</>
								) : null}
							</g>
						)
					}}
					isAnimationActive={!compact && !reduceMotion}
					stroke="var(--foreground)"
					strokeWidth={2}
					type="monotone"
				/>
			</LineChart>
		</ResponsiveContainer>
	)
}
