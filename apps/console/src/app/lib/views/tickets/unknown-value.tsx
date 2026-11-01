import { useTranslations } from 'next-intl'

interface UnknownValueProps {
	className?: string
}

// A value a ticket does not know yet: a dash, never a zero, read as words.
export function UnknownValue(props: UnknownValueProps) {
	const { className } = props

	const translate = useTranslations('tickets')

	return (
		<span className={className}>
			<span aria-hidden="true">—</span>
			<span className="sr-only">{translate('notSet')}</span>
		</span>
	)
}
