package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/urfave/cli/v2"
)

const AppName = "sorascope"

const Version = "0.0.1"

var Revision = "HEAD"

type Config struct {
	Host     string `json:"host"`
	Handle   string `json:"handle"`
	Password string `json:"password"`
	Dir      string
	Verbose  bool
	Prefix   string
}

func ConfigDir() (string, error) {
	switch runtime.GOOS {
	case "darwin":
		dir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(dir, ".config"), nil
	default:
		return os.UserConfigDir()

	}
}

func GetConfigFromCtx(cCtx *cli.Context) (cfg *Config, err error) {
	config, ok := cCtx.App.Metadata["config"]
	if !ok {
		return nil, fmt.Errorf("config not found")
	}
	cfg = config.(*Config)
	return cfg, nil
}

func GetConfigFpFromCtx(cCtx *cli.Context) (fp string, err error) {
	filePath, ok := cCtx.App.Metadata["path"]
	if !ok {
		return "", fmt.Errorf("config not found")
	}
	fp = filePath.(string)
	return fp, nil
}

func LoadConfig(profile string) (*Config, string, error) {
	dir, err := ConfigDir()
	if err != nil {
		return nil, "", err
	}
	dir = filepath.Join(dir, AppName)

	var fp string
	if profile == "" {
		fp = filepath.Join(dir, "config.json")
	} else if profile == "?" {
		names, err := filepath.Glob(filepath.Join(dir, "config-*.json"))
		if err != nil {
			return nil, "", err
		}
		for _, name := range names {
			name = filepath.Base(name)
			name = strings.TrimLeft(name[6:len(name)-5], "-")
			//fmt.Println(name)
		}
		os.Exit(0)
	} else {
		fp = filepath.Join(dir, "config-"+profile+".json")
	}
	os.MkdirAll(filepath.Dir(fp), 0700)

	b, err := os.ReadFile(fp)
	if err != nil {
		return nil, fp, fmt.Errorf("cannot load config file: %w", err)
	}
	var cfg Config
	err = json.Unmarshal(b, &cfg)
	if err != nil {
		return nil, fp, fmt.Errorf("cannot load config file: %w", err)
	}
	if cfg.Host == "" {
		cfg.Host = "https://bsky.social"
	}
	cfg.Dir = dir
	return &cfg, fp, nil
}
