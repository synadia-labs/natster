package main

import (
	"fmt"
	"os"
	"time"

	"github.com/choria-io/fisk"
	"github.com/fatih/color"
	"github.com/synadia-labs/natster/internal/models"
)

var (
	VERSION = "dev"
	COMMIT  = "none"
	DATE    = time.Now().Format(time.RFC822)

	Opts         = &models.Options{}
	HubOpts      = &models.HubOptions{}
	InitOpts     = &models.InitOptions{}
	ShareOpts    = &models.ShareOptions{}
	DlOpts       = &models.DownloadOptions{}
	ClaimOpts    = &models.ClaimOpts{}
	WebLoginOpts = &models.WebLoginOpts{}
)

func main() {
	blue := color.New(color.FgBlue).SprintFunc()
	help := fmt.Sprintf("%s\nNatster %s\n", blue(Banner), VERSION)

	ncli := fisk.New("natster", help)
	ncli.Author("Synadia Communications")
	ncli.UsageWriter(os.Stdout)
	ncli.Version(fmt.Sprintf("v%s [%s] | BuiltOn: %s", VERSION, COMMIT, DATE))
	ncli.HelpFlag.Short('h')
	ncli.WithCheats().CheatCommand.Hidden()

	ncli.Flag("timeout", "Time to wait on responses from NATS").Default("2s").Envar("NATS_TIMEOUT").PlaceHolder("DURATION").DurationVar(&Opts.Timeout)
	ncli.Flag("context", "Name of the context in which to perform the command").Default("default").StringVar(&Opts.ContextName)

	initcli := ncli.Command("init", "Initialize and configure the Natster CLI")
	initcli.Flag("token", "Synadia Cloud personal access token").Required().StringVar(&InitOpts.Token)
	//initcli.Flag("creds", "User credentials file (JWT authentication)").Envar("NATS_CREDS").PlaceHolder("FILE").Required().StringVar(&Opts.Creds)
	initcli.Action(InitNatster)

	catalog := ncli.Command("catalog", "Perform various activities related to a media catalog")
	newcat := catalog.Command("new", "Creates a new media catalog from a directory")
	newcat.Arg("name", "The name of the catalog (alphanumeric, no spaces)").Required().StringVar(&HubOpts.Name)
	newcat.Arg("description", "Description of the catalog").Required().StringVar(&HubOpts.Description)
	newcat.Arg("path", "Path to the root directory containing the catalog's media").Required().ExistingDirVar(&HubOpts.RootPath)
	newcat.Action(NewCatalog)

	inbox := catalog.Command("inbox", "Lists the catalogs that have been shared with you that have not yet been imported").
		HelpLong("This requires at least one catalog server running within the selected context. Will prompt the user to choose a catalog to import.")
	inbox.Action(RenderInbox)

	sharecat := catalog.Command("share", "Shares a catalog with a target account")
	sharecat.Arg("name", "The name of the catalog to share").Required().StringVar(&ShareOpts.Name)
	sharecat.Arg("account", "Public key of the target account").Required().StringVar(&ShareOpts.AccountKey)
	sharecat.Action(ShareCatalog)

	unsharecat := catalog.Command("unshare", "Stops sharing a catalog with a target account")
	unsharecat.Arg("name", "The name of the catalog to stop sharing").Required().StringVar(&ShareOpts.Name)
	unsharecat.Arg("account", "The public key of the target account").Required().StringVar(&ShareOpts.AccountKey)
	unsharecat.Action(UnshareCatalog)

	catimport := catalog.Command("import", "Imports a shared catalog")
	catimport.Arg("name", "Name of the catalog to import").Required().StringVar(&ShareOpts.Name)
	catimport.Arg("account", "Public key of the account from which to import").Required().StringVar(&ShareOpts.AccountKey)
	catimport.Action(ImportCatalog)

	catdl := catalog.Command("download", "Downloads a file from a catalog")
	catdl.Arg("name", "Name of the catalog from which to download the file").Required().StringVar(&ShareOpts.Name)
	catdl.Arg("hash", "SHA256 hash of the file to download. Hashes can be found in catalog metadata").Required().StringVar(&DlOpts.Hash)
	catdl.Arg("out", "Path to output file").Required().StringVar(&DlOpts.OutputPath)
	catdl.Action(DownloadFile)

	catcontents := catalog.Command("contents", "View the contents of a given catalog")
	catcontents.Arg("name", "Name of the catalog to view").Required().StringVar(&ShareOpts.Name)
	catcontents.Action(ViewCatalogItems)

	catls := catalog.Command("list", "Lists my shared catalogs and catalogs shared with me").Alias("ls")
	catls.Action(ListCatalogs)

	// TODO: reusing all these shared opts is awful. Need to isolate
	catdescribe := catalog.Command("describe", "Sets the description for a given file")
	catdescribe.Arg("name", "Name of the catalog").Required().StringVar(&ShareOpts.Name)
	catdescribe.Arg("hash", "Hash of the file to describe").Required().StringVar(&DlOpts.Hash)
	catdescribe.Arg("description", "Description of the file").Required().StringVar(&HubOpts.Description)
	catdescribe.Action(DescribeCatalogItem)

	hub_up := catalog.Command("serve", "Starts the media catalog server")
	hub_up.Arg("name", "The name of the catalog to serve").Required().StringVar(&HubOpts.Name)
	hub_up.Flag("allowall", "Disables security checks for contents and downloads - Use with caution").
		Envar("CATALOG_ALLOW_ALL").
		Default("false").
		UnNegatableBoolVar(&HubOpts.AllowAll)
	hub_up.Action(StartCatalogServer)

	auth := ncli.Command("auth", "Authenticate your local context for use with natster.io")
	weblogin := auth.Command("web", "Authenticate with one time code")
	weblogin.Flag("qrcode", "Displays QR code in terminal of login link").Default("false").UnNegatableBoolVar(&WebLoginOpts.DisplayQR)
	weblogin.Action(WebLogin)

	claim := ncli.Command("claim", "Claims an OTC code. For testing only - Can only be done from the natster.io account").Hidden()
	claim.Arg("code", "Previously generated one-time code").Required().StringVar(&ClaimOpts.Code)
	claim.Arg("identity", "OAuth identity string").Required().StringVar(&ClaimOpts.OAuthIdentity)
	claim.Action(ClaimOtc)

	ctxlookup := ncli.Command("oauthcheck", "Looks up the context bound to the given OAuth ID. Debug - can only be done from the natster.io account").Hidden()
	ctxlookup.Arg("oauthid", "OAuth identitifer to check").StringVar(&ClaimOpts.OAuthIdentity)
	ctxlookup.Action(LookupOAuthId)

	whoami := ncli.Command("whoami", "Displays information about the selected context")
	whoami.Action(DisplayContext)

	ncli.MustParseWithUsage(os.Args[1:])

	fmt.Println("") // why?? WHY??
}
