package main

import (
	"context"
	"log"
	"os"

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
