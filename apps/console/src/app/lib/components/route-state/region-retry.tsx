'use client'
import { ErrorNotice } from '@components/error-notice/error-notice'
import { useRouter } from '@i18n/navigation'

export function RegionRetry() {
	const router = useRouter()

	return <ErrorNotice error={null} retry={() => router.refresh()} />
}
