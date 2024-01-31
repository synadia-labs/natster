package main

import (
	"context"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/choria-io/fisk"
	"github.com/synadia-io/control-plane-sdk-go/syncp"
)

const (
	baseUrl              = "https://cloud.synadia.com"
	natsterGlobalAccount = "ABNZ6NGGOKLCNJOSETMDLT6KLXR5NHBFLIZKVUTXXBOZILRLTDRYH5VZ"
)

// TODO: make the import and export creation idempotent. If they're already there, skip
// those operations
func InitNatster(ctx *fisk.ParseContext) error {
	client := syncp.NewAPIClient(syncp.NewConfiguration())
	ctxx := context.WithValue(context.Background(), syncp.ContextServerVariables, map[string]string{
		"baseUrl": baseUrl,
	})
	ctxx = context.WithValue(ctxx, syncp.ContextAccessToken, InitOpts.Token)
	resp, _, err := client.SessionAPI.ListTeams(ctxx).Execute()
	if err != nil {
		return err
	}

	teamNames := make([]string, len(resp.Items))

	for i, team := range resp.Items {
		teamNames[i] = team.Name
	}

	teamPrompt := &survey.Select{
		Message: "Select a team:",
		Options: teamNames,
	}
	var selectedTeam int
	err = survey.AskOne(teamPrompt, &selectedTeam, survey.WithValidator(survey.Required))
	if err != nil {
		return err
	}

	teamId := resp.Items[selectedTeam].Id
	sysResp, _, err := client.TeamAPI.ListTeamSystems(ctxx, teamId).Execute()
	if err != nil {
		return err
	}

	systemNames := make([]string, len(sysResp.Items))
	for i, system := range sysResp.Items {
		systemNames[i] = system.Name
	}

	systemPrompt := &survey.Select{
		Message: "Select a system:",
		Options: systemNames,
	}
	var selectedSystem int
	err = survey.AskOne(systemPrompt, &selectedSystem, survey.WithValidator(survey.Required))
	if err != nil {
		return err
	}

	systemId := sysResp.Items[selectedSystem].Id

	acctResp, _, err := client.SystemAPI.ListAccounts(ctxx, systemId).Execute()
	if err != nil {
		return err
	}
	accountNames := make([]string, len(acctResp.Items))
	for i, account := range acctResp.Items {
		accountNames[i] = account.Name
	}

	accountPrompt := &survey.Select{
		Message: "Select an account:",
		Options: accountNames,
	}
	var selectedAccount int
	err = survey.AskOne(accountPrompt, &selectedAccount, survey.WithValidator(survey.Required))
	if err != nil {
		return err
	}

	accountId := acctResp.Items[selectedAccount].Id
	accountKey := acctResp.Items[selectedAccount].AccountPublicKey
	accountName := acctResp.Items[selectedAccount].Name

	jwt := syncp.Export{}

	// token position is 1-based since 0 means none
	jwt.AccountTokenPosition = syncp.Ptr(int32(1))
	jwt.Advertise = syncp.Ptr(true)
	jwt.Subject = syncp.Ptr("*.natster.catalog.>")
	jwt.Description = syncp.Ptr("Natster Catalog Service")
	jwt.Name = syncp.Ptr("natster_catalog")
	jwt.InfoUrl = syncp.Ptr("https://natster.io")
	jwt.ResponseType = syncp.Ptr(syncp.RESPONSETYPE_SINGLETON)
	jwt.Type = syncp.Ptr(syncp.EXPORTTYPE_SERVICE)
	req := syncp.SubjectExportCreateRequest{
		JwtSettings:               jwt,
		MetricsEnabled:            false,
		MetricsSamplingPercentage: 0,
	}
	_, _, err = client.AccountAPI.CreateSubjectExport(ctxx, accountId).SubjectExportCreateRequest(req).Execute()
	if err != nil {
		return err
	}

	importReq := syncp.SubjectImportCreateRequest{
		JwtSettings: syncp.Import{
			Account:      syncp.Ptr(natsterGlobalAccount),
			Subject:      syncp.Ptr(fmt.Sprintf("%s.natster.global.>", accountKey)),
			LocalSubject: syncp.Ptr("natster.global.>"),
			Name:         syncp.Ptr("natster_global"),
			Type:         syncp.Ptr(syncp.EXPORTTYPE_SERVICE),
		},
	}
	_, _, err = client.AccountAPI.CreateSubjectImport(ctxx, accountId).SubjectImportCreateRequest(importReq).Execute()
	if err != nil {
		return err
	}

	// TODO: instead of prompting for a creds file we should list users in this
	// account and prompt them and we download that user's creds automatically

	newCtx := NatsterContext{
		TeamID:           teamId,
		SystemID:         systemId,
		AccountID:        accountId,
		AccountName:      accountName,
		AccountPublicKey: accountKey,
		Token:            InitOpts.Token,
		CredsPath:        Opts.Creds,
	}
	err = writeContext(newCtx)
	if err != nil {
		return err
	}

	fmt.Printf("Congratulations! Your account (%s) is now ready to serve Natster catalogs!\n", accountName)
	fmt.Println("You can now use `natster catalog serve` to host a media catalog and `natster catalog share` to share with friends.")

	return nil
}
