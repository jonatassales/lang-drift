// The codes a demo request's fields fail with; the book-a-demo catalog
// holds a message for each.
export const bookDemoErrorCodes = [
	'emailRequired',
	'emailInvalid',
	'emailTooLong',
	'emailFreeMail',
	'companyRequired',
	'companyTooLong',
	'positionRequired',
] as const
