/** Schedule noncritical initialization after a paint, without keeping unmounted work alive. */
export function afterPaint(callback: () => void) {
	let secondFrame = 0
	let idle = 0
	let timer = 0
	const firstFrame = window.requestAnimationFrame(() => {
		secondFrame = window.requestAnimationFrame(() => {
			if (typeof window.requestIdleCallback === 'function') {
				idle = window.requestIdleCallback(callback, { timeout: 500 })
			} else {
				timer = window.setTimeout(callback, 0)
			}
		})
	})

	return () => {
		window.cancelAnimationFrame(firstFrame)
		window.cancelAnimationFrame(secondFrame)
		if (idle) window.cancelIdleCallback(idle)
		window.clearTimeout(timer)
	}
}
