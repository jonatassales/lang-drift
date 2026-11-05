import type { ProviderUsage } from '@app/api/judge/judge.types'

// What the Messages API answers for one question: the text in one text
// block, after an empty thinking block, with the usage and the stop reason.
// The route's tests answer the provider's fetch with it, so no test calls a
// provider.
export function messagesReply(
	text: string,
	usage: ProviderUsage,
	stopReason = 'end_turn',
): Response {
	return Response.json({
		content: [
			{
				signature: '',
				thinking: '',
				type: 'thinking',
			},
			{
				text,
				type: 'text',
			},
		],
		id: 'msg_test',
		model: 'claude-sonnet-5-5',
		role: 'assistant',
		stop_reason: stopReason,
		stop_sequence: null,
		type: 'message',
		usage: {
			input_tokens: usage.inputTokens,
			output_tokens: usage.outputTokens,
		},
	})
}
