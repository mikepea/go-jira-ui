package main

import (
	"fmt"
	ui "github.com/gizak/termui"
)

type TicketShowPage struct {
	selectedLine     int
	uiList           *ui.List
	displayLines     []string
	cachedResults    []string
	firstDisplayLine int
}

func (p *TicketShowPage) PreviousLine(n int) {
	p.selectedLine = p.selectedLine - n
	if p.selectedLine < 0 {
		p.selectedLine = 0
	}
	if p.selectedLine < p.firstDisplayLine {
		p.firstDisplayLine = p.selectedLine
	}
}

func (p *TicketShowPage) NextLine(n int) {
	if p.selectedLine < len(p.cachedResults)-n {
		p.selectedLine = p.selectedLine + n
	} else {
		p.selectedLine = len(p.cachedResults) - 1
	}
	if p.selectedLine > p.lastDisplayedLine() {
		p.firstDisplayLine = p.firstDisplayLine + n
	}
}

func (p *TicketShowPage) PreviousPage() {
	p.PreviousLine(p.uiList.Height - 5)
}

func (p *TicketShowPage) NextPage() {
	p.NextLine(p.uiList.Height - 5)
}

func (p *TicketShowPage) lastDisplayedLine() int {
	return lastLineDisplayed(p.uiList, p.firstDisplayLine, 5)
}

func (p *TicketShowPage) markActiveLine() {
	for i, v := range p.cachedResults {
		selected := ""
		if i == p.selectedLine {
			selected = "fg-white,bg-blue"
		}
		p.displayLines[i] = fmt.Sprintf("[%s](%s)", v, selected)
	}
}

func (p *TicketShowPage) Update() {
	ls := p.uiList
	p.markActiveLine()
	ls.Items = p.displayLines[p.firstDisplayLine:]
	ui.Render(ls)
}

func (p *TicketShowPage) Create() {
	ui.Clear()
	ls := ui.NewList()
	p.uiList = ls
	p.selectedLine = 0
	p.firstDisplayLine = 0
	ticketId := ticketListPage.GetSelectedTicketId()
	p.cachedResults = JiraTicketAsStrings(ticketId)
	p.displayLines = make([]string, len(p.cachedResults))
	ls.ItemFgColor = ui.ColorYellow
	ls.Height = ui.TermHeight()
	ls.Width = ui.TermWidth()
	ls.Overflow = "wrap"
	ls.Border = false
	ls.Y = 0
	p.Update()
}
