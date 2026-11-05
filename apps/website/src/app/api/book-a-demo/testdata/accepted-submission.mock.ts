export function acceptedSubmission() {
	return Promise.resolve(
		Response.json(
			{
				id: '701706be-7e70-4681-9608-826d3ad548b9',
			},
			{
				status: 202,
			},
		),
	)
}
