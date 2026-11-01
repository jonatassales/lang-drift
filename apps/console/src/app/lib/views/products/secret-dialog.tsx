'use client'

import { Button } from '@langdrift/web/ui/button'
import { Dialog, DialogContent, DialogTitle } from '@langdrift/web/ui/dialog'
import { useTranslations } from 'next-intl'
import { useState } from 'react'

import { useReturnFocus } from '@components/dialog/use-return-focus.hook'
import { ErrorNotice } from '@components/error-notice/error-notice'
import { fieldClass } from '@components/form-styles/form-styles'
import { FieldSkeleton } from '@components/route-state/field-skeleton'

interface SecretDialogProps {
	secret: string
	pending: boolean
	error: unknown
	onClose: () => void
}

export function SecretDialog(props: SecretDialogProps) {
	const { secret, pending, error, onClose } = props

	const translate = useTranslations('ui')
	const [copied, setCopied] = useState(false)
	const [copyFailed, setCopyFailed] = useState(false)
	const focus = useReturnFocus()
	async function copy() {
		try {
			await navigator.clipboard.writeText(secret)

			setCopied(true)
			setCopyFailed(false)
		} catch {
			setCopyFailed(true)
		}
	}

	// Escape and a press outside close the dialog, except while the secret is
	// still being generated.
	function handleOpenChange(nextOpen: boolean) {
		if (nextOpen) return
		if (pending) return

		onClose()
	}

	return (
		<Dialog open onOpenChange={handleOpenChange}>
			<DialogContent
				className="block w-[calc(100%-2rem)] max-w-lg rounded-xl border border-border bg-card p-6 text-foreground shadow-xl"
				onCloseAutoFocus={focus.restore}
				onOpenAutoFocus={focus.remember}
			>
				<DialogTitle className="text-lg font-medium">
					{translate('generateSecret')}
				</DialogTitle>
				<div className="mt-4 grid gap-4">
					{pending ? <FieldSkeleton /> : null}
					{error ? <ErrorNotice error={error} /> : null}
					{secret ? (
						<>
							<p className="text-sm text-muted-foreground">
								{translate('projectSecretOnce')}
							</p>
							<input
								className={`${fieldClass} font-mono opacity-100`}
								aria-label={translate('secretValue')}
								value={secret}
								disabled
								readOnly
							/>
							<Button onClick={copy} type="button" variant="outline">
								{translate(copied ? 'secretCopied' : 'copySecret')}
							</Button>
							{copyFailed ? (
								<p role="alert" className="text-sm text-muted-foreground">
									{translate('secretCopyFailed')}
								</p>
							) : null}
						</>
					) : null}
					<Button
						onClick={onClose}
						disabled={pending}
						type="button"
						variant="outline"
					>
						{translate('close')}
					</Button>
				</div>
			</DialogContent>
		</Dialog>
	)
}
