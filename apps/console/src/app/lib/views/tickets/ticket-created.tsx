import { useDateTime } from '@app/template/formatters/date-time.hook'

interface TicketCreatedProps {
	className?: string
	createdAt: string
	/** Set on the ticket page, so the date carries its year. */
	withYear?: boolean
}

const shortDate = {
	day: 'numeric',
	hour: 'numeric',
	minute: '2-digit',
	month: 'short',
} as const

const fullDate = {
	dateStyle: 'medium',
	timeStyle: 'short',
} as const

// When the ticket was created, as an absolute date and time: the console
// formats dates in UTC, so the server and the browser render the same text.
export function TicketCreated(props: TicketCreatedProps) {
	const { className, createdAt, withYear = false } = props

	const dateTime = useDateTime()

	return (
		<time className={className} dateTime={createdAt}>
			{dateTime(new Date(createdAt), withYear ? fullDate : shortDate)}
		</time>
	)
}
