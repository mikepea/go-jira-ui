package main

import (
	"fmt"
	"os/exec"

	"github.com/mikepea/go-jira-ui"
)

func resetTTY() {
	cmd := exec.Command("reset")
	_ = cmd.Run()
	fmt.Println()
}

func main() {
	defer resetTTY()
	jiraui.Run()
}
