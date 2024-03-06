package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/choria-io/fisk"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/synadia-labs/natster/internal/globalservice"
	"github.com/synadia-labs/natster/internal/models"
)

func DisplayContext(ctx *fisk.ParseContext) error {
	ctxx, err := loadContext()
	if err != nil {
		return err
	}
	client, err := globalservice.NewClientWithCredsPath(ctxx.CredsPath)
	if err != nil {
		return err
	}

	idString := "(unlinked)"
	initializedOn := "(never)"
	whoami, _ := client.Whoami()

	if whoami != nil {
		if whoami.OAuthIdentity != nil {
			idString = *whoami.OAuthIdentity
		}
		if whoami.Initialized > 0 {
			t := time.Unix(whoami.Initialized, 0)
			initializedOn = t.Format("2006-01-02 15:04:05")
		}
	}

	t := newTableWriter(ctxx.AccountName, "cyan")
	w := t.writer
	w.AppendRow(table.Row{"Account", ctxx.AccountPublicKey})
	w.AppendRow(table.Row{"Initialized At", initializedOn})
	w.AppendRow(table.Row{"Synadia Cloud Team", ctxx.TeamID})
	w.AppendRow(table.Row{"Synadia Cloud System", ctxx.SystemID})
	w.AppendRow(table.Row{"Synadia Cloud User", ctxx.UserID})
	w.AppendRow(table.Row{"Credentials", ctxx.CredsPath})
	w.AppendRow(table.Row{"Natster.io Login", idString})
	fmt.Println(w.Render())

	return nil
}

func loadContext() (*models.NatsterContext, error) {
	home, err := getNatsterHome()
	if err != nil {
		return nil, err
	}
	file := path.Join(home, fmt.Sprintf("%s.context", Opts.ContextName))
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

func getLocalLibraries() ([]string, error) {
	home, err := getNatsterHome()
	if err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(home)
	if err != nil {
		return nil, err
	}
	libraries := make([]string, 0)
	for _, entry := range entries {
		libraries = append(libraries, strings.ToLower(entry.Name()))
	}
	return libraries, nil
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

	file := path.Join(home, fmt.Sprintf("%s.context", Opts.ContextName))
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
