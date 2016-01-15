package jiraui

import (
	"fmt"
	ui "github.com/gizak/termui"
	"regexp"
)

const (
	defaultMaxWrapWidth = 100
)

type TicketShowPage struct {
	BaseListPage
	MaxWrapWidth int
	TicketId     string
	Template     string
	TicketTrail  []*TicketShowPage // previously viewed tickets in drill-down
	WrapWidth    int
	opts         map[string]interface{}
}

func (p *TicketShowPage) SelectItem() {
	selected := p.cachedResults[p.selectedLine]
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
	if p.selectedLine == 0 {
		return
	}
	for i := p.selectedLine - 1; i > 0; i-- {
		if ok, _ := regexp.MatchString(`^\s*$`, p.cachedResults[i]); ok {
			newDisplayLine = i
			break
		}
	}
	p.PreviousLine(p.selectedLine - newDisplayLine)
}

func (p *TicketShowPage) NextPara() {
	newDisplayLine := len(p.cachedResults) - 1
	if p.selectedLine == newDisplayLine {
		return
	}
	for i := p.selectedLine + 1; i < len(p.cachedResults); i++ {
		if ok, _ := regexp.MatchString(`^\s*$`, p.cachedResults[i]); ok {
			newDisplayLine = i
			break
		}
	}
	p.NextLine(newDisplayLine - p.selectedLine)
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

func (p *TicketShowPage) CommentTicket() {
	runJiraCmdComment(p.TicketId)
}

func (p *TicketShowPage) ticketTrailAsString() (trail string) {
	for i := len(p.TicketTrail) - 1; i >= 0; i-- {
		q := *p.TicketTrail[i]
		trail = trail + " <- " + q.Id()
	}
	return trail
}

func (p *TicketShowPage) Create() {
	p.opts = getJiraOpts()
	if p.TicketId == "" {
		p.TicketId = ticketListPage.GetSelectedTicketId()
	}
	if p.MaxWrapWidth == 0 {
		if m := p.opts["max_wrap"]; m != nil {
			p.MaxWrapWidth = int(m.(int64))
		} else {
			p.MaxWrapWidth = defaultMaxWrapWidth
		}
	}
	ui.Clear()
	ls := ui.NewList()
	p.uiList = ls
	p.selectedLine = 0
	p.firstDisplayLine = 0
	if p.Template == "" {
		if templateOpt := p.opts["template"]; templateOpt == nil {
			p.Template = "jira_ui_view"
		} else {
			p.Template = templateOpt.(string)
		}
	}
	if ui.TermWidth()-3 < p.MaxWrapWidth {
		p.WrapWidth = ui.TermWidth() - 3
	} else {
		p.WrapWidth = p.MaxWrapWidth
	}
	if len(p.cachedResults) == 0 {
		p.cachedResults = WrapText(JiraTicketAsStrings(p.TicketId, p.Template), p.WrapWidth)
	}
	p.displayLines = make([]string, len(p.cachedResults))
	ls.ItemFgColor = ui.ColorYellow
	ls.Height = ui.TermHeight()
	ls.Width = ui.TermWidth()
	ls.Border = true
	ls.BorderLabel = fmt.Sprintf("%s %s", p.TicketId, p.ticketTrailAsString())
	ls.Y = 0
	p.Update()
}
