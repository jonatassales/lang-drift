'use client'
import { useTranslations } from 'next-intl'
import { useEffect, useState } from 'react'

import { ErrorNotice } from '@components/error-notice/error-notice'
import { Button } from '@ui/button'
import { Card } from '@ui/card'

interface ConnectionInstructionsProps {
	secret?: string
	local?: boolean
}
export function ConnectionInstructions(props: ConnectionInstructionsProps) {
	const { secret, local = false } = props

	const t = useTranslations('ui')
	const [origin, setOrigin] = useState('')
	const [copied, setCopied] = useState(false)
	const [error, setError] = useState<unknown>(null)
	useEffect(() => {
		setOrigin(window.location.origin)
	}, [])

	return (
		<Card className="p-5">
			<h2 className="text-sm font-medium">{t('connectionInstructions')}</h2>
			{secret ? (
				<div className="my-4 grid gap-3">
					<p className="text-sm text-primary">{t('secretOnce')}</p>
					<code
						className="block break-all rounded-md border border-border bg-background p-3 text-xs"
						data-private="true"
					>
						{secret}
					</code>
					<Button
						onClick={async () => {
							try {
								await navigator.clipboard.writeText(secret)

								setCopied(true)
							} catch (failure) {
								setError(failure)
							}
						}}
					>
						{t(copied ? 'copied' : 'copySecret')}
					</Button>
				</div>
			) : null}
			<p className="mt-3 text-xs leading-relaxed text-muted-foreground">
				{t(local ? 'projectAgentIntro' : 'setupInstructions')}
			</p>
			{!local ? (
				<pre className="mt-4 overflow-x-auto rounded-md border border-border bg-background p-4 text-xs leading-6">
					<code>{`ex init
read -s RUNTIME_SECRET
export RUNTIME_SECRET
ex agent connect ${origin}
unset RUNTIME_SECRET
ex agent run`}</code>
				</pre>
			) : null}
			{error ? <ErrorNotice error={error} /> : null}
		</Card>
	)
}
