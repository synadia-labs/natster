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

	Opts    = &models.Options{}
	HubOpts = &models.HubOptions{}
)

func main() {
	blue := color.New(color.FgBlue).SprintFunc()
	help := fmt.Sprintf("%s\nNatster Media Hub %s\n", blue(Banner), VERSION)

	ncli := fisk.New("natster", help)
	ncli.Author("Synadia Communications")
	ncli.UsageWriter(os.Stdout)
	ncli.Version(fmt.Sprintf("v%s [%s] | Built-on: %s", VERSION, COMMIT, BUILDDATE))
	ncli.HelpFlag.Short('h')
	ncli.WithCheats().CheatCommand.Hidden()

	ncli.Flag("server", "NATS server urls").Short('s').Envar("NATS_URL").PlaceHolder("URL").StringVar(&Opts.Servers)
	ncli.Flag("user", "Username or Token").Envar("NATS_USER").PlaceHolder("USER").StringVar(&Opts.Username)
	ncli.Flag("password", "Password").Envar("NATS_PASSWORD").PlaceHolder("PASSWORD").StringVar(&Opts.Password)
	ncli.Flag("creds", "User credentials file (JWT authentication)").Envar("NATS_CREDS").PlaceHolder("FILE").StringVar(&Opts.Creds)
	ncli.Flag("nkey", "User NKEY file for single-key auth").Envar("NATS_NKEY").PlaceHolder("FILE").StringVar(&Opts.Nkey)
	ncli.Flag("tlscert", "TLS public certificate file").Envar("NATS_CERT").PlaceHolder("FILE").ExistingFileVar(&Opts.TlsCert)
	ncli.Flag("tlskey", "TLS private key file").Envar("NATS_KEY").PlaceHolder("FILE").ExistingFileVar(&Opts.TlsKey)
	ncli.Flag("tlsca", "TLS certificate authority chain file").Envar("NATS_CA").PlaceHolder("FILE").ExistingFileVar(&Opts.TlsCA)
	ncli.Flag("tlsfirst", "Perform TLS handshake before expecting the server greeting").BoolVar(&Opts.TlsFirst)
	ncli.Flag("timeout", "Time to wait on responses from NATS").Default("2s").Envar("NATS_TIMEOUT").PlaceHolder("DURATION").DurationVar(&Opts.Timeout)

	hub := ncli.Command("hub", "Interact with the media hub")
	hub_up := hub.Command("up", "Starts the media hub server")
	hub_up.Arg("name", "The name of the library (single word)").Required().StringVar(&HubOpts.Name)
	hub_up.Arg("description", "Description of the library").Required().StringVar(&HubOpts.Description)
	hub_up.Arg("path", "Path to the root directory of the library").Required().ExistingDirVar(&HubOpts.RootPath)
	hub_up.Flag("port", "HTTP port on which to run the UI/API").Default("8080").IntVar(&HubOpts.Port)
	hub_up.Action(HubUp)

	ncli.MustParseWithUsage(os.Args[1:])
}
