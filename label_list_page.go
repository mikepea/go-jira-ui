package main

import (
	"fmt"
	ui "github.com/gizak/termui"
)

type LabelListPage struct {
	BaseListPage
	labelCounts map[string]int
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
	previousPage = currentPage
	currentPage = &ticketListPage
	changePage()
}

func (p *LabelListPage) Update() {
	ls := p.uiList
	p.markActiveLine()
	ls.Items = p.displayLines[p.firstDisplayLine:]
	ui.Render(ls)
}

func (p *LabelListPage) Create() {
	ui.Clear()
	ls := ui.NewList()
	p.uiList = ls
	p.selectedLine = 0
	p.firstDisplayLine = 0
	queryName := ticketQueryPage.SelectedQuery().Name
	queryJQL := ticketQueryPage.SelectedQuery().JQL
	p.labelCounts = countLabelsFromQuery(queryJQL)
	p.cachedResults = p.labelsAsSortedList()
	p.displayLines = make([]string, len(p.cachedResults))
	ls.ItemFgColor = ui.ColorYellow
	ls.BorderLabel = fmt.Sprintf("Label view -- %s: %s", queryName, queryJQL)
	ls.Height = ui.TermHeight()
	ls.Width = ui.TermWidth()
	ls.Y = 0
	p.Update()
}
