import '@app/template/tailwind/app.css'

import { Brand } from '@langdrift/web/ui/brand'
import type { Metadata, Viewport } from 'next'
import { Head, Search } from 'nextra/components'
import { getPageMap } from 'nextra/page-map'
import { Layout, Navbar } from 'nextra-theme-docs'
import type { ReactNode } from 'react'

import { DocsFooter } from '@app/template/docs-footer'

export const metadata: Metadata = {
	applicationName: 'LangDrift',
	title: {
		default: 'LangDrift Docs',
		template: '%s – LangDrift Docs',
	},
	// Indexing the docs is the owner's decision, still open.
	robots: {
		index: false,
		follow: true,
	},
}

export const viewport: Viewport = {
	colorScheme: 'dark',
	themeColor: '#0a0a0a',
}

interface RootLayoutProps {
	children: ReactNode
}

export default async function RootLayout(props: RootLayoutProps) {
	const { children } = props

	const pageMap = await getPageMap()

	return (
		<html lang="en-US" dir="ltr" className="dark" suppressHydrationWarning>
			<Head
				backgroundColor={{
					dark: '#0a0a0a',
					light: '#0a0a0a',
				}}
				color={{
					hue: 266,
					saturation: 100,
					lightness: 76,
				}}
			/>
			<body className="min-h-screen supports-[height:100dvh]:min-h-dvh bg-background font-sans text-foreground">
				<Layout
					copyPageButton={false}
					darkMode={false}
					editLink={null}
					feedback={{ content: null }}
					footer={
						<div className="mx-auto w-full max-w-[1440px] border-t border-border px-6 py-8 sm:px-10">
							<DocsFooter />
						</div>
					}
					navbar={
						<Navbar
							logo={<Brand compact tone="dark" />}
							logoLink="/"
							projectLink="https://github.com/jonatassales"
						/>
					}
					nextThemes={{
						defaultTheme: 'dark',
						forcedTheme: 'dark',
					}}
					pageMap={pageMap}
					search={
						<Search
							aria-label="Search documentation"
							placeholder="Search docs…"
						/>
					}
					toc={{ title: 'On this page' }}
					sidebar={{ defaultMenuCollapseLevel: 2 }}
				>
					{children}
				</Layout>
			</body>
		</html>
	)
}
