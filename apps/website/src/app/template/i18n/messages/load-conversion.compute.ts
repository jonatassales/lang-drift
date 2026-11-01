import type conversionJson from './conversion.json'
import type { ConversionCopy } from './messages.types'
import { triple } from './triple.validator'

export function loadConversion(
	locale: (typeof conversionJson)['en'],
): ConversionCopy {
	return {
		...locale,
		plans: {
			...locale.plans,
			cards: {
				trial: {
					...locale.plans.cards.trial,
					items: triple(
						locale.plans.cards.trial.items,
						'plans.cards.trial.items',
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
