'use client'

import { useLocale, useTranslations } from 'next-intl'
import { useState, useTransition } from 'react'

import languages from '@i18n/messages/languages.json'
import { LanguageSwitcher } from 'web/ui/language-switcher'

import { setLocale } from './actions'
import { isLocale } from './config'
import { usePathname, useRouter } from './navigation'

interface LanguageSettingsProps {
	compact?: boolean
}

export function LanguageSettings(props: LanguageSettingsProps) {
	const { compact = false } = props

	const [pending, startTransition] = useTransition()
	const [failed, setFailed] = useState(false)

	const pathname = usePathname()
	const router = useRouter()

	const locale = useLocale()
	const translate = useTranslations('ui')

	function handleLocaleChange(next: string) {
		if (!isLocale(next) || next === locale) return
		setFailed(false)
		startTransition(async () => {
			try {
				const result = await setLocale(next)

				if (!result.ok) {
					setFailed(true)

					return
				}
				router.replace(
					`${pathname}${window.location.search}${window.location.hash}`,
					{
						locale: next,
					},
				)
			} catch {
				setFailed(true)
			}
		})
	}

	return (
		<div className={compact ? 'grid' : 'grid gap-2'}>
			<LanguageSwitcher
				label={translate('language')}
				currentLocale={locale}
				options={languages}
				disabled={pending}
				onLocaleChange={handleLocaleChange}
			/>
			<span
				aria-live="polite"
				className={compact ? 'sr-only' : 'text-xs text-muted-foreground'}
			>
				{pending ? translate('saving') : ''}
			</span>
			{failed ? (
				<p role="alert" className="text-xs text-status-danger">
					{translate('error')}
				</p>
			) : null}
		</div>
	)
}
