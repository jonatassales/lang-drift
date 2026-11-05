import type {
	JudgeRefusal,
	JudgeResponse,
	JudgeWireResponse,
	ModelPrices,
	ProviderUsage,
} from './judge.types'
import { judgeVersion } from './judge-version.data'
import { judgmentCost } from './judgment-cost.compute'
import { reportedUsage } from './reported-usage.map'

// A judgment goes out as the contract's answer. A refusal goes out with no
// text and its cause, charged at the list prices for the tokens the provider
// reported, since the call spent them.
export function judgeWireResponse(
	outcome: JudgeResponse | JudgeRefusal,
	usage: ProviderUsage,
	prices: ModelPrices,
): JudgeWireResponse {
	if ('reason' in outcome) {
		const reported = reportedUsage(usage)

		return {
			cost: judgmentCost(reported, prices),
			judgeVersion,
			reason: outcome.reason,
			refused: true,
			text: '',
			usage: reported,
		}
	}

	return {
		...outcome,
		reason: '',
		refused: false,
	}
}
