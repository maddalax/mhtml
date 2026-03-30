package tailwind

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsV4_Latest(t *testing.T) {
	t.Parallel()
	assert.True(t, IsV4("latest"))
}

func TestIsV4_Empty(t *testing.T) {
	t.Parallel()
	assert.True(t, IsV4(""))
}

func TestIsV4_V4Version(t *testing.T) {
	t.Parallel()
	assert.True(t, IsV4("v4.2.2"))
}

func TestIsV4_V4NoPrefixVersion(t *testing.T) {
	t.Parallel()
	assert.True(t, IsV4("4.2.2"))
}

func TestIsV4_V3Version(t *testing.T) {
	t.Parallel()
	assert.False(t, IsV4("v3.4.16"))
}

func TestIsV4_V3NoPrefixVersion(t *testing.T) {
	t.Parallel()
	assert.False(t, IsV4("3.4.16"))
}

func TestDownloadURL_Latest(t *testing.T) {
	t.Parallel()
	url := DownloadURL("latest", "tailwindcss-linux-x64")
	assert.Equal(t, "https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-x64", url)
}

func TestDownloadURL_Empty(t *testing.T) {
	t.Parallel()
	url := DownloadURL("", "tailwindcss-linux-x64")
	assert.Equal(t, "https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-x64", url)
}

func TestDownloadURL_Pinned(t *testing.T) {
	t.Parallel()
	url := DownloadURL("v3.4.16", "tailwindcss-linux-x64")
	assert.Equal(t, "https://github.com/tailwindlabs/tailwindcss/releases/download/v3.4.16/tailwindcss-linux-x64", url)
}

func TestDownloadURL_PinnedNoPrefix(t *testing.T) {
	t.Parallel()
	url := DownloadURL("3.4.16", "tailwindcss-linux-x64")
	assert.Equal(t, "https://github.com/tailwindlabs/tailwindcss/releases/download/v3.4.16/tailwindcss-linux-x64", url)
}

func TestDetectTailwindVersion_ExplicitConfig(t *testing.T) {
	t.Parallel()
	version := DetectTailwindVersion("v4.2.2", "")
	assert.Equal(t, "v4.2.2", version)
}

func TestDetectTailwindVersion_ExplicitLatest(t *testing.T) {
	t.Parallel()
	version := DetectTailwindVersion("latest", "")
	assert.Equal(t, "latest", version)
}

func TestDetectTailwindVersion_AutoDetectV3(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	cssDir := filepath.Join(dir, "assets", "css")
	os.MkdirAll(cssDir, 0755)
	os.WriteFile(filepath.Join(cssDir, "input.css"), []byte("@tailwind base;\n@tailwind components;\n@tailwind utilities;\n"), 0644)
	version := DetectTailwindVersion("", dir)
	assert.Equal(t, "v3.4.16", version)
}

func TestDetectTailwindVersion_AutoDetectV4(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	cssDir := filepath.Join(dir, "assets", "css")
	os.MkdirAll(cssDir, 0755)
	os.WriteFile(filepath.Join(cssDir, "input.css"), []byte("@import \"tailwindcss\";\n"), 0644)
	version := DetectTailwindVersion("", dir)
	assert.Equal(t, "latest", version)
}

func TestDetectTailwindVersion_V4WithTypographyPlugin(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	cssDir := filepath.Join(dir, "assets", "css")
	os.MkdirAll(cssDir, 0755)
	os.WriteFile(filepath.Join(cssDir, "input.css"), []byte("@import \"tailwindcss\";\n@plugin \"@tailwindcss/typography\";\n"), 0644)
	version := DetectTailwindVersion("", dir)
	assert.Equal(t, "latest", version)
}

func TestDetectTailwindVersion_NoInputCss(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	version := DetectTailwindVersion("", dir)
	assert.Equal(t, "latest", version)
}
