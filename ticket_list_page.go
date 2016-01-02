package main

import (
	"fmt"
	ui "github.com/gizak/termui"
	"strings"
)

type TicketListPage struct {
	selectedLine     int
	uiList           *ui.List
	displayLines     []string
	cachedResults    []string
	firstDisplayLine int
}

func (p *TicketListPage) PreviousLine(n int) {
	p.selectedLine = p.selectedLine - n
	if p.selectedLine < 0 {
		p.selectedLine = 0
	}
	if p.selectedLine < p.firstDisplayLine {
		p.firstDisplayLine = p.selectedLine
	}
}

func (p *TicketListPage) NextLine(n int) {
	if p.selectedLine < len(p.cachedResults)-n {
		p.selectedLine = p.selectedLine + n
	} else {
		p.selectedLine = len(p.cachedResults) - 1
	}
	if p.selectedLine > p.lastDisplayedLine() {
		p.firstDisplayLine = p.firstDisplayLine + n
	}
}

func (p *TicketListPage) PreviousPage() {
	p.PreviousLine(p.uiList.Height - 2)
}

func (p *TicketListPage) NextPage() {
	p.NextLine(p.uiList.Height - 2)
}

func (p *TicketListPage) lastDisplayedLine() int {
	return lastLineDisplayed(p.uiList, p.firstDisplayLine, 3)
}

func (p *TicketListPage) markActiveLine() {
	for i, v := range p.cachedResults {
		selected := ""
		if i == p.selectedLine {
			selected = "fg-white,bg-blue"
		}
		p.displayLines[i] = fmt.Sprintf("[%s](%s)", v, selected)
	}
}

func (p *TicketListPage) GetSelectedTicketId() string {
	return strings.Split(p.cachedResults[p.selectedLine], " ")[0]
}

func (p *TicketListPage) Update() {
	ls := p.uiList
	p.markActiveLine()
	ls.Items = p.displayLines[p.firstDisplayLine:]
	ui.Render(ls)
}

func (p *TicketListPage) Create() {
	ui.Clear()
	ls := ui.NewList()
	p.uiList = ls
	p.selectedLine = 0
	p.firstDisplayLine = 0
	queryName := ticketQueryPage.SelectedQuery().Name
	queryJQL := ticketQueryPage.SelectedQuery().JQL
	p.cachedResults = JiraQueryAsStrings(queryJQL)
	p.displayLines = make([]string, len(p.cachedResults))
	ls.ItemFgColor = ui.ColorYellow
	ls.BorderLabel = fmt.Sprintf("%s: %s", queryName, queryJQL)
	ls.Height = ui.TermHeight()
	ls.Width = ui.TermWidth()
	ls.Y = 0
	p.Update()
}
