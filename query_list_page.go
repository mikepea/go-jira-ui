package main

import (
	"fmt"
	ui "github.com/gizak/termui"
)

type Query struct {
	Name string
	JQL  string
}

var querySelected = 0
var activeQueryList *ui.List

var origQueries = []Query{
	Query{"My Assigned Tickets", "assignee = currentUser() AND resolution = Unresolved"},
	Query{"My Reported Tickets", "reporter = currentUser() AND resolution = Unresolved"},
	Query{"My Watched Tickets", "watcher = currentUser() AND resolution = Unresolved"},
	Query{"OPS unlabelled", "project = OPS AND labels IS EMPTY AND resolution = Unresolved"},
	Query{"Ops Queue", "project = OPS AND resolution = Unresolved"},
}

var displayQueries []string

func prevQuery(n int) {
	querySelected = querySelected - n
	if querySelected < 0 {
		querySelected = 0
	}
}

func nextQuery(n int) {
	if querySelected < len(origQueries)-1 {
		querySelected = querySelected + n
	}
}

func markActiveQuery() {
	for i, v := range origQueries {
		selected := ""
		if i == querySelected {
			selected = "fg-white,bg-blue"
		}
		displayQueries[i] = fmt.Sprintf("[%s](%s)", v.Name, selected)
	}
}

func updateQueryPage(ls *ui.List) {
	markActiveQuery()
	ls.Items = displayQueries
	ui.Render(ls)
}

func handleTicketQueryPage() {
	ui.Clear()
	ls := ui.NewList()
	displayQueries = make([]string, len(origQueries))
	ls.Items = displayQueries
	ls.ItemFgColor = ui.ColorYellow
	ls.BorderLabel = "List"
	ls.Height = ui.TermHeight()
	ls.Width = ui.TermWidth()
	ls.Y = 0
	activeQueryList = ls
	markActiveQuery()
	ui.Render(ls)
}
