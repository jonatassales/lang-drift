'use client'

import { usePathname } from 'next/navigation'

import '@app/template/tailwind/app.css'
import { isLocale, languageTags } from '@i18n/config'
import messages from '@i18n/messages/global-error.json'
import { Button } from '@ui/button'

interface GlobalErrorProps {
	retry: () => void
}

// It replaces the root layout when that fails, so it renders its own document
// and takes its copy, without the provider, in the locale the URL names.
export default function GlobalError(props: GlobalErrorProps) {
	const { retry } = props

	const [, prefix] = usePathname().split('/')
	const locale = isLocale(prefix) ? prefix : 'en'
	const copy = messages[locale]

	return (
		<html
			className="dark scheme-light dark:scheme-dark"
			lang={languageTags[locale]}
		>
			<body className="min-h-screen bg-background font-sans text-foreground antialiased">
				<title>{copy.documentTitle}</title>
				<main className="mx-auto flex min-h-screen max-w-xl flex-col justify-center gap-5 px-6">
					<h1 className="text-3xl font-semibold">{copy.title}</h1>
					<p className="text-muted-foreground">{copy.description}</p>
					<div className="flex flex-wrap items-center gap-4">
						<Button tone="primary" onClick={retry}>
							{copy.retry}
						</Button>
						<a
							className="text-sm text-primary underline underline-offset-4"
							href={`/${locale}/home`}
						>
							{copy.home}
						</a>
					</div>
				</main>
			</body>
		</html>
	)
}
