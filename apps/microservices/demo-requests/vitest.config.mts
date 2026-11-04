import { readFileSync } from 'node:fs'
import { resolve } from 'node:path'
import { fileURLToPath } from 'node:url'
import { defineConfig } from 'vitest/config'

const appRoot = fileURLToPath(new URL('.', import.meta.url))
const pathsConfiguration = JSON.parse(
	readFileSync(resolve(appRoot, 'tsconfig.paths.json'), 'utf8'),
) as { compilerOptions: { paths: Record<string, string[]> } }
const aliases = Object.fromEntries(
	Object.entries(pathsConfiguration.compilerOptions.paths).map(
		([alias, targets]) => [
			alias.replace(/\/\*$/, ''),
			resolve(appRoot, targets[0].replace(/\/\*$/, '')),
		],
	),
)

export default defineConfig({
	oxc: {
		decorator: {
			legacy: true,
		},
	},
	resolve: { alias: aliases },
	test: {
		include: ['src/**/*.unit.ts'],
		setupFiles: ['reflect-metadata'],
	},
})
