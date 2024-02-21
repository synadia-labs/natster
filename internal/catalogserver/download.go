package catalogserver

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nkeys"
	"github.com/synadia-labs/natster/internal/medialibrary"
	"github.com/synadia-labs/natster/internal/models"
)

const (
	chunkSizeBytes    = 5120
	headerChunkIndex  = "x-natster-chunk-idx"
	headerSenderXkey  = "x-natster-sender-xkey"
	headerTotalChunks = "x-natster-total-chunks"
)

func handleDownloadRequest(srv *CatalogServer) func(m *nats.Msg) {
	return func(m *nats.Msg) {
		// *.natster.catalog.%s.download
		tokens := strings.Split(m.Subject, ".")

		var req models.DownloadRequest
		err := json.Unmarshal(m.Data, &req)
		if err != nil {
			slog.Error("Failed to deserialize download request", err)
			_ = m.Respond(models.NewApiResultFail(err.Error(), 400))
			return
		}
		if !srv.isClientAllowed(tokens[0]) {
			slog.Debug("Request to download file from unauthorized client",
				slog.String("hash", req.Hash),
				slog.String("account", tokens[0]),
				slog.String("target_xkey", req.TargetXkey),
			)
			_ = m.Respond(models.NewApiResultFail("Forbidden", 403))
			return
		}

		slog.Info("Receiving request for file download",
			slog.String("hash", req.Hash))

		f := srv.library.FindByHash(req.Hash)
		if f == nil {
			_ = m.Respond(models.NewApiResultFail("Requested File Not Found", 404))
			return
		}

		senderKp, _ := nkeys.CreateCurveKeys()
		senderPublicKey, _ := senderKp.PublicKey()

		chunks := determineChunks(uint(f.ByteSize), chunkSizeBytes)
		resp := models.DownloadResponse{
			Acknowledged: true,
			ChunkSize:    chunkSizeBytes,
			TotalChunks:  chunks,
			SenderXKey:   senderPublicKey,
			TotalBytes:   uint(f.ByteSize),
		}
		_ = m.Respond(models.NewApiResultPass(resp))

		go srv.transmitChunkedFile(senderKp, tokens[0], req, f, chunks, resp)
	}
}

func (srv *CatalogServer) transmitChunkedFile(
	senderKp nkeys.KeyPair,
	targetAccount string,
	request models.DownloadRequest,
	entry *medialibrary.MediaEntry,
	chunks uint,
	resp models.DownloadResponse) {

	f, err := os.Open(entry.Path)
	if err != nil {
		slog.Error("Error reading file '%s': %s", entry.Path, err.Error())
	}
	r := bufio.NewReader(f)
	buf := make([]byte, 0, chunkSizeBytes)

	// Axxx.natster.media.kevvbuzz.xxxxxx
	targetSubject := fmt.Sprintf("%s.natster.media.%s.%s", targetAccount, srv.library.Name, request.Hash)

	for i := 0; i < int(chunks); i++ {
		n, err := r.Read(buf[:cap(buf)])
		buf = buf[:n]
		if n == 0 {
			if err == nil {
				continue
			}
			if err == io.EOF {
				break
			}
			slog.Error("File read error during chunk transmission", err)
		}
		if err != nil && err != io.EOF {
			slog.Error("File read error during chunk transmission", err)
			break
		}
		sealed, err := senderKp.Seal(buf, request.TargetXkey)
		if err != nil {
			slog.Error("Encryption failure", err)
			break
		}
		err = srv.transmitChunk(i, targetSubject, request, sealed, resp)
		if err != nil {
			slog.Error("Failed to transmit chunk", err)
			break
		}
	}
}

func (srv *CatalogServer) transmitChunk(index int, targetSubject string, request models.DownloadRequest, buf []byte, resp models.DownloadResponse) error {

	m := nats.NewMsg(targetSubject)
	m.Header.Add(headerChunkIndex, strconv.Itoa(index))
	m.Header.Add(headerSenderXkey, resp.SenderXKey)
	m.Header.Add(headerTotalChunks, strconv.Itoa(int(resp.TotalChunks)))
	m.Data = buf

	err := srv.nc.PublishMsg(m)
	if err != nil {
		return err
	}
	return nil
}

func determineChunks(fileSize uint, chunkSize uint) uint {
	chunks := fileSize / chunkSize
	if fileSize%chunkSize != 0 {
		chunks = chunks + 1
	}
	return chunks
}
