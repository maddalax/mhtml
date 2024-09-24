package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/maddalax/htmgo/cli/htmgo/tasks/astgen"
	"github.com/maddalax/htmgo/cli/htmgo/tasks/copyassets"
	"github.com/maddalax/htmgo/cli/htmgo/tasks/css"
	"github.com/maddalax/htmgo/cli/htmgo/tasks/downloadtemplate"
	"github.com/maddalax/htmgo/cli/htmgo/tasks/process"
	"github.com/maddalax/htmgo/cli/htmgo/tasks/reloader"
	"github.com/maddalax/htmgo/cli/htmgo/tasks/run"
	"log/slog"
	"os"
	"strings"
)

func main() {
	done := RegisterSignals()

	commandMap := make(map[string]*flag.FlagSet)
	commands := []string{"template", "run", "watch", "build", "setup", "css", "schema", "generate"}

	for _, command := range commands {
		commandMap[command] = flag.NewFlagSet(command, flag.ExitOnError)
	}

	if len(os.Args) < 2 {
		fmt.Println(fmt.Sprintf("Usage: htmgo [%s]", strings.Join(commands, " | ")))
		os.Exit(1)
	}

	c := commandMap[os.Args[1]]

	if c == nil {
		fmt.Println(fmt.Sprintf("Usage: htmgo [%s]", strings.Join(commands, " | ")))
		os.Exit(1)
		return
	}

	err := c.Parse(os.Args[2:])

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
		return
	}

	slog.SetLogLoggerLevel(getLogLevel())

	taskName := os.Args[1]

	slog.Debug("Running task:", slog.String("task", taskName))
	slog.Debug("working dir:", slog.String("dir", process.GetWorkingDir()))

	if taskName == "watch" {
		fmt.Printf("Running in watch mode\n")
		os.Setenv("ENV", "development")
		os.Setenv("WATCH_MODE", "true")
		fmt.Printf("Starting processes...")
		copyassets.CopyAssets()
		astgen.GenAst(process.ExitOnError)
		css.GenerateCss(process.ExitOnError)
		run.EntGenerate()
		go func() {
			css.GenerateCssWatch(process.ExitOnError)
		}()
		fmt.Printf("Starting server...\n")
		go func() {
			_ = run.Server()
		}()
		startWatcher(reloader.OnFileChange)
	} else {
		if taskName == "schema" {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Enter entity name:")
			text, _ := reader.ReadString('\n')
			text = strings.TrimSuffix(text, "\n")
			run.EntNewSchema(text)
		} else if taskName == "generate" {
			run.EntGenerate()
			astgen.GenAst(process.ExitOnError)
		} else if taskName == "setup" {
			run.Setup()
		} else if taskName == "css" {
			_ = css.GenerateCss(process.ExitOnError)
		} else if taskName == "ast" {
			_ = astgen.GenAst(process.ExitOnError)
		} else if taskName == "run" {
			_ = astgen.GenAst(process.ExitOnError)
			_ = css.GenerateCss(process.ExitOnError)
			_ = run.Server(process.ExitOnError)
		} else if taskName == "template" {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("What would you like to call your new app?: ")
			text, _ := reader.ReadString('\n')
			text = strings.TrimSuffix(text, "\n")
			text = strings.ReplaceAll(text, " ", "-")
			text = strings.ToLower(text)
			downloadtemplate.DownloadTemplate(fmt.Sprintf("./%s", text))
		} else if taskName == "build" {
			run.Build()
		} else if taskName == "generate" {
			astgen.GenAst(process.ExitOnError)
		} else {
			fmt.Println(fmt.Sprintf("Usage: htmgo [%s]", strings.Join(commands, " | ")))
		}
		os.Exit(0)
	}

	<-done
}
