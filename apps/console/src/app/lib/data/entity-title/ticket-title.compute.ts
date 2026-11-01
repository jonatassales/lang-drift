import 'server-only'
import { getTranslations } from 'next-intl/server'

import { readConsole } from '@data/server'
import type { TicketDetailSnapshot } from '@data/types'

export async function ticketTitle(id: string) {
	const translate = await getTranslations()

	try {
		const detail = await readConsole<TicketDetailSnapshot>(`tickets/${id}`)

		return translate('metadata.ticketTitle', {
			number: detail.ticket.number,
			title: detail.ticket.title,
		})
	} catch {
		return translate('ui.tickets')
	}
}
