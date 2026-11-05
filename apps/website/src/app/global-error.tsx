'use client'

import { usePathname } from 'next/navigation'

import messages from '@i18n/messages/route-state.json'
import { localePath, routing, type WebsiteLocale } from '@i18n/routing'

interface GlobalErrorProps {
	reset: () => void
}

export default function GlobalError(props: GlobalErrorProps) {
	const { reset } = props

	const pathname = usePathname()

	const prefix = pathname.split('/')[1].toLowerCase()
	const locale: WebsiteLocale =
		routing.locales.find((value) => value.toLowerCase() === prefix) ?? 'en'
	const copy = messages[locale]

	return (
		<html lang={locale}>
			<body>
				<main>
					<h1>{copy.error}</h1>
					<p>{copy.errorDescription}</p>
					<button onClick={reset} type="button">
						{copy.retry}
					</button>
					<p>
						<a href={localePath[locale]}>{copy.home}</a>
					</p>
				</main>
			</body>
		</html>
	)
}
