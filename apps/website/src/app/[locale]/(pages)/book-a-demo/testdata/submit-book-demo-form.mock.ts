import { fireEvent, screen } from '@testing-library/react'

import { bookDemoCopy } from './book-demo-copy.data'

export function submitBookDemoForm() {
	fireEvent.click(
		screen.getByRole('button', { name: bookDemoCopy.form.submit }),
	)
}
