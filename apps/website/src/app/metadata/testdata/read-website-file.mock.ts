import { readFileSync } from 'node:fs'
import { join } from 'node:path'

export function readWebsiteFile(path: string) {
	return readFileSync(join(process.cwd(), path), 'utf8')
}
