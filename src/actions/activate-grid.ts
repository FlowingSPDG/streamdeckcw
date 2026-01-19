import { action, KeyDownEvent, SingletonAction } from '@elgato/streamdeck'
import streamDeck from '@elgato/streamdeck'

import { createActivateGridCommand } from '../lib/characterworks/commands'
import { sendCommand } from '../lib/characterworks/client'

type ActivateGridSettings = {
	gridName?: string
	gridCell?: string // "row,column"
	host?: string
	port?: number
}

@action({ UUID: 'dev.flowingspdg.characterworks-node.activate-grid' })
export class ActivateGrid extends SingletonAction<ActivateGridSettings> {
	override async onKeyDown(ev: KeyDownEvent<ActivateGridSettings>): Promise<void> {
		const { settings } = ev.payload

		const gridName = settings.gridName ?? ''
		const gridCell = settings.gridCell ?? ''
		const host = settings.host
		const port = settings.port

		if (!host || !port) {
			streamDeck.logger.warn('CharacterWorks host/port is not configured')
			return
		}

		if (!gridName || !gridCell) {
			streamDeck.logger.warn('Grid name or cell is empty')
			return
		}

		const parts = gridCell.split(',').map((v) => v.trim())
		if (parts.length !== 2) {
			streamDeck.logger.warn('Grid cell must be in the format row,column')
			return
		}

		const row = Number.parseInt(parts[0], 10)
		const column = Number.parseInt(parts[1], 10)

		if (Number.isNaN(row) || Number.isNaN(column)) {
			streamDeck.logger.warn('Grid cell coordinates must be numbers')
			return
		}

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

