package models

// Heartbeats are emitted to the global service periodically by running natster
// catalog servers
type Heartbeat struct {
	AccountId string `json:"account_id"`
	Catalog   string `json:"catalog"`
}
