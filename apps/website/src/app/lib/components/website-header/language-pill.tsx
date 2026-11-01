'use client'

import { LanguageSwitcher } from '@langdrift/web/ui/language-switcher'
import { useLocale, useTranslations } from 'next-intl'
import type { MouseEvent } from 'react'
import { useEffect, useState } from 'react'

import languageLabels from '@i18n/messages/language-labels.json'
import { getPathname, usePathname } from '@i18n/navigation'

import { drawerFocus } from './drawer-focus'

const languages = [
	{
		locale: 'en',
		name: 'english',
	},
	{
		locale: 'pt-BR',
		name: 'portuguese',
	},
	{
		locale: 'zh-Hant',
		name: 'chinese',
	},
	{
		locale: 'ja',
		name: 'japanese',
	},
	{
		locale: 'vi',
		name: 'vietnamese',
	},
] as const

interface LanguagePillProps {
	fromDrawer?: boolean
}

export function LanguagePill(props: LanguagePillProps) {
	const { fromDrawer = false } = props

	const locale = useLocale()
	const translate = useTranslations('header')
	const pathname = usePathname()
	const [urlSuffix, setUrlSuffix] = useState('')

	useEffect(() => {
		try {
			window.localStorage.setItem('langdrift-website-locale', locale)
		} catch {
			// Locale routing does not depend on browser storage.
		}
	}, [locale])

	useEffect(() => {
		const syncUrlSuffix = () => {
			setUrlSuffix(`${window.location.search}${window.location.hash}`)
		}
		syncUrlSuffix()
		window.addEventListener('hashchange', syncUrlSuffix)
		window.addEventListener('popstate', syncUrlSuffix)

		return () => {
			window.removeEventListener('hashchange', syncUrlSuffix)
			window.removeEventListener('popstate', syncUrlSuffix)
		}
	}, [])

	function selectLanguage(
		event: MouseEvent<HTMLAnchorElement>,
		nextLocale: (typeof languages)[number]['locale'],
	) {
		if (
			event.button !== 0 ||
			event.metaKey ||
			event.ctrlKey ||
			event.shiftKey ||
			event.altKey
		) {
			return
		}

		try {
			window.localStorage.setItem('langdrift-website-locale', nextLocale)
		} catch {
			// Language navigation works even when browser storage is unavailable.
		}
		if (fromDrawer && nextLocale !== locale) {
			drawerFocus.locale = nextLocale
		}
	}

	return (
		<LanguageSwitcher
			className={
				fromDrawer
					? '[&_a]:min-h-11 [&_button]:min-h-11 [&_button]:min-w-11'
					: undefined
			}
			currentLocale={locale}
			label={translate('language')}
			options={languages.map((language) => ({
				locale: language.locale,
				label: languageLabels[language.locale],
				name: translate(`languages.${language.name}`),
				href: `${pathname}${urlSuffix}`,
			}))}
			renderLink={(language, linkProps) => (
				<a
					{...linkProps}
					href={`${getPathname({
						href: pathname,
						locale: language.locale,
					}).replace(/\/?$/, '/')}${urlSuffix}`}
					onClick={(event) => selectLanguage(event, language.locale)}
				/>
			)}
		/>
	)
}
