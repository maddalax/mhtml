package copyassets

import (
	"fmt"
	"github.com/maddalax/htmgo/cli/htmgo/internal/dirutil"
	"github.com/maddalax/htmgo/cli/htmgo/internal/tailwind"
	"github.com/maddalax/htmgo/cli/htmgo/tasks/module"
	"github.com/maddalax/htmgo/cli/htmgo/tasks/process"
	"golang.org/x/mod/modfile"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func getModuleVersion(modulePath string) (string, error) {
	// Read the go.mod file
	data, err := os.ReadFile(process.GetPathRelativeToCwd("go.mod"))
	if err != nil {
		return "", fmt.Errorf("error reading go.mod: %v", err)
	}

	// Parse the go.mod file
	modFile, err := modfile.Parse(process.GetPathRelativeToCwd("go.mod"), data, nil)
	if err != nil {
		return "", fmt.Errorf("error parsing go.mod: %v", err)
	}

	// Find the module version
	for _, req := range modFile.Require {
		if req.Mod.Path == modulePath {
			return req.Mod.Version, nil
		}
	}

	return "", fmt.Errorf("module %s not found in go.mod", modulePath)
}

func CopyAssets() {
	dirutil.CreateHtmgoDir()

	moduleName := "github.com/maddalax/htmgo/framework"
	modulePath := module.GetDependencyPath(moduleName)

	assetDir := ""
	// Is hosted version and not local version from .work file
	if strings.HasPrefix(modulePath, "github.com/") {
		version, err := getModuleVersion(modulePath)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
		dirname, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}
		assetDir = fmt.Sprintf("%s/go/pkg/mod/%s@%s/assets", dirname, modulePath, version)
	} else {
		assetDir = filepath.Join(modulePath, "assets")
	}

	assetDistDir := filepath.Join(assetDir, "dist")
	assetCssDir := filepath.Join(assetDir, "css")

	cwd := process.GetWorkingDir()

	destDir := filepath.Join(cwd, "assets")
	destDirDist := filepath.Join(destDir, "dist")
	destDirCss := filepath.Join(destDir, "css")

	err := dirutil.CopyDir(assetDistDir, destDirDist, func(path string, exists bool) bool {
		return true
	})

	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	err = dirutil.CopyDir(assetCssDir, destDirCss, func(path string, exists bool) bool {
		if strings.HasSuffix(path, "tailwind.config.js") {
			return false
		}
		return !exists
	})

	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	if dirutil.HasFileFromRoot("assets/public") {
		err = dirutil.CopyDir(filepath.Join(process.GetWorkingDir(), "assets/public"), filepath.Join(process.GetWorkingDir(), "assets/dist"),
			func(path string, exists bool) bool {
				return true
			})
	}

	cfg := dirutil.GetConfig()
	if cfg.Tailwind && !dirutil.HasFileFromRoot("tailwind.config.js") {
		if !tailwind.IsV4(tailwind.DetectTailwindVersion(cfg.TailwindVersion, process.GetWorkingDir())) {
			err = dirutil.CopyFile(
				filepath.Join(assetCssDir, "tailwind.config.js"),
				filepath.Join(process.GetWorkingDir(), "tailwind.config.js"),
			)
		}
	}

	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	cmd := fmt.Sprintf("cd %s && git add .", destDirCss)
	process.Run(process.NewRawCommand(cmd, cmd, process.Silent))
}
