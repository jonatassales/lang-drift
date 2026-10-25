'use client'

import { cn } from '@langdrift/web/utilities'
import { useEffect, useRef, useState } from 'react'
import type {
	VidroEditor,
	VidroLanguage,
	VidroTheme,
	VidroUiLocale,
} from 'vidropro'

// A read-only Vidro editor (vidro.pro): the owner's Monaco-based Web
// Component, with syntax highlighting, a titled header and a copy button. The
// server renders the plain text; the editor loads in the browser and takes its
// place once it is ready. If it cannot load, the plain text stays.

export interface CodeEditorProps {
	className?: string
	language: VidroLanguage
	/** Beyond it the editor scrolls inside its frame. */
	maxHeight?: number
	theme: VidroTheme
	title?: string
	/** Set once per document: every Vidro editor on a page shares it. */
	uiLocale: VidroUiLocale
	value: string
}

export function CodeEditor(props: CodeEditorProps) {
	const { className, language, maxHeight, theme, title, uiLocale, value } =
		props

	const hostRef = useRef<HTMLDivElement>(null)
	const [isReady, setIsReady] = useState(false)

	useEffect(() => {
		const host = hostRef.current
		if (!host) return

		let editor: VidroEditor | undefined
		let isCancelled = false

		async function mountEditor(target: HTMLDivElement) {
			try {
				await import('vidropro')

				if (isCancelled) return
				editor = document.createElement('vidro-editor')
				editor.uiLocale = uiLocale
				editor.language = language
				editor.headerTitle = title
				editor.value = value
				editor.fontSize = 12
				editor.lineHeight = 20
				editor.className = 'block w-full'
				if (maxHeight) editor.style.maxHeight = `${maxHeight}px`
				target.append(editor)

				await editor.whenReady()

				if (!isCancelled) setIsReady(true)
			} catch {
				// The plain text rendered by the server is the fallback.
				editor?.remove()
			}
		}

		void mountEditor(host)

		return () => {
			isCancelled = true
			editor?.remove()
			setIsReady(false)
		}
	}, [language, maxHeight, title, uiLocale, value])

	useEffect(() => {
		let isCancelled = false

		async function applyTheme() {
			const runtime = await import('vidropro')

			if (!isCancelled) runtime.setTheme(theme)
		}

		void applyTheme()

		return () => {
			isCancelled = true
		}
	}, [theme])

	return (
		<div className={cn('grid min-w-0', className)}>
			<pre
				className={cn(
					'col-start-1 row-start-1 m-0 overflow-auto rounded-2xl border border-border bg-card p-5 font-mono text-[12px] leading-5 text-foreground',
					isReady && 'hidden',
				)}
				style={{ maxHeight }}
			>
				<code>{value}</code>
			</pre>
			<div
				className={cn(
					'col-start-1 row-start-1 min-w-0 [--vidro-frame-border-radius:16px]',
					!isReady && 'invisible',
				)}
				ref={hostRef}
			/>
		</div>
	)
}
