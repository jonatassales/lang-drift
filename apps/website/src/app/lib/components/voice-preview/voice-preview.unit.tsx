// @vitest-environment jsdom

import { act, cleanup, fireEvent, render, screen } from '@testing-library/react'
import { NextIntlClientProvider } from 'next-intl'
import { afterEach, describe, expect, it, vi } from 'vitest'

import { getWebsiteMessages } from '@i18n/messages/index'
import type { WebsiteLocale } from '@i18n/routing'

import { VoicePreview } from './voice-preview'

// Browser checks cover the decorative dependency's actual custom-element runtime.
vi.mock('orbo-voice/browser', () => ({}))

afterEach(() => {
	cleanup()
	vi.clearAllTimers()
	vi.useRealTimers()
})

async function mount(locale: WebsiteLocale) {
	const view = render(
		<NextIntlClientProvider
			locale={locale}
			messages={getWebsiteMessages(locale)}
			timeZone="UTC"
		>
			<VoicePreview />
		</NextIntlClientProvider>,
	)
	await act(async () => {
		await vi.dynamicImportSettled()
	})
	vi.useFakeTimers()
	return view
}

describe('localized voice preview', () => {
	it.each([
		{
			locale: 'en',
			action: 'Explore the sample answer',
			disclosure:
				'Explore a scripted example of how a question can connect Product Vision, changes, people, and evidence. This preview does not use your microphone or generate a live answer.',
			listening: 'Previewing the sample question',
		},
		{
			locale: 'pt-BR',
			action: 'Ver a resposta de exemplo',
			disclosure:
				'Explore um exemplo simulado de como uma pergunta pode conectar Product Vision, mudanças, pessoas e evidências. Esta demonstração não usa o microfone nem gera respostas em tempo real.',
			listening: 'Mostrando a pergunta de exemplo',
		},
		{
			locale: 'ja',
			action: '回答例を見る',
			disclosure:
				'Product Vision、変化、関わった人、根拠を結びつけて問いかける流れを、あらかじめ用意した例で体験できます。このデモはマイクを使用せず、その場で回答を生成することもありません。',
			listening: '質問例を表示しています',
		},
		{
			locale: 'zh-Hant',
			action: '查看回答範例',
			disclosure:
				'透過預先編寫的範例，了解問題如何連結 Product Vision、變更、相關人員與證據。此預覽不會使用麥克風，也不會即時生成回答。',
			listening: '正在顯示問題範例',
		},
		{
			locale: 'vi',
			action: 'Khám phá câu trả lời mẫu',
			disclosure:
				'Khám phá ví dụ theo kịch bản về cách một câu hỏi kết nối Product Vision, thay đổi, con người và bằng chứng. Bản xem trước này không dùng micrô và không tạo câu trả lời trực tiếp.',
			listening: 'Đang xem câu hỏi mẫu',
		},
	] satisfies {
		locale: WebsiteLocale
		action: string
		disclosure: string
		listening: string
	}[])(
		'renders the real $locale copy and announces each prepared-answer state',
		async ({ locale, action, disclosure, listening }) => {
			await mount(locale)
			const copy = getWebsiteMessages(locale).home.voice
			const description = screen.getByText(disclosure)
			const button = screen.getByRole('button', { name: action })
			const status = screen.getByRole('status')
			const answer = screen.getByText(copy.answer).closest('[aria-live]')

			expect(screen.getByRole('heading', { name: copy.title })).toBeTruthy()
			expect(screen.getByText(copy.kicker)).toBeTruthy()
			expect(
				description.closest('[aria-hidden="true"], .website-sr-only'),
			).toBeNull()
			expect(status.textContent).toBe(copy.idle)
			expect(answer?.getAttribute('aria-hidden')).toBe('true')
			expect(answer?.getAttribute('aria-live')).toBe('polite')

			fireEvent.click(button)
			act(() => vi.advanceTimersByTime(0))
			expect(status.textContent).toBe(listening)
			expect(answer?.getAttribute('aria-hidden')).toBe('true')

			act(() => vi.advanceTimersByTime(850))
			expect(status.textContent).toBe(copy.thinking)
			expect(answer?.getAttribute('aria-hidden')).toBe('true')

			act(() => vi.advanceTimersByTime(850))
			expect(status.textContent).toBe(copy.speaking)
			expect(answer?.getAttribute('aria-hidden')).toBe('false')
			expect(answer?.textContent).toContain(copy.answerTitle)
			expect(answer?.textContent).toContain(copy.answer)

			act(() => vi.advanceTimersByTime(850))
			expect(status.textContent).toBe(copy.idle)
			expect(answer?.getAttribute('aria-hidden')).toBe('false')
		},
	)

	it('restarts the sample without allowing earlier timers to reveal the answer', async () => {
		await mount('en')
		const button = screen.getByRole('button', {
			name: 'Explore the sample answer',
		})
		const status = screen.getByRole('status')
		const answer = screen
			.getByText(getWebsiteMessages('en').home.voice.answer)
			.closest('[aria-live]')

		fireEvent.click(button)
		act(() => vi.advanceTimersByTime(1000))
		fireEvent.click(button)
		act(() => vi.advanceTimersByTime(700))
		expect(status.textContent).toBe('Previewing the sample question')
		expect(answer?.getAttribute('aria-hidden')).toBe('true')

		act(() => vi.advanceTimersByTime(1000))
		expect(status.textContent).toBe('Presenting the sample answer')
		expect(answer?.getAttribute('aria-hidden')).toBe('false')

		fireEvent.click(button)
		expect(answer?.getAttribute('aria-hidden')).toBe('true')
		act(() => vi.advanceTimersByTime(0))
		expect(status.textContent).toBe('Previewing the sample question')
	})

	it('cancels pending sample updates when the preview unmounts', async () => {
		const view = await mount('en')
		fireEvent.click(
			screen.getByRole('button', { name: 'Explore the sample answer' }),
		)
		act(() => vi.advanceTimersByTime(0))
		expect(vi.getTimerCount()).toBeGreaterThan(0)

		view.unmount()
		expect(vi.getTimerCount()).toBe(0)
		act(() => vi.advanceTimersByTime(3000))
		expect(screen.queryByRole('status')).toBeNull()
	})
})
