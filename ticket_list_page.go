package jiraui

import (
	"fmt"
	ui "github.com/gizak/termui"
	"strings"
)

type TicketListPage struct {
	BaseListPage
	ActiveQuery Query
}

func (p *TicketListPage) GetSelectedTicketId() string {
	return strings.Split(p.cachedResults[p.selectedLine], " ")[0]
}

func (p *TicketListPage) SelectItem() {
	if len(p.cachedResults) == 0 {
		return
	}
	q := new(TicketShowPage)
	q.TicketId = p.GetSelectedTicketId()
	q.Create()
	currentPage = q
	changePage()
}

func (p *TicketListPage) GoBack() {
	currentPage = ticketQueryPage
	changePage()
}

func (p *TicketListPage) EditTicket() {
	runJiraCmdEdit(p.GetSelectedTicketId())
}

func (p *TicketListPage) CommentTicket() {
	runJiraCmdComment(p.GetSelectedTicketId())
}

func (p *TicketListPage) Create() {
	ui.Clear()
	ls := ui.NewList()
	p.uiList = ls
	p.selectedLine = 0
	p.firstDisplayLine = 0
	if len(p.cachedResults) == 0 {
		p.cachedResults = JiraQueryAsStrings(p.ActiveQuery.JQL)
	}
	p.displayLines = make([]string, len(p.cachedResults))
	ls.ItemFgColor = ui.ColorYellow
	ls.BorderLabel = fmt.Sprintf("%s: %s", p.ActiveQuery.Name, p.ActiveQuery.JQL)
	ls.Height = ui.TermHeight()
	ls.Width = ui.TermWidth()
	ls.Y = 0
	p.Update()
}
