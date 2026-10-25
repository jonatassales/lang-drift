import { readFile, rm } from 'node:fs/promises'
import path from 'node:path'
import type { NextAdapter } from 'next'
import * as pagefind from 'pagefind'

// Framework pages that never belong in the search corpus.
const excludedPathnames = new Set([
	'/_not-found',
	'/_global-error',
	'/404',
	'/500',
])

// Nextra's search adds the base path when it loads Pagefind and when it opens
// a result, so pages are indexed under their URL without it.
function searchUrl(pathname: string, basePath: string) {
	if (!basePath || !pathname.startsWith(basePath)) return pathname

	return pathname.slice(basePath.length) || '/'
}

const adapter: NextAdapter = {
	name: 'LangDrift documentation search',
	async onBuildComplete(context) {
		const basePath = context.config.basePath ?? ''
		const documents = new Map<string, string>()
		for (const output of [
			...context.outputs.staticFiles,
			...context.outputs.prerenders.map((prerender) => ({
				pathname: prerender.pathname,
				filePath: prerender.fallback?.filePath,
			})),
		]) {
			if (!output.filePath?.endsWith('.html')) continue

			const url = searchUrl(output.pathname, basePath)
			if (excludedPathnames.has(url)) continue

			documents.set(url, output.filePath)
		}
		if (documents.size === 0) {
			throw new Error('No prerendered documentation pages found for search.')
		}

		try {
			const { index, errors } = await pagefind.createIndex()

			if (errors.length || !index) {
				throw new Error(`Unable to create search index: ${errors.join('; ')}`)
			}
			for (const [url, filePath] of documents) {
				const content = await readFile(filePath, 'utf8')

				const result = await index.addHTMLFile({
					url,
					content,
				})

				if (result.errors.length) {
					throw new Error(`Unable to index ${url}: ${result.errors.join('; ')}`)
				}
			}
			const outputPath = path.join(context.projectDir, 'public/_pagefind')

			await rm(outputPath, {
				recursive: true,
				force: true,
			})

			const result = await index.writeFiles({ outputPath })

			if (result.errors.length) {
				throw new Error(
					`Unable to write search index: ${result.errors.join('; ')}`,
				)
			}
			console.info(`Search: indexed ${documents.size} documentation pages.`)
		} finally {
			await pagefind.close()
		}
	},
}

export default adapter
