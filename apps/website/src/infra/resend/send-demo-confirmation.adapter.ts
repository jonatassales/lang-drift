import { Resend } from 'resend'

export interface DemoConfirmationInput {
	company: string
	email: string
	locale: string
	position: string
}

const confirmationSubjects: Record<string, string> = {
	en: 'We received your demo request — LangDrift',
	ja: 'デモリクエストを受け付けました — LangDrift',
	'pt-BR': 'Recebemos sua solicitação de demo — LangDrift',
	vi: 'Chúng tôi đã nhận yêu cầu demo của bạn — LangDrift',
	'zh-Hant': '我們已收到您的展示預約 — LangDrift',
}

const confirmationBodies: Record<
	string,
	(input: DemoConfirmationInput) => string
> = {
	en: (input) =>
		`Hello,\n\nThank you for requesting a demo of LangDrift for ${input.company}.\n\nOur team will be in touch shortly to schedule your walkthrough.\n\nBest regards,\nThe LangDrift Team\nhttps://langdrift.md`,
	ja: (input) =>
		`こんにちは、\n\n${input.company}様向けLangDriftデモのお申し込みをいただきありがとうございます。\n\n追って担当者よりデモのご案内を差し上げます。\n\nLangDriftチーム\nhttps://langdrift.md`,
	'pt-BR': (input) =>
		`Olá,\n\nObrigado por solicitar uma demonstração do LangDrift para a ${input.company}.\n\nNossa equipe entrará em contato em breve para agendar a sua apresentação.\n\nAtenciosamente,\nEquipe LangDrift\nhttps://langdrift.md`,
	vi: (input) =>
		`Xin chào,\n\nCảm ơn bạn đã yêu cầu bản demo LangDrift cho ${input.company}.\n\nĐội ngũ của chúng tôi sẽ sớm liên hệ để lên lịch buổi giới thiệu.\n\nTrân trọng,\nĐội ngũ LangDrift\nhttps://langdrift.md`,
	'zh-Hant': (input) =>
		`您好，\n\n感謝您為 ${input.company} 預約 LangDrift 展示。\n\n我們的團隊將很快與您聯絡以安排展示時段。\n\n祝好，\nLangDrift 團隊\nhttps://langdrift.md`,
}

export async function sendDemoConfirmation(
	input: DemoConfirmationInput,
	client?: Resend,
): Promise<{ id: string }> {
	const apiKey = process.env.RESEND_API_KEY

	if (!apiKey && !client) {
		throw new Error('Missing RESEND_API_KEY')
	}

	const resend = client ?? new Resend(apiKey)
	const subject = confirmationSubjects[input.locale] ?? confirmationSubjects.en
	const bodyFormatter =
		confirmationBodies[input.locale] ?? confirmationBodies.en
	const text = bodyFormatter(input)

	const response = await resend.emails.send({
		from: 'LangDrift <no-reply@langdrift.md>',
		subject,
		text,
		to: input.email,
	})

	if (response.error) {
		throw new Error(response.error.message)
	}

	if (!response.data?.id) {
		throw new Error('Resend confirmation failed: missing id')
	}

	return {
		id: response.data.id,
	}
}
