'use client'
import { useTransition } from 'react'

import { ErrorNotice } from '@components/error-notice/error-notice'
import { useRouter } from '@i18n/navigation'

interface PageErrorProps {
	error: Error
	reset: () => void
}
export function PageError(props: PageErrorProps) {
	const { error, reset } = props

	const router = useRouter()
	const [pending, startTransition] = useTransition()

	return (
		<div aria-busy={pending}>
			<ErrorNotice
				error={error}
				retry={() =>
					startTransition(() => {
						router.refresh()
						reset()
					})
				}
			/>
		</div>
	)
}
