'use client'

import { useTranslations } from 'next-intl'
import { useEffect, useRef, useState, useTransition } from 'react'

import { ErrorNotice } from '@components/error-notice/error-notice'
import { fieldClass, labelClass } from '@components/form-styles/form-styles'
import { PageHeading } from '@components/page-heading/page-heading'
import { mutation } from '@data/api'
import { type ReadResult, readResult } from '@data/read-result'
import type { FolderListing, Product } from '@data/types'
import { useRouter } from '@i18n/navigation'
import { Button } from 'web/ui/button'
import { Card } from 'web/ui/card'

import { FolderPicker } from './folder-picker'

interface ProductCreateProps {
	first?: boolean
	productId?: string
}

export function ProductCreate(props: ProductCreateProps) {
	const { first = false, productId } = props

	const [path, setPath] = useState('')
	const [name, setName] = useState('')
	const [picker, setPicker] = useState<Promise<
		ReadResult<FolderListing>
	> | null>(null)
	const [dropped, setDropped] = useState(false)
	const [error, setError] = useState<unknown>(null)
	const [pending, startTransition] = useTransition()
	const folderRequest = useRef<AbortController | null>(null)

	const router = useRouter()
	const translate = useTranslations('ui')

	// Abort pending folder request on unmount.
	useEffect(() => () => folderRequest.current?.abort(), [])

	function handleSelect(selectedPath: string) {
		setPath(selectedPath)
		setName(selectedPath.split('/').filter(Boolean).at(-1) || '')
		setPicker(null)
		setDropped(false)
	}

	function handleSubmit(event: React.FormEvent<HTMLFormElement>) {
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
	}

	function handleBrowseFolder() {
		folderRequest.current?.abort()
		folderRequest.current = new AbortController()
		setPicker(
			readResult<FolderListing>('local/folders', folderRequest.current.signal),
		)
	}

	function handleDragOver(event: React.DragEvent) {
		event.preventDefault()
		event.dataTransfer.dropEffect = 'copy'
	}

	function handleDrop(event: React.DragEvent) {
		event.preventDefault()
		if (pending) return
		const file = event.dataTransfer.files[0]
		setPath('')
		setName(file?.name || '')
		setDropped(true)
	}

	function handlePathChange(event: React.ChangeEvent<HTMLInputElement>) {
		setPath(event.target.value)
	}

	function handleNameChange(event: React.ChangeEvent<HTMLInputElement>) {
		setName(event.target.value)
	}

	function handleClosePicker() {
		folderRequest.current?.abort()
		setPicker(null)
	}

	let headingKey = 'createProject'
	if (productId) {
		headingKey = 'attachFolder'
	} else if (first) {
		headingKey = 'firstProject'
	}

	return (
		<>
			<PageHeading
				title={translate(headingKey)}
				description={translate('localProjectIntro')}
			/>
			<Card className="block mx-auto max-w-2xl p-7 sm:mt-10">
				<form className="grid gap-5" onSubmit={handleSubmit}>
					<button
						type="button"
						disabled={pending}
						className="grid min-h-40 place-content-center gap-3 rounded-xl border-2 border-dashed border-border p-8 text-center hover:border-primary hover:bg-muted focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-ring"
						onClick={handleBrowseFolder}
						onDragOver={handleDragOver}
						onDrop={handleDrop}
					>
						<span className="text-base font-medium">
							{translate('dropFolder')}
						</span>
						<span className="text-sm text-muted-foreground">
							{translate('browseFolder')}
						</span>
					</button>
					{dropped ? (
						<p role="status" className="text-sm text-muted-foreground">
							{translate('confirmFolderPath')}
						</p>
					) : null}
					<label className={labelClass}>
						{translate('folderPath')}
						<input
							className={fieldClass}
							value={path}
							required
							disabled={pending}
							onChange={handlePathChange}
							autoComplete="off"
						/>
					</label>
					{!productId ? (
						<label className={labelClass}>
							{translate('projectName')}
							<input
								className={fieldClass}
								value={name}
								maxLength={100}
								disabled={pending}
								onChange={handleNameChange}
							/>
						</label>
					) : null}
					<p className="text-xs text-muted-foreground">
						{translate('projectAgentIntro')}
					</p>
					{error ? <ErrorNotice error={error} /> : null}
					<Button
						type="submit"
						disabled={pending || !path.trim()}
						variant="outline"
					>
						{translate(
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
					onSelect={handleSelect}
					onClose={handleClosePicker}
				/>
			) : null}
		</>
	)
}
