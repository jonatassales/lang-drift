'use client'

import { useLocale, useTranslations } from 'next-intl'
import type { MouseEvent } from 'react'
import { useEffect, useState } from 'react'

import languageLabels from '@i18n/messages/language-labels.json'
import { Link, usePathname } from '@i18n/navigation'
import {
	LanguageSwitcher,
	type LanguageSwitcherLinkProps,
	type LanguageSwitcherOption,
} from 'web/ui/language-switcher'

import { drawerFocus } from './drawer-focus.state'
import { persistLanguagePreference } from './persist-language-preference.adapter'

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

	// Remembers the active locale for returning visitors, after every switch.
	// A refused write keeps the visitor navigating: the proxy's cookie still
	// holds the choice, and routing never reads browser storage.
	useEffect(() => {
		persistLanguagePreference(locale)
	}, [locale])

	// Listen for hash and popstate changes to preserve query and hash in language links.
	useEffect(() => {
		function syncUrlSuffix() {
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
		const modifiedClick =
			event.button !== 0 ||
			event.metaKey ||
			event.ctrlKey ||
			event.shiftKey ||
			event.altKey

		if (modifiedClick) {
			return
		}

		const nextLocale = event.currentTarget.dataset.locale as
			| (typeof languages)[number]['locale']
			| undefined

		if (!nextLocale) {
			return
		}

		const switchFromDrawer = fromDrawer && nextLocale !== locale

		if (switchFromDrawer) drawerFocus.locale = nextLocale
	}

	function renderLanguageLink(
		language: LanguageSwitcherOption & { href: string },
		linkProps: LanguageSwitcherLinkProps,
	) {
		return (
			<Link
				{...linkProps}
				data-locale={language.locale}
				href={`${pathname.replace(/\/?$/, '/')}${urlSuffix}`}
				locale={language.locale}
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
