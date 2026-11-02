'use client'

import { useEffect } from 'react'

import { useRouter } from '@i18n/navigation'

interface AutoRedirectProps {
	delayMs?: number
	destination?: string
}

// Redirects the visitor to the target destination after an elapsed duration.
export function AutoRedirect(props: AutoRedirectProps) {
	const { delayMs = 5000, destination = '/' } = props

	const router = useRouter()

	// An external browser timer initiates navigation when the countdown ends.
	useEffect(() => {
		const timer = setTimeout(() => {
			router.push(destination)
		}, delayMs)

		return () => {
			clearTimeout(timer)
		}
	}, [delayMs, destination, router])

	return null
}
