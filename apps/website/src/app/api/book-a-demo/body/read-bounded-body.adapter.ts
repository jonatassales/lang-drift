import { BodyError } from './body-error.type'
import { reportCancelFailure } from './report-cancel-failure.handler'

export async function readBoundedBody(
	body: ReadableStream<Uint8Array> | null,
	maximum: number,
	signal: AbortSignal,
): Promise<string> {
	if (!body) throw new BodyError(400)
	const reader = body.getReader()
	const decoder = new TextDecoder('utf-8', { fatal: true })
	let size = 0
	let text = ''
	function cancel() {
		void reader.cancel().catch(reportCancelFailure)
	}

	signal.addEventListener('abort', cancel, { once: true })
	try {
		for (;;) {
			signal.throwIfAborted()

			const result = await reader.read()

			signal.throwIfAborted()
			if (result.done) break
			size += result.value.byteLength
			if (size > maximum) throw new BodyError(413)
			text += decoder.decode(result.value, { stream: true })
		}
		text += decoder.decode()

		return text
	} finally {
		signal.removeEventListener('abort', cancel)
		cancel()
		reader.releaseLock()
	}
}
