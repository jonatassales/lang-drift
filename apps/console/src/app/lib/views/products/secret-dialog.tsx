'use client'

import { useTranslations } from 'next-intl'
import { useState } from 'react'

import { ErrorNotice } from '@components/error-notice/error-notice'
import { fieldClass } from '@components/form-styles/form-styles'
import { FieldSkeleton } from '@components/route-state/field-skeleton'
import { Button } from '@ui/button'
import { OverlayDialog } from '@ui/overlay-dialog'

interface SecretDialogProps {
	secret: string
	pending: boolean
	error: unknown
	onClose: () => void
}

export function SecretDialog(props: SecretDialogProps) {
	const { secret, pending, error, onClose } = props

	const t = useTranslations('ui')
	const [copied, setCopied] = useState(false)
	const [copyFailed, setCopyFailed] = useState(false)
	async function copy() {
		try {
			await navigator.clipboard.writeText(secret)

			setCopied(true)
			setCopyFailed(false)
		} catch {
			setCopyFailed(true)
		}
	}

	return (
		<OverlayDialog
			modal
			label={t('generateSecret')}
			onClose={() => {
				if (!pending) onClose()
			}}
			className="fixed inset-0 m-auto w-[calc(100%-2rem)] max-w-lg rounded-xl border border-border bg-card p-6 text-foreground shadow-xl backdrop:bg-background/80"
		>
			<h2 className="text-lg font-medium">{t('generateSecret')}</h2>
			<div className="mt-4 grid gap-4">
				{pending ? <FieldSkeleton /> : null}
				{error ? <ErrorNotice error={error} /> : null}
				{secret ? (
					<>
						<p className="text-sm text-muted-foreground">
							{t('projectSecretOnce')}
						</p>
						<input
							className={`${fieldClass} font-mono opacity-100`}
							aria-label={t('secretValue')}
							value={secret}
							disabled
							readOnly
						/>
						<Button onClick={copy}>
							{t(copied ? 'secretCopied' : 'copySecret')}
						</Button>
						{copyFailed ? (
							<p role="alert" className="text-sm text-muted-foreground">
								{t('secretCopyFailed')}
							</p>
						) : null}
					</>
				) : null}
				<Button onClick={onClose} disabled={pending}>
					{t('close')}
				</Button>
			</div>
		</OverlayDialog>
	)
}
