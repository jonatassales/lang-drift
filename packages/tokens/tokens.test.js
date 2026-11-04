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

// The theme the shared contract is read from, and its prefix.
const theme = 'website'
const prefix = 'site'

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

// The roles whose values depend on neither the theme nor the status colors.
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
const classifications = ['intentional', 'aligned', 'unexplained', 'review']
const rainbows = ['1', '2', '3', '4', '5', '6']
const drifts = ['negative', 'positive', 'zero']
const classificationTexts = ['intentional', 'review', 'unexplained']

// The theme's tokens and the shared roles that follow its brand.
const brandRolesOf = (themePrefix) => [
	`--${themePrefix}-color-primary`,
	`--${themePrefix}-color-primary-hover`,
	`--${themePrefix}-color-primary-text`,
	...[
		'primary',
		'primary-hover',
		'primary-text',
		'ring',
		'sidebar-primary',
		'sidebar-ring',
		'chart-1',
	].map((role) => `--lg-color-${role}`),
]
const brandRoles = brandRolesOf(prefix)

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

// The stylesheet an app importing `tokens/<theme>.css` from this
// copy would receive.
function stylesheet(copy, name = theme) {
	const consumer = join(sandbox, `${basename(copy)}-consumer`)
	mkdirSync(join(consumer, 'node_modules'), { recursive: true })
	rmSync(join(consumer, 'node_modules/tokens'), { force: true })
	symlinkSync(copy, join(consumer, 'node_modules/tokens'))
	const require = createRequire(join(consumer, 'index.js'))
	return readFileSync(require.resolve(`tokens/${name}.css`), 'utf8')
}

// The theme's stylesheet, built once from its own copy.
let builtStylesheet
function themeStylesheet() {
	if (!builtStylesheet) {
		const copy = copyPackage(theme)
		const result = build(copy)
		assert.equal(result.status, 0, result.stderr)
		builtStylesheet = stylesheet(copy)
	}
	return builtStylesheet
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
		css = themeStylesheet()
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

	test('declare the website fonts, and the shadow and easing kept from Exphere', () => {
		assert.deepEqual(
			pick(light, [
				'--lg-typography-font-body',
				'--lg-typography-font-mono',
				'--lg-effects-shadow-sm',
				'--lg-effects-easing-standard',
			]),
			{
				'--lg-typography-font-body':
					'ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif',
				'--lg-typography-font-mono':
					'ui-monospace, "SFMono-Regular", Consolas, "Liberation Mono", monospace',
				'--lg-effects-shadow-sm': '0px 1px 2px #00000014',
				'--lg-effects-easing-standard': 'cubic-bezier(0.2, 0, 0, 1)',
			},
		)
	})

	test('leave layout and the type scale to Tailwind', () => {
		assert.doesNotMatch(css, /--lg-dimensions-|--lg-typography-size-/)
	})

	test('declare the light roles with the website neutrals', () => {
		assert.deepEqual(pick(light, sharedRoles), {
			'--lg-color-background': '#F7F6F2',
			'--lg-color-foreground': '#212121',
			'--lg-color-card': '#FFFFFF',
			'--lg-color-card-foreground': 'var(--lg-color-foreground)',
			'--lg-color-popover': '#FFFFFF',
			'--lg-color-popover-foreground': 'var(--lg-color-foreground)',
			'--lg-color-primary-foreground': '#FFFFFF',
			'--lg-color-secondary': 'var(--lg-color-gray-100)',
			'--lg-color-secondary-foreground': 'var(--lg-color-foreground)',
			'--lg-color-muted': 'var(--lg-color-gray-100)',
			'--lg-color-muted-foreground': '#6B6B75',
			'--lg-color-accent': 'var(--lg-color-gray-100)',
			'--lg-color-accent-foreground': 'var(--lg-color-foreground)',
			'--lg-color-border': '#D9D9DD',
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

	test('declare the dark roles under .dark with the website neutrals', () => {
		assert.deepEqual(pick(dark, sharedRoles), {
			'--lg-color-background': '#0B0B0C',
			'--lg-color-foreground': '#F5F5F5',
			'--lg-color-card': '#151517',
			'--lg-color-card-foreground': 'var(--lg-color-foreground)',
			'--lg-color-popover': '#151517',
			'--lg-color-popover-foreground': 'var(--lg-color-foreground)',
			'--lg-color-primary-foreground': '#FFFFFF',
			'--lg-color-secondary': 'var(--lg-color-gray-800)',
			'--lg-color-secondary-foreground': 'var(--lg-color-foreground)',
			'--lg-color-muted': 'var(--lg-color-gray-800)',
			'--lg-color-muted-foreground': '#A1A1AA',
			'--lg-color-accent': 'var(--lg-color-gray-800)',
			'--lg-color-accent-foreground': 'var(--lg-color-foreground)',
			'--lg-color-border': '#2A2A2D',
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

	test('fix the classification colors in every color scheme', () => {
		assert.deepEqual(
			pick(
				light,
				classifications.map((name) => `--lg-classification-${name}`),
			),
			{
				'--lg-classification-intentional': '#3B82F6',
				'--lg-classification-aligned': '#16A34A',
				'--lg-classification-unexplained': '#DC2626',
				'--lg-classification-review': '#A855F7',
			},
		)
		for (const name of classifications) {
			assert.equal(dark[`--lg-classification-${name}`], undefined)
		}
	})

	test('fix the rainbow colors in every color scheme', () => {
		assert.deepEqual(
			pick(
				light,
				rainbows.map((stop) => `--lg-rainbow-${stop}`),
			),
			{
				'--lg-rainbow-1': '#F43F5E',
				'--lg-rainbow-2': '#F97316',
				'--lg-rainbow-3': '#FACC15',
				'--lg-rainbow-4': '#22C55E',
				'--lg-rainbow-5': '#3B82F6',
				'--lg-rainbow-6': '#6D0AD6',
			},
		)
		for (const stop of rainbows) {
			assert.equal(dark[`--lg-rainbow-${stop}`], undefined)
		}
	})

	test('fix the classification text tones, darker in light and lighter in dark', () => {
		const textVariables = classificationTexts.map(
			(name) => `--lg-classification-text-${name}`,
		)
		assert.deepEqual(pick(light, textVariables), {
			'--lg-classification-text-intentional': 'oklch(48.8% 0.243 264.376)',
			'--lg-classification-text-review': 'oklch(49.6% 0.265 301.924)',
			'--lg-classification-text-unexplained': 'oklch(50.5% 0.213 27.518)',
		})
		assert.deepEqual(pick(dark, textVariables), {
			'--lg-classification-text-intentional': 'oklch(70.7% 0.165 254.624)',
			'--lg-classification-text-review': 'oklch(71.4% 0.203 305.504)',
			'--lg-classification-text-unexplained': 'oklch(70.4% 0.191 22.216)',
		})
	})

	test('fix the drift colors, darker in light and lighter in dark', () => {
		const driftVariables = drifts.map((drift) => `--lg-drift-${drift}`)
		assert.deepEqual(pick(light, driftVariables), {
			'--lg-drift-negative': '#B91C1C',
			'--lg-drift-positive': '#15803D',
			'--lg-drift-zero': '#C2410C',
		})
		assert.deepEqual(pick(dark, driftVariables), {
			'--lg-drift-negative': '#F87171',
			'--lg-drift-positive': '#4ADE80',
			'--lg-drift-zero': '#FB923C',
		})
	})

	test('never let the domain colors take the brand', () => {
		for (const declared of [light, dark]) {
			for (const [name, value] of Object.entries(declared)) {
				if (/^--lg-(classification|classification-text|drift)-/.test(name)) {
					assert.doesNotMatch(value, /primary|--site-/)
				}
			}
		}
	})

	test('back the destructive role with danger', () => {
		assert.equal(light['--lg-color-destructive'], 'var(--lg-status-danger)')
		assert.equal(dark['--lg-color-destructive'], 'var(--lg-status-danger)')
	})
})

describe('website theme', () => {
	test('binds the shared primary roles to the Exphere purple', () => {
		const css = themeStylesheet()
		for (const [selector, primary, hover, text] of [
			[':root', '#6D0AD6', '#5408A5', '#6D0AD6'],
			['.dark', '#6D0AD6', '#8C3BEF', '#A970FF'],
		]) {
			assert.deepEqual(pick(declarations(css, selector), brandRoles), {
				[`--${prefix}-color-primary`]: primary,
				[`--${prefix}-color-primary-hover`]: hover,
				[`--${prefix}-color-primary-text`]: text,
				'--lg-color-primary': `var(--${prefix}-color-primary)`,
				'--lg-color-primary-hover': `var(--${prefix}-color-primary-hover)`,
				'--lg-color-primary-text': `var(--${prefix}-color-primary-text)`,
				'--lg-color-ring': 'var(--lg-color-primary)',
				'--lg-color-sidebar-primary': 'var(--lg-color-primary)',
				'--lg-color-sidebar-ring': 'var(--lg-color-primary)',
				'--lg-color-chart-1': 'var(--lg-color-primary)',
			})
		}
	})
})

describe('shadcn bridge', () => {
	test('declares every shadcn variable as a reference to its role', () => {
		const bridge = declarations(themeStylesheet(), ':root, .dark')
		for (const variable of shadcnVariables) {
			assert.equal(bridge[`--${variable}`], `var(--lg-color-${variable})`)
		}
	})

	test('leaves the radius to Tailwind', () => {
		assert.doesNotMatch(themeStylesheet(), /--radius/)
	})
})

describe('Tailwind theme', () => {
	let tailwind

	before(() => {
		tailwind = Object.entries(declarations(themeStylesheet(), '@theme inline'))
	})

	test('resets the native palette before any token color', () => {
		assert.deepEqual(tailwind[0], ['--color-*', 'initial'])
	})

	test('exposes every role, status and domain color, and no foundation', () => {
		const colors = Object.fromEntries(
			tailwind.filter(
				([name]) => name.startsWith('--color-') && name !== '--color-*',
			),
		)
		for (const variable of shadcnVariables) {
			assert.equal(colors[`--color-${variable}`], `var(--${variable})`)
		}
		assert.equal(colors['--color-primary-hover'], 'var(--primary-hover)')
		assert.equal(colors['--color-primary-text'], 'var(--primary-text)')
		for (const status of statuses) {
			assert.equal(
				colors[`--color-status-${status}`],
				`var(--lg-status-${status})`,
			)
		}
		for (const name of classifications) {
			assert.equal(
				colors[`--color-classification-${name}`],
				`var(--lg-classification-${name})`,
			)
		}
		for (const name of classificationTexts) {
			assert.equal(
				colors[`--color-classification-text-${name}`],
				`var(--lg-classification-text-${name})`,
			)
		}
		for (const drift of drifts) {
			assert.equal(colors[`--color-drift-${drift}`], `var(--lg-drift-${drift})`)
		}
		assert.equal(
			Object.keys(colors).some((name) => name.includes('gray')),
			false,
		)
	})

	test('maps the fonts, shadow and easing, and nothing else', () => {
		assert.deepEqual(
			tailwind.filter(([name]) => !name.startsWith('--color-')),
			[
				['--font-sans', 'var(--lg-typography-font-body)'],
				['--font-mono', 'var(--lg-typography-font-mono)'],
				['--shadow-sm', 'var(--lg-effects-shadow-sm)'],
				['--ease-standard', 'var(--lg-effects-easing-standard)'],
			],
		)
	})
})

describe('dark variant', () => {
	test('follows the .dark class', () => {
		assert.match(
			themeStylesheet(),
			/^@custom-variant dark \(&:where\(\.dark, \.dark \*\)\);$/m,
		)
	})
})

describe('build', () => {
	test('yields identical stylesheets from unchanged sources', () => {
		const copy = copyPackage('rebuild')
		assert.equal(build(copy).status, 0)
		const first = stylesheet(copy)
		assert.equal(build(copy).status, 0)
		assert.equal(stylesheet(copy), first)
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
