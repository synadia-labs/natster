package main

import (
	"fmt"

	"github.com/choria-io/fisk"
	"github.com/synadia-labs/natster/internal/globalservice"
)

const (
	natsterDotIoAccount = "AA2JVG74M2LCCNYYFMBAANHRNFTUAVWZJGTDAREKDBE23DRLBRWYQNLD"
)

func WebLogin(ctx *fisk.ParseContext) error {
	nctx, err := loadContext()
	if err != nil {
		return err
	}
	globalClient, err := globalservice.NewClientWithCredsPath(nctx.CredsPath)
	if err != nil {
		return err
	}

	response, err := globalClient.GenerateOneTimeCode(*nctx)
	if err != nil {
		fmt.Printf("Failed to generate one-time code: %s\n", err)
		return err
	}

	fmt.Printf("Login code generated: %s\n", response.Code)
	fmt.Printf("For the next %d minutes, you can claim this code at the following URL: %s\nAfter that you'll need to generate a new login code.\n",
		response.ValidMinutes,
		response.ClaimUrl,
	)

	return nil
}

func ClaimOtc(ctx *fisk.ParseContext) error {
	nctx, err := loadContext()
	if err != nil {
		return err
	}
	if nctx.AccountPublicKey != natsterDotIoAccount {
		fmt.Printf("Invalid account: %s\n", nctx.AccountPublicKey)
		fmt.Println("This is a debug facility only to be used within the context of the natster.io site account.")
		return nil
	}
	globalClient, err := globalservice.NewClientWithCredsPath(nctx.CredsPath)
	if err != nil {
		return err
	}
	response, err := globalClient.ClaimOneTimeCode(ClaimOpts.Code, ClaimOpts.OAuthIdentity)
	if err != nil {
		return err
	}

	fmt.Printf("DEBUG - Claimed OTC: \n\n%+v\n\n", response)
	return nil

}
