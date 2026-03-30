package tailwind

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadCachedVersion_NoFile(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	assert.Equal(t, "", ReadCachedVersion(dir))
}

func TestReadCachedVersion_Exists(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	htmgoDir := filepath.Join(dir, "__htmgo")
	os.MkdirAll(htmgoDir, 0755)
	os.WriteFile(filepath.Join(htmgoDir, "tailwind.version"), []byte("v4.2.2"), 0644)
	assert.Equal(t, "v4.2.2", ReadCachedVersion(dir))
}

func TestWriteCachedVersion(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	htmgoDir := filepath.Join(dir, "__htmgo")
	os.MkdirAll(htmgoDir, 0755)
	WriteCachedVersion(dir, "v4.2.2")
	data, err := os.ReadFile(filepath.Join(htmgoDir, "tailwind.version"))
	assert.NoError(t, err)
	assert.Equal(t, "v4.2.2", string(data))
}

func TestNeedsRedownload_NoCache(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	assert.True(t, NeedsRedownload(dir, "v4.2.2"))
}

func TestNeedsRedownload_SameVersion(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	htmgoDir := filepath.Join(dir, "__htmgo")
	os.MkdirAll(htmgoDir, 0755)
	os.WriteFile(filepath.Join(htmgoDir, "tailwind.version"), []byte("v4.2.2"), 0644)
	assert.False(t, NeedsRedownload(dir, "v4.2.2"))
}

func TestNeedsRedownload_DifferentVersion(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	htmgoDir := filepath.Join(dir, "__htmgo")
	os.MkdirAll(htmgoDir, 0755)
	os.WriteFile(filepath.Join(htmgoDir, "tailwind.version"), []byte("v3.4.16"), 0644)
	assert.True(t, NeedsRedownload(dir, "v4.2.2"))
}
