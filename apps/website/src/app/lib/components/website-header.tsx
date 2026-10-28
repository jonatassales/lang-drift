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
import { platformProducts } from './website-header/platform-products'
import { panelSurface } from './website-header/surface'

import './website-header.css'

// Entries without a destination yet read like their neighbors but do not
// navigate: Docs and GitHub.
const utilityItem =
	'inline-flex min-h-11 items-center text-[13px] font-medium whitespace-nowrap text-muted-foreground'
const drawerItem =
	'flex items-center gap-3 rounded-lg px-2 py-3.5 text-[19px] text-muted-foreground'

// The bar is transparent at the top of the page, so it reads as part of the
// background; once the page scrolls it becomes glass over what passes under
// it, or a nearly solid background without backdrop filters. The glass sits
// on a layer behind the bar's content, not on the header itself: an element
// with a backdrop filter becomes the backdrop root of everything inside it.
const glass =
	"before:pointer-events-none before:absolute before:inset-0 before:-z-10 before:bg-transparent before:content-[''] before:transition-colors data-[scrolled=true]:before:bg-background/90 supports-[backdrop-filter]:data-[scrolled=true]:before:bg-background/70 supports-[backdrop-filter]:data-[scrolled=true]:before:backdrop-blur-xl"

export function WebsiteHeader() {
	const t = useTranslations('header')
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

	const themeLabel = t(theme === 'dark' ? 'lightTheme' : 'darkTheme')

	return (
		<>
			<div
				aria-hidden="true"
				className="website-header-sentinel"
				ref={sentinelRef}
			/>
			<header className={`website-header ${glass}`} data-scrolled={scrolled}>
				<div className="website-header-inner grid-cols-[minmax(0,1fr)_auto] lg:grid-cols-[minmax(0,1fr)_auto_minmax(0,1fr)]">
					<div className="website-header-identity">
						<Link aria-label={t('home')} href="/">
							<Brand className="website-header-brand" compact tone={theme} />
						</Link>
						<div className="website-header-language">
							<LanguagePill />
						</div>
					</div>
					<nav
						aria-label={t('navigation')}
						className="website-desktop-nav hidden items-center justify-center gap-4.5 lg:flex"
					>
						{platformProducts.map((product) => (
							<Link href={product.href} key={product.label}>
								{t(product.label)}
							</Link>
						))}
						<span className={utilityItem}>{t('docs')}</span>
						<span className={utilityItem}>{t('github')}</span>
						<Link href="/pricing">{t('pricing')}</Link>
					</nav>
					<div className="website-header-controls">
						<div className="website-header-access">
							<DemoButton />
						</div>
						<button
							aria-controls="website-navigation-drawer"
							aria-expanded={drawerOpen}
							aria-label={t('openNavigation')}
							className="website-header-icon website-navigation-trigger inline-flex lg:hidden"
							onClick={() => setDrawerOpen(true)}
							ref={triggerRef}
							type="button"
						>
							<svg aria-hidden="true" fill="none" viewBox="0 0 24 24">
								<path d="M4 7h16M4 12h16M4 17h16" />
							</svg>
						</button>
						<button
							aria-label={themeLabel}
							className="website-header-icon website-theme-toggle inline-flex"
							onClick={() => setTheme(theme === 'light' ? 'dark' : 'light')}
							title={themeLabel}
							type="button"
						>
							<svg aria-hidden="true" fill="none" viewBox="0 0 24 24">
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
				aria-label={t('navigation')}
				className={`website-navigation-drawer ${panelSurface}`}
				id="website-navigation-drawer"
				onCancel={closeDrawer}
				onClose={closeDrawer}
				ref={drawerRef}
			>
				<div className="website-drawer-heading">
					<Link aria-label={t('home')} href="/" onClick={closeDrawer}>
						<Brand compact tone={theme} />
					</Link>
					<button
						aria-label={t('closeNavigation')}
						className="website-header-icon inline-flex"
						onClick={closeDrawer}
						ref={closeRef}
						type="button"
					>
						<svg aria-hidden="true" fill="none" viewBox="0 0 24 24">
							<path d="m6 6 12 12M6 18 18 6" />
						</svg>
					</button>
				</div>
				<nav aria-label={t('navigation')} className="website-drawer-nav">
					{platformProducts.map((product) => {
						const Icon = product.icon

						return (
							<Link
								className="flex items-center gap-3"
								href={product.href}
								key={product.label}
								onClick={closeDrawer}
							>
								<Icon className="size-5" />
								{t(product.label)}
							</Link>
						)
					})}
					<span className={drawerItem}>{t('docs')}</span>
					<span className={drawerItem}>{t('github')}</span>
					<Link href="/pricing" onClick={closeDrawer}>
						{t('pricing')}
					</Link>
				</nav>
				<div className="website-drawer-access">
					<LanguagePill fromDrawer />
					<div className="website-drawer-buttons">
						<DemoButton onClick={closeDrawer} />
					</div>
				</div>
			</dialog>
		</>
	)
}
