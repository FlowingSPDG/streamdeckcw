export interface CharacterWorksConfig {
	host: string
	port: number
}

export type CharacterWorksChannel = 'live1' | 'live2' | 'preview'

export type TriggerAction = 'play_motions' | 'stop_motions' | 'finish_motions' | 'pause_motions'

export interface TriggerCommand {
	action: TriggerAction
	motions: string[]
	channel: CharacterWorksChannel
}

export interface SetTextCommand {
	action: 'set_text'
	layer: string
	value: string
	channel: CharacterWorksChannel
}

export interface ActivateGridCommand {
	action: 'activate_grid_cell'
	grid: string
	cell: [number, number]
}

export type CharacterWorksCommand = TriggerCommand | SetTextCommand | ActivateGridCommand

