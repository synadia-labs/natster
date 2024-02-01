package main

import (
	"fmt"
	"os"

	"github.com/ConnectEverything/natster/internal/models"
	"github.com/choria-io/fisk"
	"github.com/fatih/color"
)

var (
	VERSION   = "development"
	COMMIT    = ""
	BUILDDATE = ""

	Opts      = &models.Options{}
	HubOpts   = &models.HubOptions{}
	InitOpts  = &models.InitOptions{}
	ShareOpts = &models.ShareOptions{}
)

func main() {
	blue := color.New(color.FgBlue).SprintFunc()
	help := fmt.Sprintf("%s\nNatster %s\n", blue(Banner), VERSION)

	ncli := fisk.New("natster", help)
	ncli.Author("Synadia Communications")
	ncli.UsageWriter(os.Stdout)
	ncli.Version(fmt.Sprintf("v%s [%s] | Built-on: %s", VERSION, COMMIT, BUILDDATE))
	ncli.HelpFlag.Short('h')
	ncli.WithCheats().CheatCommand.Hidden()

	ncli.Flag("timeout", "Time to wait on responses from NATS").Default("2s").Envar("NATS_TIMEOUT").PlaceHolder("DURATION").DurationVar(&Opts.Timeout)

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

	sharecat := catalog.Command("share", "Shares a catalog with a target account")
	sharecat.Arg("name", "The name of the catalog to share").Required().StringVar(&ShareOpts.Name)
	sharecat.Arg("account", "Public key of the target account").Required().StringVar(&ShareOpts.AccountKey)
	sharecat.Action(ShareCatalog)

	hub_up := catalog.Command("serve", "Starts the media catalog server")
	hub_up.Arg("name", "The name of the catalog to serve").Required().StringVar(&HubOpts.Name)
	hub_up.Flag("port", "HTTP port on which to run the UI/API").Default("8080").IntVar(&HubOpts.Port)
	hub_up.Action(StartCatalogServer)

	ncli.MustParseWithUsage(os.Args[1:])

	fmt.Println("") // why?? WHY??
}
