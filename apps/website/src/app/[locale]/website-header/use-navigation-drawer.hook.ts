import { useLocale } from 'next-intl'
import { type RefObject, useEffect, useRef, useState } from 'react'

import { drawerFocus } from './drawer-focus.state'

const desktopQuery = '(min-width: 64rem)'

export interface NavigationDrawer {
	closeDrawer: () => void
	closeRef: RefObject<HTMLButtonElement | null>
	drawerRef: RefObject<HTMLDialogElement | null>
	open: boolean
	openDrawer: () => void
	triggerRef: RefObject<HTMLButtonElement | null>
}

// The header's navigation drawer: a modal dialog that locks the page scroll
// while open, returns focus to its trigger when it closes and closes when the
// viewport grows to the desktop navigation. It stays open only in the locale
// it was opened in.
export function useNavigationDrawer(): NavigationDrawer {
	const [openLocale, setOpenLocale] = useState<string | null>(null)
	const drawerRef = useRef<HTMLDialogElement>(null)
	const triggerRef = useRef<HTMLButtonElement>(null)
	const closeRef = useRef<HTMLButtonElement>(null)

	const locale = useLocale()

	// A language switch from the drawer remounts the header under the new
	// locale; the replacement trigger takes the focus the old drawer held.
	useEffect(() => {
		if (drawerFocus.locale !== locale) return

		triggerRef.current?.focus({ preventScroll: true })
		drawerFocus.locale = undefined
	}, [locale])

	// The dialog element opens and closes imperatively: as a modal, with the
	// page scroll locked and the close button focused, until it closes.
	useEffect(() => {
		if (openLocale !== locale) return

		const drawer = drawerRef.current

		if (!drawer) return

		const previousOverflow = document.body.style.overflow

		drawer.showModal()
		closeRef.current?.focus()
		document.body.style.overflow = 'hidden'

		return () => {
			drawer.close()
			document.body.style.overflow = previousOverflow
			triggerRef.current?.focus()
		}
	}, [openLocale, locale])

	// Subscribes to the desktop breakpoint, where the drawer does not exist.
	useEffect(() => {
		const desktop = window.matchMedia(desktopQuery)

		function closeOnDesktop() {
			if (desktop.matches) setOpenLocale(null)
		}

		desktop.addEventListener('change', closeOnDesktop)

		return () => desktop.removeEventListener('change', closeOnDesktop)
	}, [])

	function openDrawer() {
		setOpenLocale(locale)
	}

	function closeDrawer() {
		setOpenLocale(null)
	}

	const staleLocale = openLocale !== null && openLocale !== locale

	if (staleLocale) setOpenLocale(null)

	return {
		closeDrawer,
		closeRef,
		drawerRef,
		open: openLocale === locale,
		openDrawer,
		triggerRef,
	}
}
