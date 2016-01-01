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
	selectedLine int
	uiList       *ui.List
	displayLines []string
}

var origQueries = []Query{
	Query{"My Assigned Tickets", "assignee = currentUser() AND resolution = Unresolved"},
	Query{"My Reported Tickets", "reporter = currentUser() AND resolution = Unresolved"},
	Query{"My Watched Tickets", "watcher = currentUser() AND resolution = Unresolved"},
	Query{"OPS unlabelled", "project = OPS AND labels IS EMPTY AND resolution = Unresolved"},
	Query{"Ops Queue", "project = OPS AND resolution = Unresolved"},
}

func (p *QueryPage) PreviousLine(n int) {
	p.selectedLine = p.selectedLine - n
	if p.selectedLine < 0 {
		p.selectedLine = 0
	}
}

func (p *QueryPage) NextLine(n int) {
	if p.selectedLine < len(origQueries)-n {
		p.selectedLine = p.selectedLine + n
	}
}

func (p *QueryPage) markActiveLine() {
	log.Noticef("markActiveLine: p = %s", &p)
	for i, v := range origQueries {
		log.Noticef("markActiveLine: displayLines = %s", p.displayLines)
		selected := ""
		if i == p.selectedLine {
			selected = "fg-white,bg-blue"
		}
		p.displayLines[i] = fmt.Sprintf("[%s](%s)", v.Name, selected)
	}
}

func (p *QueryPage) SelectedQuery() Query {
	return origQueries[p.selectedLine]
}

func (p *QueryPage) Update() {
	log.Noticef("Update: p = %s", &p)
	log.Noticef("Update: displayLines = %s", p.displayLines)
	ls := p.uiList
	p.markActiveLine()
	ls.Items = p.displayLines
	ui.Render(ls)
}

func (p *QueryPage) Create() {
	ui.Clear()
	ls := ui.NewList()
	p.selectedLine = 0
	p.displayLines = make([]string, len(origQueries))
	ls.Items = p.displayLines
	ls.ItemFgColor = ui.ColorYellow
	ls.BorderLabel = "List"
	ls.Height = ui.TermHeight()
	ls.Width = ui.TermWidth()
	ls.Y = 0
	p.uiList = ls
	p.markActiveLine()
	ui.Render(ls)
}
