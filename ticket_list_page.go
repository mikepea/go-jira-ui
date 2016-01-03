package main

import (
	"fmt"
	ui "github.com/gizak/termui"
	"strings"
)

type TicketListPage struct {
	BaseListPage
}

func (p *TicketListPage) GetSelectedTicketId() string {
	return strings.Split(p.cachedResults[p.selectedLine], " ")[0]
}

func (p *TicketListPage) SelectItem() {
	previousPage = currentPage
	currentPage = &ticketShowPage
	changePage()
}

func (p *TicketListPage) GoBack() {
	previousPage = currentPage
	currentPage = &ticketQueryPage
	changePage()
}

func (p *TicketListPage) Create(opts ...interface{}) {
	ui.Clear()
	var label string
	var queryJQL string
	var queryName string
	if len(opts) > 0 {
		log.Noticef("TicketListPage.Create opts: %s", opts)
		if d, ok := opts[0].(map[string]string); ok {
			label = d["label"]
		}
	}
	ls := ui.NewList()
	p.uiList = ls
	p.selectedLine = 0
	p.firstDisplayLine = 0
	if label != "" {
		queryName = ticketQueryPage.SelectedQuery().Name + "+" + label
		queryJQL = ticketQueryPage.SelectedQuery().JQL + " AND labels = " + label
	} else {
		queryName = ticketQueryPage.SelectedQuery().Name
		queryJQL = ticketQueryPage.SelectedQuery().JQL
	}
	p.cachedResults = JiraQueryAsStrings(queryJQL)
	p.displayLines = make([]string, len(p.cachedResults))
	ls.ItemFgColor = ui.ColorYellow
	ls.BorderLabel = fmt.Sprintf("%s: %s", queryName, queryJQL)
	ls.Height = ui.TermHeight()
	ls.Width = ui.TermWidth()
	ls.Y = 0
	p.Update()
}
