import type { CSSProperties } from 'react'

export interface BrandProps {
	className?: string
	compact?: boolean
	tone?: 'dark' | 'light'
}

export interface LogoMarkProps {
	size?: number
}

const wordmarkStyle: CSSProperties = {
	letterSpacing: '-0.055em',
	lineHeight: 1,
}

export function LogoMark({ size = 31 }: LogoMarkProps) {
	return (
		<svg
			aria-hidden="true"
			height={size}
			viewBox="0 0 78 56"
			width={Math.round(size * 1.39)}
		>
			<path
				d="M2 39 C20 39 27 35 35 24 C44 11 51 9 62 9"
				className="stroke-primary"
				fill="none"
				strokeLinecap="butt"
				strokeWidth="8"
			/>
			<path
				d="M2 51 C24 51 34 46 43 33 C51 21 56 18 65 18"
				className="stroke-primary"
				fill="none"
				strokeLinecap="butt"
				strokeWidth="7"
			/>
			<circle className="fill-primary" cx="68" cy="8" r="8" />
		</svg>
	)
}

export function Brand({
	className,
	compact = false,
	tone = 'light',
}: BrandProps) {
	const textColor = tone === 'dark' ? '#f5f5f5' : '#111111'

	return (
		<span
			className={className}
			style={{
				alignItems: 'center',
				color: textColor,
				display: 'inline-flex',
				fontFamily: 'var(--lg-typography-font-body)',
				fontWeight: 650,
				gap: compact ? '8px' : '11px',
			}}
		>
			<LogoMark size={compact ? 23 : 30} />
			<span
				style={{
					...wordmarkStyle,
					fontSize: compact ? 18 : 24,
				}}
			>
				LangDrift
			</span>
		</span>
	)
}
