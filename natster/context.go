package main

import (
	"encoding/json"
	"os"
	"path"

	"github.com/synadia-labs/natster/internal/models"
)

func loadContext() (*models.NatsterContext, error) {
	home, err := getNatsterHome()
	if err != nil {
		return nil, err
	}
	file := path.Join(home, ".context")
	bytes, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var context models.NatsterContext
	err = json.Unmarshal(bytes, &context)
	if err != nil {
		return nil, err
	}
	return &context, nil
}

func writeContext(ctx models.NatsterContext) error {

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
