package main

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/AlecAivazis/survey/v2"
	"github.com/ConnectEverything/natster/internal/globalservice"
	"github.com/ConnectEverything/natster/internal/models"
	"github.com/choria-io/fisk"
	"github.com/nats-io/jsm.go/natscontext"
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

	users, _, err := client.AccountAPI.ListUsers(ctxx, accountId).Execute()
	if err != nil {
		return err
	}

	usernames := make([]string, len(users.Items))
	for i, u := range users.Items {
		usernames[i] = u.Name
	}

	userPrompt := survey.Select{
		Message: "Select a user for NATS authentication:",
		Options: usernames,
	}
	var selectedUser int
	err = survey.AskOne(&userPrompt, &selectedUser, survey.WithValidator(survey.Required))
	if err != nil {
		return err
	}
	userId := users.Items[selectedUser].Id

	creds, _, err := client.NatsUserAPI.DownloadNatsUserCreds(ctxx, userId).Execute()
	if err != nil {
		return err
	}

	home, err := getNatsterHome()
	if err != nil {
		return err
	}
	credsFileName := path.Join(home, ".creds")
	err = os.WriteFile(credsFileName, []byte(creds), 0655)
	if err != nil {
		return nil
	}

	newCtx := NatsterContext{
		TeamID:           teamId,
		SystemID:         systemId,
		AccountID:        accountId,
		AccountName:      accountName,
		AccountPublicKey: accountKey,
		Token:            InitOpts.Token,
		CredsPath:        credsFileName,
	}
	err = writeContext(newCtx)
	if err != nil {
		return err
	}

	// Use the newly established context to publish an initialized event
	// on the new natster global import
	ctxOpts := []natscontext.Option{
		natscontext.WithServerURL("tls://connect.ngs.global"),
		natscontext.WithCreds(newCtx.CredsPath),
	}
	natsContext, err := natscontext.New("natster", false, ctxOpts...)
	if err != nil {
		return err
	}
	conn, err := natsContext.Connect()
	if err != nil {
		return err
	}
	globalClient := globalservice.NewClient(conn)
	err = globalClient.PublishEvent(models.NatsterInitializedEventType, "none", "none",
		models.NatsterInitializedEvent{
			AccountId:   newCtx.AccountID,
			AccountName: newCtx.AccountName,
			AccountKey:  newCtx.AccountPublicKey,
		},
	)
	if err != nil {
		fmt.Printf("Failed to contact Natster global service to post initialization event: %s", err)
		return err
	}

	fmt.Printf("Congratulations! Your account (%s) is now ready to serve Natster catalogs!\n", accountName)
	fmt.Println("You can now use `natster catalog serve` to host a media catalog and `natster catalog share` to share with friends.")

	return nil
}
