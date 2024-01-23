package models

type CatalogSummary struct {
	Name    string         `json:"name"`
	Entries []CatalogEntry `json:"entries"`
}

// It might feel like a good idea to just use the internal medialibrary.MediaEntry here
// but once we get to refactoring, that internal type will change and we'll want to insulate
// clients from that
type CatalogEntry struct {
	Path        string `json:"path"`
	Description string `json:"description"`
	MimeType    string `json:"mime_type"`
	ByteSize    int64  `json:"byte_size"`
}
