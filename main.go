package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

const AppName = "sorascope"

const version = "0.0.1"

var revision = "HEAD"

type config struct {
	Host     string `json:"host"`
	Handle   string `json:"handle"`
	Password string `json:"password"`
	dir      string
	verbose  bool
	prefix   string
}

func main() {
	app := &cli.App{
		Name:        AppName,
		Usage:       AppName,
		Version:     version,
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
				UsageText:   AppName + " get-posts",
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
				UsageText:   AppName + " get-blocks",
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
				UsageText:   AppName + " show-session",
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "json", Usage: "output JSON"},
				},
				Action: doShowSession,
			},
			{
				Name:        "login",
				Description: "Login the social",
				Usage:       "Login the social",
				UsageText:   AppName + " login [handle] [password]",
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
			cfg, fp, err := loadConfig(profile)
			cCtx.App.Metadata["path"] = fp
			if cCtx.Args().Get(0) == "login" {
				return nil
			}
			if err != nil {
				return fmt.Errorf("cannot load config file: %w", err)
			}
			cCtx.App.Metadata["config"] = cfg
			cfg.verbose = cCtx.Bool("V")
			if profile != "" {
				cfg.prefix = profile + "-"
			}
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
