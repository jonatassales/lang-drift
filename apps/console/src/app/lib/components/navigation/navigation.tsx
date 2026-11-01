'use client'

import { Brand } from '@langdrift/web/ui/brand'
import {
	CaretDown,
	ChartLineUp,
	CheckCircle,
	ClipboardText,
	Cube,
	FileText,
	FlowArrow,
	Gear,
	GitBranch,
	House,
	MagnifyingGlass,
	Ticket,
} from '@langdrift/web/ui/icons/client'
import { useTranslations } from 'next-intl'

import { useDismissibleMenu } from '@components/menu/use-dismissible-menu.hook'
import type { Product } from '@data/types'
import { Link, usePathname } from '@i18n/navigation'

import { areaOf } from './area'
import { AreaSwitcher } from './area-switcher'

interface NavigationProps {
	product?: Product
	products: Product[]
	mobile?: boolean
	onNavigate?: () => void
}

export function Navigation(props: NavigationProps) {
	const { product, products, mobile = false, onNavigate } = props

	const translate = useTranslations('ui')
	const pathname = usePathname()
	const base = product ? `/products/${product.id}` : ''
	const area = areaOf(pathname)
	const projectGroup = product
		? [
				{
					key: 'project',
					entries: [
						{
							href: base,
							key: 'runtimes',
							icon: Cube,
						},
						{
							href: `${base}/tickets`,
							key: 'tickets',
							icon: Ticket,
						},
						{
							href: `${base}/analysis`,
							key: 'analysis',
							icon: MagnifyingGlass,
						},
						{
							href: `${base}/plans`,
							key: 'plans',
							icon: ClipboardText,
						},
						{
							href: `${base}/workflows`,
							key: 'workflows',
							icon: FlowArrow,
						},
						{
							href: `${base}/executions`,
							key: 'executions',
							icon: GitBranch,
						},
						{
							href: `${base}/history`,
							key: 'history',
							icon: FileText,
						},
						{
							href: `${base}/config/secrets`,
							key: 'config',
							icon: Gear,
						},
					],
				},
			]
		: []
	const groups =
		area === 'finops'
			? [
					{
						key: 'finops',
						entries: [
							{
								href: '/finops',
								key: 'costOverview',
								icon: ChartLineUp,
							},
						],
					},
				]
			: [
					{
						key: 'microruntimes',
						entries: [
							{
								href: '/home',
								key: 'home',
								icon: House,
							},
							{
								href: '/products',
								key: 'projects',
								icon: GitBranch,
							},
						],
					},
					...projectGroup,
				]
	const {
		containerRef: selectorRef,
		menuId,
		open,
		setOpen,
		triggerRef,
	} = useDismissibleMenu()
	const settingsActive =
		pathname === '/settings' || pathname.startsWith('/settings/')
	function navigate() {
		setOpen(false)
		onNavigate?.()
	}
	const linkStyle = (active: boolean) =>
		`relative flex items-center gap-2.5 rounded-md px-3 text-[12px] no-underline transition-colors hover:bg-muted hover:text-foreground focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-ring ${active ? 'bg-muted font-medium text-foreground before:absolute before:inset-y-2.5 before:left-0 before:w-0.5 before:rounded-full before:bg-primary' : 'text-muted-foreground'} ${mobile ? 'min-h-10' : 'min-h-9'}`

	return (
		<>
			<Link
				aria-label={translate('home')}
				className="flex items-center gap-2.5 px-2.5 pt-2 pb-6 no-underline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-ring"
				href="/home"
				onClick={navigate}
			>
				<Brand className="text-foreground" compact />
			</Link>
			<AreaSwitcher onNavigate={onNavigate} />
			<nav aria-label={translate('navigation')} className="grid gap-1">
				{groups.map((group) => (
					<section
						key={group.key}
						className={
							group.key === 'project' ? 'mt-4 border-t border-border pt-4' : ''
						}
					>
						{group.key === 'project' ? (
							<h2 className="px-3 pb-2 text-[10px] font-medium tracking-wide text-muted-foreground uppercase">
								{translate(group.key)}
							</h2>
						) : null}
						{group.key === 'project' && product ? (
							<div className="relative mb-2" ref={selectorRef}>
								<button
									aria-controls={open ? menuId : undefined}
									aria-expanded={products.length > 1 ? open : undefined}
									disabled={products.length < 2}
									className="flex min-h-12 w-full cursor-pointer items-center justify-between gap-3 rounded-lg border border-border bg-background px-3 py-2.5 text-left text-foreground hover:bg-muted focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-ring"
									onClick={() => setOpen(!open)}
									ref={triggerRef}
									type="button"
								>
									<span className="grid gap-1">
										<strong className="text-[13px] font-medium">
											{product.name}
										</strong>
									</span>
									{products.length > 1 ? (
										<CaretDown aria-hidden size={13} />
									) : null}
								</button>
								{open ? (
									<nav
										aria-label={translate('project')}
										className="absolute inset-x-0 top-[calc(100%+6px)] z-50 grid rounded-lg border border-border bg-card p-1.5 shadow-lg"
										id={menuId}
									>
										{products.map((item) => (
											<Link
												className="flex min-h-10 items-center justify-between px-2.5 text-xs focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-ring"
												href={`/products/${item.id}`}
												onClick={navigate}
												key={item.id}
											>
												{item.name}
												{item.id === product.id ? (
													<CheckCircle
														aria-hidden
														className="text-primary"
														size={15}
													/>
												) : null}
											</Link>
										))}
									</nav>
								) : null}
							</div>
						) : null}
						{group.entries.map((item) => {
							const active =
								pathname === item.href ||
								(item.key === 'runtimes' &&
									(pathname.startsWith('/runtimes/') ||
										pathname.startsWith(`${base}/runtimes/`))) ||
								(item.key !== 'projects' &&
									item.key !== 'runtimes' &&
									pathname.startsWith(`${item.href}/`))
							const Icon = item.icon

							return (
								<Link
									aria-current={active ? 'page' : undefined}
									className={linkStyle(active)}
									href={item.href}
									key={item.href}
									onClick={navigate}
								>
									<Icon
										aria-hidden
										className={active ? 'text-primary' : undefined}
										size={17}
									/>
									{translate(item.key)}
								</Link>
							)
						})}
					</section>
				))}
			</nav>
			<div className="mt-auto pt-8">
				<Link
					aria-current={settingsActive ? 'page' : undefined}
					className={linkStyle(settingsActive)}
					href="/settings"
					onClick={navigate}
				>
					<Gear
						aria-hidden
						className={settingsActive ? 'text-primary' : undefined}
						size={17}
					/>
					{translate('settings')}
				</Link>
				<p className="mt-3 px-2 text-[10px] tracking-wide text-muted-foreground/60 select-none">
					{translate('copyright')}
				</p>
			</div>
		</>
	)
}
