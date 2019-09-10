package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/coryb/figtree"
	"github.com/coryb/oreo"

	"github.com/mikepea/go-jira-ui"

	"github.com/go-jira/jira/jiracli"
	"github.com/go-jira/jira/jiracmd"
	"gopkg.in/op/go-logging.v1"
)

var log = logging.MustGetLogger("jiraui")

func resetTTY() {
	cmd := exec.Command("reset")
	_ = cmd.Run()
	fmt.Println()
}

type oreoLogger struct {
	logger *logging.Logger
}

func (ol *oreoLogger) Printf(format string, args ...interface{}) {
	ol.logger.Debugf(format, args...)
}

func main() {
	defer resetTTY()

	configDir := ".jira.d"
	fig := figtree.NewFigTree(
		figtree.WithHome(jiracli.Homedir()),
		figtree.WithEnvPrefix("JIRA"),
		figtree.WithConfigDir(configDir),
	)

	if err := os.MkdirAll(filepath.Join(jiracli.Homedir(), configDir), 0755); err != nil {
		log.Errorf("%s", err)
		panic(jiracli.Exit{Code: 1})
	}

	o := oreo.New().WithCookieFile(filepath.Join(jiracli.Homedir(), configDir, "cookies.js")).WithLogger(&oreoLogger{log})

	jiracmd.RegisterAllCommands()

	app := jiracli.CommandLine(fig, o)

	jiraui.Run(app)
}
