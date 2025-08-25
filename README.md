# Character Works Stream Deck Plugin (Go)

This repository contains an Elgato Stream Deck plugin written in Go to control Character Works.

- Stream Deck Go bindings: https://github.com/FlowingSPDG/streamdeck
- Character Works Go client: https://github.com/FlowingSPDG/cwgo
- Reference implementation: https://github.com/FlowingSPDG/streamdeck-advanced-http
- CW help: https://www.chrworks.com/help/
- Stream Deck SDK WebSocket: https://docs.elgato.com/streamdeck/sdk/references/websocket/plugin

## Actions

- Update Data Field
- Take (OnAir)
- Stop (Off)
- Generic Command (HTTP passthrough)

All actions accept per-action settings including `host` and `port`.

## Project Structure

- `Source/` Go implementation
  - `actions/` individual action handlers
  - `cw/` CW client adapter
  - `manifest.json` Stream Deck plugin manifest
  - `images/` plugin icons (copy from reference repo for now)
- `Makefile` build/package helpers (mirrors reference repo)
- `Release/` packaged plugin output

## Build

Prerequisites:
- Go 1.21+
- Stream Deck Distribution Tool available as `DistributionTool` (or `DistributionTool.exe`) in repo root, or update the Makefile path

Commands:

- Build the plugin app
```
make build
```

- Package to `.streamDeckPlugin`
```
make package
```

The packaged file is written to `Release/`.

## Notes

- Icons are not included. Copy from the reference repo (advanced-http) or provide your own under `Source/images/` and ensure names match those referenced by `manifest.json`.
- The CW client currently uses an HTTP adapter with pluggable interface. Swap to the `cwgo` client when connection details are finalized.
- Default plugin UUID: `dev.flowingspdg.characterworks` 