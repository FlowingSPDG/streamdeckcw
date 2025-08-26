package main

import (
	"context"
	"log"
	"os"
	"reflect"

	cw "github.com/FlowingSPDG/cwgo"
	"github.com/FlowingSPDG/streamdeck"
	"github.com/puzpuzpuz/xsync/v3"
)

type Settings struct {
	Host     string `json:"host"`
	Channel  string `json:"channel"`
	Layer    string `json:"layer"`
	Value    string `json:"value"`
	Motion   string `json:"motion"`
	MotionID string `json:"motionId"`
}

type pluginState struct {
	clients *xsync.MapOf[string, *cw.CharacterWorks]
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
	state := &pluginState{
		clients: xsync.NewMapOf[string, *cw.CharacterWorks](),
	}

	registerAction(client, "dev.flowingspdg.characterworks.updateField", state, handleUpdateField)
	registerAction(client, "dev.flowingspdg.characterworks.take", state, handleTake)
	registerAction(client, "dev.flowingspdg.characterworks.stop", state, handleStop)
}

func registerAction(client *streamdeck.Client, actionUUID string, state *pluginState, handler func(context.Context, *streamdeck.Client, streamdeck.KeyDownPayload[Settings], *pluginState) error) {
	action := client.Action(actionUUID)
	streamdeck.OnKeyDown(action, func(ctx context.Context, client *streamdeck.Client, p streamdeck.KeyDownPayload[Settings]) error {
		return handler(ctx, client, p, state)
	})
	streamdeck.OnDidReceivePropertyInspectorMessage(action, func(ctx context.Context, client *streamdeck.Client, msg streamdeck.DidReceivePropertyInspectorMessagePayload[map[string]any]) error {
		if t, ok := msg.Message["type"].(string); ok && t == "requestInfo" {
			payload := map[string]any{"type": "info"}
			var s Settings
			if raw, ok := msg.Message["settings"].(map[string]any); ok {
				if v, ok := raw["host"].(string); ok {
					s.Host = v
				}
				if v, ok := raw["channel"].(string); ok {
					s.Channel = v
				}
				if v, ok := raw["layer"].(string); ok {
					s.Layer = v
				}
			}
			if s.Host != "" {
				cwcli, _ := getCWClient(state, s.Host)
				ch := &s.Channel
				if s.Channel == "" {
					ch = nil
				}
				if actionUUID == "dev.flowingspdg.characterworks.updateField" {
					if layers := tryListLayers(cwcli, ch); len(layers) > 0 {
						payload["layers"] = layers
					}
				}
				if actionUUID == "dev.flowingspdg.characterworks.take" || actionUUID == "dev.flowingspdg.characterworks.stop" {
					motions, ids := tryListMotions(cwcli, ch)
					if len(motions) > 0 {
						payload["motions"] = motions
					}
					if len(ids) > 0 {
						payload["motionIds"] = ids
					}
				}
			}
			return client.SendToPropertyInspector(ctx, payload)
		}
		if t, ok := msg.Message["type"].(string); ok && t == "connect" {
			var s Settings
			if raw, ok := msg.Message["settings"].(map[string]any); ok {
				if v, ok := raw["host"].(string); ok {
					s.Host = v
				}
				if v, ok := raw["channel"].(string); ok {
					s.Channel = v
				}
			}
			if s.Host == "" {
				return client.SendToPropertyInspector(ctx, map[string]any{"type": "connectResult", "ok": false, "error": "host is required"})
			}
			cwcli, _ := getCWClient(state, s.Host)
			ch := &s.Channel
			if s.Channel == "" {
				ch = nil
			}
			resp := map[string]any{"type": "connectResult", "ok": true}
			if actionUUID == "dev.flowingspdg.characterworks.updateField" {
				if layers := tryListLayers(cwcli, ch); len(layers) > 0 {
					resp["layers"] = layers
				}
			}
			if actionUUID == "dev.flowingspdg.characterworks.take" || actionUUID == "dev.flowingspdg.characterworks.stop" {
				motions, ids := tryListMotions(cwcli, ch)
				if len(motions) > 0 {
					resp["motions"] = motions
				}
				if len(ids) > 0 {
					resp["motionIds"] = ids
				}
			}
			return client.SendToPropertyInspector(ctx, resp)
		}
		if s, ok := msg.Message["settings"]; ok {
			return client.SetSettings(ctx, s)
		}
		return nil
	})
}

func handleUpdateField(ctx context.Context, client *streamdeck.Client, p streamdeck.KeyDownPayload[Settings], state *pluginState) error {
	cwcli, err := getCWClient(state, p.Settings.Host)
	if err != nil {
		return err
	}
	if p.Settings.Layer == "" {
		return nil
	}
	ch := &p.Settings.Channel
	if p.Settings.Channel == "" {
		ch = nil
	}
	_, _ = cwcli.SetText(p.Settings.Layer, nil, p.Settings.Value, ch)
	return nil
}

func handleTake(ctx context.Context, client *streamdeck.Client, p streamdeck.KeyDownPayload[Settings], state *pluginState) error {
	cwcli, err := getCWClient(state, p.Settings.Host)
	if err != nil {
		return err
	}
	ch := &p.Settings.Channel
	if p.Settings.Channel == "" {
		ch = nil
	}
	motions := []string{}
	ids := []string{}
	if p.Settings.Motion != "" {
		motions = append(motions, p.Settings.Motion)
	}
	if p.Settings.MotionID != "" {
		ids = append(ids, p.Settings.MotionID)
	}
	_, _ = cwcli.PlayMotions(motions, ids, ch)
	return nil
}

func handleStop(ctx context.Context, client *streamdeck.Client, p streamdeck.KeyDownPayload[Settings], state *pluginState) error {
	cwcli, err := getCWClient(state, p.Settings.Host)
	if err != nil {
		return err
	}
	ch := &p.Settings.Channel
	if p.Settings.Channel == "" {
		ch = nil
	}
	motions := []string{}
	ids := []string{}
	if p.Settings.Motion != "" {
		motions = append(motions, p.Settings.Motion)
	}
	if p.Settings.MotionID != "" {
		ids = append(ids, p.Settings.MotionID)
	}
	_, _ = cwcli.StopMotions(motions, ids, ch)
	return nil
}

func getCWClient(state *pluginState, host string) (*cw.CharacterWorks, error) {
	if host == "" {
		return nil, nil
	}
	client, _ := state.clients.LoadOrStore(host, cw.NewCharacterWorks(host, 10))
	return client, nil
}

// tryListLayers calls cwgo's ListLayers(*string) if available via reflection
func tryListLayers(cwcli *cw.CharacterWorks, ch *string) []string {
	if cwcli == nil {
		return nil
	}
	v := reflect.ValueOf(cwcli)
	m := v.MethodByName("ListLayers")
	if !m.IsValid() {
		return nil
	}
	args := []reflect.Value{reflect.ValueOf(ch)}
	rets := m.Call(args)
	return pickStringSlice(rets)
}

// tryListMotions calls cwgo's ListMotions(*string) if available via reflection
// It supports either ([]string, []string, error) or ([][]string, error) like returns; we coerce best-effort
func tryListMotions(cwcli *cw.CharacterWorks, ch *string) ([]string, []string) {
	if cwcli == nil {
		return nil, nil
	}
	v := reflect.ValueOf(cwcli)
	m := v.MethodByName("ListMotions")
	if !m.IsValid() {
		return nil, nil
	}
	rets := m.Call([]reflect.Value{reflect.ValueOf(ch)})
	if len(rets) == 0 {
		return nil, nil
	}
	// Common pattern: motions, ids, err
	if len(rets) >= 2 {
		motions := toStringSlice(rets[0])
		ids := toStringSlice(rets[1])
		return motions, ids
	}
	// Fallback single slice
	return toStringSlice(rets[0]), nil
}

func pickStringSlice(rets []reflect.Value) []string {
	for _, rv := range rets {
		if s := toStringSlice(rv); len(s) > 0 {
			return s
		}
	}
	return nil
}

func toStringSlice(rv reflect.Value) []string {
	if !rv.IsValid() {
		return nil
	}
	if rv.Kind() == reflect.Interface || rv.Kind() == reflect.Pointer {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Slice {
		return nil
	}
	n := rv.Len()
	res := make([]string, 0, n)
	for i := 0; i < n; i++ {
		item := rv.Index(i)
		if item.Kind() == reflect.Interface || item.Kind() == reflect.Pointer {
			item = item.Elem()
		}
		if item.Kind() == reflect.String {
			res = append(res, item.String())
		}
	}
	return res
}
