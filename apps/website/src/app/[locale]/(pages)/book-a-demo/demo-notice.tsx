import type { Ref } from 'react'

import type { NoticeCopy } from '@i18n/messages/messages.types'

interface DemoNoticeProps {
	copy: NoticeCopy
	headingRef?: Ref<HTMLHeadingElement>
}

// What stands in place of the demo form: why there is none, or what happens
// next. With a heading ref, the heading takes focus so the change is read.
export function DemoNotice(props: DemoNoticeProps) {
	const { copy, headingRef } = props

	return (
		<>
			<h2
				className="m-0 text-2xl font-medium tracking-tight outline-none"
				ref={headingRef}
				tabIndex={headingRef ? -1 : undefined}
			>
				{copy.title}
			</h2>
			<p className="mt-3 mb-0 text-base leading-relaxed text-muted-foreground">
				{copy.description}
			</p>
		</>
	)
}
