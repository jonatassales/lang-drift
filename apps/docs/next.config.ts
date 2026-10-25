import path from 'node:path'
import type { NextConfig } from 'next'
import nextra from 'nextra'

const withNextra = nextra({})

const mdxComponents = './src/mdx-components.tsx'

const nextConfig: NextConfig = {
	basePath: '/docs',
	adapterPath: path.resolve(process.cwd(), 'next.adapter.mts'),
	transpilePackages: ['@langdrift/web'],
	turbopack: {
		resolveAlias: { 'next-mdx-import-source-file': mdxComponents },
	},
	webpack(config) {
		config.resolve.alias['next-mdx-import-source-file'] = path.resolve(
			process.cwd(),
			mdxComponents,
		)

		return config
	},
}

export default withNextra(nextConfig)
