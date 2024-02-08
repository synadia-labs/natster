package medialibrary

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	log "log/slog"
	"os"
	"path"
	"path/filepath"
)

type MediaLibrary struct {
	Name        string       `json:"name"`
	RootDir     string       `json:"root_dir"`
	Description string       `json:"description"`
	Entries     []MediaEntry `json:"entries"`
	Shares      []string     `json:"shares"`
}

type MediaEntry struct {
	Path        string `json:"path"`
	Description string `json:"description"`
	MimeType    string `json:"mime_type"`
	Hash        string `json:"hash"`
	ByteSize    int64  `json:"byte_size"`
}

func New(rootDir string, name string, description string) (*MediaLibrary, error) {
	return &MediaLibrary{
		Name:        name,
		RootDir:     rootDir,
		Description: description,
		Shares:      make([]string, 0),
		Entries:     make([]MediaEntry, 0),
	}, nil
}

func Load(name string) (*MediaLibrary, error) {
	natsterHome, err := getNatsterHome()
	if err != nil {
		return nil, err
	}
	bytes, err := os.ReadFile(path.Join(natsterHome, fmt.Sprintf("%s.json", name)))
	if err != nil {
		return nil, err
	}
	var library MediaLibrary
	err = json.Unmarshal(bytes, &library)
	if err != nil {
		return nil, err
	}
	return &library, nil
}

func (library *MediaLibrary) AddShare(accountKey string) {
	library.Shares = append(library.Shares, accountKey)
}

func (library *MediaLibrary) Save() error {
	natsterHome, err := getNatsterHome()
	if err != nil {
		return err
	}
	bytes, err := json.Marshal(library)
	if err != nil {
		return err
	}
	dataFile := path.Join(natsterHome, fmt.Sprintf("%s.json", library.Name))
	err = os.WriteFile(dataFile, bytes, 0660)
	if err != nil {
		return err
	}

	return nil
}

// Attempts to locate a media entry based on the hash. If it is not found,
// the result will be nil
func (library *MediaLibrary) FindByHash(hash string) *MediaEntry {
	for _, entry := range library.Entries {
		if entry.Hash == hash {
			return &entry
		}
	}
	return nil
}

// Recursively runs through the root directory and ensures that there's at least a
// key in the library metadata for that file
func (library *MediaLibrary) Ingest() error {

	library.Entries = make([]MediaEntry, 0)

	err := filepath.Walk(library.RootDir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				hash, err := hashFile(path)
				if err != nil {
					return err
				}
				entry := MediaEntry{
					Path:        path,
					Hash:        hash,
					ByteSize:    info.Size(),
					Description: "Auto-imported library entry",
				}
				library.Entries = append(library.Entries, entry)
			}
			return nil
		})

	if err == nil {
		log.Info(
			"Imported library entries",
		)
	}

	return err
}

func (library *MediaLibrary) GetCatalog() ([]MediaEntry, error) {

	return library.Entries, nil
}

func getNatsterHome() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	natsterHome := path.Join(home, ".natster")
	err = os.MkdirAll(natsterHome, 0750)
	if err != nil {
		return "", err
	}

	return natsterHome, nil

}

func hashFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		slog.Error("Failed to open file for hashing", err)
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		slog.Error("Failed to copy file for hashing", err)
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
