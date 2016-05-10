package main

import (
	"os/exec"

	"github.com/mikepea/go-jira-ui"
)

func resetTTY() {
	cmd := exec.Command("reset")
	_ = cmd.Run()
}

func main() {
	defer resetTTY()
	jiraui.Run()
}
