'use client'

import { use, useCallback, useEffect, useRef, useState } from 'react'

import type { Page } from '@data/types'

import { consoleRequest } from './api'

/** Owns sequential page requests; route changes cancel every pending request. */
export function usePagedResource<T>(path: string, initial: Promise<Page<T>>) {
	const initialPage = use(initial)
	const [items, setItems] = useState<T[]>(initialPage.items)
	const [error, setError] = useState<unknown>(null)
	const [loading, setLoading] = useState(false)
	const [hasMore, setHasMore] = useState(initialPage.nextCursor !== null)
	const loadRef = useRef<((reset: boolean) => Promise<void>) | null>(null)
	const failedReset = useRef<boolean | null>(null)
	const next = useCallback(() => {
		void loadRef.current?.(false)
	}, [])
	const refresh = useCallback(() => {
		void loadRef.current?.(true)
	}, [])
	const retry = useCallback(() => {
		if (failedReset.current !== null)
			void loadRef.current?.(failedReset.current)
	}, [])
	useEffect(() => {
		const controller = new AbortController()
		let pending = false
		let cursor: string | null = initialPage.nextCursor
		failedReset.current = null
		setItems(initialPage.items)
		setError(null)
		setLoading(false)
		setHasMore(initialPage.nextCursor !== null)
		async function load(reset: boolean) {
			if (pending || controller.signal.aborted || (!reset && cursor === null))
				return
			pending = true
			setLoading(true)
			setError(null)
			try {
				const separator = path.includes('?') ? '&' : '?'
				const query =
					reset || !cursor
						? ''
						: `${separator}cursor=${encodeURIComponent(cursor)}`
				const page = await consoleRequest<Page<T>>(`${path}${query}`, {
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
				if (!controller.signal.aborted) setLoading(false)
			}
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
