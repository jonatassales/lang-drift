import type { ReactNode } from 'react'

interface TemplateProps {
	children: ReactNode
}

export default function Template(props: TemplateProps) {
	const { children } = props

	return children
}
