'use client'

import { ThemeToggle } from '@langdrift/web/ui/theme-toggle'
import { useTranslations } from 'next-intl'
import { useEffect, useState } from 'react'

export function ThemeControl() {
	const t = useTranslations('ui')
	const [theme, setTheme] = useState<'light' | 'dark'>('dark')
	useEffect(() => {
		try {
			const saved = localStorage.getItem('langdrift-console-theme')
			if (saved === 'dark' || saved === 'light') setTheme(saved)
		} catch {
			/* The default remains usable when browser storage is blocked. */
		}
	}, [])
	useEffect(() => {
		document.documentElement.classList.toggle('dark', theme === 'dark')
		document.documentElement.style.colorScheme = theme
	}, [theme])

	return (
		<ThemeToggle
			theme={theme}
			label={t(theme === 'dark' ? 'switchLight' : 'switchDark')}
			onToggle={() => {
				const next = theme === 'light' ? 'dark' : 'light'
				setTheme(next)
				try {
					localStorage.setItem('langdrift-console-theme', next)
				} catch {
					/* Keep the visit preference. */
				}
			}}
		/>
	)
}
