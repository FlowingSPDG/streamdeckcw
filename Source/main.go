package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strconv"
	"sync"

	cw "github.com/FlowingSPDG/cwgo"
	"github.com/FlowingSPDG/streamdeck"
)

func payloadBytes(p any) []byte {
	switch v := p.(type) {
	case json.RawMessage:
		return v
	case []byte:
		return v
	default:
		b, _ := json.Marshal(v)
		return b
	}
}

type Settings struct {
	Host       string `json:"host"`
	Port       int    `json:"port"`
	Channel    string `json:"channel"`
	Layer      string `json:"layer"`
	Field      string `json:"field"`
	Value      string `json:"value"`
	Method     string `json:"method"`
	Path       string `json:"path"`
	Body       string `json:"body"`
	Motion     string `json:"motion"`
	MotionID   string `json:"motionId"`
	ChannelInt int    `json:"channelInt"`
	LayerInt   int    `json:"layerInt"`
}

type pluginState struct {
	mu sync.RWMutex
	cw *cw.CharacterWorks
}

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context) error {
	params, err := streamdeck.ParseRegistrationParams(os.Args)
	if err != nil {
		return err
	}
	client := streamdeck.NewClient(ctx, params)
	setup(client)
	return client.Run(ctx)
}

func setup(client *streamdeck.Client) {
	state := &pluginState{}

	update := client.Action("dev.flowingspdg.characterworks.updateField")
	streamdeck.OnWillAppear[Settings](update, func(ctx context.Context, client *streamdeck.Client, p streamdeck.WillAppearPayload[Settings]) error {
		return nil
	})
	streamdeck.OnKeyDown[Settings](update, func(ctx context.Context, client *streamdeck.Client, p streamdeck.KeyDownPayload[Settings]) error {
		if err := ensureCW(state, p.Settings); err != nil {
			return err
		}
		ch := normalizeChannel(p.Settings)
		layer := normalizeLayer(p.Settings)
		if layer == "" {
			return nil
		}
		state.mu.RLock()
		cwcli := state.cw
		state.mu.RUnlock()
		if cwcli == nil {
			return nil
		}
		_, _ = cwcli.SetText(layer, nil, p.Settings.Value, ch)
		return nil
	})
	streamdeck.OnDidReceivePropertyInspectorMessage[map[string]any](update, func(ctx context.Context, client *streamdeck.Client, msg streamdeck.DidReceivePropertyInspectorMessagePayload[map[string]any]) error {
		if s, ok := msg.Message["settings"]; ok {
			return client.SetSettings(ctx, s)
		}
		return nil
	})

	take := client.Action("dev.flowingspdg.characterworks.take")
	streamdeck.OnWillAppear[Settings](take, func(ctx context.Context, client *streamdeck.Client, p streamdeck.WillAppearPayload[Settings]) error {
		return nil
	})
	streamdeck.OnKeyDown[Settings](take, func(ctx context.Context, client *streamdeck.Client, p streamdeck.KeyDownPayload[Settings]) error {
		if err := ensureCW(state, p.Settings); err != nil {
			return err
		}
		ch := normalizeChannel(p.Settings)
		motions := []string{}
		ids := []string{}
		if p.Settings.Motion != "" {
			motions = append(motions, p.Settings.Motion)
		}
		if p.Settings.MotionID != "" {
			ids = append(ids, p.Settings.MotionID)
		}
		state.mu.RLock()
		cwcli := state.cw
		state.mu.RUnlock()
		if cwcli == nil {
			return nil
		}
		_, _ = cwcli.PlayMotions(motions, ids, ch)
		return nil
	})
	streamdeck.OnDidReceivePropertyInspectorMessage[map[string]any](take, func(ctx context.Context, client *streamdeck.Client, msg streamdeck.DidReceivePropertyInspectorMessagePayload[map[string]any]) error {
		if s, ok := msg.Message["settings"]; ok {
			return client.SetSettings(ctx, s)
		}
		return nil
	})

	stop := client.Action("dev.flowingspdg.characterworks.stop")
	streamdeck.OnWillAppear[Settings](stop, func(ctx context.Context, client *streamdeck.Client, p streamdeck.WillAppearPayload[Settings]) error {
		return nil
	})
	streamdeck.OnKeyDown[Settings](stop, func(ctx context.Context, client *streamdeck.Client, p streamdeck.KeyDownPayload[Settings]) error {
		if err := ensureCW(state, p.Settings); err != nil {
			return err
		}
		ch := normalizeChannel(p.Settings)
		motions := []string{}
		ids := []string{}
		if p.Settings.Motion != "" {
			motions = append(motions, p.Settings.Motion)
		}
		if p.Settings.MotionID != "" {
			ids = append(ids, p.Settings.MotionID)
		}
		state.mu.RLock()
		cwcli := state.cw
		state.mu.RUnlock()
		if cwcli == nil {
			return nil
		}
		_, _ = cwcli.StopMotions(motions, ids, ch)
		return nil
	})
	streamdeck.OnDidReceivePropertyInspectorMessage[map[string]any](stop, func(ctx context.Context, client *streamdeck.Client, msg streamdeck.DidReceivePropertyInspectorMessagePayload[map[string]any]) error {
		if s, ok := msg.Message["settings"]; ok {
			return client.SetSettings(ctx, s)
		}
		return nil
	})

	cmd := client.Action("dev.flowingspdg.characterworks.command")
	streamdeck.OnWillAppear[Settings](cmd, func(ctx context.Context, client *streamdeck.Client, p streamdeck.WillAppearPayload[Settings]) error {
		return nil
	})
	streamdeck.OnKeyDown[Settings](cmd, func(ctx context.Context, client *streamdeck.Client, p streamdeck.KeyDownPayload[Settings]) error {
		if err := ensureCW(state, p.Settings); err != nil {
			return err
		}
		return nil
	})
	streamdeck.OnDidReceivePropertyInspectorMessage[map[string]any](cmd, func(ctx context.Context, client *streamdeck.Client, msg streamdeck.DidReceivePropertyInspectorMessagePayload[map[string]any]) error {
		if s, ok := msg.Message["settings"]; ok {
			return client.SetSettings(ctx, s)
		}
		return nil
	})
}

func ensureCW(state *pluginState, s Settings) error {
	state.mu.RLock()
	exists := state.cw != nil
	state.mu.RUnlock()
	if exists {
		return nil
	}
	if s.Host == "" {
		return nil
	}
	cli := cw.NewCharacterWorks(s.Host, 10)
	state.mu.Lock()
	state.cw = cli
	state.mu.Unlock()
	return nil
}

func normalizeChannel(s Settings) *string {
	if s.Channel != "" {
		return &s.Channel
	}
	if s.ChannelInt != 0 {
		v := strconv.Itoa(s.ChannelInt)
		return &v
	}
	return nil
}

func normalizeLayer(s Settings) string {
	if s.Layer != "" {
		return s.Layer
	}
	if s.LayerInt != 0 {
		return strconv.Itoa(s.LayerInt)
	}
	if s.Field != "" {
		return s.Field
	}
	return ""
}
