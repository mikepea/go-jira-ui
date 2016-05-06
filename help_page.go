package jiraui

import (
	ui "github.com/gizak/termui"
)

type HelpPage struct {
	BaseListPage
	CommandBarFragment
	StatusBarFragment
}

func (p *HelpPage) Search() {
	s := p.ActiveSearch
	n := len(p.cachedResults)
	if s.command == "" {
		return
	}
	increment := 1
	if s.directionUp {
		increment = -1
	}
	// we use modulo here so we can loop through every line.
	// adding 'n' means we never have '-1 % n'.
	startLine := (p.selectedLine + n + increment) % n
	for i := startLine; i != p.selectedLine; i = (i + increment + n) % n {
		if s.re.MatchString(p.cachedResults[i]) {
			p.SetSelectedLine(i)
			p.Update()
			break
		}
	}
}

func (p *HelpPage) GoBack() {
	currentPage = previousPage
	changePage()
}

func (p *HelpPage) Refresh() {
	pDeref := &p
	q := *pDeref
	q.cachedResults = make([]string, 0)
	helpPage = q
	currentPage = helpPage
	changePage()
	q.Create()
}

func (p *HelpPage) Update() {
	ls := p.uiList
	p.markActiveLine()
	ls.Items = p.displayLines[p.firstDisplayLine:]
	ui.Render(ls)
	p.statusBar.Update()
	p.commandBar.Update()
}

func (p *HelpPage) Create() {
	ui.Clear()
	ls := NewScrollableList()
	p.uiList = ls
	if p.statusBar == nil {
		p.statusBar = new(StatusBar)
	}
	if p.commandBar == nil {
		p.commandBar = commandBar
	}
	if len(p.cachedResults) == 0 {
		p.cachedResults = HelpTextAsStrings(nil, "jira_ui_help")
	}
	p.displayLines = make([]string, len(p.cachedResults))
	ls.ItemFgColor = ui.ColorYellow
	ls.BorderLabel = "Help"
	ls.Height = ui.TermHeight() - 2
	ls.Width = ui.TermWidth()
	ls.Y = 0
	p.statusBar.Create()
	p.commandBar.Create()
	p.Update()
}
