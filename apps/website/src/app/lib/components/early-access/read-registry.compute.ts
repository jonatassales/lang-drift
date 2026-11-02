import { assertCredentialFreeOrigin } from './assert-credential-free-origin.validation'
import { earlyAccessMode } from './early-access-mode.compute'

export interface EarlyAccessRegistryConfig {
	mode: 'production' | 'development' | 'test'
	registryUrl: string
	registryKey: string
	publicOrigin: string
}

export function readRegistry(): EarlyAccessRegistryConfig {
	const mode = earlyAccessMode()
	const registryEnv = process.env.EARLY_ACCESS_REGISTRY_URL
	if (!registryEnv) {
		throw new Error('EARLY_ACCESS_REGISTRY_URL is not configured')
	}
	const registry = new URL(registryEnv)

	const websiteEnv = process.env.EARLY_ACCESS_PUBLIC_ORIGIN
	if (!websiteEnv) {
		throw new Error('EARLY_ACCESS_PUBLIC_ORIGIN is not configured')
	}
	const website = new URL(websiteEnv)

	assertCredentialFreeOrigin(registry, mode)
	assertCredentialFreeOrigin(website, mode)

	const registryKey = process.env.EARLY_ACCESS_REGISTRY_KEY
	if (!registryKey || registryKey.length < 32 || /[\r\n]/.test(registryKey)) {
		throw new Error(
			'EARLY_ACCESS_REGISTRY_KEY must contain at least 32 characters',
		)
	}

	return {
		mode,
		registryUrl: registry.origin,
		registryKey,
		publicOrigin: website.origin,
	}
}
