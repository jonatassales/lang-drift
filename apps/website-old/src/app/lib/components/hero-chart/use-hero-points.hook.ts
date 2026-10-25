'use client'

import { useLocale, useTranslations } from 'next-intl'
import { useMemo } from 'react'

import { getHeroPoints } from '@data/hero-demo-data'

export function useHeroPoints() {
	const t = useTranslations('demo')
	const locale = useLocale()

	return useMemo(() => getHeroPoints(t, locale), [locale, t])
}
