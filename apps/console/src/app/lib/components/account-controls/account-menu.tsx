'use client'

import { CaretDown, Gear, UserCircle } from '@langdrift/web/ui/icons/client'
import { useTranslations } from 'next-intl'
import type { ReactNode } from 'react'

import { useDismissibleMenu } from '@components/menu/use-dismissible-menu.hook'
import { Link } from '@i18n/navigation'

interface AccountMenuProps {
	/** The sign-out item, which talks to the session itself. */
	children: ReactNode
	email: string
}

// The header's account menu, after the old dashboard's: the session email, the
// way to Settings and the way out.
export function AccountMenu(props: AccountMenuProps) {
	const { children, email } = props

	const translate = useTranslations('ui')
	const { containerRef, menuId, open, setOpen, triggerRef } =
		useDismissibleMenu()

	return (
		<div className="relative" ref={containerRef}>
			<button
				aria-controls={open ? menuId : undefined}
				aria-expanded={open}
				aria-label={translate('openAccount')}
				className="flex h-8 cursor-pointer items-center gap-1.5 rounded-full px-1.5 text-muted-foreground hover:bg-muted hover:text-foreground focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-ring"
				onClick={() => setOpen(!open)}
				ref={triggerRef}
				type="button"
			>
				<UserCircle aria-hidden size={20} />
				<span className="max-w-40 truncate text-[11px] max-sm:hidden">
					{email}
				</span>
				<CaretDown aria-hidden className="max-sm:hidden" size={11} />
			</button>
			{open ? (
				<div
					className="absolute top-[calc(100%+6px)] right-0 z-50 grid w-60 gap-0.5 rounded-lg border border-border bg-card p-1 shadow-lg"
					id={menuId}
				>
					<p className="m-0 truncate px-2.5 pt-2 pb-2.5 text-[11px] text-muted-foreground">
						{email}
					</p>
					<Link
						className="flex min-h-9 items-center gap-2.5 rounded-md px-2.5 text-[12px] text-foreground no-underline hover:bg-muted focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-ring"
						href="/settings"
						onClick={() => setOpen(false)}
					>
						<Gear aria-hidden size={15} />
						{translate('settings')}
					</Link>
					{children}
				</div>
			) : null}
		</div>
	)
}
