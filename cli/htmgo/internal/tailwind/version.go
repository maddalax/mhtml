package tailwind

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const V3Default = "v3.4.16"

// IsV4 returns true if the version string indicates Tailwind CSS v4+.
// "latest" and "" are treated as v4 since v3 is EOL.
func IsV4(version string) bool {
	if version == "latest" || version == "" {
		return true
	}
	v := strings.TrimPrefix(version, "v")
	parts := strings.Split(v, ".")
	if len(parts) == 0 {
		return true
	}
	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return true
	}
	return major >= 4
}

// DownloadURL returns the GitHub releases URL for the given Tailwind version and binary filename.
func DownloadURL(version, fileName string) string {
	if version == "latest" || version == "" {
		return fmt.Sprintf("https://github.com/tailwindlabs/tailwindcss/releases/latest/download/%s", fileName)
	}
	if !strings.HasPrefix(version, "v") {
		version = "v" + version
	}
	return fmt.Sprintf("https://github.com/tailwindlabs/tailwindcss/releases/download/%s/%s", version, fileName)
}

// DetectTailwindVersion resolves the Tailwind version to use.
// Priority: explicit configVersion > auto-detect from input.css > "latest".
func DetectTailwindVersion(configVersion, workingDir string) string {
	if configVersion != "" {
		return configVersion
	}
	inputCssPath := filepath.Join(workingDir, "assets", "css", "input.css")
	data, err := os.ReadFile(inputCssPath)
	if err != nil {
		return "latest"
	}
	// Check for v3 @tailwind directives (e.g., "@tailwind base;").
	// Must not match "@tailwindcss" which appears in v4 plugin references.
	if strings.Contains(string(data), "@tailwind ") {
		return V3Default
	}
	return "latest"
}
