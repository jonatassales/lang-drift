import { render } from '@testing-library/react'

import { BookDemoForm } from '@routes/book-a-demo/book-demo-form'

import { bookDemoCopy } from './book-demo-copy.data'

export function mountBookDemoForm() {
	return render(<BookDemoForm copy={bookDemoCopy} locale="en" />)
}
