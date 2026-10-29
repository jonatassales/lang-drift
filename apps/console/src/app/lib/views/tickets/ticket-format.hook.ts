import { useLocale, useTranslations } from 'next-intl'

import type { Resources, TicketPlan } from '@data/types'

import { compactTokens, dollars } from './ticket-format'

// The ticket widgets' copy and figures in the viewer's locale: millicores,
// MiB, the pair of them, dollars from cents, compact token counts and
// whether the team has authorized a plan.
export function useTicketFormat() {
	const t = useTranslations('tickets')
	const locale = useLocale()

	const number = new Intl.NumberFormat(locale)
	const cpu = (millis: number) =>
		t('cpuValue', { value: number.format(millis) })
	const memory = (mebibytes: number) =>
		t('memoryValue', { value: number.format(mebibytes) })
	const figures = (resources: Resources) =>
		t('figures', {
			cpu: number.format(resources.cpuMillis),
			memory: number.format(resources.memoryMiB),
		})
	const money = (cents: number) => dollars(locale, cents)
	const tokens = (thousands: number) => compactTokens(locale, thousands)
	const priceMonth = (date: string) =>
		new Intl.DateTimeFormat(locale, {
			month: 'short',
			timeZone: 'UTC',
			year: 'numeric',
		}).format(new Date(`${date}T00:00:00Z`))
	const authorization = (plan: Pick<TicketPlan, 'authorization'>) =>
		plan.authorization === 'authorized'
			? t('plan.authorized')
			: t('plan.pendingAuthorization')

	return {
		authorization,
		cpu,
		figures,
		memory,
		money,
		number,
		priceMonth,
		t,
		tokens,
	}
}

export type TicketFormat = ReturnType<typeof useTicketFormat>
