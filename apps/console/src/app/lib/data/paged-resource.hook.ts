'use client'

import {
	use,
	useCallback,
	useEffect,
	useRef,
	useState,
	useTransition,
} from 'react'

import type { Page } from '@data/types'

import { consoleRequest } from './api'

/** Owns sequential page requests; route changes cancel every pending request. */
export function usePagedResource<Item>(
	path: string,
	initial: Promise<Page<Item>>,
) {
	const initialPage = use(initial)
	const [loading, startTransition] = useTransition()
	const [items, setItems] = useState<Item[]>(initialPage.items)
	const [error, setError] = useState<unknown>(null)
	const [hasMore, setHasMore] = useState(initialPage.nextCursor !== null)
	const loadRef = useRef<((reset: boolean) => void) | null>(null)
	const failedReset = useRef<boolean | null>(null)
	const next = useCallback(() => {
		loadRef.current?.(false)
	}, [])
	const refresh = useCallback(() => {
		loadRef.current?.(true)
	}, [])
	const retry = useCallback(() => {
		if (failedReset.current !== null) loadRef.current?.(failedReset.current)
	}, [])

	// Manage paginated resource fetching with abort controller across path changes.
	useEffect(() => {
		const controller = new AbortController()
		let pending = false
		let cursor: string | null = initialPage.nextCursor
		failedReset.current = null
		setItems(initialPage.items)
		setError(null)
		setHasMore(initialPage.nextCursor !== null)
		function load(reset: boolean) {
			if (pending || controller.signal.aborted || (!reset && cursor === null))
				return
			pending = true
			setError(null)
			startTransition(async () => {
				try {
					const separator = path.includes('?') ? '&' : '?'
					const query =
						reset || !cursor
							? ''
							: `${separator}cursor=${encodeURIComponent(cursor)}`

					const page = await consoleRequest<Page<Item>>(`${path}${query}`, {
						signal: controller.signal,
					})

					if (controller.signal.aborted) return
					if (
						!Array.isArray(page.items) ||
						!(page.nextCursor === null || typeof page.nextCursor === 'string')
					)
						throw new Error('invalid_activity_page')
					if (!reset && cursor !== null && page.nextCursor === cursor)
						throw new Error('invalid_activity_cursor')
					setItems((current) =>
						reset ? page.items : [...current, ...page.items],
					)
					cursor = page.nextCursor
					failedReset.current = null
					setHasMore(cursor !== null)
				} catch (failure) {
					if (!controller.signal.aborted) {
						failedReset.current = reset
						setError(failure)
					}
				} finally {
					pending = false
				}
			})
		}
		loadRef.current = load

		return () => {
			controller.abort()
			loadRef.current = null
		}
	}, [path, initialPage])

	return {
		items,
		error,
		loading,
		hasMore,
		next,
		refresh,
		retry,
	}
}
