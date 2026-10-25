'use client'

import { useTranslations } from 'next-intl'

// The six fixed rainbow stops turn behind the button; a body in the foreground
// color covers the inside, so only a thin ring and its glow show.
const spinner =
	'aspect-square w-[200%] shrink-0 motion-safe:animate-spin bg-[conic-gradient(var(--lg-rainbow-1),var(--lg-rainbow-2),var(--lg-rainbow-3),var(--lg-rainbow-4),var(--lg-rainbow-5),var(--lg-rainbow-6),var(--lg-rainbow-1))]'

export function DemoButton() {
	const t = useTranslations('header')

	return (
		<span className="relative isolate inline-flex min-h-[38px] items-center justify-center rounded-[7px] px-[13px] py-2 text-xs font-semibold whitespace-nowrap text-background">
			<span
				aria-hidden="true"
				className="absolute -inset-1 -z-20 overflow-hidden rounded-[9px] opacity-60 blur-[7px]"
			>
				<span className="absolute inset-0 flex items-center justify-center">
					<span className={spinner} style={{ animationDuration: '6s' }} />
				</span>
			</span>
			<span
				aria-hidden="true"
				className="absolute inset-0 -z-10 rounded-[inherit] bg-foreground"
			/>
			{t('bookDemo')}
		</span>
	)
}
