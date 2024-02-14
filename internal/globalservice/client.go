package globalservice

import (
	"encoding/json"
	"log/slog"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/synadia-labs/natster/internal/models"
)

type Client struct {
	nc *nats.Conn
}

func NewClient(nc *nats.Conn) *Client {
	return &Client{
		nc: nc,
	}
}

func NewClientWithCredsPath(credsPath string) (*Client, error) {
	nc, err := nats.Connect("tls://connect.ngs.global",
		nats.UserCredentials(credsPath),
		nats.Name("natster_gsclient"),
	)
	if err != nil {
		slog.Error(
			"Failed to connect to NATS",
			slog.String("error", err.Error()),
		)
		return nil, err
	}
	return NewClient(nc), nil
}

func (c *Client) GetMyCatalogs() ([]models.CatalogShareSummary, error) {
	res, err := c.nc.Request("natster.global.my.shares", []byte{}, 1*time.Second)
	if err != nil {
		return nil, err
	}
	var apiResult models.TypedApiResult[[]models.CatalogShareSummary]
	err = json.Unmarshal(res.Data, &apiResult)
	if err != nil {
		return nil, err
	}

	return apiResult.Data, nil
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
	return c.nc.Flush()
}

func (c *Client) PublishHeartbeat(accountId string, catalog string) error {
	hb := models.Heartbeat{
		AccountId: accountId,
		Catalog:   catalog,
	}
	hbBytes, _ := json.Marshal(&hb)
	err := c.nc.Publish("natster.global.heartbeats.put", hbBytes)
	if err != nil {
		return err
	}
	return nil
}
