import type { ReactNode } from 'react'

export interface ChromaTextProps {
	children: ReactNode
	className?: string
}

export function ChromaText({ children, className = '' }: ChromaTextProps) {
	return <span className={`chroma-text ${className}`}>{children}</span>
}
