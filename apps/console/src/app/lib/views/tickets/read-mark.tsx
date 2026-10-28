import { useTranslations } from 'next-intl'

// The mark beside an artifact the agent has read: read, not approved.
export function ReadMark() {
	const t = useTranslations('tickets')

	return (
		<span className="inline-flex items-center font-mono text-[12px] text-status-success">
			<span aria-hidden="true">✓</span>
			<span className="sr-only">{t('readMark')}</span>
		</span>
	)
}
