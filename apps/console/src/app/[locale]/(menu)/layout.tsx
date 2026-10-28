import { useTranslations } from 'next-intl'
import { type ReactNode, Suspense } from 'react'

import { SessionControls } from '@components/account-controls/session-controls'
import { ThemeControl } from '@components/account-controls/theme-control'
import { MobileNavigation } from '@components/navigation/mobile-navigation'
import { NavigationRegion } from '@components/navigation/navigation-region'
import { PageFocus } from '@components/page-focus/page-focus'
import { NavigationSkeleton } from '@components/route-state/navigation-skeleton'
import { SourceNotice } from '@components/source-notice/source-notice'
import { LanguageSettings } from '@i18n/language-settings'
import { Skeleton } from '@ui/skeleton'

interface ConsoleLayoutProps {
	children: ReactNode
	breadcrumb: ReactNode
}
export default function ConsoleLayout(props: ConsoleLayoutProps) {
	const { children, breadcrumb } = props

	const t = useTranslations('ui')

	return (
		<div className="min-h-screen [&_:focus]:scroll-mt-20 max-sm:block sm:grid sm:grid-cols-[224px_minmax(0,1fr)]">
			<a
				href="#console-content"
				className="fixed top-3 left-3 z-[110] -translate-y-24 rounded-md bg-foreground px-4 py-2 text-card focus:translate-y-0"
			>
				{t('skipContent')}
			</a>
			<aside className="sticky top-0 hidden h-screen flex-col overflow-y-auto border-r border-border bg-card px-2.5 py-3.5 sm:flex">
				<Suspense fallback={<NavigationSkeleton />}>
					<NavigationRegion />
				</Suspense>
			</aside>
			<div className="relative min-w-0">
				<header className="sticky top-0 z-30 flex h-[52px] items-center justify-between border-b border-border bg-background/95 px-[18px] max-sm:h-14 max-sm:gap-2 max-sm:px-2.5">
					<div className="flex min-w-0 flex-1 items-center gap-2 overflow-hidden">
						<MobileNavigation>
							<Suspense fallback={<NavigationSkeleton />}>
								<NavigationRegion />
							</Suspense>
						</MobileNavigation>
						{breadcrumb}
					</div>
					<div className="flex shrink-0 items-center gap-2">
						<div className="relative h-8 w-14">
							<div className="absolute top-1/2 right-0 -translate-y-1/2">
								<LanguageSettings compact />
							</div>
						</div>
						<ThemeControl />
						<Suspense fallback={<Skeleton className="h-8 w-8 rounded-full" />}>
							<SessionControls />
						</Suspense>
					</div>
				</header>
				<main
					id="console-content"
					tabIndex={-1}
					className="mx-auto max-w-[1400px] px-6 py-8 outline-none max-sm:px-4 max-sm:py-6"
				>
					<PageFocus />
					<Suspense
						fallback={<Skeleton className="mb-4 h-3 w-80 max-w-full" />}
					>
						<SourceNotice />
					</Suspense>
					{children}
				</main>
			</div>
		</div>
	)
}
