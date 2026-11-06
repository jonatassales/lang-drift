// The bolt is drawn in a 100 by 100 box. It starts in the box's middle band,
// steps left in kinks, and never leaves the box's height.
const segments = 6
const startTop = 30
const startBand = 40
const largestKink = 40
const largestSway = 8
const topEdge = 2
const bottomEdge = 98

// One jagged bolt leaping from the word's edge at x 100 toward the left.
// `reach`, from 0 to 1, sets how far it travels; `random` supplies the jitter,
// so a test can fix it.
export function lightningBolt(random: () => number, reach: number): string {
	const startY = startTop + random() * startBand
	const span = reach * 100
	const points = [`M100 ${startY.toFixed(1)}`]

	let boltY = startY

	for (let step = 1; step <= segments; step += 1) {
		const boltX =
			100 - (span * step) / segments + (random() - 0.5) * largestSway

		boltY = Math.min(
			bottomEdge,
			Math.max(topEdge, boltY + (random() - 0.5) * largestKink),
		)
		points.push(`L${boltX.toFixed(1)} ${boltY.toFixed(1)}`)
	}

	return points.join(' ')
}
