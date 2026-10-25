import 'server-only'

import type { CryptoService } from './crypto.service'

export type CookieReadResult<T = unknown> =
	| { found: true; value: T }
	| { found: false }

export type SameSite = 'lax' | 'strict' | 'none'

export interface CookieWriteAttributes {
	httpOnly: true
	sameSite: SameSite
	path: string
	secure: boolean
	domain?: string
	maxAge: number
	expires: Date
}

export interface CookieDeleteAttributes {
	httpOnly: true
	sameSite: SameSite
	path: string
	secure: boolean
	domain?: string
}

export interface CookieStore {
	get(name: string): string | undefined
	set(name: string, token: string, options: CookieWriteAttributes): void
	delete(name: string, options: CookieDeleteAttributes): void
}

export interface CookieDefaults {
	httpOnly: true
	sameSite: SameSite
	path: string
	secure: boolean
	domain?: string
	ttlSeconds?: number
}

export interface CookieWriteOptions {
	ttlSeconds?: number
	domain?: string
	path?: string
	sameSite?: SameSite
	secure?: boolean
}

export interface CookieScope {
	path?: string
	domain?: string
}

const nameValueBudget = 3_800
const setCookieBudget = 4_096
const tokenBudget = 3_800
const maxDateTime = 8_640_000_000_000_000
const cookieNamePattern = /^[!#$%&'*+\-.0-9A-Z^_`a-z|~]+$/
const sameSiteLabels: Record<SameSite, string> = {
	lax: 'Lax',
	strict: 'Strict',
	none: 'None',
}

export class CookieService {
	readonly #crypto: CryptoService
	readonly #store: CookieStore
	readonly #defaults: CookieDefaults

	constructor(
		crypto: CryptoService,
		store: CookieStore,
		defaults: CookieDefaults,
	) {
		this.#crypto = crypto
		this.#store = store
		this.#defaults = defaults
	}

	set(name: string, value: unknown, options?: CookieWriteOptions): void {
		assertCookieName(name)
		assertHttpOnly(this.#defaults.httpOnly)

		const ttlSeconds = resolveTtl(
			options?.ttlSeconds ?? this.#defaults.ttlSeconds,
		)
		const expiresAt = Date.now() + ttlSeconds * 1_000
		assertExpiry(expiresAt)

		const attributes = resolveWriteAttributes(
			this.#defaults,
			options,
			ttlSeconds,
			expiresAt,
		)
		const token = this.#crypto.encrypt(
			{
				version: 1,
				value,
				expiresAt,
			},
			cookieContext(name),
		)
		assertCookieSize(name, token, attributes)
		this.#store.set(name, token, attributes)
	}

	get(name: string): CookieReadResult
	get<T>(
		name: string,
		validate: (value: unknown) => value is T,
	): CookieReadResult<T>
	get(name: string, validate?: (value: unknown) => boolean): CookieReadResult {
		assertCookieName(name)

		const token = this.#store.get(name)
		if (!token) return { found: false }
		if (Buffer.byteLength(token) > tokenBudget) return { found: false }

		const decrypted = this.#crypto.decrypt(token, cookieContext(name))
		if (!decrypted.ok) return { found: false }

		const envelope = readEnvelope(decrypted.value)
		if (!envelope) return { found: false }
		if (validate && !validate(envelope.value)) return { found: false }

		return {
			found: true,
			value: envelope.value,
		}
	}

	delete(name: string, scope?: CookieScope): void {
		assertCookieName(name)
		assertHttpOnly(this.#defaults.httpOnly)

		this.#store.delete(name, {
			httpOnly: true,
			sameSite: this.#defaults.sameSite,
			path: scope?.path ?? this.#defaults.path,
			secure: this.#defaults.secure,
			domain: scope?.domain ?? this.#defaults.domain,
		})
	}
}

function cookieContext(name: string): string {
	return `langdrift:website:cookie:${name}`
}

function assertCookieName(name: string) {
	if (!cookieNamePattern.test(name)) {
		throw new Error('Cookie name is not a valid token')
	}
}

function assertHttpOnly(httpOnly: true) {
	if (httpOnly !== true) {
		throw new Error('Encrypted cookies require HttpOnly')
	}
}

function resolveTtl(ttlSeconds: number | undefined): number {
	if (
		typeof ttlSeconds !== 'number' ||
		!Number.isSafeInteger(ttlSeconds) ||
		ttlSeconds <= 0
	) {
		throw new Error('Cookie TTL must be a positive integer number of seconds')
	}

	return ttlSeconds
}

function assertExpiry(expiresAt: number) {
	if (
		!Number.isSafeInteger(expiresAt) ||
		expiresAt <= Date.now() ||
		Math.abs(expiresAt) > maxDateTime
	) {
		throw new Error('Cookie expiry is outside the supported range')
	}
}

function resolveWriteAttributes(
	defaults: CookieDefaults,
	options: CookieWriteOptions | undefined,
	ttlSeconds: number,
	expiresAt: number,
): CookieWriteAttributes {
	const secure = options?.secure ?? defaults.secure
	if (defaults.secure && secure === false) {
		throw new Error('Production cookies cannot disable Secure')
	}

	const sameSite = options?.sameSite ?? defaults.sameSite
	if (!Object.hasOwn(sameSiteLabels, sameSite)) {
		throw new Error('Cookie SameSite is not supported')
	}

	const path = options?.path ?? defaults.path
	if (!path) {
		throw new Error('Cookie path is required')
	}

	return {
		httpOnly: true,
		sameSite,
		path,
		secure,
		domain: options?.domain ?? defaults.domain,
		maxAge: ttlSeconds,
		expires: new Date(expiresAt),
	}
}

function assertCookieSize(
	name: string,
	token: string,
	attributes: CookieWriteAttributes,
) {
	if (Buffer.byteLength(`${name}=${token}`) > nameValueBudget) {
		throw new Error('Cookie exceeds the size budget')
	}

	if (
		Buffer.byteLength(setCookieHeader(name, token, attributes)) >
		setCookieBudget
	) {
		throw new Error('Cookie exceeds the size budget')
	}
}

function setCookieHeader(
	name: string,
	token: string,
	attributes: CookieWriteAttributes,
): string {
	const parts = [
		`${name}=${token}`,
		`Max-Age=${attributes.maxAge}`,
		`Expires=${attributes.expires.toUTCString()}`,
	]
	if (attributes.domain) parts.push(`Domain=${attributes.domain}`)
	parts.push(`Path=${attributes.path}`)
	if (attributes.httpOnly) parts.push('HttpOnly')
	if (attributes.secure) parts.push('Secure')
	parts.push(`SameSite=${sameSiteLabels[attributes.sameSite]}`)

	return parts.join('; ')
}

function readEnvelope(value: unknown): { value: unknown } | undefined {
	if (typeof value !== 'object' || value === null || Array.isArray(value))
		return undefined
	if (
		!Object.hasOwn(value, 'version') ||
		!Object.hasOwn(value, 'value') ||
		!Object.hasOwn(value, 'expiresAt')
	) {
		return undefined
	}

	const record = value as {
		version: unknown
		value: unknown
		expiresAt: unknown
	}
	if (record.version !== 1) return undefined
	if (
		typeof record.expiresAt !== 'number' ||
		!Number.isSafeInteger(record.expiresAt)
	) {
		return undefined
	}
	if (
		Math.abs(record.expiresAt) > maxDateTime ||
		record.expiresAt <= Date.now()
	) {
		return undefined
	}

	return { value: record.value }
}
