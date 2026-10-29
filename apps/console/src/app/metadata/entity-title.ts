import 'server-only'
import { getTranslations } from 'next-intl/server'

import { readConsole } from '@data/server'
import type {
	ProductNavigation,
	RuntimeDetailSnapshot,
	TicketDetailSnapshot,
} from '@data/types'

// An entity screen's title names the project, runtime or ticket it shows. Each
// read is one the page or its navigation already makes, request-cached, so a
// title costs no extub call of its own, and a failed read leaves the section
// label.

export async function projectTitle(id: string, section: string) {
	const t = await getTranslations()

	const label = t(`ui.${section}`)
	try {
		const navigation = await readConsole<ProductNavigation>('products')

		const project = navigation.products.find((item) => item.id === id)
		if (!project) return label

		return t('metadata.entityTitle', {
			section: label,
			entity: project.name,
		})
	} catch {
		return label
	}
}

// A runtime's own screen takes its name alone; its sections add theirs.
export async function runtimeTitle(id: string, section?: string) {
	const t = await getTranslations()

	try {
		const detail = await readConsole<RuntimeDetailSnapshot>(`runtimes/${id}`)

		const { name } = detail.runtime
		if (!section) return name

		return t('metadata.entityTitle', {
			section: t(`ui.${section}`),
			entity: name,
		})
	} catch {
		return t(`ui.${section ?? 'runtimes'}`)
	}
}

export async function ticketTitle(id: string) {
	const t = await getTranslations()

	try {
		const detail = await readConsole<TicketDetailSnapshot>(`tickets/${id}`)

		return t('metadata.ticketTitle', {
			number: detail.ticket.number,
			title: detail.ticket.title,
		})
	} catch {
		return t('ui.tickets')
	}
}
