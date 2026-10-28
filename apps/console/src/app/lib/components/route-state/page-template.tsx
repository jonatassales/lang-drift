import type { ReactNode } from 'react'

interface PageTemplateProps {
	children: ReactNode
}
/** The route template owns fresh form, paging and interaction state on entry. */
export function PageTemplate(props: PageTemplateProps) {
	const { children } = props

	return children
}
