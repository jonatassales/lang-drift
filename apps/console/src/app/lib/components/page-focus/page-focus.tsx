'use client'
import { useEffect, useRef } from 'react'

import { usePathname } from '@i18n/navigation'

export function PageFocus() {
	const pathname = usePathname()
	const previous = useRef(pathname)
	useEffect(() => {
		if (previous.current === pathname) return
		previous.current = pathname
		document.getElementById('console-content')?.focus()
	}, [pathname])

	return null
}
