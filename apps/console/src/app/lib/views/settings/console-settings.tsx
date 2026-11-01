'use client'
import { Card } from '@langdrift/web/ui/card'
import { useTranslations } from 'next-intl'

import { PageHeading } from '@components/page-heading/page-heading'
import { LanguageSettings } from '@i18n/language-settings'

export function ConsoleSettings() {
	const t = useTranslations('ui')

	return (
		<>
			<PageHeading title={t('settings')} />
			<Card className="block max-w-2xl p-5">
				<h2 className="text-sm font-medium">{t('language')}</h2>
				<p className="mt-2 mb-4 text-xs text-muted-foreground">
					{t('languageDescription')}
				</p>
				<LanguageSettings />
			</Card>
		</>
	)
}
