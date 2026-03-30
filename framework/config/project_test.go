package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"testing"
)

func TestDefaultProjectConfig(t *testing.T) {
	t.Parallel()
	cfg := DefaultProjectConfig()
	assert.Equal(t, true, cfg.Tailwind)
	assert.Equal(t, 4, len(cfg.WatchIgnore))
	assert.Equal(t, 8, len(cfg.WatchFiles))
}

func TestNoConfigFileUsesDefault(t *testing.T) {
	t.Parallel()
	cfg := FromConfigFile("non-existing-dir")
	assert.Equal(t, true, cfg.Tailwind)
	assert.Equal(t, 4, len(cfg.WatchIgnore))
	assert.Equal(t, 8, len(cfg.WatchFiles))
}

func TestPartialConfigMerges(t *testing.T) {
	t.Parallel()
	dir := writeConfigFile(t, "tailwind: false")
	cfg := FromConfigFile(dir)
	assert.Equal(t, false, cfg.Tailwind)
	assert.Equal(t, 4, len(cfg.WatchIgnore))
	assert.Equal(t, 8, len(cfg.WatchFiles))
}

func TestShouldNotSetTailwindTrue(t *testing.T) {
	t.Parallel()
	dir := writeConfigFile(t, "someValue: true")
	cfg := FromConfigFile(dir)
	assert.Equal(t, false, cfg.Tailwind)
	assert.Equal(t, 4, len(cfg.WatchIgnore))
	assert.Equal(t, 8, len(cfg.WatchFiles))
}

func TestShouldPrefixAutomaticPageRoutingIgnore(t *testing.T) {
	t.Parallel()
	cfg := DefaultProjectConfig()
	cfg.AutomaticPageRoutingIgnore = []string{"somefile"}
	cfg.Enhance()
	assert.Equal(t, []string{"pages/somefile"}, cfg.AutomaticPageRoutingIgnore)
}

func TestShouldPrefixAutomaticPageRoutingIgnore_1(t *testing.T) {
	t.Parallel()
	cfg := DefaultProjectConfig()
	cfg.AutomaticPageRoutingIgnore = []string{"pages/somefile/*"}
	cfg.Enhance()
	assert.Equal(t, []string{"pages/somefile/*"}, cfg.AutomaticPageRoutingIgnore)
}

func TestShouldPrefixAutomaticPartialRoutingIgnore(t *testing.T) {
	t.Parallel()
	cfg := DefaultProjectConfig()
	cfg.AutomaticPartialRoutingIgnore = []string{"somefile/*"}
	cfg.Enhance()
	assert.Equal(t, []string{"partials/somefile/*"}, cfg.AutomaticPartialRoutingIgnore)
}

func TestShouldPrefixAutomaticPartialRoutingIgnore_1(t *testing.T) {
	t.Parallel()
	cfg := DefaultProjectConfig()
	cfg.AutomaticPartialRoutingIgnore = []string{"partials/somefile/*"}
	cfg.Enhance()
	assert.Equal(t, []string{"partials/somefile/*"}, cfg.AutomaticPartialRoutingIgnore)
}

func TestPublicAssetPath(t *testing.T) {
	t.Parallel()
	cfg := DefaultProjectConfig()
	assert.Equal(t, "/public", cfg.PublicAssetPath)

	cfg.PublicAssetPath = "/assets"
	assert.Equal(t, "/assets", cfg.PublicAssetPath)
}

func TestConfigGet(t *testing.T) {
	t.Parallel()
	cfg := Get()
	assert.Equal(t, "/public", cfg.PublicAssetPath)
}

func TestTailwindVersionDefault(t *testing.T) {
	t.Parallel()
	cfg := DefaultProjectConfig()
	assert.Equal(t, "", cfg.TailwindVersion)
}

func TestTailwindVersionFromConfig(t *testing.T) {
	t.Parallel()
	dir := writeConfigFile(t, "tailwind: true\ntailwind_version: v3.4.16")
	cfg := FromConfigFile(dir)
	assert.Equal(t, "v3.4.16", cfg.TailwindVersion)
}

func TestTailwindVersionLatest(t *testing.T) {
	t.Parallel()
	dir := writeConfigFile(t, "tailwind: true\ntailwind_version: latest")
	cfg := FromConfigFile(dir)
	assert.Equal(t, "latest", cfg.TailwindVersion)
}

func writeConfigFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	err := os.WriteFile(path.Join(dir, "htmgo.yml"), []byte(content), 0644)
	assert.Nil(t, err)
	return dir
}
