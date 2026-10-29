import { readFileSync } from 'node:fs'
import { resolve } from 'node:path'
import { fileURLToPath } from 'node:url'
import { defineConfig } from 'vitest/config'

const websiteRoot = fileURLToPath(new URL('.', import.meta.url))
const pathsConfiguration = JSON.parse(
	readFileSync(resolve(websiteRoot, 'tsconfig.paths.json'), 'utf8'),
) as { compilerOptions: { paths: Record<string, string[]> } }
const aliases = Object.fromEntries(
	Object.entries(pathsConfiguration.compilerOptions.paths).map(
		([alias, targets]) => [
			alias.replace(/\/\*$/, ''),
			resolve(websiteRoot, targets[0].replace(/\/\*$/, '')),
		],
	),
)

export default defineConfig({
	oxc: { jsx: { runtime: 'automatic' } },
	resolve: { alias: aliases },
	test: {
		include: ['src/**/*.unit.{ts,tsx}', 'next.config.unit.ts'],
		server: { deps: { inline: ['next-intl'] } },
	},
})
