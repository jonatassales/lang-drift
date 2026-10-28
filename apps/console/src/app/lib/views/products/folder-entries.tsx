'use client'

import { useTranslations } from 'next-intl'
import { use } from 'react'

import { ErrorNotice } from '@components/error-notice/error-notice'
import type { ReadResult } from '@data/read-result'
import type { FolderListing } from '@data/types'
import { Button } from '@ui/button'

interface FolderEntriesProps {
	listing: Promise<ReadResult<FolderListing>>
	navigate: (path: string) => void
	onSelect: (path: string) => void
}

export function FolderEntries(props: FolderEntriesProps) {
	const { listing, navigate, onSelect } = props

	const result = use(listing)
	const t = useTranslations('ui')
	if (!result.data) return <ErrorNotice error={result.error} />
	const data = result.data

	return (
		<>
			<p className="mb-3 break-all text-xs text-muted-foreground">
				{data.path}
			</p>
			{data.parent ? (
				<Button onClick={() => navigate(data.parent ?? data.path)}>
					{t('parentFolder')}
				</Button>
			) : null}
			<ul
				className="my-4 grid max-h-64 gap-1 overflow-y-auto"
				aria-label={t('folders')}
			>
				{data.folders.map((folder) => (
					<li key={folder.path}>
						<button
							type="button"
							className="w-full rounded-md px-3 py-2 text-left text-sm hover:bg-muted"
							onClick={() => navigate(folder.path)}
						>
							{folder.name}
						</button>
					</li>
				))}
			</ul>
			<Button onClick={() => onSelect(data.path)}>{t('selectFolder')}</Button>
		</>
	)
}
