'use client'

import { FinopsIcon, MicroruntimesIcon } from '@langdrift/web/ui/icons/platform'
import { useTranslations } from 'next-intl'
import { useId } from 'react'

import { Link, usePathname } from '@i18n/navigation'
import { CaretDown, CheckCircle } from '@ui/icons'

import { type Area, areaHome, areaOf, areas } from './area'
import { useDismissableMenu } from './use-dismissable-menu.hook'

interface AreaSwitcherProps {
	onNavigate?: () => void
}

const areaIcon: Record<Area, typeof FinopsIcon> = {
	microruntimes: MicroruntimesIcon,
	finops: FinopsIcon,
}

// The sidebar's first control, after the old dashboard's product switcher: a
// badge, a small label and the current area, opening a menu of both areas.
export function AreaSwitcher(props: AreaSwitcherProps) {
	const { onNavigate } = props

	const t = useTranslations('ui')
	const pathname = usePathname()
	const menuId = useId()
	const { containerRef, open, setOpen, triggerRef } = useDismissableMenu()
	const current = areaOf(pathname)
	const CurrentIcon = areaIcon[current]

	function choose() {
		setOpen(false)
		onNavigate?.()
	}

	return (
		<div className="relative mb-4" ref={containerRef}>
			<button
				aria-controls={open ? menuId : undefined}
				aria-expanded={open}
				className="grid min-h-11 w-full cursor-pointer grid-cols-[28px_minmax(0,1fr)_auto] items-center gap-2 rounded-lg border border-border bg-muted px-2 py-1.5 text-left text-foreground hover:bg-background"
				onClick={() => setOpen(!open)}
				ref={triggerRef}
				type="button"
			>
				<span
					aria-hidden="true"
					className="inline-flex size-7 items-center justify-center rounded-md bg-primary text-primary-foreground"
				>
					<CurrentIcon className="size-4" />
				</span>
				<span className="grid min-w-0 gap-px">
					<small className="text-[9px] tracking-wide text-muted-foreground uppercase">
						{t('area')}
					</small>
					<strong className="truncate text-[12px] font-medium">
						{t(current)}
					</strong>
				</span>
				<CaretDown aria-hidden size={12} />
			</button>
			{open ? (
				<nav
					aria-label={t('areas')}
					className="absolute inset-x-0 top-[calc(100%+5px)] z-50 grid gap-0.5 rounded-lg border border-border bg-card p-1 shadow-lg"
					id={menuId}
				>
					{areas.map((area) => {
						const Icon = areaIcon[area]

						return (
							<Link
								aria-current={area === current ? 'true' : undefined}
								className="flex min-h-10 items-center gap-2.5 rounded-md px-2 text-[12px] no-underline hover:bg-muted"
								href={areaHome[area]}
								key={area}
								onClick={choose}
							>
								<Icon className="size-4" />
								<span className="flex-1">{t(area)}</span>
								{area === current ? (
									<CheckCircle aria-hidden className="text-primary" size={15} />
								) : null}
							</Link>
						)
					})}
				</nav>
			) : null}
		</div>
	)
}
