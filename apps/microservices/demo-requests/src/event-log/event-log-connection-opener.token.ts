// Opens one connection to the event log; the producer and the reader each
// open the connections they need through it.
export const eventLogConnectionOpenerToken = Symbol(
	'event-log-connection-opener',
)
