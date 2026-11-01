import { Brand } from '@langdrift/web/ui/brand'
import { ElevatedFrame } from '@langdrift/web/ui/elevated-frame'
import { HeroShader } from '@langdrift/web/ui/hero-shader'
import type { ReactNode } from 'react'

import { ThemeControl } from '@components/account-controls/theme-control'
import { LanguageSettings } from '@i18n/language-settings'

interface AuthShellProps {
	children: ReactNode
}

// The console's way in, on the website's ground: the shader in its
// muted tone across the screen, the language at the top left and the
// theme at the top right as in the website's header, the LangDrift brand alone
// above the form, and the form in the elevated frame.
export function AuthShell(props: AuthShellProps) {
	const { children } = props

	return (
		<section className="relative isolate min-h-screen overflow-hidden bg-background">
			<HeroShader tone="muted" />
			<div
				aria-hidden="true"
				className="pointer-events-none absolute inset-x-0 bottom-0 h-48 bg-linear-to-b from-transparent to-background"
			/>
			<div className="absolute inset-x-0 top-0 z-10 flex items-center justify-between px-5 py-4 sm:px-7">
				<div className="relative h-8 w-14">
					<div className="absolute top-1/2 left-0 -translate-y-1/2">
						<LanguageSettings compact />
					</div>
				</div>
				<ThemeControl />
			</div>
			<div className="relative mx-auto flex min-h-screen max-w-md flex-col justify-center gap-6 px-5 py-10">
				<div className="flex justify-center">
					<Brand className="text-foreground" />
				</div>
				<ElevatedFrame className="p-6">{children}</ElevatedFrame>
			</div>
		</section>
	)
}
