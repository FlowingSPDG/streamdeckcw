package main

import (
	"context"
	"log"
	"os"

	"dev.flowingspdg.characterworks/handlers"
	"github.com/FlowingSPDG/streamdeck"
)

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
	h := handlers.NewHandlers()
	h.RegisterAllActions(client)
}
