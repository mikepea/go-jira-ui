package jiraui

import (
	"fmt"
	"regexp"

	ui "github.com/gizak/termui"
)

const (
	defaultMaxWrapWidth = 100
)

type TicketShowPage struct {
	BaseListPage
	CommandBarFragment
	StatusBarFragment
	MaxWrapWidth uint
	TicketId     string
	Template     string
	apiBody      interface{}
	TicketTrail  []*TicketShowPage // previously viewed tickets in drill-down
	WrapWidth    uint
	opts         map[string]interface{}
}

func (p *TicketShowPage) Search() {
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
			p.SetSelectedLine(i)
			p.Update()
			break
		}
	}
}

func (p *TicketShowPage) SelectItem() {
	selected := p.cachedResults[p.uiList.Cursor]
	if ok, _ := regexp.MatchString(`^epic_links:`, selected); ok {
		q := new(TicketListPage)
		q.ActiveQuery.Name = fmt.Sprintf("Open Tasks in Epic %s", p.TicketId)
		q.ActiveQuery.JQL = fmt.Sprintf("\"Epic Link\" = %s AND resolution = Unresolved", p.TicketId)
		currentPage = q
	} else {
		newTicketId := findTicketIdInString(selected)
		if newTicketId == "" {
			return
		} else if newTicketId == p.TicketId {
			return
		}
		q := new(TicketShowPage)
		q.TicketId = newTicketId
		q.TicketTrail = append(p.TicketTrail, p)
		currentPage = q
	}
	changePage()
}

func (p *TicketShowPage) Id() string {
	return p.TicketId
}

func (p *TicketShowPage) PreviousPara() {
	newDisplayLine := 0
	sl := p.uiList.Cursor
	if sl == 0 {
		return
	}
	for i := sl - 1; i > 0; i-- {
		if ok, _ := regexp.MatchString(`^\s*$`, p.cachedResults[i]); ok {
			newDisplayLine = i
			break
		}
	}
	p.PreviousLine(sl - newDisplayLine)
}

func (p *TicketShowPage) NextPara() {
	newDisplayLine := len(p.cachedResults) - 1
	sl := p.uiList.Cursor
	if sl == newDisplayLine {
		return
	}
	for i := sl + 1; i < len(p.cachedResults); i++ {
		if ok, _ := regexp.MatchString(`^\s*$`, p.cachedResults[i]); ok {
			newDisplayLine = i
			break
		}
	}
	p.NextLine(newDisplayLine - sl)
}

func (p *TicketShowPage) GoBack() {
	if len(p.TicketTrail) == 0 {
		if ticketListPage != nil {
			currentPage = ticketListPage
		} else {
			currentPage = ticketQueryPage
		}
	} else {
		last := len(p.TicketTrail) - 1
		currentPage = p.TicketTrail[last]
	}
	changePage()
}

func (p *TicketShowPage) EditTicket() {
	runJiraCmdEdit(p.TicketId)
}

func (p *TicketShowPage) ActiveTicketId() string {
	return p.TicketId
}

func (p *TicketShowPage) ticketTrailAsString() (trail string) {
	for i := len(p.TicketTrail) - 1; i >= 0; i-- {
		q := *p.TicketTrail[i]
		trail = trail + " <- " + q.Id()
	}
	return trail
}

func (p *TicketShowPage) Refresh() {
	pDeref := &p
	q := *pDeref
	q.cachedResults = make([]string, 0)
	q.apiBody = nil
	currentPage = q
	changePage()
	q.Create()
}

func (p *TicketShowPage) Update() {
	ls := p.uiList
	log.Debugf("TicketShowPage.Update(): self:        %s (%p), ls: (%p)", p.Id(), p, ls)
	p.markActiveLine()
	ls.Items = p.displayLines[p.firstDisplayLine:]
	ui.Render(ls)
	p.statusBar.Update()
	p.commandBar.Update()
}

func (p *TicketShowPage) Create() {
	log.Debugf("TicketShowPage.Create(): self:        %s (%p)", p.Id(), p)
	log.Debugf("TicketShowPage.Create(): currentPage: %s (%p)", currentPage.Id(), currentPage)
	p.opts = getJiraOpts()
	if p.TicketId == "" {
		p.TicketId = ticketListPage.GetSelectedTicketId()
	}
	if p.MaxWrapWidth == 0 {
		if m := p.opts["max_wrap"]; m != nil {
			p.MaxWrapWidth = uint(m.(int64))
		} else {
			p.MaxWrapWidth = defaultMaxWrapWidth
		}
	}
	ui.Clear()
	ls := NewScrollableList()
	if p.statusBar == nil {
		p.statusBar = new(StatusBar)
	}
	if p.commandBar == nil {
		p.commandBar = commandBar
	}
	p.uiList = ls
	if p.Template == "" {
		if templateOpt := p.opts["template"]; templateOpt == nil {
			p.Template = "jira_ui_view"
		} else {
			p.Template = templateOpt.(string)
		}
	}
	innerWidth := uint(ui.TermWidth()) - 3
	if innerWidth < p.MaxWrapWidth {
		p.WrapWidth = innerWidth
	} else {
		p.WrapWidth = p.MaxWrapWidth
	}
	if p.apiBody == nil {
		p.apiBody, _ = FetchJiraTicket(p.TicketId)
	}
	p.cachedResults = WrapText(JiraTicketAsStrings(p.apiBody, p.Template), p.WrapWidth)
	p.displayLines = make([]string, len(p.cachedResults))
	if p.uiList.Cursor >= len(p.cachedResults) {
		p.uiList.Cursor = len(p.cachedResults) - 1
	}
	ls.ItemFgColor = ui.ColorYellow
	ls.Height = ui.TermHeight() - 2
	ls.Width = ui.TermWidth()
	ls.Border = true
	ls.BorderLabel = fmt.Sprintf("%s %s", p.TicketId, p.ticketTrailAsString())
	ls.Y = 0
	p.statusBar.Create()
	p.commandBar.Create()
	p.Update()
}
