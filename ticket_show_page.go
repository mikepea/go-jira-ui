package main

import (
	ui "github.com/gizak/termui"
)

type TicketShowPage struct {
	BaseListPage
}

func (p *TicketShowPage) PreviousPage() {
	p.PreviousLine(p.uiList.Height - 5)
}

func (p *TicketShowPage) NextPage() {
	p.NextLine(p.uiList.Height - 5)
}

func (p *TicketShowPage) GoBack() {
	previousPage = currentPage
	currentPage = &ticketListPage
	changePage()
}

func (p *TicketShowPage) lastDisplayedLine() int {
	return lastLineDisplayed(p.uiList, p.firstDisplayLine, 5)
}

func (p *TicketShowPage) Create(opts ...interface{}) {
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
