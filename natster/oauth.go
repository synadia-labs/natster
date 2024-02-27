package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"

	"github.com/choria-io/fisk"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/synadia-labs/natster/internal/globalservice"
	"github.com/synadia-labs/natster/internal/models"
	"golang.org/x/oauth2"
)

const (
	AUTH0_DOMAIN        = "natster.us.auth0.com"
	AUTH0_CLIENT_ID     = "96EEzrbHVLHzgY2JVbgtb9b7ecOjYQsO"
	AUTH0_CLIENT_SECRET = ""
	AUTH0_CALLBACK_URL  = "http://127.0.0.1:8088/callback"
)

func OauthLogin(ctx *fisk.ParseContext) error {
	state, err := func() (string, error) {
		b := make([]byte, 32)
		_, err := rand.Read(b)
		if err != nil {
			return "", err
		}

		state := base64.StdEncoding.EncodeToString(b)

		return state, nil
	}()
	if err != nil {
		return err
	}

	auth, err := New(state)
	if err != nil {
		return errors.New("Failed to initialize the authenticator: " + err.Error())
	}

	srv := &http.Server{Addr: ":8088"}
	go func() {
		http.HandleFunc("/callback", auth.CallbackHandler)
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	fmt.Println("Open your browser and navigate to: ", auth.AuthCodeURL(auth.state))
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", auth.AuthCodeURL(auth.state)).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", auth.AuthCodeURL(auth.state)).Start()
	case "darwin":
		err = exec.Command("open", auth.AuthCodeURL(auth.state)).Start()
	}
	if err != nil {
		fmt.Println("Failed to automatically open browser. Please do so manually.")
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		auth.loginComplete <- struct{}{}
	}()

	// Wait for the login to complete
	<-auth.loginComplete
	err = srv.Shutdown(context.TODO())
	if err != nil {
		return err
	}

	return nil
}

type Authenticator struct {
	state         string
	natsCtx       *models.NatsterContext
	loginComplete chan struct{}
	globalservice *globalservice.Client
	*oidc.Provider
	oauth2.Config
}

func New(state string) (*Authenticator, error) {
	provider, err := oidc.NewProvider(
		context.Background(),
		"https://"+AUTH0_DOMAIN+"/",
	)
	if err != nil {
		return nil, err
	}

	nctx, err := loadContext()
	if err != nil {
		return nil, err
	}

	globalClient, err := globalservice.NewClientWithCredsPath(nctx.CredsPath)
	if err != nil {
		return nil, err
	}

	conf := oauth2.Config{
		ClientID:     AUTH0_CLIENT_ID,
		ClientSecret: AUTH0_CLIENT_SECRET,
		RedirectURL:  AUTH0_CALLBACK_URL,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile"},
	}

	return &Authenticator{
		state:         state,
		loginComplete: make(chan struct{}, 1),
		natsCtx:       nctx,
		globalservice: globalClient,
		Provider:      provider,
		Config:        conf,
	}, nil
}

// VerifyIDToken verifies that an *oauth2.Token is a valid *oidc.IDToken.
func (a *Authenticator) VerifyIDToken(ctx context.Context, token *oauth2.Token) (*oidc.IDToken, error) {
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, errors.New("no id_token field in oauth2 token")
	}

	oidcConfig := &oidc.Config{
		ClientID: a.ClientID,
	}

	return a.Verifier(oidcConfig).Verify(ctx, rawIDToken)
}

// func CallbackHandler(auth *authenticator.Authenticator) gin.HandlerFunc {
func (a Authenticator) CallbackHandler(w http.ResponseWriter, r *http.Request) {
	// Exchange an authorization code for a token.
	queryCode := r.URL.Query().Get("code")
	token, err := a.Exchange(r.Context(), queryCode)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("Failed to exchange an authorization code for a token."))
		return
	}

	idToken, err := a.VerifyIDToken(r.Context(), token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Failed to verify ID Token."))
		return
	}

	var profile map[string]interface{}
	if err := idToken.Claims(&profile); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	oauthProfile, ok := profile["sub"].(string)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Auth0 did not return a valid oauth_id to bind to"))
		return
	}

	cbe := models.ContextBoundEvent{
		OAuthIdentity: oauthProfile,
		BoundContext:  *a.natsCtx,
	}

	r_cbe, err := json.Marshal(cbe)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	err = a.globalservice.PublishEvent(
		models.ContextBoundEventType,
		"none",
		"none",
		r_cbe,
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	fmt.Printf("\nBound Synadia Cloud account [%s | %s...] to oauth_id -> %s\nAttempting to redirect to natster.io",
		a.natsCtx.AccountName,
		a.natsCtx.AccountPublicKey[:8],
		profile["sub"])

	a.loginComplete <- struct{}{}
	http.Redirect(w, r, "https://natster.io/#/library", http.StatusPermanentRedirect)
}
