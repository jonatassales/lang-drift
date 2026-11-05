export function stubClipboard(writeText: (text: string) => Promise<void>) {
	Object.defineProperty(navigator, 'clipboard', {
		configurable: true,
		value: {
			writeText,
		},
	})
}
