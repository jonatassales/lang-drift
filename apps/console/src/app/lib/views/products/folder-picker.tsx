'use client'

import { Button } from '@langdrift/web/ui/button'
import { Dialog, DialogContent, DialogTitle } from '@langdrift/web/ui/dialog'
import { useTranslations } from 'next-intl'
import { Suspense, useEffect, useRef, useState } from 'react'

import { useMobileViewport } from '@components/dialog/use-mobile-viewport.hook'
import { useReturnFocus } from '@components/dialog/use-return-focus.hook'
import { fieldClass, labelClass } from '@components/form-styles/form-styles'
import { FolderListSkeleton } from '@components/route-state/folder-list-skeleton'
import { type ReadResult, readResult } from '@data/read-result'
import type { FolderListing } from '@data/types'

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
	const mobile = useMobileViewport()
	const focus = useReturnFocus()
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

	function handleOpenChange(nextOpen: boolean) {
		if (!nextOpen) onClose()
	}

	// The picker is modal on small screens only: on wider ones the page stays
	// usable around it, so a press or focus outside leaves it open.
	function handleInteractOutside(event: Event) {
		if (!mobile) event.preventDefault()
	}

	return (
		<Dialog modal={mobile} open onOpenChange={handleOpenChange}>
			<DialogContent
				className="block max-h-[80dvh] w-[min(560px,94vw)] max-w-none overflow-y-auto rounded-xl border border-border bg-card p-6 text-foreground shadow-xl sm:max-w-none"
				onCloseAutoFocus={focus.restore}
				onInteractOutside={handleInteractOutside}
				onOpenAutoFocus={focus.remember}
			>
				<DialogTitle className="mb-4 text-lg font-medium">
					{t('chooseFolder')}
				</DialogTitle>
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
					<Button type="submit" disabled={!input.trim()} variant="outline">
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
				<Button
					className="mt-3"
					onClick={onClose}
					type="button"
					variant="outline"
				>
					{t('cancel')}
				</Button>
			</DialogContent>
		</Dialog>
	)
}
