'use client'

import {
	CodeEditor,
	type CodeEditorProps,
} from '@langdrift/web/vendors/vidro/code-editor'
import { useAtomValue } from 'jotai'

import type { WebsiteLocale } from '@i18n/routing'
import { themeAtom } from '@state'

interface RuntimeFileEditorProps {
	locale: WebsiteLocale
	path: string
	value: string
}

// Vidro has no Traditional Chinese interface, and Simplified Chinese never
// stands in for it, so a zh-Hant page gets the editor's English controls.
const editorLocales: Record<WebsiteLocale, CodeEditorProps['uiLocale']> = {
	en: 'en',
	ja: 'ja',
	'pt-BR': 'pt-BR',
	vi: 'vi',
	'zh-Hant': 'en',
}

// The runtime file's langdrift fence in the Vidro editor, read only, with JSON
// highlighting and the site's color scheme.
export function RuntimeFileEditor(props: RuntimeFileEditorProps) {
	const { locale, path, value } = props

	const theme = useAtomValue(themeAtom)

	return (
		<CodeEditor
			className="[&_button]:whitespace-nowrap"
			language="json"
			maxHeight={640}
			theme={theme}
			title={path}
			uiLocale={editorLocales[locale]}
			value={value}
		/>
	)
}
