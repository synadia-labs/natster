package main

import (
	"encoding/json"
	"os"
	"path"
)

type NatsterContext struct {
	TeamID           string `json:"team_id"`
	SystemID         string `json:"system_id"`
	AccountID        string `json:"account_id"`
	AccountName      string `json:"account_name"`
	AccountPublicKey string `json:"account_public_key"`
	Token            string `json:"access_token"`
	CredsPath        string `json:"creds"`
}

func loadContext() (*NatsterContext, error) {
	home, err := getNatsterHome()
	if err != nil {
		return nil, err
	}
	file := path.Join(home, ".context")
	bytes, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var context NatsterContext
	err = json.Unmarshal(bytes, &context)
	if err != nil {
		return nil, err
	}
	return &context, nil
}

func writeContext(ctx NatsterContext) error {

	home, err := getNatsterHome()
	if err != nil {
		return err
	}

	bytes, err := json.Marshal(ctx)
	if err != nil {
		return err
	}

	file := path.Join(home, ".context")
	err = os.WriteFile(file, bytes, 0644)
	if err != nil {
		return err
	}

	return nil
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
