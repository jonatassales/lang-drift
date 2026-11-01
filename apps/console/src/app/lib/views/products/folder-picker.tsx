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

	const [input, setInput] = useState('')
	const [listing, setListing] = useState(initial)
	const controller = useRef<AbortController | null>(null)

	const mobile = useMobileViewport()
	const focus = useReturnFocus()

	const translate = useTranslations('ui')

	// Abort pending directory listing request on unmount.
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

	function handleInputChange(event: React.ChangeEvent<HTMLInputElement>) {
		setInput(event.target.value)
	}

	function handleSubmit(event: React.FormEvent<HTMLFormElement>) {
		event.preventDefault()
		navigate(input)
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
					{translate('chooseFolder')}
				</DialogTitle>
				<form className="mb-4 flex items-end gap-2" onSubmit={handleSubmit}>
					<label className={`${labelClass} min-w-0 flex-1`}>
						{translate('folderPath')}
						<input
							className={fieldClass}
							value={input}
							onChange={handleInputChange}
						/>
					</label>
					<Button type="submit" disabled={!input.trim()} variant="outline">
						{translate('openFolder')}
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
					{translate('cancel')}
				</Button>
			</DialogContent>
		</Dialog>
	)
}
