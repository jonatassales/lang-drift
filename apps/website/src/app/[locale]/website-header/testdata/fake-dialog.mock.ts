import { type Mock, vi } from 'vitest'

export interface FakeDialog {
	close: Mock<() => void>
	showModal: Mock<() => void>
}

// Stands in for the drawer's dialog element: jsdom has no modal dialogs, and
// the hook only opens and closes it.
export function fakeDialog(): FakeDialog {
	return {
		close: vi.fn(),
		showModal: vi.fn(),
	}
}
