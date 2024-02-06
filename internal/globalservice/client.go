package globalservice

import (
	"encoding/json"

	"github.com/ConnectEverything/natster/internal/models"
	"github.com/nats-io/nats.go"
)

type Client struct {
	nc *nats.Conn
}

func NewClient(nc *nats.Conn) *Client {
	return &Client{
		nc: nc,
	}
}

func (c *Client) PublishEvent(eventType string, catalog string, target string, rawData interface{}) error {
	bytes, err := json.Marshal(models.NatsterEvent{
		Catalog:   catalog,
		Target:    target,
		Data:      rawData,
		EventType: eventType,
	})
	if err != nil {
		return err
	}
	err = c.nc.Publish("natster.global.events.put", bytes)
	if err != nil {
		return err
	}
	return nil
}
