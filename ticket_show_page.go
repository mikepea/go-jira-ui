package main

import (
	"fmt"
	ui "github.com/gizak/termui"
)

type TicketShowPage struct {
	BaseListPage
	TicketId string
}

func (p *TicketShowPage) PreviousPage() {
	p.PreviousLine(p.uiList.Height - 5)
}

func (p *TicketShowPage) NextPage() {
	p.NextLine(p.uiList.Height - 5)
}

func (p *TicketShowPage) SelectItem() {
	newTicketId := findTicketIdInString(p.cachedResults[p.selectedLine])
	log.Noticef("TicketShowPage: Found ticket %s", newTicketId)
	if newTicketId == "" {
		return
	} else if newTicketId == p.TicketId {
		return
	}
	q := new(TicketShowPage)
	q.TicketId = newTicketId
	currentPage = q
	changePage()
}

func (p *TicketShowPage) Id() string {
	return p.TicketId
}

func (p *TicketShowPage) GoBack() {
	previousPage = currentPage
	currentPage = &ticketListPage
	changePage()
}

func (p *TicketShowPage) EditTicket() {
	runJiraCmdEdit(p.TicketId)
}

func (p *TicketShowPage) CommentTicket() {
	runJiraCmdComment(p.TicketId)
}

func (p *TicketShowPage) lastDisplayedLine() int {
	return lastLineDisplayed(p.uiList, p.firstDisplayLine, 5)
}

func (p *TicketShowPage) Create(opts ...interface{}) {
	if p.TicketId == "" {
		p.TicketId = ticketListPage.GetSelectedTicketId()
	}
	ui.Clear()
	ls := ui.NewList()
	p.uiList = ls
	p.selectedLine = 0
	p.firstDisplayLine = 0
	if len(p.cachedResults) == 0 {
		p.cachedResults = JiraTicketAsStrings(p.TicketId)
	}
	p.displayLines = make([]string, len(p.cachedResults))
	ls.ItemFgColor = ui.ColorYellow
	ls.Height = ui.TermHeight()
	ls.Width = ui.TermWidth()
	ls.Overflow = "wrap"
	ls.Border = true
	ls.BorderLabel = fmt.Sprintf("%s", p.TicketId)
	ls.Y = 0
	p.Update()
}
