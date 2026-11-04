// What a send leaves behind for the record: the provider's id and the
// recipient.
export interface DeliveredEmail {
	id: string
	to: string
}
