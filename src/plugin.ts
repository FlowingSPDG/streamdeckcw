import streamDeck from '@elgato/streamdeck'

import { TriggerCharacterWorks } from './actions/trigger-cw'
import { SetText } from './actions/set-text'
import { ActivateGrid } from './actions/activate-grid'

// Enable trace logging so that all messages between the Stream Deck, and the plugin are recorded.
streamDeck.logger.setLevel('trace')

// Register actions.
streamDeck.actions.registerAction(new TriggerCharacterWorks())
streamDeck.actions.registerAction(new SetText())
streamDeck.actions.registerAction(new ActivateGrid())

// Finally, connect to the Stream Deck.
streamDeck.connect()
