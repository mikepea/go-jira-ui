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
	selectedLine     int
	uiList           *ui.List
	displayLines     []string
	cachedResults    []Query
	firstDisplayLine int
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
	if p.selectedLine < p.firstDisplayLine {
		p.firstDisplayLine = p.selectedLine
	}
}

func (p *QueryPage) NextLine(n int) {
	if p.selectedLine < len(p.cachedResults)-n {
		p.selectedLine = p.selectedLine + n
	} else {
		p.selectedLine = len(p.cachedResults) - 1
	}
	if p.selectedLine > p.lastDisplayedLine() {
		p.firstDisplayLine = p.firstDisplayLine + n
	}
}

func (p *QueryPage) PreviousPage() {
	p.PreviousLine(p.uiList.Height - 2)
}

func (p *QueryPage) NextPage() {
	p.NextLine(p.uiList.Height - 2)
}

func (p *QueryPage) lastDisplayedLine() int {
	return lastLineDisplayed(p.uiList, p.firstDisplayLine, 3)
}

func (p *QueryPage) markActiveLine() {
	for i, v := range origQueries {
		selected := ""
		if i == p.selectedLine {
			selected = "fg-white,bg-blue"
		}
		p.displayLines[i] = fmt.Sprintf("[%-20s %s](%s)", v.Name, v.JQL, selected)
	}
}

func (p *QueryPage) SelectedQuery() Query {
	return origQueries[p.selectedLine]
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
