package medialibrary

import (
	"encoding/json"
	"fmt"
	log "log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/nats-io/nats.go"
)

type MediaLibrary struct {
	name           string
	rootDir        string
	description    string
	nc             *nats.Conn
	metadataBucket nats.KeyValue
}

type MediaEntry struct {
	Path        string `json:"path"`
	Description string `json:"description"`
	MimeType    string `json:"mime_type"`
	ByteSize    int64  `json:"byte_size"`
}

func New(nc *nats.Conn, rootDir string, name string, description string) (*MediaLibrary, error) {
	bucket, err := createOrLocateLibraryMetadata(nc, name)
	if err != nil {
		return nil, err
	}

	return &MediaLibrary{
		name:           name,
		nc:             nc,
		rootDir:        rootDir,
		description:    description,
		metadataBucket: bucket,
	}, nil
}

// Recursively runs through the root directory and ensures that there's at least a
// key in the library metadata for that file
func (library *MediaLibrary) Ingest() error {
	newCount := 0

	err := filepath.Walk(library.rootDir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				err, exists :=
					library.putEntry(MediaEntry{
						Path:        path,
						ByteSize:    info.Size(),
						Description: "Auto-ingested library entry",
					})

				if err != nil {
					return err
				}
				if !exists {
					newCount += 1
				}

			}
			return nil
		})

	if err == nil {
		log.Info(
			"Ingested new library entries",
			"count", newCount,
		)
	}

	return err
}

func (library *MediaLibrary) putEntry(entry MediaEntry) (error, bool) {
	j, err := json.Marshal(entry)
	if err != nil {
		return err, false
	}
	key := strings.ReplaceAll(entry.Path, "/", "_")
	_, err = library.metadataBucket.Get(key)
	if err != nil {
		// Key doesn't exist, so create it
		_, err = library.metadataBucket.Put(key, j)
		if err != nil {
			log.Error(
				"Failed to write media library entry",
				"key", key,
				"error", err,
			)
			return err, false
		}
		return nil, false
	}

	return nil, true
}

func createOrLocateLibraryMetadata(nc *nats.Conn, name string) (nats.KeyValue, error) {
	opts := []nats.JSOpt{}
	js, err := nc.JetStream(opts...)
	if err != nil {
		return nil, err
	}

	bucketName := fmt.Sprintf("%s_MD", strings.ToUpper(name))
	bucket, err := js.KeyValue(bucketName)
	if err != nil {
		bucket, err = js.CreateKeyValue(&nats.KeyValueConfig{
			Bucket:      bucketName,
			Description: fmt.Sprintf("Library metadata for %s", name),
		})
		if err != nil {
			return nil, err
		}
	}

	log.Info(
		"Bound to library metadata KV store",
		"name", bucketName,
	)

	return bucket, nil
}
