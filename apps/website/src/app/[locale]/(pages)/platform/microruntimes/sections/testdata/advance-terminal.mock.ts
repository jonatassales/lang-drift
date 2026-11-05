import { act } from '@testing-library/react'
import { vi } from 'vitest'

export function advanceTerminal(milliseconds: number) {
	for (let elapsed = 0; elapsed < milliseconds; elapsed += 40) {
		act(() => {
			vi.advanceTimersByTime(40)
		})
	}
}
