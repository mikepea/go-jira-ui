package jiraui

import (
	"fmt"
	"regexp"

	//"gopkg.in/Netflix-Skunkworks/go-jira.v1"
	ui "gopkg.in/gizak/termui.v2"
)

type TicketListPage struct {
	BaseListPage
	CommandBarFragment
	StatusBarFragment
	ActiveQuery     Query
	ActiveSort      Sort
	RankingTicketId string
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
	startLine := (p.uiList.Cursor + n + increment) % n
	for i := startLine; i != p.uiList.Cursor; i = (i + increment + n) % n {
		if s.re.MatchString(p.cachedResults[i]) {
			p.uiList.SetCursorLine(i)
			p.Update()
			break
		}
	}
}

func (p *TicketListPage) ActiveTicketId() string {
	return p.GetSelectedTicketId()
}

func (p *TicketListPage) GetSelectedTicketId() string {
	return findTicketIdInString(p.cachedResults[p.uiList.Cursor])
}

func (p *TicketListPage) MarkItemForRanking() {
	p.RankingTicketId = p.GetSelectedTicketId()
}

func (p *TicketListPage) PreviousLine(n int) {
	if p.RankingTicketId != "" {
		p.uiList.MoveUp(n)
	} else {
		p.uiList.CursorUpLines(n)
	}
}

func (p *TicketListPage) NextLine(n int) {
	if p.RankingTicketId != "" {
		p.uiList.MoveDown(n)
	} else {
		p.uiList.CursorDownLines(n)
	}
}

func (p *TicketListPage) SelectItem() {
	if p.RankingTicketId != "" {
		log.Debugf("Setting Rank for %s", p.RankingTicketId)
		log.Info("TODO: Reenable RANKAFTER")
		//order := jira.RANKAFTER
		var targetId string
		if p.uiList.Cursor == 0 {
			log.Info("TODO: Reenable RANKBEFORE")
			//order = jira.RANKBEFORE
			targetId = findTicketIdInString(p.cachedResults[p.uiList.Cursor+1])
		} else {
			targetId = findTicketIdInString(p.cachedResults[p.uiList.Cursor-1])
		}
		//runJiraCmdRank(p.RankingTicketId, targetId, order)
		log.Infof("TODO: Reenable runJiraCmdRank: %#v", targetId)
		p.RankingTicketId = ""
		p.Refresh()
		return
	}

	if len(p.cachedResults) == 0 {
		return
	}
	q := new(TicketShowPage)
	q.TicketId = p.GetSelectedTicketId()
	previousPages = append(previousPages, currentPage)
	currentPage = q
	q.Create()
	changePage()
}

func (p *TicketListPage) GoBack() {
	if len(previousPages) == 0 {
		currentPage = ticketQueryPage
	} else {
		currentPage, previousPages = previousPages[len(previousPages)-1], previousPages[:len(previousPages)-1]
	}
	changePage()
}

func (p *TicketListPage) EditTicket() {
	//runJiraCmdEdit(p.GetSelectedTicketId())
	log.Infof("TODO: Reenable runJiraCmdEdit")
}

func (p *TicketListPage) Update() {
	ui.Render(p.uiList)
	p.statusBar.Update()
	p.commandBar.Update()
}

func (p *TicketListPage) Refresh() {
	pDeref := &p
	q := *pDeref
	q.cachedResults = make([]string, 0)
	currentPage = q
	changePage()
	q.Create()
}

func (p *TicketListPage) Create() {
	log.Debugf("TicketListPage.Create(): self:        %s (%p)", p.Id(), p)
	log.Debugf("TicketListPage.Create(): currentPage: %s (%p)", currentPage.Id(), currentPage)
	ui.Clear()
	if p.uiList == nil {
		p.uiList = NewScrollableList()
	}
	if p.statusBar == nil {
		p.statusBar = new(StatusBar)
	}
	if p.commandBar == nil {
		p.commandBar = commandBar
	}
	query := p.ActiveQuery.JQL
	if sort := p.ActiveSort.JQL; sort != "" {
		re := regexp.MustCompile(`(?i)\s+ORDER\s+BY.+$`)
		query = re.ReplaceAllString(query, ``) + " " + sort
	}
	if len(p.cachedResults) == 0 {
		p.cachedResults = JiraQueryAsStrings(query, p.ActiveQuery.Template)
	}
	if p.uiList.Cursor >= len(p.cachedResults) {
		p.uiList.Cursor = len(p.cachedResults) - 1
	}
	p.uiList.Items = p.cachedResults
	p.uiList.ItemFgColor = ui.ColorYellow
	p.uiList.BorderLabel = fmt.Sprintf("%s: %s", p.ActiveQuery.Name, p.ActiveQuery.JQL)
	p.uiList.Height = ui.TermHeight() - 2
	p.uiList.Width = ui.TermWidth()
	p.uiList.Y = 0
	p.statusBar.Create()
	p.commandBar.Create()
	p.Update()
}
