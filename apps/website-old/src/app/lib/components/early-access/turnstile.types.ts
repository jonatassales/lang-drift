export interface TurnstileOptions {
	sitekey: string
	action: 'early-access'
	execution: 'execute'
	'response-field': false
	retry: 'never'
	'refresh-expired': 'never'
	callback: (token: string) => void
	'expired-callback': () => void
	'error-callback': (code: string) => void
	'unsupported-callback': () => void
}

export interface TurnstileApi {
	render: (element: HTMLElement, options: TurnstileOptions) => string
	execute: (id: string) => void
	isExpired: (id: string) => boolean
	reset: (id: string) => void
	remove: (id: string) => void
}

export interface SubmissionChallenge {
	execute: () => Promise<void>
	takeToken: () => string
}

declare global {
	interface Window {
		turnstile?: TurnstileApi
	}
}
