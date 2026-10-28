'use client'

import Image from 'next/image'
import { useTranslations } from 'next-intl'
import { useEffect, useId, useRef, useState } from 'react'

import type { Product } from '@data/types'
import { Link, usePathname } from '@i18n/navigation'
import {
	CaretDown,
	CheckCircle,
	Cube,
	FileText,
	Gear,
	GitBranch,
	House,
} from '@ui/icons'

interface NavigationProps {
	product?: Product
	products: Product[]
	mobile?: boolean
	onNavigate?: () => void
}

export function Navigation(props: NavigationProps) {
	const { product, products, mobile = false, onNavigate } = props

	const t = useTranslations('ui')
	const pathname = usePathname()
	const base = product ? `/products/${product.id}` : ''
	const groups = [
		{
			key: 'account',
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
				{
					href: '/settings',
					key: 'settings',
					icon: Gear,
				},
			],
		},
		...(product
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
			: []),
	]
	const [open, setOpen] = useState(false)
	const selectorRef = useRef<HTMLDivElement>(null)
	const triggerRef = useRef<HTMLButtonElement>(null)
	const menuId = useId()
	useEffect(() => {
		if (!open) return
		function dismiss(event: Event) {
			if (!selectorRef.current?.contains(event.target as Node)) setOpen(false)
		}
		function handleEscape(event: KeyboardEvent) {
			if (event.key === 'Escape') {
				setOpen(false)
				triggerRef.current?.focus()
			}
		}
		document.addEventListener('pointerdown', dismiss)
		document.addEventListener('focusin', dismiss)
		document.addEventListener('keydown', handleEscape)

		return () => {
			document.removeEventListener('pointerdown', dismiss)
			document.removeEventListener('focusin', dismiss)
			document.removeEventListener('keydown', handleEscape)
		}
	}, [open])
	function navigate() {
		setOpen(false)
		onNavigate?.()
	}
	const linkStyle = (active: boolean) =>
		`relative flex items-center gap-2.5 rounded-md px-3 text-[12px] no-underline transition-colors hover:bg-muted hover:text-foreground ${active ? 'bg-muted font-medium text-foreground before:absolute before:inset-y-2.5 before:left-0 before:w-0.5 before:rounded-full before:bg-primary' : 'text-muted-foreground'} ${mobile ? 'min-h-10' : 'min-h-9'}`

	return (
		<>
			<Link
				aria-label={t('home')}
				className="flex items-center gap-2.5 px-2.5 pt-2 pb-6 no-underline"
				href="/home"
				onClick={navigate}
			>
				<Image
					src="/brand/exphere-logo.svg"
					alt="Exphere"
					width={132}
					height={31}
					priority
					className="h-auto w-[132px] dark:invert"
				/>
			</Link>
			<nav aria-label={t('navigation')} className="grid gap-1">
				{groups.map((group) => (
					<section
						key={group.key}
						className={
							group.key === 'project' ? 'mt-4 border-t border-border pt-4' : ''
						}
					>
						<h2 className="px-3 pb-2 text-[10px] font-medium tracking-wide text-muted-foreground uppercase">
							{t(group.key)}
						</h2>
						{group.key === 'project' && product ? (
							<div className="relative mb-2" ref={selectorRef}>
								<button
									aria-controls={open ? menuId : undefined}
									aria-expanded={products.length > 1 ? open : undefined}
									disabled={products.length < 2}
									className="flex min-h-12 w-full cursor-pointer items-center justify-between gap-3 rounded-lg border border-border bg-background px-3 py-2.5 text-left text-foreground hover:bg-muted"
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
										aria-label={t('project')}
										className="absolute inset-x-0 top-[calc(100%+6px)] z-50 grid rounded-lg border border-border bg-card p-1.5 shadow-lg"
										id={menuId}
									>
										{products.map((item) => (
											<Link
												className="flex min-h-10 items-center justify-between px-2.5 text-xs"
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
									{t(item.key)}
								</Link>
							)
						})}
					</section>
				))}
			</nav>
			<div className="mt-auto flex items-center justify-between gap-2 px-3 pt-8 pb-4 text-[10px]">
				{product ? (
					<>
						<span className="text-muted-foreground">{t('project')}</span>
						<span className="truncate font-medium">{product.name}</span>
					</>
				) : null}
			</div>
		</>
	)
}
