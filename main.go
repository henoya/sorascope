package main

import (
	"fmt"
	"github.com/henoya/sorascope/config"
	"github.com/henoya/sorascope/user"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:        config.AppName,
		Usage:       config.AppName,
		Version:     config.Version,
		Description: "A cli application for sorascope",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "a", Usage: "profile name"},
			&cli.BoolFlag{Name: "V", Usage: "verbose"},
		},
		Commands: []*cli.Command{
			{
				Name:        "get-posts",
				Description: "Get account's all posts",
				Usage:       "get-posts -H handle [-n number] [--json]",
				UsageText:   config.AppName + " get-posts",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "handle", Aliases: []string{"H"}, Required: true, Usage: "user handle or DID"},
					&cli.IntFlag{Name: "n", Value: 30, Usage: "number of items"},
					&cli.BoolFlag{Name: "json", Usage: "output JSON"},
				},
				Action: doGetPosts,
			},
			{
				Name:        "get-blocks",
				Description: "get-blocks",
				Usage:       "get-blocks",
				UsageText:   config.AppName + " get-blocks",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "handle", Aliases: []string{"H"}, Value: "", Usage: "user handle"},
					&cli.IntFlag{Name: "n", Value: 30, Usage: "number of items"},
					&cli.BoolFlag{Name: "json", Usage: "output JSON"},
				},
				Action: doGetBlocks,
			},
			{
				Name:        "show-session",
				Description: "Show session",
				Usage:       "Show session",
				UsageText:   config.AppName + " show-session",
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "json", Usage: "output JSON"},
				},
				Action: doShowSession,
			},
			{
				Name:        "add-user",
				Description: "add user",
				Usage:       "add-user",
				UsageText:   config.AppName + " add-user [handle] [password]",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "host", Value: "https://bsky.social"},
					&cli.StringFlag{Name: "handle", Aliases: []string{"H"}, Value: ""},
					&cli.StringFlag{Name: "did", Aliases: []string{"D"}, Value: ""},
					&cli.StringFlag{Name: "app-pass", Aliases: []string{"P"}, Required: true},
				},
				HelpName: "add-user",
				Action:   user.DoAddUser,
				Before: func(cCtx *cli.Context) error {
					handle := cCtx.String("handle")
					did := cCtx.String("did")
					if handle == "" && did == "" {
						return fmt.Errorf("Need handle or did parameter")
					}
					return nil
				},
			},
			{
				Name:        "login",
				Description: "Login the social",
				Usage:       "Login the social",
				UsageText:   config.AppName + " login [handle] [password]",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "host", Value: "https://bsky.social"},
				},
				HelpName: "login",
				Action:   doLogin,
			},
		},
		Metadata: map[string]any{},
		Before: func(cCtx *cli.Context) error {
			profile := cCtx.String("a")
			cfg, fp, err := config.LoadConfig(profile)
			cCtx.App.Metadata["path"] = fp

			if cCtx.Args().Get(0) == "login" {
				return nil
			}
			if err != nil {
				return fmt.Errorf("cannot load config file: %w", err)
			}
			cCtx.App.Metadata["config"] = cfg
			cfg.Verbose = cCtx.Bool("V")
			if profile != "" {
				cfg.Prefix = profile + "-"
			}
			db, err := InitDBConnection()
			if err != nil {
				return err
			}
			cCtx.App.Metadata["db"] = db
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
