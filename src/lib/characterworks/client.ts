import type { CharacterWorksCommand, CharacterWorksConfig } from './types'

/**
 * Send a command to CharacterWorks over HTTP using Node.js 20 global fetch.
 * This mirrors the Companion module behaviour (HTTP POST with JSON body).
 */
export async function sendCommand(
	command: CharacterWorksCommand,
	config: CharacterWorksConfig
): Promise<Response> {
	const url = `http://${config.host}:${config.port}/`

	const controller = new AbortController()
	const timeoutId = setTimeout(() => controller.abort(), 10_000) // 10s timeout

	try {
		const response = await fetch(url, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
			},
			body: JSON.stringify(command),
			signal: controller.signal,
		})

		clearTimeout(timeoutId)
		return response
	} catch (error) {
		clearTimeout(timeoutId)
		throw error
	}
}

