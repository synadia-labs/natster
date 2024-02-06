package models

const (
	NatsterInitializedEventType = "natster_initialized"
)

// Heartbeats are emitted to the global service periodically by running natster
// catalog servers
type Heartbeat struct {
	AccountId string `json:"account_id"`
	Catalog   string `json:"catalog"`
}

// Events are emitted by the natster server process
type NatsterEvent struct {
	Catalog   string      `json:"catalog"`
	Target    string      `json:"target"`
	EventType string      `json:"event_type"`
	Data      interface{} `json:"data,omitempty"`
}

type NatsterInitializedEvent struct {
	AccountId   string `json:"account_id"`
	AccountName string `json:"account_name"`
	AccountKey  string `json:"account_key"`
}

type CommunityStats struct {
	TotalInitialized uint64 `json:"total_initialized"`
	RunningCatalogs  uint64 `json:"running_catalogs"`
}
