'use client'

import { useLocale, useTranslations } from 'next-intl'
import type { MouseEvent } from 'react'
import { useEffect, useState } from 'react'

import languageLabels from '@i18n/messages/language-labels.json'
import { getPathname, usePathname } from '@i18n/navigation'
import {
	LanguageSwitcher,
	type LanguageSwitcherLinkProps,
	type LanguageSwitcherOption,
} from 'web/ui/language-switcher'

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

	const [urlSuffix, setUrlSuffix] = useState('')
	const pathname = usePathname()
	const locale = useLocale()
	const translate = useTranslations('header')

	// Sync the active locale to local storage for returning visitors.
	useEffect(() => {
		try {
			window.localStorage.setItem('langdrift-website-locale', locale)
		} catch {
			// Locale routing does not depend on browser storage.
		}
	}, [locale])

	// Listen for hash and popstate changes to preserve query and hash in language links.
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

	function handleLanguageClick(event: MouseEvent<HTMLAnchorElement>) {
		if (
			event.button !== 0 ||
			event.metaKey ||
			event.ctrlKey ||
			event.shiftKey ||
			event.altKey
		) {
			return
		}

		const nextLocale = event.currentTarget.dataset.locale as
			| (typeof languages)[number]['locale']
			| undefined

		if (!nextLocale) {
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

	function renderLanguageLink(
		language: LanguageSwitcherOption & { href: string },
		linkProps: LanguageSwitcherLinkProps,
	) {
		return (
			<a
				{...linkProps}
				data-locale={language.locale}
				href={`${getPathname({
					href: pathname,
					locale: language.locale,
				}).replace(/\/?$/, '/')}${urlSuffix}`}
				onClick={handleLanguageClick}
			/>
		)
	}

	const options = languages.map((language) => ({
		href: `${pathname}${urlSuffix}`,
		label: languageLabels[language.locale],
		locale: language.locale,
		name: translate(`languages.${language.name}`),
	}))

	return (
		<LanguageSwitcher
			className={
				fromDrawer
					? '[&_a]:min-h-11 [&_button]:min-h-11 [&_button]:min-w-11'
					: undefined
			}
			currentLocale={locale}
			label={translate('language')}
			options={options}
			renderLink={renderLanguageLink}
		/>
	)
}
