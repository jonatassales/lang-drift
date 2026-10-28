'use client'

import { useTranslations } from 'next-intl'
import { Suspense, useEffect, useRef, useState } from 'react'

import { fieldClass, labelClass } from '@components/form-styles/form-styles'
import { FolderListSkeleton } from '@components/route-state/folder-list-skeleton'
import { type ReadResult, readResult } from '@data/read-result'
import type { FolderListing } from '@data/types'
import { Button } from '@ui/button'
import { OverlayDialog } from '@ui/overlay-dialog'

import { FolderEntries } from './folder-entries'

interface FolderPickerProps {
	initial: Promise<ReadResult<FolderListing>>
	onSelect: (path: string) => void
	onClose: () => void
}

export function FolderPicker(props: FolderPickerProps) {
	const { onSelect, onClose, initial } = props

	const t = useTranslations('ui')
	const [input, setInput] = useState('')
	const [listing, setListing] = useState(initial)
	const controller = useRef<AbortController | null>(null)
	useEffect(() => () => controller.current?.abort(), [])
	function navigate(path: string) {
		controller.current?.abort()
		controller.current = new AbortController()
		setListing(
			readResult<FolderListing>(
				`local/folders?path=${encodeURIComponent(path)}`,
				controller.current.signal,
			),
		)
	}

	return (
		<OverlayDialog
			label={t('chooseFolder')}
			onClose={onClose}
			className="fixed inset-0 z-50 m-auto max-h-[80dvh] w-[min(560px,94vw)] overflow-y-auto rounded-xl border border-border bg-card p-6 text-foreground shadow-xl backdrop:bg-background/80"
		>
			<h2 className="mb-4 text-lg font-medium">{t('chooseFolder')}</h2>
			<form
				className="mb-4 flex items-end gap-2"
				onSubmit={(event) => {
					event.preventDefault()
					navigate(input)
				}}
			>
				<label className={`${labelClass} min-w-0 flex-1`}>
					{t('folderPath')}
					<input
						className={fieldClass}
						value={input}
						onChange={(event) => setInput(event.target.value)}
					/>
				</label>
				<Button type="submit" disabled={!input.trim()}>
					{t('openFolder')}
				</Button>
			</form>
			<Suspense fallback={<FolderListSkeleton />}>
				<FolderEntries
					listing={listing}
					navigate={navigate}
					onSelect={onSelect}
				/>
			</Suspense>
			<Button className="mt-3" onClick={onClose}>
				{t('cancel')}
			</Button>
		</OverlayDialog>
	)
}
