package catalogserver

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strconv"
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

// Queries the contents of a given calog. Note that the contents supplied are
// summary items, including only the path and hash
func (c *Client) GetCatalogItems(catalog string) ([]models.CatalogEntry, error) {
	reqSubject := fmt.Sprintf("natster.catalog.%s.get", catalog)
	res, err := c.nc.Request(reqSubject, []byte{}, 1*time.Second)
	if err != nil {
		return nil, err
	}
	var resp models.TypedApiResult[models.CatalogSummary]
	err = json.Unmarshal(res.Data, &resp)
	if err != nil {
		fmt.Printf("Deserialization failure getting catalog items: %s\n", err.Error())
		return nil, err
	}
	if resp.Error != nil {
		fmt.Printf("%s (%d)\n", *resp.Error, resp.Code)
		return nil, errors.New(*resp.Error)
	}

	return resp.Data.Entries, nil
}

// Submits a request to download a file containing the hash of the file in question
// and a public Xkey to be used to encrypt chunks
// Subscribes to natster.media.{catalog}.{hash} for encrypted chunks
func (c *Client) DownloadFile(catalog string, hash string, targetPath string) error {
	targetKp, _ := nkeys.CreateCurveKeys()
	targetPublic, _ := targetKp.PublicKey()
	reqSubject := fmt.Sprintf("natster.catalog.%s.download", catalog)
	subscribeSubject := fmt.Sprintf("natster.media.%s.%s", catalog, hash)

	ch := make(chan []byte)

	sub, err := c.nc.Subscribe(subscribeSubject, func(m *nats.Msg) {
		senderXKey := m.Header.Get("x-natster-sender-xkey")

		chunkIdx, err := strconv.Atoi(m.Header.Get("x-natster-chunk-idx"))
		if err != nil {
			slog.Error("Failed to parse x-natster-chunk-idx header", err,
				slog.String("sender_key", senderXKey),
			)
		}

		totalChunks, err := strconv.Atoi(m.Header.Get("x-natster-total-chunks"))
		if err != nil {
			slog.Error("Failed to parse x-natster-chunk-idx header", err,
				slog.String("sender_key", senderXKey),
			)
		}

		decrypted, err := targetKp.Open(m.Data, senderXKey)
		if err != nil {
			slog.Error("Failed to decrypt chunk", err,
				slog.Int("chunk_idx", chunkIdx),
				slog.String("sender_key", senderXKey),
			)
		}

		ch <- decrypted

		fmt.Printf("Received chunk %d (%d bytes)\n", chunkIdx, len(decrypted))
		if chunkIdx == totalChunks-1 {
			close(ch)
		}
	})
	if err != nil {
		return err
	}

	defer func() {
		_ = sub.Unsubscribe()
	}()

	dlRequest := models.DownloadRequest{
		Hash:       hash,
		TargetXkey: targetPublic,
	}
	reqBytes, _ := json.Marshal(&dlRequest)
	resp, err := c.nc.Request(reqSubject, reqBytes, 1*time.Second)
	if err != nil {
		return err
	}

	var newResponse models.TypedApiResult[models.DownloadResponse]
	err = json.Unmarshal(resp.Data, &newResponse)
	if err != nil {
		return err
	}

	fmt.Printf("File download request acknowledged: %d bytes (%d chunks of %d bytes each.) from %s\n",
		newResponse.Data.TotalBytes,
		newResponse.Data.TotalChunks,
		newResponse.Data.ChunkSize,
		newResponse.Data.SenderXKey,
	)

	f, err := os.Create(targetPath)
	if err != nil {
		return err
	}

	writer := bufio.NewWriter(f)
	for buf := range ch {
		writer.Write(buf)
	}
	_ = writer.Flush()
	_ = f.Close()

	return nil
}
