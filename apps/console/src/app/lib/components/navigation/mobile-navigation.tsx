'use client'
import { Dialog, DialogContent } from '@langdrift/web/ui/dialog'
import { List, X } from '@langdrift/web/ui/icons/client'
import { useTranslations } from 'next-intl'
import { type ReactNode, useEffect, useState } from 'react'

import { useMobileViewport } from '@components/dialog/use-mobile-viewport.hook'
import { useReturnFocus } from '@components/dialog/use-return-focus.hook'
import { usePathname } from '@i18n/navigation'

interface MobileNavigationProps {
	children: ReactNode
}
export function MobileNavigation(props: MobileNavigationProps) {
	const { children } = props

	const [openPath, setOpenPath] = useState<string | null>(null)
	const mobile = useMobileViewport()
	const focus = useReturnFocus()
	const translate = useTranslations('ui')
	const pathname = usePathname()
	useEffect(() => {
		if (openPath !== pathname) setOpenPath(null)
	}, [pathname, openPath])

	function openNavigation() {
		setOpenPath(pathname)
	}

	function closeNavigation() {
		setOpenPath(null)
	}

	function handleOpenChange(nextOpen: boolean) {
		if (!nextOpen) closeNavigation()
	}

	// The drawer exists below sm only: growing past it closes the drawer for
	// good, so it does not reopen when the viewport narrows again.
	if (!mobile && openPath) setOpenPath(null)

	const open = openPath === pathname

	return (
		<>
			<button
				aria-controls="console-mobile-navigation"
				aria-expanded={open}
				aria-label={translate('openNavigation')}
				className="inline-flex h-8 w-8 shrink-0 items-center justify-center rounded-md border border-border bg-card focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-ring sm:hidden"
				type="button"
				onClick={openNavigation}
			>
				<List aria-hidden size={18} />
			</button>
			<Dialog open={open} onOpenChange={handleOpenChange}>
				<DialogContent
					aria-label={translate('navigation')}
					className="fixed inset-y-0 left-0 m-0 flex h-dvh max-h-none w-[min(320px,88vw)] max-w-none translate-x-0 translate-y-0 flex-col gap-0 overflow-y-auto rounded-none border-0 border-r border-border bg-card px-3 pt-[calc(14px+env(safe-area-inset-top,0))] pb-[calc(16px+env(safe-area-inset-bottom,0))] text-foreground shadow-xl"
					id="console-mobile-navigation"
					onCloseAutoFocus={focus.restore}
					onOpenAutoFocus={focus.remember}
				>
					<button
						aria-label={translate('closeNavigation')}
						className="mb-2 inline-flex h-8 w-8 items-center justify-center self-end rounded-full border border-border focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-ring"
						type="button"
						onClick={closeNavigation}
					>
						<X aria-hidden size={17} />
					</button>
					{children}
				</DialogContent>
			</Dialog>
		</>
	)
}
