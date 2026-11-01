'use client'

import { Brand } from '@langdrift/web/ui/brand'
import { useAtom } from 'jotai'
import { useLocale, useTranslations } from 'next-intl'
import { useEffect, useRef, useState } from 'react'

import { Link } from '@i18n/navigation'
import { themeAtom } from '@state'

import { DemoButton } from './website-header/demo-button'
import { drawerFocus } from './website-header/drawer-focus'
import { LanguagePill } from './website-header/language-pill'
import { PlatformMenu } from './website-header/platform-menu'
import { platformProducts } from './website-header/platform-products'
import { panelSurface } from './website-header/surface'

// Entries without a destination yet read like their neighbors but do not
// navigate: Docs, and the products without a page.
const utilityItem =
	'inline-flex min-h-11 items-center text-[13px] font-medium whitespace-nowrap text-muted-foreground'
const drawerItem =
	'flex items-center gap-3 rounded-lg px-2 py-3.5 text-[19px] text-muted-foreground'
const focusRing =
	'focus-visible:outline-2 focus-visible:outline-offset-3 focus-visible:outline-primary'
const desktopLink = `${utilityItem} hover:text-foreground ${focusRing}`
const drawerLink = `rounded-lg px-2 py-3.5 text-[19px] text-foreground hover:bg-card ${focusRing}`
const brandLink = `inline-flex ${focusRing}`
const iconButton = `inline-flex size-11 shrink-0 basis-11 cursor-pointer items-center justify-center rounded-full border border-transparent bg-transparent p-2.5 whitespace-nowrap text-foreground hover:border-border hover:bg-card/78 ${focusRing}`

// The bar is transparent at the top of the page, so it reads as part of the
// background; once the page scrolls it becomes glass over what passes under
// it, or a nearly solid background without backdrop filters. The glass sits
// on a layer behind the bar's content, not on the header itself: an element
// with a backdrop filter becomes the backdrop root of everything inside it, so
// the Platform dropdown, which hangs below the bar, could only blur the bar and
// would show the page through it sharp.
const glass =
	"before:pointer-events-none before:absolute before:inset-0 before:-z-10 before:bg-transparent before:content-[''] before:transition-colors data-[scrolled=true]:before:bg-background/60 supports-[backdrop-filter]:data-[scrolled=true]:before:bg-background/25 supports-[backdrop-filter]:data-[scrolled=true]:before:backdrop-blur-2xl"

export function WebsiteHeader() {
	const translate = useTranslations('header')
	const locale = useLocale()
	const [theme, setTheme] = useAtom(themeAtom)
	const [scrolled, setScrolled] = useState(false)
	const [drawerOpen, setDrawerOpen] = useState(false)
	const sentinelRef = useRef<HTMLDivElement>(null)
	const drawerRef = useRef<HTMLDialogElement>(null)
	const triggerRef = useRef<HTMLButtonElement>(null)
	const closeRef = useRef<HTMLButtonElement>(null)
	const previousLocale = useRef(locale)

	useEffect(() => {
		document.documentElement.classList.toggle('dark', theme === 'dark')
		document.documentElement.style.colorScheme = theme
		document
			.querySelector('meta[name="theme-color"]')
			?.setAttribute('content', theme === 'dark' ? '#0b0b0c' : '#f7f6f2')
	}, [theme])

	useEffect(() => {
		const sentinel = sentinelRef.current
		if (!sentinel) return
		const observer = new IntersectionObserver(([entry]) => {
			setScrolled(!entry.isIntersecting)
		})
		observer.observe(sentinel)

		return () => observer.disconnect()
	}, [])

	useEffect(() => {
		if (previousLocale.current !== locale) {
			previousLocale.current = locale
			setDrawerOpen(false)
		}
	}, [locale])

	useEffect(() => {
		if (drawerFocus.locale !== locale) return
		const frame = requestAnimationFrame(() => {
			if (drawerFocus.locale !== locale) return
			triggerRef.current?.focus({ preventScroll: true })
			drawerFocus.locale = undefined
		})

		return () => cancelAnimationFrame(frame)
	}, [locale])

	useEffect(() => {
		if (!drawerOpen) return
		const drawer = drawerRef.current
		if (!drawer) return
		const previousOverflow = document.body.style.overflow
		drawer.showModal()
		closeRef.current?.focus()
		document.body.style.overflow = 'hidden'

		return () => {
			drawer.close()
			document.body.style.overflow = previousOverflow
			triggerRef.current?.focus()
		}
	}, [drawerOpen])

	useEffect(() => {
		const desktop = window.matchMedia('(min-width: 64rem)')
		const closeOnDesktop = () => {
			if (desktop.matches) setDrawerOpen(false)
		}
		desktop.addEventListener('change', closeOnDesktop)

		return () => desktop.removeEventListener('change', closeOnDesktop)
	}, [])

	function closeDrawer() {
		setDrawerOpen(false)
	}

	const themeLabel = translate(theme === 'dark' ? 'lightTheme' : 'darkTheme')

	return (
		<>
			<div
				aria-hidden="true"
				className="pointer-events-none -mb-px h-px"
				ref={sentinelRef}
			/>
			<header
				className={`sticky top-0 isolate z-80 border-b border-transparent transition-[background,border-color,box-shadow] duration-220 ease-[ease] data-[scrolled=true]:border-b-border/60 motion-reduce:transition-none ${glass}`}
				data-scrolled={scrolled}
			>
				<div className="mx-auto grid min-h-17.5 max-w-[1280px] grid-cols-[minmax(0,1fr)_auto] items-center gap-3 px-[clamp(18px,3vw,36px)] md:min-h-20 md:gap-5 lg:grid-cols-[minmax(0,1fr)_auto_minmax(0,1fr)]">
					<div className="flex items-center gap-5 justify-self-start">
						<Link aria-label={translate('home')} className={brandLink} href="/">
							<Brand className="website-header-brand" compact tone={theme} />
						</Link>
						<div className="website-header-language">
							<LanguagePill />
						</div>
					</div>
					<nav
						aria-label={translate('navigation')}
						className="hidden items-center justify-center gap-4.5 lg:flex"
					>
						<PlatformMenu />
						<span className={utilityItem}>{translate('docs')}</span>
						<a
							className={desktopLink}
							href="https://github.com/jonatassales/lang-drift"
							rel="noreferrer"
							target="_blank"
						>
							{translate('github')}
						</a>
						<Link className={desktopLink} href="/pricing">
							{translate('pricing')}
						</Link>
					</nav>
					<div className="flex items-center gap-1 justify-self-end md:gap-2.5">
						<div className="hidden items-center gap-2.5 md:flex">
							<DemoButton />
						</div>
						<button
							aria-controls="website-navigation-drawer"
							aria-expanded={drawerOpen}
							aria-label={translate('openNavigation')}
							className={`${iconButton} website-navigation-trigger lg:hidden`}
							onClick={() => setDrawerOpen(true)}
							ref={triggerRef}
							type="button"
						>
							<svg
								aria-hidden="true"
								className="size-4.75"
								fill="none"
								stroke="currentColor"
								strokeLinecap="round"
								strokeLinejoin="round"
								strokeWidth="1.6"
								viewBox="0 0 24 24"
							>
								<path d="M4 7h16M4 12h16M4 17h16" />
							</svg>
						</button>
						<button
							aria-label={themeLabel}
							className={`${iconButton} website-theme-toggle`}
							onClick={() => setTheme(theme === 'light' ? 'dark' : 'light')}
							title={themeLabel}
							type="button"
						>
							<svg
								aria-hidden="true"
								className="size-4.75"
								fill="none"
								stroke="currentColor"
								strokeLinecap="round"
								strokeLinejoin="round"
								strokeWidth="1.6"
								viewBox="0 0 24 24"
							>
								{theme === 'light' ? (
									<path d="M20.8 13.1A9 9 0 0 1 10.9 3.2a9 9 0 1 0 9.9 9.9Z" />
								) : (
									<>
										<circle cx="12" cy="12" r="4" />
										<path d="M12 2v2m0 16v2M2 12h2m16 0h2M4.9 4.9l1.4 1.4m11.4 11.4 1.4 1.4M4.9 19.1l1.4-1.4M17.7 6.3l1.4-1.4" />
									</>
								)}
							</svg>
						</button>
					</div>
				</div>
			</header>
			<dialog
				aria-label={translate('navigation')}
				className={`my-3 mr-3 ml-auto max-h-[calc(100dvh-24px)] w-97.5 max-w-[calc(100vw-24px)] px-5.5 pt-4.5 pb-7.5 backdrop:bg-foreground/38 backdrop:backdrop-blur-[5px] ${panelSurface}`}
				id="website-navigation-drawer"
				onCancel={closeDrawer}
				onClose={closeDrawer}
				ref={drawerRef}
			>
				<div className="flex items-center justify-between gap-2.5">
					<Link
						aria-label={translate('home')}
						className={brandLink}
						href="/"
						onClick={closeDrawer}
					>
						<Brand compact tone={theme} />
					</Link>
					<button
						aria-label={translate('closeNavigation')}
						className={iconButton}
						onClick={closeDrawer}
						ref={closeRef}
						type="button"
					>
						<svg
							aria-hidden="true"
							className="size-4.75"
							fill="none"
							stroke="currentColor"
							strokeLinecap="round"
							strokeLinejoin="round"
							strokeWidth="1.6"
							viewBox="0 0 24 24"
						>
							<path d="m6 6 12 12M6 18 18 6" />
						</svg>
					</button>
				</div>
				<nav aria-label={translate('navigation')} className="grid gap-1 py-5.5">
					<p className="px-2 pt-2 text-xs font-medium tracking-wide text-muted-foreground uppercase">
						{translate('platform')}
					</p>
					{platformProducts.map((product) => {
						const Icon = product.icon

						return 'href' in product ? (
							<Link
								className={`flex items-center gap-3 ${drawerLink}`}
								href={product.href}
								key={product.label}
								onClick={closeDrawer}
							>
								<Icon className="size-5" />
								{translate(product.label)}
							</Link>
						) : (
							<span className={drawerItem} key={product.label}>
								<Icon className="size-5" />
								{translate(product.label)}
							</span>
						)
					})}
					<span className={drawerItem}>{translate('docs')}</span>
					<a
						className={drawerLink}
						href="https://github.com/jonatassales/lang-drift"
						onClick={closeDrawer}
						rel="noreferrer"
						target="_blank"
					>
						{translate('github')}
					</a>
					<Link className={drawerLink} href="/pricing" onClick={closeDrawer}>
						{translate('pricing')}
					</Link>
				</nav>
				<div className="grid justify-items-start gap-6 border-t border-border pt-6">
					<LanguagePill fromDrawer />
					<div className="flex flex-wrap gap-3.5">
						<DemoButton className="min-h-11" onClick={closeDrawer} />
					</div>
				</div>
			</dialog>
		</>
	)
}
