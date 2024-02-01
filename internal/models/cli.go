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
	Port        int
}

type ShareOptions struct {
	Name       string
	AccountKey string
}

// Options configure the CLI
type Options struct {
	// Creds is nats credentials to authenticate with
	Creds string
	// Timeout is how long to wait for operations
	Timeout time.Duration
}
