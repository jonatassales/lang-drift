// @vitest-environment jsdom
import { act, cleanup, fireEvent, render, screen } from '@testing-library/react'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

import {
	type FakeDialog,
	fakeDialog,
} from '@layout/website-header/testdata/fake-dialog.mock'
import { NavigationDrawerHarness } from '@layout/website-header/testdata/navigation-drawer-harness.mock'
import {
	type DesktopQuery,
	stubDesktopQuery,
} from '@layout/website-header/testdata/stub-desktop-query.mock'

import { drawerFocus } from './drawer-focus'

const intl = vi.hoisted(() => ({
	locale: 'en',
}))

vi.mock('next-intl', () => ({
	useLocale: () => intl.locale,
}))

let dialog: FakeDialog
let desktop: DesktopQuery

beforeEach(() => {
	dialog = fakeDialog()
	desktop = stubDesktopQuery()
})

afterEach(() => {
	cleanup()
	vi.unstubAllGlobals()
	drawerFocus.locale = undefined
	document.body.style.overflow = ''
	intl.locale = 'en'
})

describe('website/website-header/use-navigation-drawer.hook', () => {
	it('opens the dialog as a modal, locks the page scroll and focuses the close button', () => {
		render(<NavigationDrawerHarness dialog={dialog} />)

		fireEvent.click(screen.getByText('Open navigation'))

		expect(screen.getByRole('status').textContent).toBe('open')
		expect(dialog.showModal).toHaveBeenCalledTimes(1)
		expect(document.body.style.overflow).toBe('hidden')
		expect(document.activeElement).toBe(screen.getByText('Close navigation'))
	})

	it('closes the dialog, restores the page scroll and returns focus to the trigger', () => {
		document.body.style.overflow = 'auto'
		render(<NavigationDrawerHarness dialog={dialog} />)

		fireEvent.click(screen.getByText('Open navigation'))
		fireEvent.click(screen.getByText('Close navigation'))

		expect(screen.getByRole('status').textContent).toBe('closed')
		expect(dialog.close).toHaveBeenCalledTimes(1)
		expect(document.body.style.overflow).toBe('auto')
		expect(document.activeElement).toBe(screen.getByText('Open navigation'))
	})

	it('closes when the locale changes and stays closed on the way back', () => {
		const { rerender } = render(<NavigationDrawerHarness dialog={dialog} />)

		fireEvent.click(screen.getByText('Open navigation'))
		intl.locale = 'pt-BR'
		rerender(<NavigationDrawerHarness dialog={dialog} />)

		expect(screen.getByRole('status').textContent).toBe('closed')
		expect(dialog.close).toHaveBeenCalledTimes(1)

		intl.locale = 'en'
		rerender(<NavigationDrawerHarness dialog={dialog} />)

		expect(screen.getByRole('status').textContent).toBe('closed')
		expect(dialog.showModal).toHaveBeenCalledTimes(1)
	})

	it('closes when the viewport grows to the desktop navigation', () => {
		render(<NavigationDrawerHarness dialog={dialog} />)

		fireEvent.click(screen.getByText('Open navigation'))
		act(() => desktop.grow())

		expect(screen.getByRole('status').textContent).toBe('closed')
		expect(dialog.close).toHaveBeenCalledTimes(1)
	})

	it('gives the trigger the focus a language switch from the drawer left', () => {
		drawerFocus.locale = 'pt-BR'
		intl.locale = 'pt-BR'

		render(<NavigationDrawerHarness dialog={dialog} />)

		expect(document.activeElement).toBe(screen.getByText('Open navigation'))
		expect(drawerFocus.locale).toBeUndefined()
	})

	it('leaves the focus alone after a switch made outside the drawer', () => {
		intl.locale = 'pt-BR'

		render(<NavigationDrawerHarness dialog={dialog} />)

		expect(document.activeElement).toBe(document.body)
	})

	it('drops the breakpoint subscription when it unmounts', () => {
		const { unmount } = render(<NavigationDrawerHarness dialog={dialog} />)

		expect(desktop.listeners.size).toBe(1)

		unmount()

		expect(desktop.listeners.size).toBe(0)
	})
})
