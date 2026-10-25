import { cleanup, fireEvent, render, screen } from '@testing-library/react'
import { afterEach, describe, expect, it, vi } from 'vitest'

import { LanguageSwitcher } from './language-switcher'

afterEach(cleanup)

const options = [
	{
		locale: 'en',
		label: 'EN',
		name: 'English',
		href: '/settings?view=all#language',
	},
	{
		locale: 'pt-BR',
		label: 'PT-BR',
		name: 'Português',
		href: '/pt-BR/settings?view=all#language',
	},
] as const

describe('LanguageSwitcher', () => {
	it('reveals alternatives on hover and hides them when the pointer leaves', () => {
		render(
			<LanguageSwitcher
				currentLocale="en"
				label="Language"
				options={options}
			/>,
		)

		const pill = screen.getByRole('navigation', { name: 'Language' })
		const trigger = screen.getByRole('button', { name: 'Language: English' })
		expect(trigger.textContent).toBe('EN')
		expect(trigger.getAttribute('aria-expanded')).toBe('false')
		expect(screen.queryByRole('link')).toBeNull()

		fireEvent.pointerEnter(pill)
		expect(trigger.getAttribute('aria-expanded')).toBe('true')
		expect(screen.getByRole('link', { name: 'Português' })).toBeTruthy()

		fireEvent.pointerLeave(pill)
		expect(trigger.getAttribute('aria-expanded')).toBe('false')
		expect(screen.queryByRole('link')).toBeNull()
	})

	it('supports activation and Escape without losing focus in hidden options', () => {
		render(
			<LanguageSwitcher
				currentLocale="en"
				label="Language"
				options={options}
			/>,
		)

		const trigger = screen.getByRole('button', { name: 'Language: English' })
		fireEvent.click(trigger)
		const portuguese = screen.getByRole('link', { name: 'Português' })
		portuguese.focus()
		fireEvent.keyDown(portuguese, { key: 'Escape' })

		expect(trigger.getAttribute('aria-expanded')).toBe('false')
		expect(document.activeElement).toBe(trigger)
		expect(screen.queryByRole('link')).toBeNull()

		fireEvent.click(trigger)
		fireEvent.pointerDown(document.body)
		expect(trigger.getAttribute('aria-expanded')).toBe('false')
	})

	it('closes after a modified link click while preserving the navigation event', () => {
		const navigate = vi.fn()
		render(
			<LanguageSwitcher
				currentLocale="en"
				label="Language"
				options={options}
				renderLink={(option, linkProps) => (
					<a
						{...linkProps}
						href={linkProps.href}
						onClick={(event) => {
							navigate(option.locale, event.ctrlKey, event.defaultPrevented)
							event.preventDefault()
						}}
					/>
				)}
			/>,
		)

		const trigger = screen.getByRole('button', { name: 'Language: English' })
		fireEvent.click(trigger)
		fireEvent.click(screen.getByRole('link', { name: 'Português' }), {
			ctrlKey: true,
		})
		expect(navigate).toHaveBeenCalledWith('pt-BR', true, false)
		expect(trigger.getAttribute('aria-expanded')).toBe('false')
	})

	it('keeps pill instances independent and updates the displayed current language', () => {
		const onLocaleChange = vi.fn()
		const { rerender } = render(
			<>
				<LanguageSwitcher
					currentLocale="en"
					label="Website"
					options={options}
				/>
				<LanguageSwitcher
					currentLocale="pt-BR"
					label="Drawer"
					options={options}
				/>
			</>,
		)
		fireEvent.click(screen.getByRole('button', { name: 'Website: English' }))
		expect(
			screen
				.getByRole('button', { name: 'Drawer: Português' })
				.getAttribute('aria-expanded'),
		).toBe('false')

		fireEvent.keyDown(
			screen.getByRole('button', { name: 'Website: English' }),
			{
				key: 'Escape',
			},
		)
		rerender(
			<LanguageSwitcher
				currentLocale="en"
				label="Language"
				onLocaleChange={onLocaleChange}
				options={options}
			/>,
		)
		fireEvent.click(screen.getByRole('button', { name: 'Language: English' }))
		fireEvent.click(screen.getByRole('button', { name: 'Português' }))
		expect(onLocaleChange).toHaveBeenCalledWith('pt-BR')

		rerender(
			<LanguageSwitcher
				currentLocale="pt-BR"
				label="Language"
				onLocaleChange={onLocaleChange}
				options={options}
			/>,
		)
		const trigger = screen.getByRole('button', { name: 'Language: Português' })
		expect(trigger.textContent).toBe('PT-BR')
		expect(trigger.getAttribute('aria-expanded')).toBe('false')
	})

	it('supports a disabled preference control without inventing language URLs', () => {
		const onLocaleChange = vi.fn()
		const buttonOptions = options.map(({ locale, label, name }) => ({
			locale,
			label,
			name,
		}))
		const { rerender } = render(
			<LanguageSwitcher
				currentLocale="en"
				label="Language"
				onLocaleChange={onLocaleChange}
				options={buttonOptions}
			/>,
		)

		expect(screen.getByRole('group', { name: 'Language' })).toBeTruthy()
		expect(screen.queryByRole('link')).toBeNull()
		const trigger = screen.getByRole('button', { name: 'Language: English' })
		fireEvent.click(trigger)
		expect(
			screen
				.getByRole('button', { name: 'Português' })
				.getAttribute('aria-pressed'),
		).toBe('false')
		fireEvent.click(screen.getByRole('button', { name: 'Português' }))
		expect(onLocaleChange).toHaveBeenCalledWith('pt-BR')

		rerender(
			<LanguageSwitcher
				currentLocale="en"
				disabled
				label="Language"
				onLocaleChange={onLocaleChange}
				options={buttonOptions}
			/>,
		)
		fireEvent.click(trigger)
		fireEvent.pointerEnter(screen.getByRole('group', { name: 'Language' }))
		expect(trigger.hasAttribute('disabled')).toBe(true)
		expect(trigger.getAttribute('aria-expanded')).toBe('false')
		expect(onLocaleChange).toHaveBeenCalledTimes(1)
	})

	it('exposes native alternative links and names the current language', () => {
		render(
			<LanguageSwitcher
				currentLocale="pt-BR"
				label="Language"
				options={options}
			/>,
		)

		expect(screen.getByRole('navigation', { name: 'Language' })).toBeTruthy()
		const portuguese = screen.getByRole('button', {
			name: 'Language: Português',
		})
		fireEvent.click(portuguese)
		const english = screen.getByRole('link', { name: 'English' })
		expect(english.getAttribute('href')).toBe(options[0].href)
		expect(english.getAttribute('aria-current')).toBeNull()
		expect(english.getAttribute('hrefLang')).toBe('en')
		expect(portuguese.textContent).toBe('PT-BR')
	})

	it('lets a router adapter preserve anchor semantics and own navigation events', () => {
		const navigate = vi.fn()
		render(
			<LanguageSwitcher
				currentLocale="en"
				label="Language"
				options={options}
				renderLink={(option, linkProps) => (
					<a
						{...linkProps}
						href={linkProps.href}
						onClick={(event) => {
							navigate(option.locale, event.ctrlKey, event.defaultPrevented)
							event.preventDefault()
						}}
					/>
				)}
			/>,
		)

		const trigger = screen.getByRole('button', { name: 'Language: English' })
		fireEvent.click(trigger)
		const portuguese = screen.getByRole('link', { name: 'Português' })
		fireEvent.click(portuguese, { ctrlKey: true })
		expect(navigate).toHaveBeenCalledWith('pt-BR', true, false)
		expect(portuguese.getAttribute('href')).toBe(options[1].href)
		expect(trigger.textContent).toBe('EN')
		expect(trigger.getAttribute('aria-expanded')).toBe('false')
	})
})
