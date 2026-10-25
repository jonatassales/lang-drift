import { readPalette, type ShaderPalette } from './palette'
import { fragmentSource, vertexSource } from './shaders'

export interface SurfaceOptions {
	paused: boolean
	reducedMotion: boolean
}

export interface SurfaceController {
	dispose(): void
	recolor(): void
	update(options: SurfaceOptions): void
}

// The vec3 uniforms the fragment shader declares, by palette key.
const colorUniforms = [
	'background',
	'surface',
	'primary',
	'primaryHover',
	'shadeDeep',
	'shadeSoft',
	'glow',
	'violet',
] as const

// About 1920×1080 pixels, so a large hero stays cheap on dense screens.
const pixelBudget = 2073600

function compileShader(
	gl: WebGL2RenderingContext,
	type: number,
	source: string,
) {
	const shader = gl.createShader(type)
	if (!shader) return null
	gl.shaderSource(shader, source)
	gl.compileShader(shader)
	if (!gl.getShaderParameter(shader, gl.COMPILE_STATUS)) {
		gl.deleteShader(shader)

		return null
	}

	return shader
}

function createProgram(gl: WebGL2RenderingContext) {
	const vertex = compileShader(gl, gl.VERTEX_SHADER, vertexSource)
	const fragment = compileShader(gl, gl.FRAGMENT_SHADER, fragmentSource)
	const program = gl.createProgram()
	if (!vertex || !fragment || !program) {
		gl.deleteShader(vertex)
		gl.deleteShader(fragment)
		gl.deleteProgram(program)

		return null
	}
	gl.attachShader(program, vertex)
	gl.attachShader(program, fragment)
	gl.linkProgram(program)
	gl.deleteShader(vertex)
	gl.deleteShader(fragment)
	if (!gl.getProgramParameter(program, gl.LINK_STATUS)) {
		gl.deleteProgram(program)

		return null
	}

	return program
}

// Draws the scene on the canvas until disposed. Returns nothing when the page
// gives no palette, the browser has no WebGL2 or the program does not build;
// the caller's fallback then stays visible.
export function mountSurface(
	canvas: HTMLCanvasElement,
	initial: SurfaceOptions,
): SurfaceController | undefined {
	const palette = readPalette(canvas)
	if (!palette) return
	const context = canvas.getContext('webgl2', {
		alpha: true,
		antialias: false,
		depth: false,
		stencil: false,
		powerPreference: 'low-power',
	})
	if (!context) return
	const gl: WebGL2RenderingContext = context
	const program = createProgram(gl)
	if (!program) return

	gl.useProgram(program)
	const colorLocations = new Map(
		colorUniforms.map((name) => [name, gl.getUniformLocation(program, name)]),
	)
	const schemeLocation = gl.getUniformLocation(program, 'scheme')
	const resolutionLocation = gl.getUniformLocation(program, 'resolution')
	const elapsedLocation = gl.getUniformLocation(program, 'elapsed')
	const motionPreference = matchMedia('(prefers-reduced-motion: reduce)')
	let options = initial
	let frame = 0
	let previousTime: number | null = null
	let activeTime = 0
	let visible = true
	let failed = false
	let disposed = false
	let hasSize = false

	function uploadPalette(next: ShaderPalette) {
		for (const name of colorUniforms) {
			gl.uniform3fv(colorLocations.get(name) ?? null, next[name])
		}
		gl.uniform1f(schemeLocation, next.scheme)
	}

	function stop() {
		cancelAnimationFrame(frame)
		frame = 0
		previousTime = null
	}

	function canAnimate() {
		const alive = !disposed && !failed && hasSize
		const shown = visible && !document.hidden
		const moving =
			!motionPreference.matches && !options.paused && !options.reducedMotion

		return alive && shown && moving
	}

	function draw() {
		if (disposed || failed || !hasSize) return
		gl.uniform1f(elapsedLocation, activeTime)
		gl.drawArrays(gl.TRIANGLES, 0, 3)
	}

	function animate(timestamp: number) {
		frame = 0
		if (!canAnimate()) {
			previousTime = null

			return
		}
		if (previousTime !== null) {
			activeTime += Math.min((timestamp - previousTime) / 1000, 0.05)
		}
		previousTime = timestamp
		draw()
		frame = requestAnimationFrame(animate)
	}

	// Runs the loop when it may run; otherwise leaves one still frame on screen.
	function reconcile() {
		stop()
		if (canAnimate()) {
			frame = requestAnimationFrame(animate)

			return
		}
		if (!document.hidden && visible) draw()
	}

	function resize() {
		if (disposed || failed) return
		const bounds = canvas.getBoundingClientRect()
		hasSize = bounds.width > 0 && bounds.height > 0
		if (!hasSize) {
			stop()

			return
		}
		const pixelRatio = Math.min(window.devicePixelRatio || 1, 2)
		const scale = Math.min(
			pixelRatio,
			Math.sqrt(pixelBudget / (bounds.width * bounds.height)),
		)
		canvas.width = Math.max(1, Math.floor(bounds.width * scale))
		canvas.height = Math.max(1, Math.floor(bounds.height * scale))
		gl.viewport(0, 0, canvas.width, canvas.height)
		gl.uniform2f(resolutionLocation, canvas.width, canvas.height)
		draw()
		reconcile()
	}

	function contextLost(event: Event) {
		event.preventDefault()
		failed = true
		stop()
		canvas.dataset.unavailable = ''
	}

	// The page's theme toggles `.dark` on <html>; the scene follows it.
	function recolor() {
		if (disposed || failed) return
		const next = readPalette(canvas)
		if (!next) return
		uploadPalette(next)
		draw()
	}

	const resizeObserver = new ResizeObserver(resize)
	const intersectionObserver = new IntersectionObserver(([entry]) => {
		visible = entry?.isIntersecting ?? false
		reconcile()
	})
	const themeObserver = new MutationObserver(recolor)
	delete canvas.dataset.unavailable
	uploadPalette(palette)
	resizeObserver.observe(canvas)
	intersectionObserver.observe(canvas)
	themeObserver.observe(document.documentElement, {
		attributes: true,
		attributeFilter: ['class'],
	})
	motionPreference.addEventListener('change', reconcile)
	document.addEventListener('visibilitychange', reconcile)
	window.addEventListener('resize', resize)
	canvas.addEventListener('webglcontextlost', contextLost)
	resize()

	return {
		dispose() {
			disposed = true
			stop()
			resizeObserver.disconnect()
			intersectionObserver.disconnect()
			themeObserver.disconnect()
			motionPreference.removeEventListener('change', reconcile)
			document.removeEventListener('visibilitychange', reconcile)
			window.removeEventListener('resize', resize)
			canvas.removeEventListener('webglcontextlost', contextLost)
			gl.useProgram(null)
			gl.deleteProgram(program)
			gl.clearColor(0, 0, 0, 0)
			gl.clear(gl.COLOR_BUFFER_BIT)
			canvas.width = 1
			canvas.height = 1
		},
		recolor,
		update(next) {
			options = next
			reconcile()
		},
	}
}
