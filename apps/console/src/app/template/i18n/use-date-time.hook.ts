'use client'
import { type DateTimeFormatOptions, useFormatter } from 'next-intl'

/** ICU versions vary in date/time separator width; preserve locale text with stable spaces. */
export function useDateTime() {
	const format = useFormatter()

	return (value: Date | number, options?: DateTimeFormatOptions) =>
		format.dateTime(value, options).replace(/[\u00a0\u2009\u202f]/g, ' ')
}
