package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/henoya/sorascope/util"
	"io/ioutil"
	"os"

	"github.com/henoya/sorascope/config"

	comatproto "github.com/bluesky-social/indigo/api/atproto"

	"github.com/urfave/cli/v2"
)

func doLogin(cCtx *cli.Context) error {
	fp, _ := cCtx.App.Metadata["path"].(string)
	var cfg config.Config
	cfg.Host = cCtx.String("host")
	cfg.Handle = cCtx.Args().Get(0)
	cfg.Password = cCtx.Args().Get(1)
	if cfg.Handle == "" || cfg.Password == "" {
		cli.ShowSubcommandHelpAndExit(cCtx, 1)
	}
	b, err := json.MarshalIndent(&cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("cannot make config file: %w", err)
	}
	err = ioutil.WriteFile(fp, b, 0644)
	if err != nil {
		return fmt.Errorf("cannot write config file: %w", err)
	}
	return nil
}

func doShowSession(cCtx *cli.Context) error {
	xrpcc, err := util.MakeXRPCC(cCtx)
	if err != nil {
		return fmt.Errorf("cannot create client: %w", err)
	}

	session, err := comatproto.ServerGetSession(context.TODO(), xrpcc)
	if err != nil {
		return err
	}

	if cCtx.Bool("json") {
		json.NewEncoder(os.Stdout).Encode(session)
		return nil
	}

	fmt.Printf("Did: %s\n", session.Did)
	fmt.Printf("Email: %s\n", util.Stringp(session.Email))
	fmt.Printf("Handle: %s\n", session.Handle)
	return nil
}
