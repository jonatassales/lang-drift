// The JSON-LD graph a page renders, as the structured-data test reads it.
export interface GraphNode {
	'@id': string
	'@type': string
	[property: string]: unknown
}

export interface JsonLdGraph {
	'@context': string
	'@graph': GraphNode[]
}
