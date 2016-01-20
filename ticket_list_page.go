package jiraui

import (
	"fmt"
	ui "github.com/gizak/termui"
	"regexp"
)

type Search struct {
	command     string
	directionUp bool
	re          *regexp.Regexp
}

type TicketListPage struct {
	BaseListPage
	CommandBarFragment
	StatusBarFragment
	ActiveQuery  Query
	ActiveSort   Sort
	ActiveSearch Search
}

func (p *TicketListPage) SetSearch(searchCommand string) {
	if len(searchCommand) < 2 {
		// must be '/a' minimum
		return
	}
	direction := []byte(searchCommand)[0]
	regex := string([]byte(searchCommand)[1:])
	s := new(Search)
	s.command = searchCommand
	if direction == '?' {
		s.directionUp = true
	} else if direction == '/' {
		s.directionUp = false
	} else {
		// bad command
		return
	}
	if re, err := regexp.Compile(regex); err != nil {
		// bad regex
		return
	} else {
		s.re = re
		p.ActiveSearch = *s
	}
}

func (p *TicketListPage) Search() {
	return
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

func (p *TicketListPage) CommentTicket() {
	runJiraCmdComment(p.GetSelectedTicketId())
}

func (p *TicketListPage) Update() {
	ls := p.uiList
	p.markActiveLine()
	ls.Items = p.displayLines[p.firstDisplayLine:]
	ui.Render(ls)
	p.statusBar.Update()
	p.commandBar.Update()
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
