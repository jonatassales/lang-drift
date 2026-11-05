import { createTranslator } from 'next-intl'

import type { WebsiteLocale } from '@i18n/routing'
import type {
	MicroruntimesMessages,
	Resources,
} from '@routes/platform/microruntimes/content/microruntimes-content.data'
import { priceDate } from '@routes/platform/microruntimes/content/microruntimes-workflow.data'
import { tokenCount } from '@routes/platform/microruntimes/content/token-count.fmt'
import { workflowDollars } from '@routes/platform/microruntimes/content/workflow-dollars.fmt'

export type TicketCardCopy = MicroruntimesMessages['session']['card']

// The cards' figures in the visitor's locale: millicores, MiB, the pair of
// them, dollars from cents, token counts, the price date, and the cards'
// ICU messages.
export function ticketCardFormat(copy: TicketCardCopy, locale: WebsiteLocale) {
	const message = createTranslator({
		locale,
		messages: copy,
	})
	const formatNumber = new Intl.NumberFormat(locale)
	const dollars = workflowDollars(locale)
	const cpu = (millis: number) =>
		message('cpuValue', { value: formatNumber.format(millis) })
	const memory = (mebibytes: number) =>
		message('memoryValue', { value: formatNumber.format(mebibytes) })
	const figures = (resources: Resources) =>
		message('figures', {
			cpu: formatNumber.format(resources.cpuMillis),
			memory: formatNumber.format(resources.memoryMiB),
		})
	const money = (cents: number) => dollars.format(cents / 100)
	const tokens = (thousands: number) => tokenCount(locale, thousands)
	const priceMonth = new Intl.DateTimeFormat(locale, {
		month: 'short',
		timeZone: 'UTC',
		year: 'numeric',
	}).format(priceDate)

	return {
		cpu,
		figures,
		formatNumber,
		memory,
		message,
		money,
		priceMonth,
		tokens,
	}
}

export type TicketCardFormat = ReturnType<typeof ticketCardFormat>
