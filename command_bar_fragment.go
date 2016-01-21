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
	case ":":
		log.Debugf("Command: %q", command)
		handleCommand(command)
	}
}

func handleCommand(command string) {
	if len(command) < 2 {
		// must be :something
		return
	}
	fields := strings.Fields(string(command[1:]))
	action := fields[0]
	var args []string
	if len(fields) > 1 {
		args = fields[1:]
	}
	log.Debugf("handleCommand: action %q, args %s", action, args)
	switch {
	case action == "q" || action == "quit":
		handleQuit()
	case action == "label" || action == "labels":
		handleLabelCommand(args)
	case action == "assign":
		handleAssignCommand(args[0])
	case action == "unassign":
		handleAssignCommand("-1")
	case action == "take":
		opts := getJiraOpts()
		handleAssignCommand(opts["user"].(string))
	case action == "comment":
		handleCommentCommand(string(command[9:]))
	}
}

func handleLabelCommand(args []string) {
	log.Debugf("handleLabelCommand: args %s", args)
	if obj, ok := currentPage.(TicketCommander); ok {
		ticketId := obj.ActiveTicketId()
		if ticketId == "" || args == nil {
			return
		}
		action := "add"
		var labels []string
		switch args[0] {
		case "add":
			action = "add"
			if len(args) > 1 {
				labels = args[1:]
			}
		case "remove":
			action = "remove"
			if len(args) > 1 {
				labels = args[1:]
			}
		default:
			labels = args
		}
		runJiraCmdLabels(ticketId, action, labels)
		obj.Refresh()
	}
}

func handleCommentCommand(comment string) {
	log.Debugf("handleCommentCommand: comment %s", comment)
	if obj, ok := currentPage.(TicketCommander); ok {
		ticketId := obj.ActiveTicketId()
		if ticketId == "" || comment == "" {
			return
		}
		log.Debugf("handleCommentCommand: ticket: %s, comment %s", ticketId, comment)
		runJiraCmdCommentNoEditor(ticketId, comment)
		obj.Refresh()
	}
}

func handleAssignCommand(user string) {
	log.Debugf("handleAssignCommand: user %s", user)
	if obj, ok := currentPage.(TicketCommander); ok {
		ticketId := obj.ActiveTicketId()
		if ticketId == "" || user == "" {
			return
		}
		log.Debugf("handleAssignCommand: ticket: %s, user %s", ticketId, user)
		runJiraCmdAssign(ticketId, user)
		obj.Refresh()
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
