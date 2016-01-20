package jiraui

import (
	"strings"
)

type CommandBarFragment struct {
	commandBar  *CommandBar
	commandMode bool
}

func (p *CommandBarFragment) ExecuteCommand() {
	command := string(p.commandBar.text)
	if command == "" {
		return
	}
	commandMode := string([]rune(command)[0])
	switch commandMode {
	case "/":
		log.Debugf("Search down: %q", command)
		if obj, ok := currentPage.(Searcher); ok {
			obj.SetSearch(command)
			obj.Search()
		}
	case "?":
		log.Debugf("Search up: %q", command)
		if obj, ok := currentPage.(Searcher); ok {
			obj.SetSearch(command)
			obj.Search()
		}
	}
}

func (p *CommandBarFragment) SetCommandMode(mode bool) {
	p.commandMode = mode
}

func (p *CommandBarFragment) CommandMode() bool {
	return p.commandMode
}

func (p *CommandBarFragment) CommandBar() *CommandBar {
	return p.commandBar
}
