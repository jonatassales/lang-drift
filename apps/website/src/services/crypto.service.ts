import 'server-only'

import {
	createCipheriv,
	createDecipheriv,
	createSecretKey,
	type KeyObject,
	randomBytes,
} from 'node:crypto'
import { deserialize, serialize } from 'node:v8'

export interface DecryptSuccess {
	ok: true
	value: unknown
}

export interface DecryptFailure {
	ok: false
}

export type DecryptResult = DecryptSuccess | DecryptFailure

const secretPattern = /^[A-Za-z0-9+/]{43}=$/
const version = 'v2'
const codec = 'v8'
const defaultContext = 'value'
const ivBytes = 12
const tagBytes = 16

// V8 serialization drops symbol keys and non-enumerable properties, runs visited
// getters, and does not keep prototypes, methods, descriptors, or object identity.
// This token is Node-only. It is not a protocol for Go or browsers.
export class CryptoService {
	readonly #key: KeyObject

	constructor(secretBase64: string) {
		if (!secretPattern.test(secretBase64)) {
			throw new Error(
				'CRYPTO_SECRET must be canonical base64 encoding 32 bytes',
			)
		}

		const decoded = Buffer.from(secretBase64, 'base64')
		if (decoded.length !== 32) {
			throw new Error('CRYPTO_SECRET must decode to exactly 32 bytes')
		}

		this.#key = createSecretKey(decoded)
	}

	encrypt(value: unknown, context = defaultContext): string {
		this.#assertContext(context)

		const plaintext = serialize(value)
		const initializationVector = randomBytes(ivBytes)
		const cipher = createCipheriv(
			'aes-256-gcm',
			this.#key,
			initializationVector,
			{
				authTagLength: tagBytes,
			},
		)
		cipher.setAAD(this.#associatedData(context))
		const ciphertext = Buffer.concat([cipher.update(plaintext), cipher.final()])
		const authTag = cipher.getAuthTag()

		return [
			version,
			codec,
			initializationVector.toString('base64url'),
			ciphertext.toString('base64url'),
			authTag.toString('base64url'),
		].join('.')
	}

	decrypt(token: string, context = defaultContext): DecryptResult {
		this.#assertContext(context)

		const parsed = this.#parseToken(token)
		if (!parsed) return { ok: false }

		try {
			const decipher = createDecipheriv(
				'aes-256-gcm',
				this.#key,
				parsed.initializationVector,
				{
					authTagLength: tagBytes,
				},
			)
			decipher.setAAD(this.#associatedData(context))
			decipher.setAuthTag(parsed.authTag)
			const plaintext = Buffer.concat([
				decipher.update(parsed.ciphertext),
				decipher.final(),
			])

			return {
				ok: true,
				value: deserialize(plaintext),
			}
		} catch {
			return { ok: false }
		}
	}

	#assertContext(context: string): void {
		if (typeof context !== 'string') {
			throw new Error('Crypto context must be a string')
		}
	}

	#associatedData(context: string): Buffer {
		return Buffer.from(
			JSON.stringify(['langdrift-crypto', version, codec, context]),
			'utf8',
		)
	}

	#parseToken(token: string):
		| {
				initializationVector: Buffer
				ciphertext: Buffer
				authTag: Buffer
		  }
		| undefined {
		if (typeof token !== 'string') return undefined

		const segments = token.split('.')
		if (segments.length !== 5) return undefined

		const [tokenVersion, tokenCodec, encodedIv, encodedCiphertext, encodedTag] =
			segments
		const unsupportedVersion = tokenVersion !== version || tokenCodec !== codec
		if (unsupportedVersion) return undefined

		const initializationVector = this.#decodeCanonicalBase64Url(
			encodedIv,
			ivBytes,
		)
		const ciphertext = this.#decodeCanonicalBase64Url(encodedCiphertext)
		const authTag = this.#decodeCanonicalBase64Url(encodedTag, tagBytes)
		const malformedSegments = !initializationVector || !ciphertext || !authTag
		if (malformedSegments) return undefined

		return {
			initializationVector,
			ciphertext,
			authTag,
		}
	}

	#decodeCanonicalBase64Url(
		segment: string,
		expectedLength?: number,
	): Buffer | undefined {
		const invalidBase64Url = segment !== '' && !/^[A-Za-z0-9_-]+$/.test(segment)
		if (invalidBase64Url) return undefined

		const decoded = Buffer.from(segment, 'base64url')
		if (decoded.toString('base64url') !== segment) return undefined

		const lengthMismatch =
			expectedLength !== undefined && decoded.length !== expectedLength
		if (lengthMismatch) return undefined

		return decoded
	}
}
