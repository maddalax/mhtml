package css

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"

	"github.com/maddalax/htmgo/cli/htmgo/internal/dirutil"
	"github.com/maddalax/htmgo/cli/htmgo/internal/tailwind"
	"github.com/maddalax/htmgo/cli/htmgo/tasks/copyassets"
	"github.com/maddalax/htmgo/cli/htmgo/tasks/process"
)

func resolveVersion() string {
	cfg := dirutil.GetConfig()
	return tailwind.DetectTailwindVersion(cfg.TailwindVersion, process.GetWorkingDir())
}

func IsTailwindEnabled() bool {
	cfg := dirutil.GetConfig()
	if !cfg.Tailwind {
		return false
	}
	version := resolveVersion()
	if tailwind.IsV4(version) {
		return true
	}
	return dirutil.HasFileFromRoot("tailwind.config.js")
}

func Setup() bool {
	if !IsTailwindEnabled() {
		slog.Debug("Tailwind is not enabled. Skipping CSS generation.")
		return false
	}
	downloadTailwindCli()

	if !dirutil.HasFileFromRoot("assets/css/input.css") {
		copyassets.CopyAssets()
	}

	return true
}

func GetTailwindExecutableName() string {
	if runtime.GOOS == "windows" {
		return "./__htmgo/tailwind.exe"
	}
	return "./__htmgo/tailwind"
}

func GenerateCss(flags ...process.RunFlag) error {
	if !Setup() {
		return nil
	}
	exec := GetTailwindExecutableName()
	version := resolveVersion()
	var cmd string
	if tailwind.IsV4(version) {
		cmd = fmt.Sprintf("%s -i ./assets/css/input.css -o ./assets/dist/main.css", exec)
	} else {
		cmd = fmt.Sprintf("%s -i ./assets/css/input.css -o ./assets/dist/main.css -c ./tailwind.config.js", exec)
	}
	return process.Run(process.NewRawCommand("tailwind", cmd, append(flags, process.Silent)...))
}

func downloadTailwindCli() {
	version := resolveVersion()
	workingDir := process.GetWorkingDir()

	binaryExists := dirutil.HasFileFromRoot(GetTailwindExecutableName())

	if binaryExists && !tailwind.NeedsRedownload(workingDir, version) {
		slog.Debug("Tailwind CLI already exists and version matches. Skipping download.")
		return
	}

	if binaryExists {
		slog.Debug("Tailwind version changed, re-downloading.", slog.String("version", version))
		execPath := filepath.Join(workingDir, GetTailwindExecutableName())
		os.Remove(execPath)
	}

	if !IsTailwindEnabled() {
		slog.Debug("Tailwind is not enabled. Skipping tailwind cli download.")
		return
	}

	distro := ""
	goos := runtime.GOOS
	arch := runtime.GOARCH
	switch {
	case goos == "darwin" && arch == "arm64":
		distro = "macos-arm64"
	case goos == "darwin" && arch == "amd64":
		distro = "macos-x64"
	case goos == "linux" && arch == "arm64":
		distro = "linux-arm64"
	case goos == "linux" && arch == "amd64":
		distro = "linux-x64"
	case goos == "windows" && arch == "amd64":
		distro = "windows-x64.exe"
	case goos == "windows" && arch == "arm64":
		if tailwind.IsV4(version) {
			log.Fatal("Tailwind CSS v4 does not provide a windows-arm64 binary. Please use tailwind_version: v3.4.16 or run on x64.")
		}
		distro = "windows-arm64.exe"
	default:
		log.Fatal(fmt.Sprintf("Unsupported OS/ARCH: %s/%s", goos, arch))
	}

	fileName := fmt.Sprintf("tailwindcss-%s", distro)
	url := tailwind.DownloadURL(version, fileName)

	cmd := fmt.Sprintf("curl -LO %s", url)
	process.Run(process.NewRawCommand("tailwind-cli-download", cmd, process.ExitOnError))

	outputFileName := GetTailwindExecutableName()
	newPath := filepath.Join(workingDir, outputFileName)

	err := dirutil.MoveFile(
		filepath.Join(workingDir, fileName),
		newPath)

	if err != nil {
		log.Fatalf("Error moving file: %s\n", err.Error())
	}

	if goos != "windows" {
		err = process.Run(process.NewRawCommand("chmod-tailwind-cli",
			fmt.Sprintf("chmod +x %s", newPath),
			process.ExitOnError))
	}

	if err != nil {
		log.Fatalf("Error setting executable permission: %s\n", err.Error())
	}

	tailwind.WriteCachedVersion(workingDir, version)
	slog.Debug("Successfully downloaded Tailwind CLI", slog.String("url", url), slog.String("version", version))
}
