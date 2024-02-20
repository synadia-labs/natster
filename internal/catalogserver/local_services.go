package catalogserver

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/nats-io/nats.go"
	"github.com/synadia-io/control-plane-sdk-go/syncp"
	"github.com/synadia-labs/natster/internal/models"
)

const (
	synadiaCloudBaseUrl = "https://cloud.synadia.com"
)

type inboxResponse struct {
	UnimportedShares []models.CatalogShareSummary `json:"unimported_shares"`
}

func handleLocalServiceRequest(srv *CatalogServer) func(m *nats.Msg) {
	return func(m *nats.Msg) {
		// - natster.local.inbox (req)

		tokens := strings.Split(m.Subject, ".")
		// shouldn't ever happen, but better safe than sorry
		if len(tokens) < 3 {
			_ = m.Respond(models.NewApiResultFail("Bad request", 400))
			return
		}
		scpClient, scpContext := getSynadiaCloudClient(srv.nctx)
		if tokens[2] == "inbox" {
			inboxData, err := srv.getInbox(scpClient, scpContext)
			if err != nil {
				_ = m.Respond(models.NewApiResultFail(err.Error(), 500))
			} else {
				_ = m.Respond(models.NewApiResultPass(inboxData))
			}
		}
	}
}

func (srv *CatalogServer) getInbox(scpClient *syncp.APIClient, scpContext context.Context) (*inboxResponse, error) {
	cats, err := srv.globalServiceClient.GetMyCatalogs()
	if err != nil {
		return nil, err
	}

	resp, _, err := scpClient.AccountAPI.ListSubjectImports(scpContext, srv.nctx.AccountID).Execute()
	if err != nil {
		return nil, err
	}
	importedSubjects := make([]string, len(resp.Items))
	for i, subject := range resp.Items {
		importedSubjects[i] = subject.Name
	}

	inboxCatalogs := make([]models.CatalogShareSummary, 0)
	for _, catalog := range cats {
		target := fmt.Sprintf("natster_%s", strings.ToLower(catalog.Catalog))
		if !slices.Contains(importedSubjects, target) &&
			catalog.FromAccount != srv.nctx.AccountPublicKey {
			inboxCatalogs = append(inboxCatalogs, catalog)
		}
	}
	return &inboxResponse{
		UnimportedShares: inboxCatalogs,
	}, nil
}

func getSynadiaCloudClient(nctx *models.NatsterContext) (*syncp.APIClient, context.Context) {
	client := syncp.NewAPIClient(syncp.NewConfiguration())
	ctxx := context.WithValue(context.Background(), syncp.ContextServerVariables, map[string]string{
		"baseUrl": synadiaCloudBaseUrl,
	})
	ctxx = context.WithValue(ctxx, syncp.ContextAccessToken, nctx.Token)

	return client, ctxx
}
