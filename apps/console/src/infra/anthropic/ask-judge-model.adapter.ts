import Anthropic from '@anthropic-ai/sdk'

import type { JudgeQuestion, ProviderResult } from '@app/api/judge/judge.types'

import { JudgeModelUnavailable } from './judge-model-unavailable.type'

// Below the route's 120-second budget, so a slow answer still leaves time to
// say the model is unavailable.
const requestTimeoutMs = 110_000

// Asks the platform's model one question with LangDrift's own key, read at
// the call: without the key no call starts. A paid call is never retried, and
// a failure to get an answer is the model being unavailable, never a judgment.
export async function askJudgeModel(
	question: JudgeQuestion,
	signal: AbortSignal,
): Promise<ProviderResult> {
	const apiKey = process.env.ANTHROPIC_API_KEY

	if (!apiKey) throw new JudgeModelUnavailable('platform-key-missing')

	const client = new Anthropic({
		apiKey,
		maxRetries: 0,
		timeout: requestTimeoutMs,
	})

	let message: Anthropic.Message

	try {
		message = await client.messages.create(
			{
				max_tokens: question.maxTokens,
				messages: [
					{
						content: question.prompt,
						role: 'user',
					},
				],
				model: question.model,
				system: question.system,
			},
			{ signal },
		)
	} catch (error) {
		if (error instanceof Anthropic.APIError) {
			throw new JudgeModelUnavailable('provider-unavailable', {
				cause: error,
			})
		}

		throw error
	}

	return {
		stopReason: message.stop_reason,
		text: message.content
			.map((block) => (block.type === 'text' ? block.text : ''))
			.join(''),
		usage: {
			inputTokens: message.usage.input_tokens,
			outputTokens: message.usage.output_tokens,
		},
	}
}
