import type conversionJson from '@i18n/messages/conversion.json'
import type { ConversionCopy } from '@i18n/messages/messages.types'

import { triple } from './triple.validator'

export function loadConversion(
	locale: (typeof conversionJson)['en'],
): ConversionCopy {
	return {
		...locale,
		plans: {
			...locale.plans,
			cards: {
				starter: {
					...locale.plans.cards.starter,
					items: triple(
						locale.plans.cards.starter.items,
						'plans.cards.starter.items',
					),
				},
				plus: {
					...locale.plans.cards.plus,
					items: triple(
						locale.plans.cards.plus.items,
						'plans.cards.plus.items',
					),
				},
				pro: {
					...locale.plans.cards.pro,
					items: triple(locale.plans.cards.pro.items, 'plans.cards.pro.items'),
				},
			},
		},
	}
}
