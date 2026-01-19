import { action, KeyDownEvent, SingletonAction } from '@elgato/streamdeck'
import streamDeck from '@elgato/streamdeck'
import { z } from 'zod'

import { createActivateGridCommand, sendCommand } from '../lib/characterworks'

const ActivateGridSettingsSchema = z.object({
	gridName: z.string().trim().min(1, 'Grid name is required'),
	gridCell: z
		.string()
		.trim()
		.regex(/^\d+\s*,\s*\d+$/, 'Grid cell must be in the format row,column'),
	host: z.string().trim().min(1, 'Host is required'),
	port: z.coerce.number().int().min(1).max(65535),
})

type ActivateGridSettings = z.infer<typeof ActivateGridSettingsSchema>

@action({ UUID: 'dev.flowingspdg.characterworks-node.activate-grid' })
export class ActivateGrid extends SingletonAction<ActivateGridSettings> {
	override async onKeyDown(ev: KeyDownEvent<ActivateGridSettings>): Promise<void> {
		let parsed: ActivateGridSettings

		try {
			parsed = ActivateGridSettingsSchema.parse(ev.payload.settings ?? {})
		} catch (error) {
			streamDeck.logger.warn(
				`activate-grid settings validation failed: ${(error as Error).message}`
			)
			return
		}

		const { gridName, gridCell, host, port } = parsed

		const [rowStr, columnStr] = gridCell.split(',').map((v) => v.trim())
		const row = Number.parseInt(rowStr, 10)
		const column = Number.parseInt(columnStr, 10)

		try {
			const command = createActivateGridCommand(gridName, row, column)
			const response = await sendCommand(command, { host, port })

			streamDeck.logger.debug(
				`activate-grid: status=${response.status} statusText=${response.statusText}`
			)
		} catch (error) {
			streamDeck.logger.error(`activate-grid error: ${(error as Error).message}`)
		}
	}
}

