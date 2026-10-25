'use client'

import {
	Bar,
	BarChart,
	CartesianGrid,
	Cell,
	ResponsiveContainer,
	Tooltip,
	XAxis,
	YAxis,
} from 'recharts'

export type DistributionClass = 'Expected' | 'Unexpected' | 'unclassified'
export interface DistributionTarget {
	targetId: string
	label: string
	Expected: number
	Unexpected: number
	unclassified: number
}

const series: DistributionClass[] = ['Expected', 'Unexpected', 'unclassified']
const colors = {
	Expected: 'var(--review-expected, #4f79ba)',
	Unexpected: 'var(--review-unexpected, #9272ba)',
	unclassified: 'var(--review-unclassified, #97979f)',
}

/** Counts stay stable while the parent selects a subset of events. */
export function TargetDistributionChart({
	data,
	labels,
	selectedTarget,
	selectedClass,
	onSelect,
	locale,
}: {
	data: DistributionTarget[]
	labels: Record<DistributionClass, string>
	selectedTarget: string | null
	selectedClass: DistributionClass | null
	onSelect: (targetId: string, classification: DistributionClass) => void
	locale: string
}) {
	const number = new Intl.NumberFormat(locale)
	const max = Math.max(
		1,
		...data.map((row) => row.Expected + row.Unexpected + row.unclassified),
	)
	return (
		<div
			className="review-distribution-chart"
			style={{
				height: 270,
				minWidth: 0,
			}}
		>
			<ResponsiveContainer
				width="100%"
				height="100%"
				minWidth={0}
				initialDimension={{
					width: 360,
					height: 270,
				}}
			>
				<BarChart
					accessibilityLayer
					data={data}
					layout="vertical"
					margin={{
						top: 12,
						bottom: 8,
						left: 0,
						right: 14,
					}}
				>
					<CartesianGrid horizontal={false} stroke="var(--border)" />
					<XAxis
						type="number"
						domain={[0, max]}
						allowDecimals={false}
						ticks={Array.from({ length: max + 1 }, (_, i) => i)}
						tickFormatter={number.format}
						axisLine={false}
						tickLine={false}
						tick={{
							fill: 'var(--muted-foreground)',
							fontSize: 13,
						}}
					/>
					<YAxis
						type="category"
						dataKey="label"
						width={110}
						axisLine={false}
						tickLine={false}
						tick={{
							fill: 'var(--foreground)',
							fontSize: 13,
						}}
					/>
					<Tooltip
						cursor={{
							fill: 'var(--border)',
							fillOpacity: 0.3,
						}}
						isAnimationActive={false}
						formatter={(value, name) => [number.format(Number(value)), name]}
						contentStyle={{
							background: 'var(--card)',
							border: '1px solid var(--border)',
							color: 'var(--foreground)',
							borderRadius: 8,
							fontSize: 13,
						}}
					/>
					{series.map((classification) => (
						<Bar
							key={classification}
							dataKey={classification}
							name={labels[classification]}
							stackId="events"
							barSize={32}
							fill={colors[classification]}
							isAnimationActive={false}
							onClick={(_, index) => {
								const target = data[index]
								if (target) onSelect(target.targetId, classification)
							}}
						>
							{data.map((target) => (
								<Cell
									key={target.targetId}
									cursor="pointer"
									stroke={
										(selectedTarget !== null || selectedClass !== null) &&
										(selectedTarget === null ||
											selectedTarget === target.targetId) &&
										(!selectedClass || selectedClass === classification)
											? 'var(--primary)'
											: 'transparent'
									}
									strokeWidth={2}
								/>
							))}
						</Bar>
					))}
				</BarChart>
			</ResponsiveContainer>
		</div>
	)
}
