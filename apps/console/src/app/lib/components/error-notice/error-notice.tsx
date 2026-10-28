'use client'
import { useTranslations } from 'next-intl'

import { ConsoleError } from '@data/api'
import { Link } from '@i18n/navigation'
import { Button } from '@ui/button'

interface ErrorNoticeProps {
	error: unknown
	retry?: () => void
}
export function ErrorNotice(props: ErrorNoticeProps) {
	const { error, retry } = props

	const t = useTranslations('ui')
	const messages = useTranslations('errors')
	const code = error instanceof ConsoleError ? error.code : 'operation_failed'

	return (
		<div
			role="alert"
			className="flex flex-wrap items-center gap-3 rounded-lg border border-status-danger/30 bg-status-danger/5 p-4 text-sm text-status-danger"
		>
			<span>{messages.has(code) ? messages(code) : t('error')}</span>
			{code === 'authentication_required' ? (
				<Link href="/login" className="underline">
					{t('login')}
				</Link>
			) : retry ? (
				<Button onClick={retry}>{t('retry')}</Button>
			) : null}
		</div>
	)
}
