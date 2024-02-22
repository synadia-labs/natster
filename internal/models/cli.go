package models

import "time"

type InitOptions struct {
	// Personal access token created in synadia cloud
	Token string
}

type HubOptions struct {
	Name        string
	Description string
	RootPath    string
	AllowAll    bool
}

type DownloadOptions struct {
	Hash       string
	OutputPath string
}

type ClaimOpts struct {
	Code          string
	OAuthIdentity string
}

type ShareOptions struct {
	Name       string
	AccountKey string
}

type WebLoginOpts struct {
	DisplayQR bool
}

// Options configure the CLI
type Options struct {
	// Creds is nats credentials to authenticate with
	Creds string
	// Timeout is how long to wait for operations
	Timeout time.Duration
	// Context in which action is to be performed
	ContextName string
}

type NatsterContext struct {
	TeamID           string `json:"team_id"`
	SystemID         string `json:"system_id"`
	AccountID        string `json:"account_id"`
	AccountName      string `json:"account_name"`
	AccountPublicKey string `json:"account_public_key"`
	Token            string `json:"access_token"`
	UserID           string `json:"user_id"`
	CredsPath        string `json:"creds"`
}
