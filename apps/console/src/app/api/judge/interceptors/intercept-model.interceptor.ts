import type { ModelPrices } from '@app/api/judge/judge.types'
import { judgeReply } from '@app/api/judge/judge-reply.map'
import { modelPrices } from '@app/api/judge/model-prices.data'
import { offeredModels } from '@app/api/judge/offered-models.compute'

export interface OfferedModel {
	prices: ModelPrices
}

// The console judges only with the models its environment offers, each at
// the list prices it holds; another model is refused with the offered set.
export function interceptModel(model: string): Response | OfferedModel {
	const offered = offeredModels(process.env.JUDGE_MODELS)
	const prices = modelPrices[model]

	if (offered.length === 0) {
		return judgeReply(503, {
			code: 'UNAVAILABLE',
			reason: 'judge-not-configured',
		})
	}

	if (!offered.includes(model)) {
		return judgeReply(400, {
			code: 'MODEL_NOT_OFFERED',
			offered,
		})
	}

	if (!prices) {
		return judgeReply(503, {
			code: 'UNAVAILABLE',
			reason: 'judge-not-configured',
		})
	}

	return {
		prices,
	}
}
