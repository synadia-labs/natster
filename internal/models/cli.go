package models

import "time"

type HubOptions struct {
	Name        string
	Description string
	RootPath    string
	Port        int
}

// Options configure the CLI
type Options struct {
	// List of server nodes
	Servers string
	// Creds is nats credentials to authenticate with
	Creds string
	// TlsCert is the TLS Public Certificate
	TlsCert string
	// TlsKey is the TLS Private Key
	TlsKey string
	// TlsCA is the certificate authority to verify the connection with
	TlsCA string
	// Timeout is how long to wait for operations
	Timeout time.Duration
	// ConnectionName is the name to use for the underlying NATS connection
	ConnectionName string
	// Username is the username or token to connect with
	Username string
	// Password is the password to connect with
	Password string
	// Nkey is the file holding a nkey to connect with
	Nkey string
	// Trace enables verbose debug logging
	Trace bool
	// SocksProxy is a SOCKS5 proxy to use for NATS connections
	SocksProxy string
	// TlsFirst configures theTLSHandshakeFirst behavior in nats.go
	TlsFirst bool
}
