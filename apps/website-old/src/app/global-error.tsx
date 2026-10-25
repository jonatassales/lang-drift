'use client'

import { useEffect, useState } from 'react'

import messages from '@i18n/messages/route-state.json'
import { localePath, routing, type WebsiteLocale } from '@i18n/routing'

export default function GlobalError({ reset }: { reset: () => void }) {
	const [locale, setLocale] = useState<WebsiteLocale>('en')
	useEffect(() => {
		const prefix = window.location.pathname.split('/')[1].toLowerCase()
		setLocale(
			routing.locales.find((value) => value.toLowerCase() === prefix) ?? 'en',
		)
	}, [])
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
