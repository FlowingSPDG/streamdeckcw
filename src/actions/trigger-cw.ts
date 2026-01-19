import { action, KeyDownEvent, SingletonAction } from '@elgato/streamdeck'
import streamDeck from '@elgato/streamdeck'

import { createTriggerCommand } from '../lib/characterworks/commands'
import type { CharacterWorksChannel, TriggerAction } from '../lib/characterworks/types'
import { sendCommand } from '../lib/characterworks/client'

type TriggerCWSettings = {
	action?: TriggerAction
	motionName?: string
	channel?: CharacterWorksChannel
	host?: string
	port?: number
}

@action({ UUID: 'dev.flowingspdg.characterworks-node.trigger-cw' })
export class TriggerCharacterWorks extends SingletonAction<TriggerCWSettings> {
	override async onKeyDown(ev: KeyDownEvent<TriggerCWSettings>): Promise<void> {
		const { settings } = ev.payload

		const action = settings.action ?? 'play_motions'
		const motionName = settings.motionName ?? ''
		const channel = settings.channel ?? 'preview'
		const host = settings.host
		const port = settings.port

		if (!host || !port) {
			streamDeck.logger.warn('CharacterWorks host/port is not configured')
			return
		}

		if (!motionName) {
			streamDeck.logger.warn('Motion name is empty')
			return
		}

		try {
			const command = createTriggerCommand(action, motionName, channel)
			const response = await sendCommand(command, { host, port })

			streamDeck.logger.debug(
				`trigger-cw: status=${response.status} statusText=${response.statusText}`
			)
		} catch (error) {
			streamDeck.logger.error(`trigger-cw error: ${(error as Error).message}`)
		}
	}
}

