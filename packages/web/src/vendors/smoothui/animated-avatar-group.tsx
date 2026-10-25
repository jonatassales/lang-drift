'use client'

import { useState } from 'react'

export interface AvatarPerson {
	initials: string
	name: string
	role?: string
	src?: string
}

export interface AnimatedAvatarGroupProps {
	people: AvatarPerson[]
	size?: number
}

export function AnimatedAvatarGroup({
	people,
	size = 28,
}: AnimatedAvatarGroupProps) {
	const [expanded, setExpanded] = useState(false)

	return (
		<fieldset
			aria-label={people.map((person) => person.name).join(', ')}
			className="avatar-group"
			onBlur={() => setExpanded(false)}
			onFocus={() => setExpanded(true)}
			onMouseEnter={() => setExpanded(true)}
			onMouseLeave={() => setExpanded(false)}
			style={{
				border: 0,
				margin: 0,
				minInlineSize: 0,
				padding: 0,
			}}
		>
			{people.map((person, index) => (
				<span
					className="inline-flex items-center justify-center overflow-hidden rounded-[999px] border-2 border-card bg-foreground font-bold text-[9px] text-card transition-[margin-inline-start] duration-180 ease-[ease] motion-reduce:transition-none"
					key={person.name}
					style={{
						height: size,
						marginLeft: index === 0 || expanded ? 0 : -size * 0.28,
						width: size,
						zIndex: people.length - index,
					}}
					title={`${person.name}${person.role ? ` · ${person.role}` : ''}`}
				>
					{person.src ? (
						<img
							alt=""
							aria-hidden="true"
							className="block h-full w-full object-cover"
							src={person.src}
						/>
					) : (
						person.initials
					)}
				</span>
			))}
		</fieldset>
	)
}
