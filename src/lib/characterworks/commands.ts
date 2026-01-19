import type {
	ActivateGridCommand,
	CharacterWorksChannel,
	SetTextCommand,
	TriggerAction,
	TriggerCommand,
} from './types'

export function createTriggerCommand(
	action: TriggerAction,
	motionNames: string,
	channel: CharacterWorksChannel
): TriggerCommand {
	const motions = motionNames
		.split(',')
		.map((m) => m.trim())
		.filter((m) => m.length > 0)

	return {
		action,
		motions,
		channel,
	}
}

export function createSetTextCommand(
	motionName: string,
	textLayer: string,
	value: string,
	channel: CharacterWorksChannel
): SetTextCommand {
	const motion = motionName.trim()
	const layer = textLayer.trim()

	return {
		action: 'set_text',
		layer: `${motion}\\${layer}`,
		value,
		channel,
	}
}

export function createActivateGridCommand(
	gridName: string,
	row: number,
	column: number
): ActivateGridCommand {
	return {
		action: 'activate_grid_cell',
		grid: gridName,
		cell: [row, column],
	}
}

