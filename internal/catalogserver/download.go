package catalogserver

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nkeys"
	"github.com/synadia-labs/natster/internal/medialibrary"
	"github.com/synadia-labs/natster/internal/models"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

const (
	chunkSizeBytes    = 4096
	headerChunkIndex  = "x-natster-chunk-idx"
	headerSenderXkey  = "x-natster-sender-xkey"
	headerTotalChunks = "x-natster-total-chunks"
	headerTranscoding = "x-natster-transcoding"

	mimeTypeVideoMP4     = "video/mp4"
	natsterFragmentedKey = "io.natster.fragmented"
)

func handleDownloadRequest(srv *CatalogServer, local bool) func(m *nats.Msg) {
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
		if !local {
			if !srv.isClientAllowed(tokens[0]) {
				slog.Debug("Request to download file from unauthorized client",
					slog.String("hash", req.Hash),
					slog.String("account", tokens[0]),
					slog.String("target_xkey", req.TargetXkey),
				)
				_ = m.Respond(models.NewApiResultFail("Forbidden", 403))
				return
			}
		}

		slog.Info("Received request for file download",
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

		go srv.transmitChunkedFile(senderKp, tokens[0], req, f, chunks, resp, req.Transcode, local)
	}
}

func (srv *CatalogServer) transmitChunkedFile(
	senderKp nkeys.KeyPair,
	targetAccount string,
	request models.DownloadRequest,
	entry *medialibrary.MediaEntry,
	chunks uint,
	resp models.DownloadResponse,
	transcode bool,
	local bool) {

	path := filepath.Join(srv.library.RootDir, entry.Path)

	transcoding := false // this means we are streaming, not that transcoding is actively running
	transcodingInProgress := false
	var fileInfo os.FileInfo

	if strings.EqualFold(strings.ToLower(entry.MimeType), mimeTypeVideoMP4) {
		cmd := exec.Command("ffprobe", "-show_format", "-of", "json", path)
		out, err := cmd.CombinedOutput()
		if err != nil {
			slog.Error("Error reading mp4 metadata", slog.String("path", entry.Path), err.Error())
		} else {
			transcode = !strings.Contains(string(out), natsterFragmentedKey)
		}

		if transcode {
			id, _ := uuid.NewUUID()
			tmppath := filepath.Join(os.TempDir(), fmt.Sprintf("%s.mp4", id))
			chunks = chunks + uint(math.RoundToEven(float64(chunks)*.05)) // HACK expand total possible chunks by 5% so we iterate enough to read the entire fragmented file
			transcoding = true

			slog.Info("Transcoding mp4", slog.Uint64("chunks", uint64(chunks)), slog.String("path", path))

			go func() {
				transcodingInProgress = true
				err := ffmpeg.Input(path).
					Output(tmppath, ffmpeg.KwArgs{"movflags": "frag_keyframe+empty_moov+default_base_moof"}).
					Run()
				if err != nil {
					slog.Error("Error transcoding mp4", slog.String("path", entry.Path), err.Error())
				}
				transcodingInProgress = false

				slog.Info("Completed transcoding", "path", path)
			}()

			defer func() {
				_ = os.Remove(tmppath)
			}()

			time.Sleep(time.Millisecond * 5)
			path = tmppath

			var err error
			fileInfo, err = os.Stat(path)
			for err != nil {
				fileInfo, err = os.Stat(path)
			}
		}
	}

	f, err := os.Open(path)
	if err != nil {
		slog.Error("Error reading file", slog.String("path": path), err.Error())
	}
	r := bufio.NewReader(f)
	buf := make([]byte, 0, chunkSizeBytes)

	// Axxx.natster.media.kevvbuzz.xxxxxx or natster.media.kevbuzz.xxxx
	targetSubject := ""
	if local {
		targetSubject = fmt.Sprintf("natster.media.%s.%s", srv.library.Name, request.Hash)
	} else {
		targetSubject = fmt.Sprintf("%s.natster.media.%s.%s", targetAccount, srv.library.Name, request.Hash)
	}

	x := 0
	z := 0

	for i := 0; i < int(chunks) || transcoding; i++ {
		if transcoding && transcodingInProgress {
			fi, _ := os.Stat(path)
			if fi.Size() == fileInfo.Size() && z == 0 {
				slog.Info("size is the same as the last attempt...", slog.Int64("filesize", fi.Size()))

				if transcodingInProgress {
					time.Sleep(time.Millisecond * 50)
				}
				fileInfo = fi
				z++
				continue
			} else {
				z = 0
			}

			fileInfo = fi

			if fileInfo.Size() < chunkSizeBytes {
				slog.Info("cannot read chunk yet", slog.Int64("filesize", fi.Size()))
				continue
			}
			slog.Info("transcoding still in progress...", slog.Int64("filesize", fi.Size()))
		}

		n, err := r.Read(buf[:cap(buf)])
		buf = buf[:n]

		x += n
		slog.Info("read bytes of chunk...", slog.Int("bytes", n), slog.Int("total_bytes", x))

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

		err = srv.transmitChunk(i, transcoding, targetSubject, sealed, resp)
		if err != nil {
			slog.Error("Failed to transmit chunk", err)
			break
		}
	}
}

func (srv *CatalogServer) transmitChunk(index int, transcoding bool, targetSubject string, buf []byte, resp models.DownloadResponse) error {

	m := nats.NewMsg(targetSubject)
	m.Header.Add(headerChunkIndex, strconv.Itoa(index))
	m.Header.Add(headerSenderXkey, resp.SenderXKey)
	m.Header.Add(headerTotalChunks, strconv.Itoa(int(resp.TotalChunks)))

	if transcoding {
		m.Header.Add(headerTranscoding, strconv.FormatBool(transcoding))
	}

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
