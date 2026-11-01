'use client'

import { Button } from '@langdrift/web/ui/button'
import { useTranslations } from 'next-intl'
import { use } from 'react'

import { ErrorNotice } from '@components/error-notice/error-notice'
import type { ReadResult } from '@data/read-result'
import type { FolderListing } from '@data/types'

interface FolderEntriesProps {
	listing: Promise<ReadResult<FolderListing>>
	navigate: (path: string) => void
	onSelect: (path: string) => void
}

export function FolderEntries(props: FolderEntriesProps) {
	const { listing, navigate, onSelect } = props

	const result = use(listing)
	const translate = useTranslations('ui')

	if (!result.data) return <ErrorNotice error={result.error} />
	const data = result.data

	function handleNavigateParent() {
		navigate(data.parent ?? data.path)
	}

	function handleSelectFolder(event: React.MouseEvent<HTMLButtonElement>) {
		const targetPath = event.currentTarget.dataset.path
		if (targetPath) navigate(targetPath)
	}

	function handleSelectCurrent() {
		onSelect(data.path)
	}

	return (
		<>
			<p className="mb-3 break-all text-xs text-muted-foreground">
				{data.path}
			</p>
			{data.parent ? (
				<Button onClick={handleNavigateParent} type="button" variant="outline">
					{translate('parentFolder')}
				</Button>
			) : null}
			<ul
				className="my-4 grid max-h-64 gap-1 overflow-y-auto"
				aria-label={translate('folders')}
			>
				{data.folders.map((folder) => (
					<li key={folder.path}>
						<button
							data-path={folder.path}
							type="button"
							className="w-full rounded-md px-3 py-2 text-left text-sm hover:bg-muted focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-ring"
							onClick={handleSelectFolder}
						>
							{folder.name}
						</button>
					</li>
				))}
			</ul>
			<Button onClick={handleSelectCurrent} type="button" variant="outline">
				{translate('selectFolder')}
			</Button>
		</>
	)
}
