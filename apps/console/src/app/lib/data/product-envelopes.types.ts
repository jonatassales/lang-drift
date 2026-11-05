import type { Product, ProductSecret } from './types'

// What the console API wraps a product read and a product command's reply in.

export interface ProductEnvelope {
	product: Product
}

export interface ProductsEnvelope {
	products: Product[]
}

export interface ProductSecretsEnvelope {
	secrets: ProductSecret[]
}

export interface IssuedSecretEnvelope {
	secret: string
}
