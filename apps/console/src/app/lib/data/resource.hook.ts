'use client'
import { use, useCallback, useEffect, useRef, useState } from 'react'

import { consoleRequest } from './api'
export function useResource<T>(
	path: string | null,
	interval: number,
	initial: Promise<T>,
) {
	const initialData = use(initial)
	const [data, setData] = useState<T | null>(initialData ?? null)
	const [error, setError] = useState<unknown>(null)
	const loadRef = useRef<(() => Promise<void>) | null>(null)
	const refresh = useCallback(() => {
		void loadRef.current?.()
	}, [])
	useEffect(() => {
		if (!path) {
			setData(null)
			return
		}
		const controller = new AbortController()
		let pending = false
		setData(initialData ?? null)
		setError(null)
		async function load() {
			if (pending || controller.signal.aborted || !path) return
			pending = true
			try {
				const result = await consoleRequest<T>(path, {
					signal: controller.signal,
				})

				if (!controller.signal.aborted) {
					setData(result)
					setError(null)
				}
			} catch (failure) {
				if (!controller.signal.aborted) setError(failure)
			} finally {
				pending = false
			}
		}
		loadRef.current = load
		const timer = interval
			? window.setInterval(() => {
					if (document.visibilityState === 'visible') void load()
				}, interval)
			: undefined

		return () => {
			controller.abort()
			loadRef.current = null
			if (timer) window.clearInterval(timer)
		}
	}, [path, interval, initialData])

	return {
		data,
		error,
		refresh,
	}
}
