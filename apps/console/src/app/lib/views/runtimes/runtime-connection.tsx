'use client'

import { use } from 'react'

import type { Product } from '@data/types'

import { ConnectionInstructions } from './connection-instructions'

interface RuntimeConnectionProps {
	initial: Promise<{ product: Product }>
}

export function RuntimeConnection(props: RuntimeConnectionProps) {
	const { initial } = props

	const { product } = use(initial)

	return <ConnectionInstructions local={!!product.source} />
}
