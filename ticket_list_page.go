package jiraui

import (
	"fmt"
	ui "github.com/gizak/termui"
	"regexp"
)

type TicketListPage struct {
	BaseListPage
	CommandBarFragment
	StatusBarFragment
	ActiveQuery Query
	ActiveSort  Sort
}

func (p *TicketListPage) Search() {
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

func (p *TicketListPage) ActiveTicketId() string {
	return p.GetSelectedTicketId()
}

func (p *TicketListPage) GetSelectedTicketId() string {
	return findTicketIdInString(p.cachedResults[p.selectedLine])
}

func (p *TicketListPage) SelectItem() {
	if len(p.cachedResults) == 0 {
		return
	}
	q := new(TicketShowPage)
	q.TicketId = p.GetSelectedTicketId()
	currentPage = q
	q.Create()
	changePage()
}

func (p *TicketListPage) GoBack() {
	currentPage = ticketQueryPage
	changePage()
}

func (p *TicketListPage) EditTicket() {
	runJiraCmdEdit(p.GetSelectedTicketId())
}

func (p *TicketListPage) Update() {
	ls := p.uiList
	p.markActiveLine()
	ls.Items = p.displayLines[p.firstDisplayLine:]
	ui.Render(ls)
	p.statusBar.Update()
	p.commandBar.Update()
}

func (p *TicketListPage) Refresh() {
	pDeref := &p
	q := *pDeref
	q.cachedResults = make([]string, 0)
	ticketListPage = q
	changePage()
	q.Create()
}

func (p *TicketListPage) Create() {
	ui.Clear()
	ls := ui.NewList()
	p.uiList = ls
	if p.statusBar == nil {
		p.statusBar = new(StatusBar)
	}
	if p.commandBar == nil {
		p.commandBar = new(CommandBar)
	}
	query := p.ActiveQuery.JQL
	if sort := p.ActiveSort.JQL; sort != "" {
		re := regexp.MustCompile(`(?i)\s+ORDER\s+BY.+$`)
		query = re.ReplaceAllString(query, ``) + " " + sort
	}
	if len(p.cachedResults) == 0 {
		p.cachedResults = JiraQueryAsStrings(query, p.ActiveQuery.Template)
	}
	p.displayLines = make([]string, len(p.cachedResults))
	ls.ItemFgColor = ui.ColorYellow
	ls.BorderLabel = fmt.Sprintf("%s: %s", p.ActiveQuery.Name, p.ActiveQuery.JQL)
	ls.Height = ui.TermHeight() - 2
	ls.Width = ui.TermWidth()
	ls.Y = 0
	p.statusBar.Create()
	p.commandBar.Create()
	p.Update()
}
