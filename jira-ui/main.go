package main

import (
	"fmt"
	"os/exec"

	"github.com/mikepea/go-jira-ui"

	"gopkg.in/op/go-logging.v1"
)

var log = logging.MustGetLogger("jiraui")

func resetTTY() {
	cmd := exec.Command("reset")
	_ = cmd.Run()
	fmt.Println()
}

func main() {
	defer resetTTY()
	jiraui.Run()
}
