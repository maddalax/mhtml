package config

import (
	"gopkg.in/yaml.v3"
	"log/slog"
	"os"
	"path"
	"strings"
)

type ProjectConfig struct {
	Tailwind                      bool     `yaml:"tailwind"`
	TailwindVersion               string   `yaml:"tailwind_version"`
	WatchIgnore                   []string `yaml:"watch_ignore"`
	WatchFiles                    []string `yaml:"watch_files"`
	AutomaticPageRoutingIgnore    []string `yaml:"automatic_page_routing_ignore"`
	AutomaticPartialRoutingIgnore []string `yaml:"automatic_partial_routing_ignore"`
	PublicAssetPath               string   `yaml:"public_asset_path"`
}

func DefaultProjectConfig() *ProjectConfig {
	return &ProjectConfig{
		Tailwind: true,
		WatchIgnore: []string{
			"node_modules", ".git", ".idea", "assets/dist",
		},
		WatchFiles: []string{
			"**/*.go", "**/*.html", "**/*.css", "**/*.js", "**/*.json", "**/*.yaml", "**/*.yml", "**/*.md",
		},
		PublicAssetPath: "/public",
	}
}

func (cfg *ProjectConfig) Enhance() *ProjectConfig {
	defaultCfg := DefaultProjectConfig()
	if len(cfg.WatchFiles) == 0 {
		cfg.WatchFiles = defaultCfg.WatchFiles
	}
	if len(cfg.WatchIgnore) == 0 {
		cfg.WatchIgnore = defaultCfg.WatchIgnore
	}

	for i, s := range cfg.AutomaticPartialRoutingIgnore {
		parts := strings.Split(s, string(os.PathSeparator))
		if len(parts) == 0 {
			continue
		}
		if parts[0] != "partials" {
			cfg.AutomaticPartialRoutingIgnore[i] = path.Join("partials", s)
		}
	}

	for i, s := range cfg.AutomaticPageRoutingIgnore {
		parts := strings.Split(s, string(os.PathSeparator))
		if len(parts) == 0 {
			continue
		}
		if parts[0] != "pages" {
			cfg.AutomaticPageRoutingIgnore[i] = path.Join("pages", s)
		}
	}

	if cfg.PublicAssetPath == "" {
		cfg.PublicAssetPath = "/public"
	}

	return cfg
}

func Get() *ProjectConfig {
	cwd, err := os.Getwd()
	if err != nil {
		return DefaultProjectConfig()
	}
	config := FromConfigFile(cwd)
	return config
}

func FromConfigFile(workingDir string) *ProjectConfig {
	defaultCfg := DefaultProjectConfig()
	names := []string{"htmgo.yaml", "htmgo.yml", "_htmgo.yaml", "_htmgo.yml"}
	for _, name := range names {
		filePath := path.Join(workingDir, name)
		if _, err := os.Stat(filePath); err == nil {
			cfg := &ProjectConfig{}
			bytes, err := os.ReadFile(filePath)
			if err == nil {
				err = yaml.Unmarshal(bytes, cfg)
				if err != nil {
					slog.Error("Error parsing config file", slog.String("file", filePath), slog.String("error", err.Error()))
					os.Exit(1)
				}
				return cfg.Enhance()
			}
		}
	}
	return defaultCfg
}
