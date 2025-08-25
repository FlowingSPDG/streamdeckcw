package cw

import (
	"context"
	"fmt"
)

type Client interface {
	Take(ctx context.Context, channel int, layer int) error
	Stop(ctx context.Context, channel int, layer int) error
	UpdateField(ctx context.Context, channel int, layer int, field string, value string) error
	Command(ctx context.Context, method string, path string, body string) (int, string, error)
}

type HTTPClient struct {
	host string
	port int
}

func NewHTTPClient(host string, port int) *HTTPClient { return &HTTPClient{host: host, port: port} }

func (c *HTTPClient) Take(ctx context.Context, channel int, layer int) error {
	// TODO: implement using cwgo once endpoint confirmed
	return nil
}

func (c *HTTPClient) Stop(ctx context.Context, channel int, layer int) error {
	// TODO: implement using cwgo once endpoint confirmed
	return nil
}

func (c *HTTPClient) UpdateField(ctx context.Context, channel int, layer int, field string, value string) error {
	// TODO: implement using cwgo once endpoint confirmed
	return nil
}

func (c *HTTPClient) Command(ctx context.Context, method string, path string, body string) (int, string, error) {
	return 0, "", fmt.Errorf("not implemented")
}
