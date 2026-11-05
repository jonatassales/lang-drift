import { fireEvent, screen } from '@testing-library/react'

import { bookDemoCopy } from './book-demo-copy.data'

interface BookDemoFields {
	company?: string
	email?: string
	position?: string
}

export function fillBookDemoForm(fields: BookDemoFields) {
	if (fields.email !== undefined) {
		fireEvent.change(screen.getByLabelText(bookDemoCopy.form.email), {
			target: { value: fields.email },
		})
	}
	if (fields.company !== undefined) {
		fireEvent.change(screen.getByLabelText(bookDemoCopy.form.company), {
			target: { value: fields.company },
		})
	}
	if (fields.position !== undefined) {
		fireEvent.change(screen.getByLabelText(bookDemoCopy.form.position), {
			target: { value: fields.position },
		})
	}
}
