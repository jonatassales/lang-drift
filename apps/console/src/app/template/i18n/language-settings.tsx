'use client'

import { LanguageSwitcher } from '@langdrift/web/ui/language-switcher'
import { useLocale, useTranslations } from 'next-intl'
import { useState, useTransition } from 'react'

import { setLocale } from './actions'
import { isLocale } from './config'
import languages from './messages/languages.json'
import { usePathname, useRouter } from './navigation'

interface LanguageSettingsProps {
	compact?: boolean
}

export function LanguageSettings(props: LanguageSettingsProps) {
	const { compact = false } = props

	const locale = useLocale()
	const translate = useTranslations('ui')
	const pathname = usePathname()
	const router = useRouter()
	const [pending, startTransition] = useTransition()
	const [failed, setFailed] = useState(false)

	return (
		<div className={compact ? 'grid' : 'grid gap-2'}>
			<LanguageSwitcher
				label={translate('language')}
				currentLocale={locale}
				options={languages}
				disabled={pending}
				onLocaleChange={(next) => {
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
				}}
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
