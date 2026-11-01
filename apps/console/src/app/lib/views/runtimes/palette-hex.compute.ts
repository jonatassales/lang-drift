/** Sinapsi accepts hex tokens; canvas resolves the browser's computed color space. */
export function paletteHex(color: string): string {
	const canvas = window.document.createElement('canvas')
	canvas.width = canvas.height = 1
	const context = canvas.getContext('2d')
	if (!context) return '#6d0ad6'
	context.fillStyle = color
	context.fillRect(0, 0, 1, 1)
	const rgba = context.getImageData(0, 0, 1, 1).data

	return `#${[rgba[0], rgba[1], rgba[2]].map((value) => value.toString(16).padStart(2, '0')).join('')}`
}
