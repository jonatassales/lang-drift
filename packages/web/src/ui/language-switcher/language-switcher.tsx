'use client'

import { cn } from '@langdrift/web/utilities'
import type { ComponentPropsWithoutRef, ReactNode } from 'react'
import { Fragment, useEffect, useId, useRef, useState } from 'react'

export interface LanguageSwitcherOption<Locale extends string = string> {
	locale: Locale
	label: string
	name: string
	href?: string
}

export type LanguageSwitcherLinkProps = Omit<
	ComponentPropsWithoutRef<'a'>,
	'href'
> & { href: string }

export type LanguageSwitcherProps<Locale extends string = string> = {
	label: string
	currentLocale: string
	className?: string
} & (
	| {
			options: readonly (LanguageSwitcherOption<Locale> & { href: string })[]
			onLocaleChange?: never
			disabled?: never
			renderLink?: (
				option: LanguageSwitcherOption<Locale> & { href: string },
				linkProps: LanguageSwitcherLinkProps,
			) => ReactNode
	  }
	| {
			options: readonly LanguageSwitcherOption<Locale>[]
			onLocaleChange: (locale: Locale) => void
			disabled?: boolean
			renderLink?: never
	  }
)

const optionClassName =
	'inline-flex min-h-[1.7rem] cursor-pointer items-center justify-center whitespace-nowrap rounded-full border-0 bg-transparent px-[0.7rem] text-[0.72rem]! font-[650]! leading-normal! tracking-[0.04em] text-[var(--muted-foreground)] no-underline transition-colors duration-180 hover:text-[var(--foreground)] focus-visible:text-[var(--foreground)] focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[var(--foreground)]! aria-[current=page]:bg-[color-mix(in_srgb,var(--background)_88%,transparent)] aria-[current=page]:text-[var(--foreground)] aria-pressed:bg-[color-mix(in_srgb,var(--background)_88%,transparent)] aria-pressed:text-[var(--foreground)] disabled:cursor-wait disabled:opacity-60 motion-reduce:transition-none'

export function LanguageSwitcher<Locale extends string>(
	props: LanguageSwitcherProps<Locale>,
) {
	const { label, currentLocale, className } = props

	const [open, setOpen] = useState(false)
	const containerRef = useRef<HTMLElement | null>(null)
	const triggerRef = useRef<HTMLButtonElement>(null)
	const optionsRef = useRef<HTMLDivElement>(null)
	const optionsId = useId()
	const currentOption = props.options.find(
		(option) => option.locale === currentLocale,
	)
	const Container = props.onLocaleChange ? 'div' : 'nav'

	useEffect(() => {
		if (!open) return

		function dismissOutside(event: PointerEvent) {
			if (!containerRef.current?.contains(event.target as Node)) {
				setOpen(false)
			}
		}

		document.addEventListener('pointerdown', dismissOutside)

		return () => document.removeEventListener('pointerdown', dismissOutside)
	}, [open])

	function closePill() {
		if (optionsRef.current?.contains(document.activeElement)) {
			triggerRef.current?.focus({ preventScroll: true })
		}
		setOpen(false)
	}

	const choices = props.onLocaleChange
		? props.options
				.filter((option) => option.locale !== currentLocale)
				.map((option) => (
					<button
						aria-label={option.name}
						aria-pressed={currentLocale === option.locale}
						className={cn(
							optionClassName,
							'focus-visible:outline-offset-[-2px]',
						)}
						disabled={props.disabled}
						key={option.locale}
						onClick={() => props.onLocaleChange(option.locale)}
						type="button"
					>
						{option.label}
					</button>
				))
		: props.options
				.filter((option) => option.locale !== currentLocale)
				.map((option) => {
					const linkProps: LanguageSwitcherLinkProps = {
						'aria-current':
							currentLocale === option.locale ? 'page' : undefined,
						'aria-label': option.name,
						className: cn(
							optionClassName,
							'focus-visible:outline-offset-[-2px]',
						),
						href: option.href,
						hrefLang: option.locale,
						children: option.label,
					}

					return (
						<Fragment key={option.locale}>
							{props.renderLink ? (
								props.renderLink(option, linkProps)
							) : (
								<a {...linkProps} />
							)}
						</Fragment>
					)
				})

	return (
		<Container
			aria-label={label}
			className={cn(
				'inline-flex w-max max-w-full shrink-0 flex-nowrap items-center rounded-full border border-[color-mix(in_srgb,var(--border)_80%,transparent)] bg-[color-mix(in_srgb,var(--card)_72%,transparent)] p-[0.2rem] shadow-[0_8px_32px_color-mix(in_srgb,var(--foreground)_8%,transparent)] backdrop-blur-lg',
				className,
			)}
			data-state={open ? 'open' : 'idle'}
			onBlur={(event) => {
				if (!event.currentTarget.contains(event.relatedTarget)) {
					setOpen(false)
				}
			}}
			onKeyDown={(event) => {
				if (event.key === 'Escape' && open) {
					event.preventDefault()
					event.stopPropagation()
					closePill()
				}
			}}
			onPointerEnter={(event) => {
				if (event.pointerType !== 'touch' && !props.disabled) {
					setOpen(true)
				}
			}}
			onPointerLeave={(event) => {
				if (event.pointerType !== 'touch') closePill()
			}}
			ref={(element) => {
				containerRef.current = element
			}}
			role={props.onLocaleChange ? 'group' : undefined}
		>
			<button
				aria-controls={optionsId}
				aria-expanded={open}
				aria-label={`${label}: ${currentOption?.name ?? currentLocale}`}
				className={cn(
					optionClassName,
					'shrink-0 bg-[color-mix(in_srgb,var(--background)_88%,transparent)] text-[var(--foreground)]',
				)}
				disabled={props.disabled}
				onClick={() => setOpen(!open)}
				ref={triggerRef}
				type="button"
			>
				{currentOption?.label ?? currentLocale}
			</button>
			<div
				aria-hidden={!open}
				className={cn(
					'grid grid-cols-[0fr] opacity-0 transition-[grid-template-columns,opacity] duration-250 ease-out motion-reduce:transition-none',
					open && 'grid-cols-[1fr] opacity-100',
				)}
				id={optionsId}
				inert={!open}
				onClickCapture={closePill}
				ref={optionsRef}
			>
				<div className="min-w-0 overflow-hidden">
					<div className="flex w-max items-center gap-[0.15rem] pl-[0.15rem]">
						{choices}
					</div>
				</div>
			</div>
		</Container>
	)
}
