package catalogserver

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nkeys"
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
		nats.Name("natster_catalogclient"),
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

// Submits a request to download a file containing the hash of the file in question
// and a public Xkey to be used to encrypt chunks
// Subscribes to natster.media.{catalog}.{hash} for encrypted chunks
func (c *Client) DownloadFile(catalog string, hash string, targetPath string) error {
	targetKp, _ := nkeys.CreateCurveKeys()
	targetPublic, _ := targetKp.PublicKey()
	reqSubject := fmt.Sprintf("natster.catalog.%s.download", catalog)
	subscribeSubject := fmt.Sprintf("natster.media.%s.%s", catalog, hash)

	var response models.TypedApiResult[models.DownloadResponse]
	ch := make(chan bool)

	f, err := os.Create(targetPath)
	if err != nil {
		return err
	}
	defer f.Close()

	writer := bufio.NewWriter(f)
	chunkCount := 0
	c.nc.Subscribe(subscribeSubject, func(m *nats.Msg) {
		decrypted, err := targetKp.Open(m.Data, response.Data.SenderXKey)
		if err != nil {
			slog.Error("Failed to decrypt chunk", err)
		}
		_, _ = writer.Write(decrypted)

		chunkCount := chunkCount + 1
		if chunkCount == int(response.Data.TotalChunks) {
			ch <- true
		}
	})

	dlRequest := models.DownloadRequest{
		Hash:       hash,
		TargetXkey: targetPublic,
	}
	reqBytes, _ := json.Marshal(&dlRequest)
	resp, err := c.nc.Request(reqSubject, reqBytes, 1*time.Second)
	if err != nil {
		return err
	}

	err = json.Unmarshal(resp.Data, &response)
	if err != nil {
		return err
	}
	<-ch

	return nil
}
