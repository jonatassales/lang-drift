'use client'

import { Button } from '@langdrift/web/ui/button'
import { Card } from '@langdrift/web/ui/card'
import { useTranslations } from 'next-intl'
import { useState, useTransition } from 'react'

import { ErrorNotice } from '@components/error-notice/error-notice'
import { PageHeading } from '@components/page-heading/page-heading'
import { mutation } from '@data/api'
import { useResource } from '@data/resource.hook'
import type { ProductSecret } from '@data/types'

import { SecretDialog } from './secret-dialog'

interface ProductSecretsProps {
	initial: Promise<{ secrets: ProductSecret[] }>
	productId: string
}

export function ProductSecrets(props: ProductSecretsProps) {
	const { productId, initial } = props

	const translate = useTranslations('ui')
	const resource = useResource<{ secrets: ProductSecret[] }>(
		`products/${productId}/secrets`,
		0,
		initial,
	)
	const [open, setOpen] = useState(false)
	const [secret, setSecret] = useState('')
	const [error, setError] = useState<unknown>(null)
	const [pending, startTransition] = useTransition()

	return (
		<>
			<PageHeading
				title={translate('secrets')}
				description={translate('projectSecretsIntro')}
			>
				<Button
					disabled={pending || open}
					onClick={() => {
						setOpen(true)
						setError(null)
						setSecret('')
						startTransition(async () => {
							try {
								const result = await mutation<{ secret: string }>(
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
					}}
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
					onClose={() => {
						setOpen(false)
						setSecret('')
						setError(null)
					}}
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
										disabled={pending}
										onClick={() => {
											if (
												!window.confirm(translate('revokeProjectSecretConfirm'))
											)
												return
											startTransition(async () => {
												setError(null)

												try {
													await mutation(
														`products/${productId}/secrets/${credential.id}`,
														'DELETE',
													)

													setSecret('')
													resource.refresh()
												} catch (failure) {
													setError(failure)
												}
											})
										}}
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
