import { action, KeyDownEvent, SingletonAction } from '@elgato/streamdeck'
import streamDeck from '@elgato/streamdeck'
import { z } from 'zod'

import { createSetTextCommand, sendCommand, type CharacterWorksChannel } from 'characterworks-client'

const ChannelEnum = z.enum(['live1', 'live2', 'preview'])

const SetTextSettingsSchema = z.object({
	motionName: z.string().trim().min(1, 'Motion name is required'),
	textLayer: z.string().trim().min(1, 'Text layer is required'),
	textValue: z.string().trim().default(''),
	channel: ChannelEnum.default('preview'),
	host: z.string().trim().min(1, 'Host is required'),
	port: z.coerce.number().int().min(1).max(65535),
})

type SetTextSettings = z.infer<typeof SetTextSettingsSchema>

@action({ UUID: 'dev.flowingspdg.characterworks-node.set-text' })
export class SetText extends SingletonAction<SetTextSettings> {
	override async onKeyDown(ev: KeyDownEvent<SetTextSettings>): Promise<void> {
		let parsed: SetTextSettings

		try {
			parsed = SetTextSettingsSchema.parse(ev.payload.settings ?? {})
		} catch (error) {
			streamDeck.logger.warn(`set-text settings validation failed: ${(error as Error).message}`)
			return
		}

		const { motionName, textLayer, textValue, channel, host, port } = parsed

		try {
			const command = createSetTextCommand(
				motionName,
				textLayer,
				textValue,
				channel as CharacterWorksChannel
			)
			const response = await sendCommand(command, { host, port })

			streamDeck.logger.debug(
				`set-text: status=${response.status} statusText=${response.statusText}`
			)
		} catch (error) {
			streamDeck.logger.error(`set-text error: ${(error as Error).message}`)
		}
	}
}

