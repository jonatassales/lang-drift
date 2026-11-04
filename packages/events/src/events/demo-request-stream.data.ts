// Where demo-request events live in Redis and who consumes them.
export const demoRequestAggregateType = 'demo-request'
export const demoRequestStream = 'demo-requests:events'
export const demoRequestDeadLetterStream = 'demo-requests:events:dead-letter'
export const demoRequestArchiveGroup = 'archive'
export const demoRequestMailerGroup = 'mailer'
