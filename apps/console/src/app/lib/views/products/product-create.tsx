'use client'

import { Button } from '@langdrift/web/ui/button'
import { Card } from '@langdrift/web/ui/card'
import { useTranslations } from 'next-intl'
import { useEffect, useRef, useState, useTransition } from 'react'

import { ErrorNotice } from '@components/error-notice/error-notice'
import { fieldClass, labelClass } from '@components/form-styles/form-styles'
import { PageHeading } from '@components/page-heading/page-heading'
import { mutation } from '@data/api'
import { type ReadResult, readResult } from '@data/read-result'
import type { FolderListing, Product } from '@data/types'
import { useRouter } from '@i18n/navigation'

import { FolderPicker } from './folder-picker'

interface ProductCreateProps {
	first?: boolean
	productId?: string
}

export function ProductCreate(props: ProductCreateProps) {
	const { first = false, productId } = props

	const t = useTranslations('ui')
	const router = useRouter()
	const [path, setPath] = useState('')
	const [name, setName] = useState('')
	const [picker, setPicker] = useState<Promise<
		ReadResult<FolderListing>
	> | null>(null)
	const folderRequest = useRef<AbortController | null>(null)
	useEffect(() => () => folderRequest.current?.abort(), [])
	const [dropped, setDropped] = useState(false)
	const [error, setError] = useState<unknown>(null)
	const [pending, startTransition] = useTransition()
	function select(path: string) {
		setPath(path)
		setName(path.split('/').filter(Boolean).at(-1) || '')
		setPicker(null)
		setDropped(false)
	}

	return (
		<>
			<PageHeading
				title={t(
					productId ? 'attachFolder' : first ? 'firstProject' : 'createProject',
				)}
				description={t('localProjectIntro')}
			/>
			<Card className="block mx-auto max-w-2xl p-7 sm:mt-10">
				<form
					className="grid gap-5"
					onSubmit={(event) => {
						event.preventDefault()
						setError(null)
						startTransition(async () => {
							try {
								const result = await mutation<{ product: Product }>(
									productId ? `products/${productId}/folder` : 'products',
									'POST',
									{
										path,
										name: name.trim() || undefined,
									},
								)

								router.push(`/products/${result.product.id}`)
								router.refresh()
							} catch (failure) {
								setError(failure)
							}
						})
					}}
				>
					<button
						type="button"
						disabled={pending}
						className="grid min-h-40 place-content-center gap-3 rounded-xl border-2 border-dashed border-border p-8 text-center hover:border-primary hover:bg-muted focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-ring"
						onClick={() => {
							folderRequest.current?.abort()
							folderRequest.current = new AbortController()
							setPicker(
								readResult<FolderListing>(
									'local/folders',
									folderRequest.current.signal,
								),
							)
						}}
						onDragOver={(event) => {
							event.preventDefault()
							event.dataTransfer.dropEffect = 'copy'
						}}
						onDrop={(event) => {
							event.preventDefault()
							if (pending) return
							const file = event.dataTransfer.files[0]
							setPath('')
							setName(file?.name || '')
							setDropped(true)
						}}
					>
						<span className="text-base font-medium">{t('dropFolder')}</span>
						<span className="text-sm text-muted-foreground">
							{t('browseFolder')}
						</span>
					</button>
					{dropped ? (
						<p role="status" className="text-sm text-muted-foreground">
							{t('confirmFolderPath')}
						</p>
					) : null}
					<label className={labelClass}>
						{t('folderPath')}
						<input
							className={fieldClass}
							value={path}
							required
							disabled={pending}
							onChange={(event) => setPath(event.target.value)}
							autoComplete="off"
						/>
					</label>
					{!productId ? (
						<label className={labelClass}>
							{t('projectName')}
							<input
								className={fieldClass}
								value={name}
								maxLength={100}
								disabled={pending}
								onChange={(event) => setName(event.target.value)}
							/>
						</label>
					) : null}
					<p className="text-xs text-muted-foreground">
						{t('projectAgentIntro')}
					</p>
					{error ? <ErrorNotice error={error} /> : null}
					<Button
						type="submit"
						disabled={pending || !path.trim()}
						variant="outline"
					>
						{t(
							pending
								? 'preparingProject'
								: productId
									? 'attachFolder'
									: 'createProject',
						)}
					</Button>
				</form>
			</Card>
			{picker ? (
				<FolderPicker
					initial={picker}
					onSelect={select}
					onClose={() => {
						folderRequest.current?.abort()
						setPicker(null)
					}}
				/>
			) : null}
		</>
	)
}
