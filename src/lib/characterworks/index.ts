/**
 * CharacterWorks HTTP Client Library
 *
 * This module provides a pure client library for communicating with CharacterWorks
 * over HTTP. It is designed to be platform-agnostic and can be used by Stream Deck
 * plugins, Companion modules, or any other Node.js application.
 *
 * @module characterworks
 */

// Re-export all types
export type {
	CharacterWorksConfig,
	CharacterWorksChannel,
	TriggerAction,
	TriggerCommand,
	SetTextCommand,
	ActivateGridCommand,
	CharacterWorksCommand,
} from './types'

// Re-export all command creation functions
export {
	createTriggerCommand,
	createSetTextCommand,
	createActivateGridCommand,
} from './commands'

// Re-export the HTTP client function
export { sendCommand } from './client'
