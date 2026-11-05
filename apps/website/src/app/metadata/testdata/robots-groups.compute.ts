export function robotsGroups(text: string) {
	return text
		.trim()
		.split(/\n\s*\n/)
		.map((block) => block.split('\n'))
		.filter((lines) => lines[0].startsWith('User-agent:'))
		.map((lines) => ({
			agents: lines
				.filter((line) => line.startsWith('User-agent:'))
				.map((line) => line.replace('User-agent:', '').trim()),
			rules: lines.filter((line) => !line.startsWith('User-agent:')),
		}))
}
