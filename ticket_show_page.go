package main

import (
	"fmt"
	ui "github.com/gizak/termui"
)

type TicketShowPage struct {
	BaseListPage
	TicketId    string
	TicketTrail []*TicketShowPage // previously viewed tickets in drill-down
}

func (p *TicketShowPage) PreviousPage() {
	p.PreviousLine(p.uiList.Height - 5)
}

func (p *TicketShowPage) NextPage() {
	p.NextLine(p.uiList.Height - 5)
}

func (p *TicketShowPage) BottomOfPage() {
	p.selectedLine = len(p.cachedResults) - 1
	firstLine := p.selectedLine - (p.uiList.Height - 10)
	if firstLine > 0 {
		p.firstDisplayLine = firstLine
	} else {
		p.firstDisplayLine = 0
	}
}

func (p *TicketShowPage) SelectItem() {
	newTicketId := findTicketIdInString(p.cachedResults[p.selectedLine])
	if newTicketId == "" {
		return
	} else if newTicketId == p.TicketId {
		return
	}
	q := new(TicketShowPage)
	q.TicketId = newTicketId
	q.TicketTrail = append(p.TicketTrail, p)
	currentPage = q
	changePage()
}

func (p *TicketShowPage) Id() string {
	return p.TicketId
}

func (p *TicketShowPage) GoBack() {
	if len(p.TicketTrail) == 0 {
		previousPage = currentPage
		currentPage = &ticketListPage
	} else {
		last := len(p.TicketTrail) - 1
		currentPage = p.TicketTrail[last]
	}
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

func (p *TicketShowPage) ticketTrailAsString() (trail string) {
	for i := len(p.TicketTrail) - 1; i >= 0; i-- {
		q := *p.TicketTrail[i]
		trail = trail + " <- " + q.Id()
	}
	return trail
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
	ls.BorderLabel = fmt.Sprintf("%s %s", p.TicketId, p.ticketTrailAsString())
	ls.Y = 0
	p.Update()
}
