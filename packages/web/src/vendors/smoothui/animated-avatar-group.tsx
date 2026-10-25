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
					className="avatar-chip"
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
						<img alt="" aria-hidden="true" src={person.src} />
					) : (
						person.initials
					)}
				</span>
			))}
		</fieldset>
	)
}
