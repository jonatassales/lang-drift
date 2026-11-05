'use client'

import { useTranslations } from 'next-intl'
import { useState, useTransition } from 'react'

import { ErrorNotice } from '@components/error-notice/error-notice'
import { PageHeading } from '@components/page-heading/page-heading'
import { mutation } from '@data/api'
import type {
	IssuedSecretEnvelope,
	ProductSecretsEnvelope,
} from '@data/read-envelopes.types'
import { useResource } from '@data/resource.hook'
import { Button } from 'web/ui/button'
import { Card } from 'web/ui/card'

import { SecretDialog } from './secret-dialog'

interface ProductSecretsProps {
	initial: Promise<ProductSecretsEnvelope>
	productId: string
}

export function ProductSecrets(props: ProductSecretsProps) {
	const { productId, initial } = props

	const [open, setOpen] = useState(false)
	const [secret, setSecret] = useState('')
	const [error, setError] = useState<unknown>(null)
	const [pending, startTransition] = useTransition()

	const resource = useResource<ProductSecretsEnvelope>(
		`products/${productId}/secrets`,
		0,
		initial,
	)
	const translate = useTranslations('ui')

	function handleGenerateSecret() {
		setOpen(true)
		setError(null)
		setSecret('')
		startTransition(async () => {
			try {
				const result = await mutation<IssuedSecretEnvelope>(
					`products/${productId}/secrets`,
					'POST',
					{},
				)

				setSecret(result.secret)
				resource.refresh()
			} catch (failure) {
				setError(failure)
			}
		})
	}

	function handleCloseDialog() {
		setOpen(false)
		setSecret('')
		setError(null)
	}

	function handleRevokeSecret(event: React.MouseEvent<HTMLButtonElement>) {
		const credentialId = event.currentTarget.dataset.credentialId
		if (!credentialId) return
		if (!window.confirm(translate('revokeProjectSecretConfirm'))) return
		startTransition(async () => {
			setError(null)
			try {
				await mutation(
					`products/${productId}/secrets/${credentialId}`,
					'DELETE',
				)

				setSecret('')
				resource.refresh()
			} catch (failure) {
				setError(failure)
			}
		})
	}

	return (
		<>
			<PageHeading
				title={translate('secrets')}
				description={translate('projectSecretsIntro')}
			>
				<Button
					disabled={pending || open}
					onClick={handleGenerateSecret}
					type="button"
					variant="outline"
				>
					{translate('generateSecret')}
				</Button>
			</PageHeading>
			{open ? (
				<SecretDialog
					secret={secret}
					error={error}
					pending={pending}
					onClose={handleCloseDialog}
				/>
			) : null}
			{!open && error ? <ErrorNotice error={error} /> : null}
			{resource.error ? (
				<ErrorNotice error={resource.error} retry={resource.refresh} />
			) : (
				<ul className="grid max-w-2xl gap-3">
					{resource.data?.secrets.map((credential) => (
						<li key={credential.id}>
							<Card className="flex flex-row flex-wrap items-center justify-between gap-3 p-5">
								<span className="text-sm">{credential.name}</span>
								{credential.revokedAt ? (
									<span className="text-xs text-muted-foreground">
										{translate('revoked')}
									</span>
								) : (
									<Button
										data-credential-id={credential.id}
										disabled={pending}
										onClick={handleRevokeSecret}
										type="button"
										variant="outline"
									>
										{translate('revokeSecret')}
									</Button>
								)}
							</Card>
						</li>
					))}
				</ul>
			)}
		</>
	)
}
