'use client'

import { useTranslations } from 'next-intl'
import { useEffect, useState } from 'react'

import { ThemeToggle } from 'web/ui/theme-toggle'

export function ThemeControl() {
	const [theme, setTheme] = useState<'light' | 'dark'>('dark')

	const translate = useTranslations('ui')

	// Restore saved theme from local storage on mount.
	useEffect(() => {
		try {
			const saved = localStorage.getItem('langdrift-console-theme')
			if (saved === 'dark' || saved === 'light') setTheme(saved)
		} catch {
			/* The default remains usable when browser storage is blocked. */
		}
	}, [])

	// Synchronize active theme with document element classes and color scheme.
	useEffect(() => {
		document.documentElement.classList.toggle('dark', theme === 'dark')
		document.documentElement.style.colorScheme = theme
	}, [theme])

	function handleToggle() {
		const next = theme === 'light' ? 'dark' : 'light'
		setTheme(next)
		try {
			localStorage.setItem('langdrift-console-theme', next)
		} catch {
			/* Keep the visit preference. */
		}
	}

	return (
		<ThemeToggle
			theme={theme}
			label={translate(theme === 'dark' ? 'switchLight' : 'switchDark')}
			onToggle={handleToggle}
		/>
	)
}
