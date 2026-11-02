import { Resend } from 'resend'

export interface DemoLeadNotificationInput {
	company: string
	email: string
	locale: string
	position: string
	submittedAt?: string
}

export async function sendDemoLeadNotification(
	input: DemoLeadNotificationInput,
	client?: Resend,
): Promise<{ id: string }> {
	const apiKey = process.env.RESEND_API_KEY

	if (!apiKey && !client) {
		throw new Error('Missing RESEND_API_KEY')
	}

	const resend = client ?? new Resend(apiKey)
	const timestamp = input.submittedAt ?? new Date().toISOString()
	const text = [
		'New demo request submitted on langdrift.md:',
		'',
		`Company: ${input.company}`,
		`Contact Email: ${input.email}`,
		`Position: ${input.position}`,
		`Locale: ${input.locale}`,
		`Submitted At: ${timestamp}`,
	].join('\n')

	const response = await resend.emails.send({
		from: 'LangDrift <no-reply@langdrift.md>',
		subject: `[LangDrift] New Demo Request: ${input.company}`,
		text,
		to: 'jonatas.langdrift@gmail.com',
	})

	if (response.error) {
		throw new Error(response.error.message)
	}

	if (!response.data?.id) {
		throw new Error('Resend notification failed: missing id')
	}

	return {
		id: response.data.id,
	}
}
