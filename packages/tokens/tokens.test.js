import assert from 'node:assert/strict'
import { spawnSync } from 'node:child_process'
import {
	cpSync,
	mkdirSync,
	mkdtempSync,
	readFileSync,
	rmSync,
	symlinkSync,
	writeFileSync,
} from 'node:fs'
import { createRequire } from 'node:module'
import { tmpdir } from 'node:os'
import { basename, join } from 'node:path'
import { after, before, describe, test } from 'node:test'

// The contract is the published stylesheet, read the way an app reads it:
// built from an isolated copy and resolved through the package's exports.

// Each product, its product theme prefix and the other product's prefix.
const products = [
	['vision', 'viz', 'fp'],
	['finops', 'fp', 'viz'],
]

// The variables shadcn and SmoothUI components read, without the leading `--`.
const shadcnVariables = [
	'background',
	'foreground',
	'card',
	'card-foreground',
	'popover',
	'popover-foreground',
	'primary',
	'primary-foreground',
	'secondary',
	'secondary-foreground',
	'muted',
	'muted-foreground',
	'accent',
	'accent-foreground',
	'destructive',
	'border',
	'input',
	'ring',
	'chart-1',
	'chart-2',
	'chart-3',
	'chart-4',
	'chart-5',
	'sidebar',
	'sidebar-foreground',
	'sidebar-primary',
	'sidebar-primary-foreground',
	'sidebar-accent',
	'sidebar-accent-foreground',
	'sidebar-border',
	'sidebar-ring',
]

// The roles whose values depend on neither the product theme nor the status
// colors.
const sharedRoles = [
	'background',
	'foreground',
	'card',
	'card-foreground',
	'popover',
	'popover-foreground',
	'primary-foreground',
	'secondary',
	'secondary-foreground',
	'muted',
	'muted-foreground',
	'accent',
	'accent-foreground',
	'border',
	'input',
	'chart-2',
	'chart-3',
	'chart-4',
	'chart-5',
	'sidebar',
	'sidebar-foreground',
	'sidebar-primary-foreground',
	'sidebar-accent',
	'sidebar-accent-foreground',
	'sidebar-border',
].map((role) => `--lg-color-${role}`)

const statuses = ['success', 'warning', 'danger', 'info']

// A product theme's tokens and the shared roles that follow its brand.
function brandRoles(prefix) {
	return [
		`--${prefix}-color-primary`,
		`--${prefix}-color-primary-hover`,
		...[
			'primary',
			'primary-hover',
			'ring',
			'sidebar-primary',
			'sidebar-ring',
			'chart-1',
		].map((role) => `--lg-color-${role}`),
	]
}

const sandbox = mkdtempSync(join(tmpdir(), 'langdrift-tokens-'))
after(() =>
	rmSync(sandbox, {
		recursive: true,
		force: true,
	}),
)

function copyPackage(name) {
	const copy = join(sandbox, name)
	cpSync(import.meta.dirname, copy, {
		recursive: true,
		filter: (path) => !/\/(node_modules|dist)(\/|$)/.test(path),
	})
	symlinkSync(
		join(import.meta.dirname, 'node_modules'),
		join(copy, 'node_modules'),
	)
	return copy
}

function build(copy) {
	return spawnSync('pnpm', ['run', '--silent', 'build'], {
		cwd: copy,
		encoding: 'utf8',
	})
}

// Builds a copy whose source file has been changed by `edit`.
function buildEdited(name, sourceFile, edit) {
	const copy = copyPackage(name)
	const source = join(copy, sourceFile)
	const tokens = JSON.parse(readFileSync(source, 'utf8'))
	edit(tokens)
	writeFileSync(source, JSON.stringify(tokens))
	return build(copy)
}

// The stylesheet an app importing `@langdrift/tokens/<product>.css` from
// this copy would receive.
function stylesheet(copy, product) {
	const consumer = join(sandbox, `${basename(copy)}-consumer`)
	mkdirSync(join(consumer, 'node_modules/@langdrift'), { recursive: true })
	rmSync(join(consumer, 'node_modules/@langdrift/tokens'), { force: true })
	symlinkSync(copy, join(consumer, 'node_modules/@langdrift/tokens'))
	const require = createRequire(join(consumer, 'index.js'))
	return readFileSync(
		require.resolve(`@langdrift/tokens/${product}.css`),
		'utf8',
	)
}

// Each product's stylesheet, built once from its own copy.
const productStylesheets = new Map()
function productStylesheet(product) {
	if (!productStylesheets.has(product)) {
		const copy = copyPackage(product)
		const result = build(copy)
		assert.equal(result.status, 0, result.stderr)
		productStylesheets.set(product, stylesheet(copy, product))
	}
	return productStylesheets.get(product)
}

// The declarations of the top-level rule whose selector is `selector`.
function declarations(css, selector) {
	const escaped = selector.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')
	const start = css.search(new RegExp(`^${escaped} \\{$`, 'm'))
	assert.notEqual(start, -1, `missing rule ${selector}`)
	const body = css.slice(css.indexOf('{', start) + 1, css.indexOf('\n}', start))
	return Object.fromEntries(
		body
			.split(';')
			.map((line) => line.trim())
			.filter(Boolean)
			.map((line) => {
				const colon = line.indexOf(':')
				return [line.slice(0, colon).trim(), line.slice(colon + 1).trim()]
			}),
	)
}

function pick(declared, names) {
	return Object.fromEntries(names.map((name) => [name, declared[name]]))
}

describe('shared tokens', () => {
	let css
	let light
	let dark

	before(() => {
		css = productStylesheet('vision')
		light = declarations(css, ':root')
		dark = declarations(css, '.dark')
	})

	test('declare the gray foundation, including gray 400', () => {
		assert.deepEqual(
			Object.entries(light).filter(([name]) =>
				name.startsWith('--lg-color-gray-'),
			),
			[
				['--lg-color-gray-50', '#fafafa'],
				['--lg-color-gray-100', '#f5f5f5'],
				['--lg-color-gray-200', '#e5e5e5'],
				['--lg-color-gray-400', '#a3a3a3'],
				['--lg-color-gray-800', '#262626'],
				['--lg-color-gray-900', '#171717'],
				['--lg-color-gray-950', '#0a0a0a'],
			],
		)
	})

	test('declare the font family, shadow and easing kept from Exphere', () => {
		assert.deepEqual(
			pick(light, [
				'--lg-typography-font-body',
				'--lg-effects-shadow-sm',
				'--lg-effects-easing-standard',
			]),
			{
				'--lg-typography-font-body': 'system-ui, sans-serif',
				'--lg-effects-shadow-sm': '0px 1px 2px #00000014',
				'--lg-effects-easing-standard': 'cubic-bezier(0.2, 0, 0, 1)',
			},
		)
	})

	test('leave layout and the type scale to Tailwind', () => {
		assert.doesNotMatch(css, /--lg-dimensions-|--lg-typography-size-/)
	})

	test('declare the light roles, keeping the Exphere palette', () => {
		assert.deepEqual(pick(light, sharedRoles), {
			'--lg-color-background': '#FFFFFF',
			'--lg-color-foreground': '#27262A',
			'--lg-color-card': '#F9F9F9',
			'--lg-color-card-foreground': 'var(--lg-color-foreground)',
			'--lg-color-popover': '#FFFFFF',
			'--lg-color-popover-foreground': 'var(--lg-color-foreground)',
			'--lg-color-primary-foreground': '#FFFFFF',
			'--lg-color-secondary': 'var(--lg-color-gray-100)',
			'--lg-color-secondary-foreground': 'var(--lg-color-foreground)',
			'--lg-color-muted': 'var(--lg-color-gray-100)',
			'--lg-color-muted-foreground': '#535257',
			'--lg-color-accent': 'var(--lg-color-gray-100)',
			'--lg-color-accent-foreground': 'var(--lg-color-foreground)',
			'--lg-color-border': '#DDDDDD',
			'--lg-color-input': 'var(--lg-color-border)',
			'--lg-color-chart-2': '#0D9488',
			'--lg-color-chart-3': '#D97706',
			'--lg-color-chart-4': '#0284C7',
			'--lg-color-chart-5': '#E11D48',
			'--lg-color-sidebar': 'var(--lg-color-card)',
			'--lg-color-sidebar-foreground': 'var(--lg-color-foreground)',
			'--lg-color-sidebar-primary-foreground':
				'var(--lg-color-primary-foreground)',
			'--lg-color-sidebar-accent': 'var(--lg-color-accent)',
			'--lg-color-sidebar-accent-foreground':
				'var(--lg-color-accent-foreground)',
			'--lg-color-sidebar-border': 'var(--lg-color-border)',
		})
	})

	test('declare the dark roles under .dark', () => {
		assert.deepEqual(pick(dark, sharedRoles), {
			'--lg-color-background': 'var(--lg-color-gray-950)',
			'--lg-color-foreground': 'var(--lg-color-gray-50)',
			'--lg-color-card': 'var(--lg-color-gray-900)',
			'--lg-color-card-foreground': 'var(--lg-color-foreground)',
			'--lg-color-popover': 'var(--lg-color-gray-900)',
			'--lg-color-popover-foreground': 'var(--lg-color-foreground)',
			'--lg-color-primary-foreground': 'var(--lg-color-gray-950)',
			'--lg-color-secondary': 'var(--lg-color-gray-800)',
			'--lg-color-secondary-foreground': 'var(--lg-color-foreground)',
			'--lg-color-muted': 'var(--lg-color-gray-800)',
			'--lg-color-muted-foreground': 'var(--lg-color-gray-400)',
			'--lg-color-accent': 'var(--lg-color-gray-800)',
			'--lg-color-accent-foreground': 'var(--lg-color-foreground)',
			'--lg-color-border': 'var(--lg-color-gray-800)',
			'--lg-color-input': 'var(--lg-color-border)',
			'--lg-color-chart-2': '#2DD4BF',
			'--lg-color-chart-3': '#FBBF24',
			'--lg-color-chart-4': '#38BDF8',
			'--lg-color-chart-5': '#FB7185',
			'--lg-color-sidebar': 'var(--lg-color-card)',
			'--lg-color-sidebar-foreground': 'var(--lg-color-foreground)',
			'--lg-color-sidebar-primary-foreground':
				'var(--lg-color-primary-foreground)',
			'--lg-color-sidebar-accent': 'var(--lg-color-accent)',
			'--lg-color-sidebar-accent-foreground':
				'var(--lg-color-accent-foreground)',
			'--lg-color-sidebar-border': 'var(--lg-color-border)',
		})
	})

	test('fix the status colors, darker in light and lighter in dark', () => {
		const statusVariables = statuses.map((status) => `--lg-status-${status}`)
		assert.deepEqual(pick(light, statusVariables), {
			'--lg-status-success': '#15803D',
			'--lg-status-warning': '#B45309',
			'--lg-status-danger': '#B91C1C',
			'--lg-status-info': '#1D4ED8',
		})
		assert.deepEqual(pick(dark, statusVariables), {
			'--lg-status-success': '#4ADE80',
			'--lg-status-warning': '#FBBF24',
			'--lg-status-danger': '#F87171',
			'--lg-status-info': '#60A5FA',
		})
	})

	test('back the destructive role with danger', () => {
		assert.equal(light['--lg-color-destructive'], 'var(--lg-status-danger)')
		assert.equal(dark['--lg-color-destructive'], 'var(--lg-status-danger)')
	})
})

describe('product themes', () => {
	for (const [product, prefix, otherPrefix] of products) {
		test(`bind the shared primary roles to ${product}'s brand`, () => {
			const css = productStylesheet(product)
			for (const [selector, primary, hover] of [
				[':root', '#6D0AD6', '#5408A5'],
				['.dark', '#A66BF0', '#B98AF4'],
			]) {
				assert.deepEqual(
					pick(declarations(css, selector), brandRoles(prefix)),
					{
						[`--${prefix}-color-primary`]: primary,
						[`--${prefix}-color-primary-hover`]: hover,
						'--lg-color-primary': `var(--${prefix}-color-primary)`,
						'--lg-color-primary-hover': `var(--${prefix}-color-primary-hover)`,
						'--lg-color-ring': 'var(--lg-color-primary)',
						'--lg-color-sidebar-primary': 'var(--lg-color-primary)',
						'--lg-color-sidebar-ring': 'var(--lg-color-primary)',
						'--lg-color-chart-1': 'var(--lg-color-primary)',
					},
				)
			}
		})

		test(`keep ${product}'s stylesheet free of the other product theme`, () => {
			assert.doesNotMatch(
				productStylesheet(product),
				new RegExp(`--${otherPrefix}-`),
			)
		})
	}
})

describe('shadcn bridge', () => {
	test('declares every shadcn variable as a reference to its role', () => {
		const bridge = declarations(productStylesheet('vision'), ':root, .dark')
		for (const variable of shadcnVariables) {
			assert.equal(bridge[`--${variable}`], `var(--lg-color-${variable})`)
		}
	})

	test('leaves the radius to Tailwind', () => {
		assert.doesNotMatch(productStylesheet('vision'), /--radius/)
	})
})

describe('Tailwind theme', () => {
	let theme

	before(() => {
		theme = Object.entries(
			declarations(productStylesheet('vision'), '@theme inline'),
		)
	})

	test('resets the native palette before any token color', () => {
		assert.deepEqual(theme[0], ['--color-*', 'initial'])
	})

	test('exposes every role and status color, and no foundation', () => {
		const colors = Object.fromEntries(
			theme.filter(
				([name]) => name.startsWith('--color-') && name !== '--color-*',
			),
		)
		for (const variable of shadcnVariables) {
			assert.equal(colors[`--color-${variable}`], `var(--${variable})`)
		}
		assert.equal(colors['--color-primary-hover'], 'var(--primary-hover)')
		for (const status of statuses) {
			assert.equal(
				colors[`--color-status-${status}`],
				`var(--lg-status-${status})`,
			)
		}
		assert.equal(
			Object.keys(colors).some((name) => name.includes('gray')),
			false,
		)
	})

	test('maps the font, shadow and easing, and nothing else', () => {
		assert.deepEqual(
			theme.filter(([name]) => !name.startsWith('--color-')),
			[
				['--font-sans', 'var(--lg-typography-font-body)'],
				['--shadow-sm', 'var(--lg-effects-shadow-sm)'],
				['--ease-standard', 'var(--lg-effects-easing-standard)'],
			],
		)
	})
})

describe('build', () => {
	test('yields identical stylesheets from unchanged sources', () => {
		const copy = copyPackage('rebuild')
		assert.equal(build(copy).status, 0)
		const first = products.map(([product]) => stylesheet(copy, product))
		assert.equal(build(copy).status, 0)
		assert.deepEqual(
			products.map(([product]) => stylesheet(copy, product)),
			first,
		)
	})

	test('fails when a token references a missing token', () => {
		const result = buildEdited(
			'missing-alias',
			'shared/light.json',
			(tokens) => {
				tokens.light.lg.color.card.$value = '{lg.color.gray.300}'
			},
		)
		assert.notEqual(result.status, 0)
		assert.match(result.stderr, /lg\.color\.gray\.300/)
	})

	test('fails when references form a cycle', () => {
		const result = buildEdited(
			'circular-alias',
			'shared/light.json',
			(tokens) => {
				tokens.light.lg.color.foreground.$value =
					'{light.lg.color.card-foreground}'
			},
		)
		assert.notEqual(result.status, 0)
		assert.match(result.stderr, /Circular definition/)
	})
})
