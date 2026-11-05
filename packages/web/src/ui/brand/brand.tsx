import { cn } from 'web/utilities'

import { LogoMark } from './logo-mark'

// The LangDrift mark and wordmark, in the color of the text around them.
export interface BrandProps {
	className?: string
	compact?: boolean
}

export function Brand(props: BrandProps) {
	const { className, compact = false } = props

	return (
		<span
			className={cn(
				'inline-flex items-center font-sans font-[650]',
				compact ? 'gap-2' : 'gap-2.75',
				className,
			)}
		>
			<LogoMark size={compact ? 23 : 30} />
			<span
				className={cn(
					'leading-none tracking-[-0.055em]',
					compact ? 'text-lg' : 'text-2xl',
				)}
			>
				LangDrift
			</span>
		</span>
	)
}
