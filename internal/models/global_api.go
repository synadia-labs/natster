package models

const (
	NatsterInitializedEventType = "natster_initialized"
	CatalogSharedEventType      = "catalog_shared"
	ContextBoundEventType       = "context_bound"
)

// Heartbeats are emitted to the global service periodically by running natster
// catalog servers
type Heartbeat struct {
	Catalog    string `json:"catalog"`
	AccountKey string `json:"account_key"`
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

// Occurs when someone with an OAuth ID claims a one-time code, resulting in a binding
// between the context originally submitted with the code and the OAuth ID, which in turn
// allows that web user to access their natster resources
type ContextBoundEvent struct {
	OAuthIdentity string         `json:"oauth_id"`
	BoundContext  NatsterContext `json:"bound_context"`
}

type CommunityStats struct {
	TotalInitialized uint64 `json:"total_initialized"`
	RunningCatalogs  uint64 `json:"running_catalogs"`
	// Total number of catalogs shared with others
	SharedCatalogs uint64 `json:"share_count"`
}

type CatalogShareSummary struct {
	FromAccount   string `json:"from_account"`
	ToAccount     string `json:"to_account"`
	Catalog       string `json:"catalog"`
	CatalogOnline bool   `json:"catalog_online"`
	// TODO: add a timestamp here
}

type WhoamiResponse struct {
	AccountKey    string  `json:"account_key"`
	OAuthIdentity *string `json:"oauth_id,omitempty"`
}

type ContextQueryResponse struct {
	Context   NatsterContext `json:"context"`
	FullCreds string         `json:"full_creds"`
}

type OtcGenerateResponse struct {
	Code         string `json:"code"`
	ValidMinutes int    `json:"valid_minutes"`
	ClaimUrl     string `json:"claim_url"`
}

type OtcClaimRequest struct {
	Code          string `json:"code"`
	OAuthIdentity string `json:"oauth_id"`
}
