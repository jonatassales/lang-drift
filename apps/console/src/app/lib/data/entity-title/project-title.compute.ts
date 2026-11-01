import 'server-only'
import { getTranslations } from 'next-intl/server'

import { readConsole } from '@data/server'
import type { ProductNavigation } from '@data/types'

export async function projectTitle(id: string, section: string) {
	const translate = await getTranslations()

	const label = translate(`ui.${section}`)
	try {
		const navigation = await readConsole<ProductNavigation>('products')

		const project = navigation.products.find((item) => item.id === id)
		if (!project) return label

		return translate('metadata.entityTitle', {
			section: label,
			entity: project.name,
		})
	} catch {
		return label
	}
}
