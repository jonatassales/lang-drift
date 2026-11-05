'use client'

import { use } from 'react'

import type { ProductEnvelope } from '@data/product-envelopes.types'

import { ConnectionInstructions } from './connection-instructions'

interface RuntimeConnectionProps {
	initial: Promise<ProductEnvelope>
}

export function RuntimeConnection(props: RuntimeConnectionProps) {
	const { initial } = props

	const { product } = use(initial)

	return <ConnectionInstructions local={!!product.source} />
}
