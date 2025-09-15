package handlers

import (
	"context"

	cw "github.com/FlowingSPDG/cwgo"
	"github.com/FlowingSPDG/streamdeck"
	"github.com/puzpuzpuz/xsync/v3"
)

// Settings represents the plugin settings structure
type Settings struct {
	Host     string `json:"host"`
	Channel  string `json:"channel"`
	Layer    string `json:"layer"`
	Value    string `json:"value"`
	Motion   string `json:"motion"`
	MotionID string `json:"motionId"`
}

// PluginState manages the plugin's global state
type PluginState struct {
	clients *xsync.MapOf[string, *cw.CharacterWorks]
}

// NewPluginState creates a new plugin state instance
func NewPluginState() *PluginState {
	return &PluginState{
		clients: xsync.NewMapOf[string, *cw.CharacterWorks](),
	}
}

// Handlers manages all action handlers
type Handlers struct {
	state *PluginState
}

// NewHandlers creates a new handlers instance
func NewHandlers() *Handlers {
	return &Handlers{
		state: NewPluginState(),
	}
}

// PropertyInspectorMessage represents messages from Property Inspector
type PropertyInspectorMessage struct {
	Type     string   `json:"type"`
	Settings Settings `json:"settings,omitempty"`
}

// PropertyInspectorResponse represents responses to Property Inspector
type PropertyInspectorResponse struct {
	Type      string   `json:"type"`
	OK        bool     `json:"ok,omitempty"`
	Error     string   `json:"error,omitempty"`
	Layers    []string `json:"layers,omitempty"`
	Motions   []string `json:"motions,omitempty"`
	MotionIds []string `json:"motionIds,omitempty"`
}

// getCWClient retrieves or creates a CharacterWorks client for the given host
func (h *Handlers) getCWClient(host string) (*cw.CharacterWorks, error) {
	if host == "" {
		return nil, nil
	}
	client, _ := h.state.clients.LoadOrStore(host, cw.NewCharacterWorks(host, 10))
	return client, nil
}

// listLayers retrieves available layers from CharacterWorks
func (h *Handlers) listLayers(ctx context.Context, cwcli *cw.CharacterWorks, channel *string) []string {
	if cwcli == nil {
		return nil
	}

	resp, err := cwcli.ListLayers(ctx, nil, channel)
	if err != nil {
		return nil
	}

	if resp == nil || resp.Layers == nil {
		return nil
	}

	// Convert Layer structs to strings
	var layers []string
	for _, layer := range resp.Layers {
		layers = append(layers, layer.Name)
	}

	return layers
}

// listMotions retrieves available motions from CharacterWorks
func (h *Handlers) listMotions(ctx context.Context, cwcli *cw.CharacterWorks, channel *string) ([]string, []string) {
	if cwcli == nil {
		return nil, nil
	}

	// Get motions with IDs
	resp, err := cwcli.ListMotionsWithIDs(ctx)
	if err != nil {
		return nil, nil
	}

	if resp == nil {
		return nil, nil
	}

	var motions []string
	var ids []string

	if resp.Motions != nil {
		for _, motion := range resp.Motions {
			motions = append(motions, motion.Name)
			ids = append(ids, motion.ID)
		}
	}

	return motions, ids
}

// handlePropertyInspectorMessage handles messages from Property Inspector
func (h *Handlers) handlePropertyInspectorMessage(ctx context.Context, client *streamdeck.Client, msg PropertyInspectorMessage, actionUUID string) error {
	switch msg.Type {
	case "requestInfo":
		return h.handleRequestInfo(ctx, client, msg, actionUUID)
	case "connect":
		return h.handleConnect(ctx, client, msg, actionUUID)
	default:
		// Handle settings update
		return client.SetSettings(ctx, msg.Settings)
	}
}

// handleRequestInfo handles requestInfo messages from Property Inspector
func (h *Handlers) handleRequestInfo(ctx context.Context, client *streamdeck.Client, msg PropertyInspectorMessage, actionUUID string) error {
	response := PropertyInspectorResponse{Type: "info"}

	if msg.Settings.Host != "" {
		cwcli, _ := h.getCWClient(msg.Settings.Host)
		channel := &msg.Settings.Channel
		if msg.Settings.Channel == "" {
			channel = nil
		}

		// Add action-specific data
		switch actionUUID {
		case "dev.flowingspdg.characterworks.updateField":
			if layers := h.listLayers(ctx, cwcli, channel); len(layers) > 0 {
				response.Layers = layers
			}
		case "dev.flowingspdg.characterworks.take", "dev.flowingspdg.characterworks.stop":
			if motions, ids := h.listMotions(ctx, cwcli, channel); len(motions) > 0 || len(ids) > 0 {
				response.Motions = motions
				response.MotionIds = ids
			}
		}
	}

	return client.SendToPropertyInspector(ctx, response)
}

// handleConnect handles connect messages from Property Inspector
func (h *Handlers) handleConnect(ctx context.Context, client *streamdeck.Client, msg PropertyInspectorMessage, actionUUID string) error {
	if msg.Settings.Host == "" {
		response := PropertyInspectorResponse{
			Type:  "connectResult",
			OK:    false,
			Error: "host is required",
		}
		return client.SendToPropertyInspector(ctx, response)
	}

	cwcli, _ := h.getCWClient(msg.Settings.Host)
	channel := &msg.Settings.Channel
	if msg.Settings.Channel == "" {
		channel = nil
	}

	response := PropertyInspectorResponse{Type: "connectResult", OK: true}

	// Add action-specific data
	switch actionUUID {
	case "dev.flowingspdg.characterworks.updateField":
		if layers := h.listLayers(ctx, cwcli, channel); len(layers) > 0 {
			response.Layers = layers
		}
	case "dev.flowingspdg.characterworks.take", "dev.flowingspdg.characterworks.stop":
		if motions, ids := h.listMotions(ctx, cwcli, channel); len(motions) > 0 || len(ids) > 0 {
			response.Motions = motions
			response.MotionIds = ids
		}
	}

	return client.SendToPropertyInspector(ctx, response)
}

// registerUpdateFieldAction registers the updateField action handler
func (h *Handlers) registerUpdateFieldAction(client *streamdeck.Client) {
	action := client.Action("dev.flowingspdg.characterworks.updateField")

	// Register KeyDown handler
	streamdeck.OnKeyDown(action, func(ctx context.Context, client *streamdeck.Client, p streamdeck.KeyDownPayload[Settings]) error {
		return h.handleUpdateField(ctx, client, p)
	})

	// Register Property Inspector message handler
	streamdeck.OnDidReceivePropertyInspectorMessage(action, func(ctx context.Context, client *streamdeck.Client, msg streamdeck.DidReceivePropertyInspectorMessagePayload[PropertyInspectorMessage]) error {
		return h.handlePropertyInspectorMessage(ctx, client, msg.Message, "dev.flowingspdg.characterworks.updateField")
	})
}

// registerTakeAction registers the take action handler
func (h *Handlers) registerTakeAction(client *streamdeck.Client) {
	action := client.Action("dev.flowingspdg.characterworks.take")

	// Register KeyDown handler
	streamdeck.OnKeyDown(action, func(ctx context.Context, client *streamdeck.Client, p streamdeck.KeyDownPayload[Settings]) error {
		return h.handleTake(ctx, client, p)
	})

	// Register Property Inspector message handler
	streamdeck.OnDidReceivePropertyInspectorMessage(action, func(ctx context.Context, client *streamdeck.Client, msg streamdeck.DidReceivePropertyInspectorMessagePayload[PropertyInspectorMessage]) error {
		return h.handlePropertyInspectorMessage(ctx, client, msg.Message, "dev.flowingspdg.characterworks.take")
	})
}

// registerStopAction registers the stop action handler
func (h *Handlers) registerStopAction(client *streamdeck.Client) {
	action := client.Action("dev.flowingspdg.characterworks.stop")

	// Register KeyDown handler
	streamdeck.OnKeyDown(action, func(ctx context.Context, client *streamdeck.Client, p streamdeck.KeyDownPayload[Settings]) error {
		return h.handleStop(ctx, client, p)
	})

	// Register Property Inspector message handler
	streamdeck.OnDidReceivePropertyInspectorMessage(action, func(ctx context.Context, client *streamdeck.Client, msg streamdeck.DidReceivePropertyInspectorMessagePayload[PropertyInspectorMessage]) error {
		return h.handlePropertyInspectorMessage(ctx, client, msg.Message, "dev.flowingspdg.characterworks.stop")
	})
}

// handleUpdateField handles the updateField action
func (h *Handlers) handleUpdateField(ctx context.Context, client *streamdeck.Client, p streamdeck.KeyDownPayload[Settings]) error {
	cwcli, err := h.getCWClient(p.Settings.Host)
	if err != nil {
		return err
	}
	if cwcli == nil || p.Settings.Layer == "" {
		return nil
	}

	channel := &p.Settings.Channel
	if p.Settings.Channel == "" {
		channel = nil
	}

	_, _ = cwcli.SetText(ctx, p.Settings.Layer, nil, p.Settings.Value, channel)
	return nil
}

// handleTake handles the take action
func (h *Handlers) handleTake(ctx context.Context, client *streamdeck.Client, p streamdeck.KeyDownPayload[Settings]) error {
	cwcli, err := h.getCWClient(p.Settings.Host)
	if err != nil {
		return err
	}
	if cwcli == nil {
		return nil
	}

	channel := &p.Settings.Channel
	if p.Settings.Channel == "" {
		channel = nil
	}

	var motions []string
	var ids []string
	if p.Settings.Motion != "" {
		motions = append(motions, p.Settings.Motion)
	}
	if p.Settings.MotionID != "" {
		ids = append(ids, p.Settings.MotionID)
	}

	_, _ = cwcli.PlayMotions(ctx, motions, ids, channel)
	return nil
}

// handleStop handles the stop action
func (h *Handlers) handleStop(ctx context.Context, client *streamdeck.Client, p streamdeck.KeyDownPayload[Settings]) error {
	cwcli, err := h.getCWClient(p.Settings.Host)
	if err != nil {
		return err
	}
	if cwcli == nil {
		return nil
	}

	channel := &p.Settings.Channel
	if p.Settings.Channel == "" {
		channel = nil
	}

	var motions []string
	var ids []string
	if p.Settings.Motion != "" {
		motions = append(motions, p.Settings.Motion)
	}
	if p.Settings.MotionID != "" {
		ids = append(ids, p.Settings.MotionID)
	}

	_, _ = cwcli.StopMotions(ctx, motions, ids, channel)
	return nil
}

// RegisterAllActions registers all action handlers with the StreamDeck client
func (h *Handlers) RegisterAllActions(client *streamdeck.Client) {
	h.registerUpdateFieldAction(client)
	h.registerTakeAction(client)
	h.registerStopAction(client)
}
