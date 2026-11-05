import requestSample from '@langdrift/core/judge/request.sample.json' with {
	type: 'json',
}

import { judgeTestSecret } from './judge-secret.data'

export interface JudgeRequestOptions {
	body?: string
	headers?: Record<string, string>
	method?: string
}

// A request to the judge route as the CLI sends it: the contract's sample
// question, as JSON, with the deployment secret as a bearer. Each option
// replaces one part of it.
export function judgeRequest(options: JudgeRequestOptions = {}): Request {
	const {
		body = JSON.stringify(requestSample),
		headers = {},
		method = 'POST',
	} = options

	return new Request('https://console.test/api/judge', {
		body: method === 'GET' ? undefined : body,
		headers: {
			authorization: `Bearer ${judgeTestSecret}`,
			'content-type': 'application/json',
			...headers,
		},
		method,
	})
}
