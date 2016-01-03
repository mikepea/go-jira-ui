package main

import (
	"fmt"
	ui "github.com/gizak/termui"
)

type Query struct {
	Name string
	JQL  string
}

type QueryPage struct {
	BaseListPage
	cachedResults []Query
}

var origQueries = []Query{
	Query{"My Assigned Tickets", "assignee = currentUser() AND resolution = Unresolved"},
	Query{"My Reported Tickets", "reporter = currentUser() AND resolution = Unresolved"},
	Query{"My Watched Tickets", "watcher = currentUser() AND resolution = Unresolved"},
	Query{"OPS unlabelled", "project = OPS AND labels IS EMPTY AND resolution = Unresolved"},
	Query{"Ops Queue", "project = OPS AND resolution = Unresolved"},
}

func (p *QueryPage) markActiveLine() {
	for i, v := range p.cachedResults {
		selected := ""
		if i == p.selectedLine {
			selected = "fg-white,bg-blue"
		}
		p.displayLines[i] = fmt.Sprintf("[%-30s -- %s](%s)", v.Name, v.JQL, selected)
	}
}

func (p *QueryPage) SelectedQuery() Query {
	return p.cachedResults[p.selectedLine]
}

func (p *QueryPage) SelectItem() {
	previousPage = currentPage
	currentPage = &ticketListPage
	changePage()
}

func (p *QueryPage) Update() {
	ls := p.uiList
	p.markActiveLine()
	ls.Items = p.displayLines[p.firstDisplayLine:]
	ui.Render(ls)
}

func (p *QueryPage) Create() {
	ui.Clear()
	ls := ui.NewList()
	p.uiList = ls
	p.selectedLine = 0
	p.firstDisplayLine = 0
	p.cachedResults = origQueries
	p.displayLines = make([]string, len(p.cachedResults))
	ls.ItemFgColor = ui.ColorYellow
	ls.BorderLabel = "Queries"
	ls.Height = ui.TermHeight()
	ls.Width = ui.TermWidth()
	ls.Y = 0
	p.Update()
}
