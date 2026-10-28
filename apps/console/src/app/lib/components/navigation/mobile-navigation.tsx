'use client'
import { useTranslations } from 'next-intl'
import { type ReactNode, useEffect, useState } from 'react'

import { usePathname } from '@i18n/navigation'
import { List, X } from '@ui/icons'
import { OverlayDialog } from '@ui/overlay-dialog'

interface MobileNavigationProps {
	children: ReactNode
}
export function MobileNavigation(props: MobileNavigationProps) {
	const { children } = props

	const t = useTranslations('ui')
	const pathname = usePathname()
	const [openPath, setOpenPath] = useState<string | null>(null)
	const open = openPath === pathname
	useEffect(() => {
		if (openPath !== pathname) setOpenPath(null)
	}, [pathname, openPath])

	return (
		<>
			<button
				aria-controls="console-mobile-navigation"
				aria-expanded={open}
				aria-label={t('openNavigation')}
				className="inline-flex h-8 w-8 shrink-0 items-center justify-center rounded-md border border-border bg-card sm:hidden"
				type="button"
				onClick={() => setOpenPath(pathname)}
			>
				<List aria-hidden size={18} />
			</button>
			{open ? (
				<OverlayDialog
					label={t('navigation')}
					mobileOnly
					onClose={() => setOpenPath(null)}
					className="fixed inset-y-0 left-0 m-0 flex h-dvh max-h-none w-[min(320px,88vw)] max-w-none flex-col overflow-y-auto border-0 border-r border-border bg-card px-3 pt-[calc(14px+env(safe-area-inset-top,0))] pb-[calc(16px+env(safe-area-inset-bottom,0))] text-foreground shadow-xl backdrop:bg-background/80"
					id="console-mobile-navigation"
				>
					<button
						aria-label={t('closeNavigation')}
						className="mb-2 inline-flex h-8 w-8 items-center justify-center self-end rounded-full border border-border"
						type="button"
						onClick={() => setOpenPath(null)}
					>
						<X aria-hidden size={17} />
					</button>
					{children}
				</OverlayDialog>
			) : null}
		</>
	)
}
