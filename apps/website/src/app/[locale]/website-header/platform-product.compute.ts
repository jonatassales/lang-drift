import {
	type PlatformProduct,
	platformProducts,
} from './platform-products.data'

// The product a section of the home presents, by its message key; a key
// the menu does not list is a programming error.
export function platformProduct(label: PlatformProduct['label']) {
	const product = platformProducts.find((entry) => entry.label === label)

	if (!product) throw new Error(`Unknown platform product: ${label}`)

	return product
}
