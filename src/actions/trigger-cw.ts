import { action, KeyDownEvent, SendToPluginEvent, SingletonAction } from '@elgato/streamdeck'
import streamDeck from '@elgato/streamdeck'
import { z } from 'zod'

import { createTriggerCommand, sendCommand, type CharacterWorksChannel, type TriggerAction } from '../lib/characterworks'

const TriggerActionEnum = z.enum(['play_motions', 'stop_motions', 'finish_motions', 'pause_motions'])
const ChannelEnum = z.enum(['live1', 'live2', 'preview'])

const TriggerCWSettingsSchema = z.object({
	action: TriggerActionEnum.default('play_motions'),
	motionName: z.string().trim().min(1, 'Motion name is required'),
	channel: ChannelEnum.default('preview'),
	host: z.string().trim().min(1, 'Host is required'),
	// `port` may come from the Property Inspector as string or number.
	port: z.coerce.number().int().min(1).max(65535),
})

type TriggerCWSettings = z.infer<typeof TriggerCWSettingsSchema>

@action({ UUID: 'dev.flowingspdg.characterworks-node.trigger-cw' })
export class TriggerCharacterWorks extends SingletonAction<TriggerCWSettings> {
	override async onKeyDown(ev: KeyDownEvent<TriggerCWSettings>): Promise<void> {
		let parsed: TriggerCWSettings

		try {
			parsed = TriggerCWSettingsSchema.parse(ev.payload.settings ?? {})
		} catch (error) {
			streamDeck.logger.warn(`trigger-cw settings validation failed: ${(error as Error).message}`)
			return
		}

		const { action, motionName, channel, host, port } = parsed

		try {
			const command = createTriggerCommand(action as TriggerAction, motionName, channel as CharacterWorksChannel)
			const response = await sendCommand(command, { host, port })

			streamDeck.logger.debug(
				`trigger-cw: status=${response.status} statusText=${response.statusText}`
			)
		} catch (error) {
			streamDeck.logger.error(`trigger-cw error: ${(error as Error).message}`)
		}
	}

	/**
	 * Handle requests from the Property Inspector (e.g. data source for motionName).
	 *
	 * The sdpi-components data source helper will send a `sendToPlugin` message
	 * when the Property Inspector is opened or refreshed. We respond with a
	 * standardized payload containing `event` and `items`, as documented in:
	 * https://sdpi-components.dev/docs/helpers/data-source
	 */
	override async onSendToPlugin(ev: SendToPluginEvent<{ event?: string; isRefresh?: boolean }, TriggerCWSettings>): Promise<void> {
		const dataEvent = ev.payload?.event

		// Data source for the Motion Name select.
		if (dataEvent === 'getCharacterWorksMotions') {
			// TODO: Replace this placeholder with a real CharacterWorks CWML-based
			// query when the concrete API for listing motions is defined.
			// See: https://chrworks.com/help/?url=doc/CWML.html
			const motions = ['motion_1', 'motion_2', 'motion_3']

			await streamDeck.ui.sendToPropertyInspector({
				event: dataEvent,
				items: motions.map((name) => ({
					label: name,
					value: name,
				})),
			})
		}
	}
}

