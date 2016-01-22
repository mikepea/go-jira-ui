package jiraui

import (
	"fmt"
	ui "github.com/gizak/termui"
)

type LabelListPage struct {
	BaseListPage
	CommandBarFragment
	StatusBarFragment
	labelCounts map[string]int
	ActiveQuery Query
}

func (p *LabelListPage) Search() {
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

func (p *LabelListPage) labelsAsSortedList() []string {
	return sortedKeys(p.labelCounts)
}

func (p *LabelListPage) labelsAsSortedListWithCounts() []string {
	data := p.labelsAsSortedList()
	ret := make([]string, len(data))
	for i, v := range data {
		ret[i] = fmt.Sprintf("%s (%d found)", v, p.labelCounts[v])
	}
	return ret
}

func (p *LabelListPage) SelectItem() {
	label := p.cachedResults[p.selectedLine]
	q := new(TicketListPage)
	q.ActiveQuery.Name = ticketListPage.ActiveQuery.Name + "+" + label
	q.ActiveQuery.JQL = ticketListPage.ActiveQuery.JQL + " AND labels = " + label
	currentPage = q
	changePage()
}

func (p *LabelListPage) markActiveLine() {
	for i, v := range p.cachedResults {
		selected := ""
		if i == p.selectedLine {
			selected = "fg-white,bg-blue"
		}
		p.displayLines[i] = fmt.Sprintf("[%-40s -- %d tickets](%s)", v, p.labelCounts[v], selected)
	}
}

func (p *LabelListPage) GoBack() {
	currentPage = ticketListPage
	changePage()
}

func (p *LabelListPage) Update() {
	ls := p.uiList
	p.markActiveLine()
	ls.Items = p.displayLines[p.firstDisplayLine:]
	ui.Render(ls)
	p.statusBar.Update()
	p.commandBar.Update()
}

func (p *LabelListPage) Create() {
	ui.Clear()
	ls := ui.NewList()
	p.uiList = ls
	if p.statusBar == nil {
		p.statusBar = new(StatusBar)
	}
	if p.commandBar == nil {
		p.commandBar = new(CommandBar)
	}
	queryName := p.ActiveQuery.Name
	queryJQL := p.ActiveQuery.JQL
	p.labelCounts = countLabelsFromQuery(queryJQL)
	p.cachedResults = p.labelsAsSortedList()
	p.isPopulated = true
	p.displayLines = make([]string, len(p.cachedResults))
	ls.ItemFgColor = ui.ColorYellow
	ls.BorderLabel = fmt.Sprintf("Label view -- %s: %s", queryName, queryJQL)
	ls.Height = ui.TermHeight() - 2
	ls.Width = ui.TermWidth()
	ls.Y = 0
	p.statusBar.Create()
	p.commandBar.Create()
	p.Update()
}
