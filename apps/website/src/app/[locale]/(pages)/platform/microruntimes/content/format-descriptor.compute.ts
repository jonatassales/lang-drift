import type { runtimeDescriptor } from './microruntimes-content'

// Renders the descriptor as the indented JSON a reader would see in the
// file, with short arrays kept on one line.
export function formatDescriptor(descriptor: typeof runtimeDescriptor): string {
	return JSON.stringify(descriptor, null, 2).replace(
		/\[\n\s+([^\]{}]*?)\n\s+\]/g,
		(_match, items: string) => `[${items.replace(/,\n\s+/g, ', ')}]`,
	)
}
