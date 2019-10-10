package jiraui

import (
	ui "gopkg.in/gizak/termui.v2"
)

type DebugPage struct {
	BaseListPage
	CommandBarFragment
	StatusBarFragment
}

func (p *DebugPage) Search() {
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
	startLine := (p.uiList.Cursor + n + increment) % n
	for i := startLine; i != p.uiList.Cursor; i = (i + increment + n) % n {
		if s.re.MatchString(p.cachedResults[i]) {
			p.uiList.SetCursorLine(i)
			p.Update()
			break
		}
	}
}

func (p *DebugPage) GoBack() {
	currentPage, previousPages = previousPages[len(previousPages)-1], previousPages[:len(previousPages)-1]
	changePage()
}

func (p *DebugPage) Refresh() {
	pDeref := &p
	q := *pDeref
	q.cachedResults = make([]string, 0)
	debugPage = q
	currentPage = debugPage
	changePage()
	q.Create()
}

func (p *DebugPage) Update() {
	ls := p.uiList
	ui.Render(ls)
	p.statusBar.Update()
	p.commandBar.Update()
}

func (p *DebugPage) Create() {
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
		// logBuffer contains output of all log.* calls
		p.cachedResults = logBuffer
	}
	ls.Items = p.cachedResults
	ls.ItemFgColor = ui.ColorYellow
	ls.BorderLabel = "Debug"
	ls.Height = ui.TermHeight() - 2
	ls.Width = ui.TermWidth()
	ls.Y = 0
	p.statusBar.Create()
	p.commandBar.Create()
	p.Update()
}
