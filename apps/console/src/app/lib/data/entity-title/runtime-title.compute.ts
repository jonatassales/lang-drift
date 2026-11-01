import 'server-only'
import { getTranslations } from 'next-intl/server'

import { readConsole } from '@data/server'
import type { RuntimeDetailSnapshot } from '@data/types'

export async function runtimeTitle(id: string, section?: string) {
	const translate = await getTranslations()

	try {
		const detail = await readConsole<RuntimeDetailSnapshot>(`runtimes/${id}`)

		const { name } = detail.runtime
		if (!section) return name

		return translate('metadata.entityTitle', {
			section: translate(`ui.${section}`),
			entity: name,
		})
	} catch {
		return translate(`ui.${section ?? 'runtimes'}`)
	}
}
