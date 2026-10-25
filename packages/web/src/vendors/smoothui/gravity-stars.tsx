'use client'

import { cn } from '@langdrift/web/utilities'
import { useReducedMotion } from 'motion/react'
import type { CSSProperties, ReactNode } from 'react'
import { useEffect, useRef, useState } from 'react'

// Adapted from SmoothUI's Gravity Stars (smoothui.dev): a starfield on a 2D
// canvas with a pointer gravity well, a still frame under reduced motion and
// colours read from the page, so a token such as `--foreground` follows the
// theme.

export interface GravityStarsProps {
	/** Foreground content, always painted above the decorative starfield. */
	children?: ReactNode
	/** Extra classes for the positioned wrapper element. */
	className?: string
	/**
	 * Base star colour. Accepts any CSS colour, a `var(--token)` expression, or a
	 * bare custom property name such as `--color-foreground`.
	 */
	color?: string
	/** Draws thin, distance-faded links between stars that are close together. */
	connect?: boolean
	/** Maximum distance in CSS pixels at which two stars are linked. */
	connectDistance?: number
	/** Number of stars in the pool. Reduced automatically on narrow viewports. */
	count?: number
	/** Velocity damping per frame, from 0 (frictionless) to 1 (frozen). */
	drag?: number
	/** Halo radius multiplier for the glow sprite. 1 is a bare dot, 8 is a bloom. */
	glow?: number
	/** Strength of the pointer gravity well. */
	gravity?: number
	/** Freezes the animation on the current frame. */
	paused?: boolean
	/** Ambient drift speed multiplier. 1 is the calibrated default. */
	speed?: number
	/** Base star radius in CSS pixels, before the depth multiplier. */
	starSize?: number
	/**
	 * How far individual stars are tinted away from `color` towards a cool
	 * blue-white and a warm amber, from 0 (one flat colour) to 1 (full spread).
	 * A small amount is what makes the field read as a sky rather than grey dots.
	 */
	tint?: number
	/** Brightness oscillation amount, from 0 (steady) to 1 (strong). */
	twinkle?: number
}

type Rgb = [number, number, number]

interface Palette {
	base: Rgb
	cool: Rgb
	warm: Rgb
}

interface StarSettings {
	connect: boolean
	connectDistance: number
	count: number
	drag: number
	glow: number
	gravity: number
	palette: Palette
	speed: number
	starSize: number
	still: boolean
	tint: number
	twinkle: number
}

interface StarController {
	destroy: () => void
	render: () => void
	resize: () => void
	setPointer: (nextX: number, nextY: number, active: boolean) => void
	setRunning: (running: boolean) => void
	setSettings: (settings: StarSettings) => void
}

// Without a 2D canvas, two tiled dot grids in the star colour stand in.
const fallbackField =
	'pointer-events-none absolute inset-0 opacity-60 bg-[radial-gradient(circle_at_30%_30%,var(--star)_0px,transparent_1.5px),radial-gradient(circle_at_70%_65%,var(--star)_0px,transparent_1px)] bg-size-[96px_96px,53px_53px]'

const MAX_DPR = 2
const MS_PER_SECOND = 1000
const TAU = Math.PI * 2

const DEFAULT_COUNT = 220
const DEFAULT_STAR_SIZE = 1.5
const DEFAULT_GRAVITY = 1
const DEFAULT_DRAG = 0.03
const DEFAULT_CONNECT_DISTANCE = 130
const DEFAULT_TWINKLE = 0.55
const DEFAULT_GLOW = 5.5
const DEFAULT_SPEED = 1
const DEFAULT_TINT = 0.5
const DEFAULT_COLOR = 'var(--foreground, oklch(0.93 0.004 265))'

/**
 * Stellar tint endpoints. Real skies are near-neutral with a slight hue
 * spread — hot stars lean blue-white, cool ones amber. Lightness is held
 * constant across both so the spread reads as colour, never as brightness, and
 * each hue gets a comparable fraction of its own sRGB chroma ceiling.
 */
const TINT_COOL = 'oklch(0.92 0.028 255)'
const TINT_WARM = 'oklch(0.92 0.055 75)'

/** oklch(0.93 0.004 265) — the neutral used when the probe cannot resolve. */
const FALLBACK_BASE: Rgb = [232, 233, 238]
/** oklch(0.92 0.028 255) */
const FALLBACK_COOL: Rgb = [221, 233, 247]
/** oklch(0.92 0.055 75) */
const FALLBACK_WARM: Rgb = [243, 226, 194]

const FALLBACK_PALETTE: Palette = {
	base: FALLBACK_BASE,
	cool: FALLBACK_COOL,
	warm: FALLBACK_WARM,
}

const MAX_COUNT = 600
const MIN_COUNT = 8
const SMALL_VIEWPORT_WIDTH = 640
const SMALL_VIEWPORT_RATIO = 0.6

/** Gravity well tuning. Softening keeps the inverse square finite at r → 0. */
const POINTER_RANGE = 260
const POINTER_ACCEL = 4_000_000
const SOFTENING_SQ = 676
const MAX_ACCEL = 5000
const MAX_SPEED = 1200
const SWIRL = 0.85
const DRIFT_SPEED = 12

const MAX_DELTA = 0.05
const REFERENCE_FPS = 60
const MIN_GRID_CELL = 24
const MAX_LINKS_PER_STAR = 8

/** Depth shaping: `u ** DEPTH_BIAS` skews the pool towards small, far stars. */
const DEPTH_BIAS = 1.7
const NEAR_SIZE_GAIN = 1.9
const FAR_SIZE_FLOOR = 0.42
const NEAR_ALPHA_GAIN = 0.82
const FAR_ALPHA_FLOOR = 0.18
const NEAR_HALO_GAIN = 0.7
const FAR_HALO_FLOOR = 0.6
const NEAR_REACH_GAIN = 0.8
const FAR_REACH_FLOOR = 0.6
const NEAR_PULL_GAIN = 1
const FAR_PULL_FLOOR = 0.35
const MIN_VISIBLE_ALPHA = 0.01

const LINE_LEVELS = 6
const LINE_BASE_ALPHA = 0.3
const LINE_MIN_WIDTH = 0.3
const LINE_WIDTH_STEP = 0.12
const LINE_CLOSENESS_CURVE = 1.5
const LINE_DEPTH_FLOOR = 0.25
const LINE_DEPTH_GAIN = 0.75

const TWINKLE_RATE = 1.15
const TWINKLE_RATE_SPREAD = 0.8
const TWINKLE_DEPTH = 0.5
const SIZE_JITTER = 0.5

const SPRITE_PX = 128
const SPRITE_LEVELS = 5
const SEED = 987_654_321
const LEHMER_MULTIPLIER = 16_807
const LEHMER_MODULUS = 2_147_483_647

/**
 * Half the neighbourhood as flat `(dCol, dRow)` pairs, so every pair of cells is
 * compared exactly once and the hot loop allocates nothing.
 */
const FORWARD_NEIGHBOURS = [1, 0, -1, 1, 0, 1, 1, 1]

const clamp = (value: number, min: number, max: number) =>
	Math.max(min, Math.min(max, value))

const mix = (from: number, to: number, amount: number) =>
	from + (to - from) * amount

const toCssColor = (input: string) =>
	input.trim().startsWith('--') ? `var(${input.trim()})` : input.trim()

/**
 * Reads real computed colours out of the document, so tokens such as
 * `--color-foreground` follow the active theme instead of being duplicated in a
 * parallel JS palette. A hidden probe span resolves the cascade, then a 1×1
 * canvas converts whatever notation came back into channels.
 */
const resolveCssColors = (inputs: string[], host: HTMLElement): Rgb[] => {
	const probe = document.createElement('span')
	probe.style.display = 'none'
	host.append(probe)

	const canvas = document.createElement('canvas')
	canvas.width = 1
	canvas.height = 1
	const context = canvas.getContext('2d', { willReadFrequently: true })

	const fallbacks = [FALLBACK_BASE, FALLBACK_COOL, FALLBACK_WARM]
	const resolved = inputs.map((input, index) => {
		const fallback = fallbacks[index] ?? FALLBACK_BASE
		if (!context) {
			return fallback
		}
		probe.style.color = ''
		probe.style.color = toCssColor(input)
		const computed = window.getComputedStyle(probe).color
		if (!computed) {
			return fallback
		}
		context.clearRect(0, 0, 1, 1)
		context.fillStyle = '#000000'
		context.fillStyle = computed
		context.fillRect(0, 0, 1, 1)
		const { data } = context.getImageData(0, 0, 1, 1)

		return [data[0], data[1], data[2]] as Rgb
	})

	probe.remove()

	return resolved
}

/**
 * Deterministic Lehmer PRNG so the seeded layout is identical on every render
 * and on both sides of hydration.
 */
const createRandom = (seed: number) => {
	let state = seed % LEHMER_MODULUS || 1

	return () => {
		state = (state * LEHMER_MULTIPLIER) % LEHMER_MODULUS

		return state / LEHMER_MODULUS
	}
}

/**
 * Radial gradient sprite: a hot core dissolving into a wide, low-alpha halo.
 * Built once per tint and blitted per star, which is far cheaper — and far
 * softer — than filling a path or paying `shadowBlur` for every particle. The
 * falloff is deliberately wide and faint so it reads astronomical, not neon.
 */
const paintSprite = (sprite: HTMLCanvasElement, color: Rgb) => {
	const context = sprite.getContext('2d')
	if (!context) {
		return
	}
	const [red, green, blue] = color
	const center = SPRITE_PX / 2
	context.clearRect(0, 0, SPRITE_PX, SPRITE_PX)
	const gradient = context.createRadialGradient(
		center,
		center,
		0,
		center,
		center,
		center,
	)
	gradient.addColorStop(0, `rgba(${red}, ${green}, ${blue}, 1)`)
	gradient.addColorStop(0.06, `rgba(${red}, ${green}, ${blue}, 0.95)`)
	gradient.addColorStop(0.14, `rgba(${red}, ${green}, ${blue}, 0.5)`)
	gradient.addColorStop(0.28, `rgba(${red}, ${green}, ${blue}, 0.16)`)
	gradient.addColorStop(0.55, `rgba(${red}, ${green}, ${blue}, 0.04)`)
	gradient.addColorStop(1, `rgba(${red}, ${green}, ${blue}, 0)`)
	context.fillStyle = gradient
	context.fillRect(0, 0, SPRITE_PX, SPRITE_PX)
}

/**
 * One sprite per tint step, from the cool end to the warm end. `tint` controls
 * how far each step travels away from the base colour, so 0 collapses the whole
 * ramp back onto a single flat hue.
 */
const paintSprites = (
	sprites: HTMLCanvasElement[],
	palette: Palette,
	tint: number,
) => {
	const last = Math.max(sprites.length - 1, 1)
	for (let level = 0; level < sprites.length; level++) {
		const ramp = level / last
		const color: Rgb = [
			Math.round(
				mix(palette.base[0], mix(palette.cool[0], palette.warm[0], ramp), tint),
			),
			Math.round(
				mix(palette.base[1], mix(palette.cool[1], palette.warm[1], ramp), tint),
			),
			Math.round(
				mix(palette.base[2], mix(palette.cool[2], palette.warm[2], ramp), tint),
			),
		]
		paintSprite(sprites[level], color)
	}
}

const samePalette = (first: Palette, second: Palette) =>
	first.base[0] === second.base[0] &&
	first.base[1] === second.base[1] &&
	first.base[2] === second.base[2] &&
	first.cool[0] === second.cool[0] &&
	first.warm[0] === second.warm[0]

const createStarController = (
	canvas: HTMLCanvasElement,
): StarController | null => {
	const context = canvas.getContext('2d')
	if (!context) {
		return null
	}

	const sprites: HTMLCanvasElement[] = []
	for (let level = 0; level < SPRITE_LEVELS; level++) {
		const sprite = document.createElement('canvas')
		sprite.width = SPRITE_PX
		sprite.height = SPRITE_PX
		sprites.push(sprite)
	}

	const posX = new Float32Array(MAX_COUNT)
	const posY = new Float32Array(MAX_COUNT)
	const velX = new Float32Array(MAX_COUNT)
	const velY = new Float32Array(MAX_COUNT)
	const driftX = new Float32Array(MAX_COUNT)
	const driftY = new Float32Array(MAX_COUNT)
	const depth = new Float32Array(MAX_COUNT)
	const phase = new Float32Array(MAX_COUNT)
	const rate = new Float32Array(MAX_COUNT)
	const jitter = new Float32Array(MAX_COUNT)
	const seedX = new Float32Array(MAX_COUNT)
	const seedY = new Float32Array(MAX_COUNT)
	const tintLevel = new Uint8Array(MAX_COUNT)

	const random = createRandom(SEED)
	for (let star = 0; star < MAX_COUNT; star++) {
		seedX[star] = random()
		seedY[star] = random()
		const starDepth = random() ** DEPTH_BIAS
		depth[star] = starDepth
		phase[star] = random() * TAU
		rate[star] = TWINKLE_RATE * (1 + (random() - 0.5) * TWINKLE_RATE_SPREAD)
		jitter[star] = 1 + (random() - 0.5) * SIZE_JITTER
		// Averaging two draws gives a centre-weighted spread, so most stars sit
		// near the neutral base and only a few reach the blue or amber ends.
		const hue = (random() + random()) / 2
		tintLevel[star] = Math.min(
			SPRITE_LEVELS - 1,
			Math.floor(hue * SPRITE_LEVELS),
		)
		const angle = random() * TAU
		// Nearer stars drift faster: the parallax that sells the depth.
		const magnitude = DRIFT_SPEED * (0.15 + starDepth)
		driftX[star] = Math.cos(angle) * magnitude
		driftY[star] = Math.sin(angle) * magnitude
	}

	const maxSegments = MAX_COUNT * MAX_LINKS_PER_STAR
	const segments = new Float32Array(maxSegments * 4)
	const segmentLevel = new Uint8Array(maxSegments)

	let cellStart = new Int32Array(1)
	let cellCursor = new Int32Array(1)
	let cellItems = new Int32Array(MAX_COUNT)

	const startedAt = performance.now()
	let settings: StarSettings = {
		connect: true,
		connectDistance: DEFAULT_CONNECT_DISTANCE,
		count: DEFAULT_COUNT,
		drag: DEFAULT_DRAG,
		glow: DEFAULT_GLOW,
		gravity: DEFAULT_GRAVITY,
		palette: FALLBACK_PALETTE,
		speed: DEFAULT_SPEED,
		starSize: DEFAULT_STAR_SIZE,
		still: false,
		tint: DEFAULT_TINT,
		twinkle: DEFAULT_TWINKLE,
	}
	paintSprites(sprites, settings.palette, settings.tint)

	let width = 1
	let height = 1
	let activeCount = 0
	let seeded = false
	let pointerX = 0
	let pointerY = 0
	let pointerActive = false
	let lastFrameAt = startedAt
	let frame = 0
	let running = false
	let destroyed = false

	const effectiveCount = () => {
		const base = clamp(Math.round(settings.count), MIN_COUNT, MAX_COUNT)

		return width < SMALL_VIEWPORT_WIDTH
			? Math.max(MIN_COUNT, Math.round(base * SMALL_VIEWPORT_RATIO))
			: base
	}

	const seedPositions = (next: number) => {
		for (let star = 0; star < next; star++) {
			posX[star] = seedX[star] * width
			posY[star] = seedY[star] * height
			velX[star] = 0
			velY[star] = 0
		}
		activeCount = next
		seeded = true
	}

	const resize = () => {
		const rect = canvas.getBoundingClientRect()
		const dpr = Math.min(window.devicePixelRatio || 1, MAX_DPR)
		const previousWidth = width
		const previousHeight = height
		width = Math.max(1, rect.width)
		height = Math.max(1, rect.height)
		const pixelWidth = Math.max(1, Math.round(width * dpr))
		const pixelHeight = Math.max(1, Math.round(height * dpr))
		if (canvas.width !== pixelWidth || canvas.height !== pixelHeight) {
			canvas.width = pixelWidth
			canvas.height = pixelHeight
		}
		context.setTransform(dpr, 0, 0, dpr, 0, 0)

		const next = effectiveCount()
		if (!seeded || next !== activeCount) {
			seedPositions(next)

			return
		}
		const scaleX = width / previousWidth
		const scaleY = height / previousHeight
		for (let star = 0; star < activeCount; star++) {
			posX[star] *= scaleX
			posY[star] *= scaleY
		}
	}

	const gridGeometry = () => {
		const cell = Math.max(settings.connectDistance, MIN_GRID_CELL)

		return {
			cell,
			cols: Math.max(1, Math.ceil(width / cell)),
			rows: Math.max(1, Math.ceil(height / cell)),
		}
	}

	/** Counting sort of the particle pool into a uniform grid. */
	const sortIntoGrid = (cell: number, cols: number, rows: number) => {
		const total = cols * rows
		if (cellStart.length < total + 1) {
			cellStart = new Int32Array(total + 1)
			cellCursor = new Int32Array(total + 1)
		}
		cellStart.fill(0, 0, total + 1)
		if (cellItems.length < activeCount) {
			cellItems = new Int32Array(activeCount)
		}

		for (let star = 0; star < activeCount; star++) {
			const col = clamp(Math.floor(posX[star] / cell), 0, cols - 1)
			const row = clamp(Math.floor(posY[star] / cell), 0, rows - 1)
			cellStart[row * cols + col + 1]++
		}
		for (let cellIndex = 0; cellIndex < total; cellIndex++) {
			cellStart[cellIndex + 1] += cellStart[cellIndex]
			cellCursor[cellIndex] = cellStart[cellIndex]
		}
		for (let star = 0; star < activeCount; star++) {
			const col = clamp(Math.floor(posX[star] / cell), 0, cols - 1)
			const row = clamp(Math.floor(posY[star] / cell), 0, rows - 1)
			const index = row * cols + col
			cellItems[cellCursor[index]] = star
			cellCursor[index]++
		}
	}

	/**
	 * Softened inverse-square attraction with a tangential term, so stars fall
	 * in, swirl past the pointer and are flung back out carrying momentum.
	 * Nearer stars feel a wider, stronger well than distant ones.
	 */
	const applyPointer = (star: number, out: Float32Array) => {
		const starDepth = depth[star]
		const reach =
			POINTER_RANGE * (FAR_REACH_FLOOR + starDepth * NEAR_REACH_GAIN)
		const dx = pointerX - posX[star]
		const dy = pointerY - posY[star]
		const distanceSq = dx * dx + dy * dy
		if (distanceSq > reach * reach) {
			return
		}
		const distance = Math.sqrt(distanceSq)
		const softened = Math.sqrt(distanceSq + SOFTENING_SQ)
		const fade = 1 - distance / reach
		const pull = FAR_PULL_FLOOR + starDepth * NEAR_PULL_GAIN
		const strength =
			Math.min(
				(settings.gravity * POINTER_ACCEL * pull) / (softened * softened),
				MAX_ACCEL,
			) *
			fade *
			fade
		const nx = dx / softened
		const ny = dy / softened
		out[0] += nx * strength - ny * strength * SWIRL
		out[1] += ny * strength + nx * strength * SWIRL
	}

	const accel = new Float32Array(2)

	const step = (delta: number) => {
		const damping = (1 - settings.drag) ** (delta * REFERENCE_FPS)
		const { speed } = settings

		for (let star = 0; star < activeCount; star++) {
			accel[0] = 0
			accel[1] = 0
			if (pointerActive) {
				applyPointer(star, accel)
			}

			let vx = (velX[star] + accel[0] * delta) * damping
			let vy = (velY[star] + accel[1] * delta) * damping
			const speedSq = vx * vx + vy * vy
			if (speedSq > MAX_SPEED * MAX_SPEED) {
				const scale = MAX_SPEED / Math.sqrt(speedSq)
				vx *= scale
				vy *= scale
			}
			velX[star] = vx
			velY[star] = vy

			// Ambient drift is added outside the damped term, so a flung star bleeds
			// its momentum away and settles back into the calm parallax layer.
			posX[star] += (vx + driftX[star] * speed) * delta
			posY[star] += (vy + driftY[star] * speed) * delta

			if (posX[star] < 0) {
				posX[star] += width
			} else if (posX[star] > width) {
				posX[star] -= width
			}
			if (posY[star] < 0) {
				posY[star] += height
			} else if (posY[star] > height) {
				posY[star] -= height
			}
		}
	}

	const buildSegments = (cols: number, rows: number) => {
		const limit = settings.connectDistance * settings.connectDistance
		let segmentTotal = 0

		const consider = (first: number, second: number) => {
			if (segmentTotal >= maxSegments) {
				return
			}
			const dx = posX[first] - posX[second]
			const dy = posY[first] - posY[second]
			const d2 = dx * dx + dy * dy
			if (d2 > limit || d2 === 0) {
				return
			}
			// Alpha falls off with distance and with how far back the endpoints sit,
			// so the constellation reads as depth rather than a web of grey sticks.
			const closeness =
				(1 - Math.sqrt(d2) / settings.connectDistance) ** LINE_CLOSENESS_CURVE
			const depthWeight =
				LINE_DEPTH_FLOOR +
				((depth[first] + depth[second]) / 2) * LINE_DEPTH_GAIN
			const level = Math.min(
				LINE_LEVELS - 1,
				Math.floor(closeness * depthWeight * LINE_LEVELS),
			)
			const offset = segmentTotal * 4
			segments[offset] = posX[first]
			segments[offset + 1] = posY[first]
			segments[offset + 2] = posX[second]
			segments[offset + 3] = posY[second]
			segmentLevel[segmentTotal] = level
			segmentTotal++
		}

		// Only forward neighbours are visited, so every pair is tested once.
		for (let row = 0; row < rows; row++) {
			for (let col = 0; col < cols; col++) {
				const index = row * cols + col
				const start = cellStart[index]
				const end = cellStart[index + 1]
				for (let first = start; first < end; first++) {
					const star = cellItems[first]
					for (let second = first + 1; second < end; second++) {
						consider(star, cellItems[second])
					}
					for (
						let neighbour = 0;
						neighbour < FORWARD_NEIGHBOURS.length;
						neighbour += 2
					) {
						const nc = col + FORWARD_NEIGHBOURS[neighbour]
						const nr = row + FORWARD_NEIGHBOURS[neighbour + 1]
						if (nc < 0 || nc >= cols || nr >= rows) {
							continue
						}
						const nIndex = nr * cols + nc
						const nStart = cellStart[nIndex]
						const nEnd = cellStart[nIndex + 1]
						for (let second = nStart; second < nEnd; second++) {
							consider(star, cellItems[second])
						}
					}
				}
			}
		}

		return segmentTotal
	}

	/**
	 * One `beginPath`/`stroke` pair per alpha bucket. Stroking each pair
	 * individually would blow the frame budget on its own.
	 */
	const drawLinks = (segmentTotal: number) => {
		const [red, green, blue] = settings.palette.base
		for (let level = 0; level < LINE_LEVELS; level++) {
			const alpha = (LINE_BASE_ALPHA * (level + 1)) / LINE_LEVELS
			context.strokeStyle = `rgba(${red}, ${green}, ${blue}, ${alpha})`
			context.lineWidth = LINE_MIN_WIDTH + level * LINE_WIDTH_STEP
			context.beginPath()
			let drew = false
			for (let segmentIndex = 0; segmentIndex < segmentTotal; segmentIndex++) {
				if (segmentLevel[segmentIndex] !== level) {
					continue
				}
				const offset = segmentIndex * 4
				context.moveTo(segments[offset], segments[offset + 1])
				context.lineTo(segments[offset + 2], segments[offset + 3])
				drew = true
			}
			if (drew) {
				context.stroke()
			}
		}
	}

	const drawStars = () => {
		const elapsed = (performance.now() - startedAt) / MS_PER_SECOND
		const { glow, starSize, still, twinkle } = settings

		for (let star = 0; star < activeCount; star++) {
			const starDepth = depth[star]
			// Each star carries its own phase and rate, so the field never pulses in
			// unison. When still, the wave is frozen at the seeded phase.
			const wave = still
				? Math.sin(phase[star])
				: Math.sin(elapsed * rate[star] + phase[star])
			const brightness =
				(FAR_ALPHA_FLOOR + starDepth * NEAR_ALPHA_GAIN) *
				(1 - twinkle * TWINKLE_DEPTH * (1 - wave))
			const alpha = clamp(brightness, 0, 1)
			if (alpha <= MIN_VISIBLE_ALPHA) {
				continue
			}
			const core =
				starSize * (FAR_SIZE_FLOOR + starDepth * NEAR_SIZE_GAIN) * jitter[star]
			const halo = core * glow * (FAR_HALO_FLOOR + starDepth * NEAR_HALO_GAIN)
			context.globalAlpha = alpha
			context.drawImage(
				sprites[tintLevel[star]],
				posX[star] - halo,
				posY[star] - halo,
				halo * 2,
				halo * 2,
			)
		}
		context.globalAlpha = 1
	}

	const draw = (segmentTotal: number) => {
		context.clearRect(0, 0, width, height)
		if (settings.connect && segmentTotal > 0) {
			drawLinks(segmentTotal)
		}
		drawStars()
	}

	const renderStatic = () => {
		const { cell, cols, rows } = gridGeometry()
		sortIntoGrid(cell, cols, rows)
		const segmentTotal = settings.connect ? buildSegments(cols, rows) : 0
		draw(segmentTotal)
	}

	const tick = () => {
		if (destroyed || !running) {
			return
		}
		const now = performance.now()
		const delta = Math.min((now - lastFrameAt) / MS_PER_SECOND, MAX_DELTA)
		lastFrameAt = now
		step(delta)
		let segmentTotal = 0
		if (settings.connect) {
			const { cell, cols, rows } = gridGeometry()
			sortIntoGrid(cell, cols, rows)
			segmentTotal = buildSegments(cols, rows)
		}
		draw(segmentTotal)
		frame = requestAnimationFrame(tick)
	}

	resize()

	return {
		destroy: () => {
			destroyed = true
			running = false
			cancelAnimationFrame(frame)
			context.setTransform(1, 0, 0, 1, 0, 0)
			context.clearRect(0, 0, canvas.width, canvas.height)
			canvas.width = 0
			canvas.height = 0
			for (const sprite of sprites) {
				sprite.width = 0
				sprite.height = 0
			}
		},
		render: () => {
			resize()
			renderStatic()
		},
		resize: () => {
			resize()
			if (!running) {
				renderStatic()
			}
		},
		setPointer: (nextX: number, nextY: number, active: boolean) => {
			pointerX = nextX
			pointerY = nextY
			pointerActive = active
		},
		setRunning: (next: boolean) => {
			if (destroyed || running === next) {
				return
			}
			running = next
			if (next) {
				lastFrameAt = performance.now()
				frame = requestAnimationFrame(tick)
			} else {
				cancelAnimationFrame(frame)
			}
		},
		setSettings: (next: StarSettings) => {
			const changedCount = next.count !== settings.count
			const changedSprites =
				next.tint !== settings.tint ||
				!samePalette(next.palette, settings.palette)
			settings = next
			if (changedSprites) {
				paintSprites(sprites, next.palette, next.tint)
			}
			if (changedCount) {
				seedPositions(effectiveCount())
			}
		},
	}
}

export function GravityStars(props: GravityStarsProps) {
	const {
		children,
		className,
		color = DEFAULT_COLOR,
		connect = true,
		connectDistance = DEFAULT_CONNECT_DISTANCE,
		count = DEFAULT_COUNT,
		drag = DEFAULT_DRAG,
		glow = DEFAULT_GLOW,
		gravity = DEFAULT_GRAVITY,
		paused = false,
		speed = DEFAULT_SPEED,
		starSize = DEFAULT_STAR_SIZE,
		tint = DEFAULT_TINT,
		twinkle = DEFAULT_TWINKLE,
	} = props

	const shouldReduceMotion = useReducedMotion()
	const hostRef = useRef<HTMLDivElement>(null)
	const canvasRef = useRef<HTMLCanvasElement>(null)
	const controllerRef = useRef<StarController | null>(null)
	const [isSupported, setIsSupported] = useState(true)
	const [isActive, setIsActive] = useState(false)
	const [palette, setPalette] = useState<Palette>(FALLBACK_PALETTE)

	useEffect(() => {
		const host = hostRef.current
		if (!host) {
			return
		}
		const resolve = () => {
			const [base, cool, warm] = resolveCssColors(
				[color, TINT_COOL, TINT_WARM],
				host,
			)
			setPalette({
				base: base ?? FALLBACK_BASE,
				cool: cool ?? FALLBACK_COOL,
				warm: warm ?? FALLBACK_WARM,
			})
		}
		resolve()
		// A theme switch changes what the colour token resolves to; the page
		// toggles its scheme through the class on <html>.
		const observer = new MutationObserver(resolve)
		observer.observe(document.documentElement, {
			attributeFilter: ['class'],
			attributes: true,
		})

		return () => observer.disconnect()
	}, [color])

	useEffect(() => {
		const canvas = canvasRef.current
		if (!canvas) {
			return
		}
		const controller = createStarController(canvas)
		controllerRef.current = controller
		setIsSupported(controller !== null)
		if (!controller) {
			return
		}

		const observer = new ResizeObserver(() => controller.resize())
		observer.observe(canvas)

		return () => {
			observer.disconnect()
			controller.destroy()
			controllerRef.current = null
		}
	}, [])

	useEffect(() => {
		const host = hostRef.current
		if (!(host && isSupported)) {
			return
		}

		let isOnScreen = true
		const sync = () => {
			setIsActive(isOnScreen && document.visibilityState === 'visible')
		}

		const observer = new IntersectionObserver((entries) => {
			isOnScreen = entries.some((entry) => entry.isIntersecting)
			sync()
		})
		observer.observe(host)
		document.addEventListener('visibilitychange', sync)
		sync()

		return () => {
			observer.disconnect()
			document.removeEventListener('visibilitychange', sync)
		}
	}, [isSupported])

	useEffect(() => {
		const controller = controllerRef.current
		if (!controller) {
			return
		}
		controller.setSettings({
			connect,
			connectDistance: Math.max(connectDistance, 1),
			count: clamp(Math.round(count), MIN_COUNT, MAX_COUNT),
			drag: clamp(drag, 0, 1),
			glow: Math.max(glow, 1),
			gravity: Math.max(gravity, 0),
			palette,
			speed: Math.max(speed, 0),
			starSize: Math.max(starSize, 0.2),
			still: shouldReduceMotion === true,
			tint: clamp(tint, 0, 1),
			twinkle: clamp(twinkle, 0, 1),
		})
		controller.render()
	}, [
		connect,
		connectDistance,
		count,
		drag,
		glow,
		gravity,
		palette,
		shouldReduceMotion,
		speed,
		starSize,
		tint,
		twinkle,
	])

	useEffect(() => {
		const controller = controllerRef.current
		if (!controller) {
			return
		}
		const shouldRun = isActive && !paused && !shouldReduceMotion
		controller.setRunning(shouldRun)
		if (!shouldRun) {
			controller.render()
		}
	}, [isActive, paused, shouldReduceMotion])

	useEffect(() => {
		const host = hostRef.current
		const canvas = canvasRef.current
		if (!(host && canvas) || shouldReduceMotion) {
			return
		}

		// Pointer events, not hover: the well has to answer a dragged finger as
		// well as a mouse.
		const handleMove = (event: PointerEvent) => {
			const rect = canvas.getBoundingClientRect()
			controllerRef.current?.setPointer(
				event.clientX - rect.left,
				event.clientY - rect.top,
				true,
			)
		}
		const handleLeave = () => {
			controllerRef.current?.setPointer(0, 0, false)
		}

		host.addEventListener('pointermove', handleMove, { passive: true })
		host.addEventListener('pointerleave', handleLeave)
		host.addEventListener('pointercancel', handleLeave)

		return () => {
			host.removeEventListener('pointermove', handleMove)
			host.removeEventListener('pointerleave', handleLeave)
			host.removeEventListener('pointercancel', handleLeave)
		}
	}, [shouldReduceMotion])

	const fallbackColor = toCssColor(color)

	return (
		<div
			className={cn('relative isolate overflow-hidden', className)}
			ref={hostRef}
		>
			{isSupported ? (
				<div
					aria-hidden="true"
					className="pointer-events-none absolute inset-0"
				>
					<canvas className="h-full w-full" ref={canvasRef} />
				</div>
			) : (
				<div
					aria-hidden="true"
					className={fallbackField}
					style={{ '--star': fallbackColor } as CSSProperties}
				/>
			)}
			<div className="relative z-10 h-full">{children}</div>
		</div>
	)
}
