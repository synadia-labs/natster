package models

import (
	"encoding/json"

	"github.com/synadia-io/control-plane-sdk-go/syncp"
)

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
	Hash        string `json:"hash"`
}

type ApiResult struct {
	Error *string     `json:"error,omitempty"`
	Code  int         `json:"code"`
	Data  interface{} `json:"data"`
}

type TypedApiResult[T any] struct {
	Error *string `json:"error,omitempty"`
	Code  int     `json:"code"`
	Data  T       `json:"data"`
}

func NewApiResultPass(data interface{}) []byte {
	res := ApiResult{
		Data: data,
		Code: 200,
	}
	bytes, _ := json.Marshal(res)
	return bytes
}

func NewApiResultFail(msg string, code int) []byte {
	res := ApiResult{
		Error: syncp.Ptr(msg),
		Code:  code,
	}
	bytes, _ := json.Marshal(res)
	return bytes
}
