package globalservice

import (
	"encoding/json"
	"errors"
	"log/slog"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/synadia-labs/natster/internal/medialibrary"
	"github.com/synadia-labs/natster/internal/models"
)

const (
	defaultClientTimeout = 2 * time.Second
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

func (c *Client) ValidateCatalogName(catalog string) (*models.CatalogNameValidationResult, error) {
	res, err := c.nc.Request("natster.global.catalogs.validatename", []byte(catalog), defaultClientTimeout)
	if err != nil {
		return nil, err
	}
	var apiResult models.TypedApiResult[models.CatalogNameValidationResult]
	err = json.Unmarshal(res.Data, &apiResult)
	if err != nil {
		return nil, err
	}
	if apiResult.Error != nil {
		return &models.CatalogNameValidationResult{
			Valid:   false,
			Message: *apiResult.Error,
		}, nil
	}
	return apiResult.Data, nil
}

func (c *Client) Whoami() (*models.WhoamiResponse, error) {
	res, err := c.nc.Request("natster.global.whoami", []byte{}, defaultClientTimeout)
	if err != nil {
		return nil, err
	}
	var apiResult models.TypedApiResult[models.WhoamiResponse]
	err = json.Unmarshal(res.Data, &apiResult)
	if err != nil {
		return nil, err
	}
	if apiResult.Code == 404 {
		return nil, nil
	}

	if apiResult.Error != nil {
		return nil, errors.New(*apiResult.Error)
	}
	return apiResult.Data, nil
}

// This is to be called by the natster.io site when someone logs into the site, which
// will provide an OAuth identifier. If this OAuth ID has been bound to a context, we
// should be able to download that context and the corresponding credentials
func (c *Client) GetBoundContextByOAuth(oauthId string) (*models.ContextQueryResponse, error) {
	res, err := c.nc.Request("natster.global.context.get", []byte(oauthId), defaultClientTimeout)
	if err != nil {
		return nil, err
	}
	var apiResult models.TypedApiResult[models.ContextQueryResponse]
	err = json.Unmarshal(res.Data, &apiResult)
	if err != nil {
		return nil, err
	}
	if apiResult.Error != nil {
		return nil, errors.New(*apiResult.Error)
	}
	return apiResult.Data, nil
}

func (c *Client) GenerateOneTimeCode(context models.NatsterContext) (*models.OtcGenerateResponse, error) {
	ctxBytes, err := json.Marshal(context)
	if err != nil {
		return nil, err
	}

	res, err := c.nc.Request("natster.global.otc.generate", ctxBytes, defaultClientTimeout)
	if err != nil {
		return nil, err
	}
	var apiResult models.TypedApiResult[models.OtcGenerateResponse]
	err = json.Unmarshal(res.Data, &apiResult)
	if err != nil {
		return nil, err
	}
	if apiResult.Error != nil {
		return nil, errors.New(*apiResult.Error)
	}
	return apiResult.Data, nil
}

func (c *Client) ClaimOneTimeCode(code string, oauthIdentifier string) (*models.NatsterContext, error) {
	request := models.OtcClaimRequest{
		Code:          code,
		OAuthIdentity: oauthIdentifier,
	}
	bytes, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	res, err := c.nc.Request("natster.global.otc.claim", bytes, defaultClientTimeout)
	if err != nil {
		return nil, err
	}
	var apiResult models.TypedApiResult[models.NatsterContext]
	err = json.Unmarshal(res.Data, &apiResult)
	if err != nil {
		return nil, err
	}
	if apiResult.Code != 200 {
		return nil, errors.New(*apiResult.Error)
	}
	return apiResult.Data, nil
}

func (c *Client) GetMyCatalogs() (*[]models.CatalogShareSummary, error) {
	res, err := c.nc.Request("natster.global.my.shares", []byte{}, defaultClientTimeout)
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

func (c *Client) PublishEvent(eventType string, catalog string, target string, rawData json.RawMessage) error {
	bytes, err := json.Marshal(models.NatsterEvent{
		Catalog:   catalog,
		Target:    target,
		Data:      rawData,
		EventType: eventType,
	})
	if err != nil {
		return err
	}
	resp, err := c.nc.Request("natster.global.events.put", bytes, defaultClientTimeout)
	if err != nil {
		return err
	}
	var apiResponse models.ApiResult
	err = json.Unmarshal(resp.Data, &apiResponse)
	if err != nil {
		return err
	}
	if apiResponse.Error != nil {
		return errors.New(*apiResponse.Error)
	}
	return nil
}

func (c *Client) PublishHeartbeat(nctx *models.NatsterContext, library *medialibrary.MediaLibrary) error {
	hb := models.Heartbeat{
		AccountKey: nctx.AccountPublicKey,
		Catalog:    library.Name,
		Revision:   library.LastModified,
	}
	hbBytes, _ := json.Marshal(&hb)
	err := c.nc.Publish("natster.global.heartbeats.put", hbBytes)
	if err != nil {
		return err
	}
	return c.nc.Flush()
}

func (c *Client) Drain() {
	_ = c.nc.Drain()
}
