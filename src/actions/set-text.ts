import { action, KeyDownEvent, SingletonAction } from '@elgato/streamdeck'
import streamDeck from '@elgato/streamdeck'

import { createSetTextCommand } from '../lib/characterworks/commands'
import type { CharacterWorksChannel } from '../lib/characterworks/types'
import { sendCommand } from '../lib/characterworks/client'

type SetTextSettings = {
	motionName?: string
	textLayer?: string
	textValue?: string
	channel?: CharacterWorksChannel
	host?: string
	port?: number
}

@action({ UUID: 'dev.flowingspdg.characterworks-node.set-text' })
export class SetText extends SingletonAction<SetTextSettings> {
	override async onKeyDown(ev: KeyDownEvent<SetTextSettings>): Promise<void> {
		const { settings } = ev.payload

		const motionName = settings.motionName ?? ''
		const textLayer = settings.textLayer ?? ''
		const textValue = settings.textValue ?? ''
		const channel = settings.channel ?? 'preview'
		const host = settings.host
		const port = settings.port

		if (!host || !port) {
			streamDeck.logger.warn('CharacterWorks host/port is not configured')
			return
		}

		if (!motionName || !textLayer) {
			streamDeck.logger.warn('Motion name or text layer is empty')
			return
		}

		try {
			const command = createSetTextCommand(motionName, textLayer, textValue, channel)
			const response = await sendCommand(command, { host, port })

			streamDeck.logger.debug(
				`set-text: status=${response.status} statusText=${response.statusText}`
			)
		} catch (error) {
			streamDeck.logger.error(`set-text error: ${(error as Error).message}`)
		}
	}
}

