package tailwind

import (
	"os"
	"path/filepath"
	"strings"
)

func versionFilePath(workingDir string) string {
	return filepath.Join(workingDir, "__htmgo", "tailwind.version")
}

// ReadCachedVersion reads the cached Tailwind version from the sidecar file.
// Returns empty string if file doesn't exist.
func ReadCachedVersion(workingDir string) string {
	data, err := os.ReadFile(versionFilePath(workingDir))
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}

// WriteCachedVersion writes the resolved version to the sidecar file.
func WriteCachedVersion(workingDir, version string) {
	os.WriteFile(versionFilePath(workingDir), []byte(version), 0644)
}

// NeedsRedownload returns true if the cached version doesn't match the resolved version.
func NeedsRedownload(workingDir, resolvedVersion string) bool {
	cached := ReadCachedVersion(workingDir)
	return cached != resolvedVersion
}
